package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	costGetAsset         sdk.Gas = 10
	costCreateAsset      sdk.Gas = 10
	costSetAsset         sdk.Gas = 10
	costHasAsset         sdk.Gas = 10
	costSubtractQuantity sdk.Gas = 10
	costAddQuantity      sdk.Gas = 10
	costUpdateProperties sdk.Gas = 10
	costCreateReporter   sdk.Gas = 10
	costRevokeReporter   sdk.Gas = 10
	costAddMaterials     sdk.Gas = 10
	costFinalize         sdk.Gas = 10
	costTransfer         sdk.Gas = 10
)

// Keeper ...
type Keeper struct {
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc      *wire.Codec
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

// Register register new asset
func (k Keeper) CreateAsset(ctx sdk.Context, msg MsgCreateAsset) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costCreateAsset, "createAsset")
	if k.has(ctx, msg.AssetID) {
		return nil, ErrInvalidTransaction(fmt.Sprintf("Asset {%s} already exists", msg.AssetID))
	}

	tags := sdk.NewTags(
		TagAsset, []byte(msg.AssetID),
		TagSender, []byte(msg.Sender.String()),
	)

	newAsset := Asset{
		ID:       msg.AssetID,
		Name:     msg.Name,
		Owner:    msg.Sender,
		Quantity: msg.Quantity,
		Parent:   msg.Parent,
		Final:    false,
		Height:   ctx.BlockHeight(),
		Created:  ctx.BlockHeader().Time,
		Unit:     msg.Unit,
	}

	if len(msg.Parent) > 0 {
		// get asset to check quantity and check authorized
		parent, found := k.GetAsset(ctx, msg.Parent)
		if !found {
			return nil, ErrAssetNotFound(msg.Parent)
		}
		if err := parent.ValidateAddChildren(msg.Sender, msg.Quantity); err != nil {
			return nil, err
		}
		parent.Quantity = parent.Quantity.Sub(msg.Quantity)

		if len(parent.Root) != 0 && parent.Quantity.IsZero() {
			parent.Final = true
		}

		if len(parent.Root) > 0 {
			newAsset.Root = parent.Root
		} else {
			newAsset.Root = parent.ID
		}

		tags = tags.AppendTag(TagAsset, []byte(parent.ID))
		newAsset.Unit = parent.Unit
		k.setAsset(ctx, parent)
	}

	if len(msg.Properties) > 0 {
		newAsset.Properties = msg.Properties.Sort()
	}

	// update asset info
	k.SetAsset(ctx, newAsset)
	k.setAssetByAccountIndex(ctx, newAsset.ID, newAsset.Owner)

	if len(newAsset.Parent) > 0 {
		// index by parent
		k.setAssetByParentIndex(ctx, newAsset)
	}

	return tags, nil
}

// set the main record holding asset details
func (k Keeper) setAsset(ctx sdk.Context, asset Asset) {
	ctx.GasMeter().ConsumeGas(costSetAsset, "setAsset")
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(asset)
	store.Set(GetAssetKey(asset.ID), bz)
}

func (k Keeper) setAssetByAccountIndex(ctx sdk.Context, assetID string, recipient sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(assetID)
	store.Set(GetAccountAssetKey(recipient, assetID), bz)
}

func (k Keeper) removeAssetByAccountIndex(ctx sdk.Context, assetID string, recipient sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetAccountAssetKey(recipient, assetID))
}

func (k Keeper) setAssetByParentIndex(ctx sdk.Context, asset Asset) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(asset.ID)
	store.Set(GetAssetChildrenKey(asset.Parent, asset.ID), bz)
}

// set the main record holding asset details
func (k Keeper) SetAsset(ctx sdk.Context, asset Asset) {
	k.setAsset(ctx, asset)
}

// Has asset
func (k Keeper) has(ctx sdk.Context, assetID string) bool {
	ctx.GasMeter().ConsumeGas(costHasAsset, "hasAsset")
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey(assetID)
	return store.Has(assetKey)
}

// GetAsset get asset by IDS
func (k Keeper) GetAsset(ctx sdk.Context, assetID string) (asset Asset, found bool) {
	ctx.GasMeter().ConsumeGas(costGetAsset, "getAsset")
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetAssetKey(assetID))
	if b == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinary(b, &asset)
	return asset, true
}

// AddQuantity ...
func (k Keeper) AddQuantity(ctx sdk.Context, msg MsgAddQuantity) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddQuantity, "addQuantity")
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}

	if err := asset.ValidateAddQuantity(msg.Sender); err != nil {
		return nil, err
	}

	// add inventory
	k.addInventory(ctx, asset.Owner, sdk.Coin{
		Denom:  asset.GetRoot(),
		Amount: msg.Quantity,
	})

	asset.Quantity = asset.Quantity.Add(msg.Quantity)
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}

// SubtractQuantity  subtract quantity of the asset
func (k Keeper) SubtractQuantity(ctx sdk.Context, msg MsgSubtractQuantity) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costSubtractQuantity, "subtractQuantity")
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}

	if err := asset.ValidateSubtractQuantity(msg.Sender, msg.Quantity); err != nil {
		return nil, err
	}

	// subtract inventory
	k.subtractInventory(ctx, asset.Owner, sdk.Coin{
		Denom:  asset.GetRoot(),
		Amount: msg.Quantity,
	})

	asset.Quantity = asset.Quantity.Sub(msg.Quantity)
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagRecipient, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}

// Send ...
func (k Keeper) Finalize(ctx sdk.Context, msg MsgFinalize) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costFinalize, "finalizeAsset")
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}
	if err := asset.ValidateFinalize(msg.Sender); err != nil {
		return nil, err
	}
	asset.Final = true
	k.removeAssetByAccountIndex(ctx, asset.ID, asset.Owner)

	// delete all index for reporter
	for _, reporter := range asset.Reporters {
		k.removeAssetByAccountIndex(ctx, asset.ID, reporter.Addr)
	}

	// subtract inventory
	k.subtractInventory(ctx, asset.Owner, sdk.Coin{
		Denom:  asset.GetRoot(),
		Amount: asset.Quantity,
	})

	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagAsset, []byte(msg.AssetID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}
