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
		"asset_id", []byte(msg.AssetID),
		"sender", []byte(msg.Sender.String()),
	)
	var assetRoot string
	if len(msg.Parent) > 0 {
		// get asset to check quantity and check authorized
		parent, found := k.GetAsset(ctx, msg.Parent)
		if !found {
			return nil, ErrAssetNotFound(msg.Parent)
		}
		if parent.Final {
			return nil, ErrAssetAlreadyFinal(parent.ID)
		}

		if !parent.IsOwner(msg.Sender) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("Address {%v} not unauthorized to create asset", msg.Sender))
		}

		if parent.Quantity < msg.Quantity {
			return nil, ErrInvalidAssetQuantity(parent.ID)
		}
		parent.Quantity -= msg.Quantity

		if len(parent.Root) != 0 && parent.Quantity == 0 {
			parent.Final = true
		}

		// save parent asset to store
		k.setAsset(ctx, parent)
		if len(parent.Root) > 0 {
			assetRoot = parent.Root
		} else {
			assetRoot = parent.ID
		}

		tags = tags.AppendTag("asset_id", []byte(parent.ID))
	}

	asset := Asset{
		ID:        msg.AssetID,
		Name:      msg.Name,
		Owner:     msg.Sender,
		Quantity:  msg.Quantity,
		Root:      assetRoot,
		Parent:    msg.Parent,
		Final:     false,
		Precision: msg.Precision,
		Height:    ctx.BlockHeight(),
	}

	if len(msg.Properties) > 0 {
		asset.Properties = msg.Properties.Sort()
	}

	// update asset info
	k.setAsset(ctx, asset)
	k.setAssetByAccountIndex(ctx, asset)
	return tags, nil
}

// set the main record holding asset details
func (k Keeper) setAsset(ctx sdk.Context, asset Asset) {
	ctx.GasMeter().ConsumeGas(costSetAsset, "setAsset")
	store := ctx.KVStore(k.storeKey)
	// set main store
	bz := k.cdc.MustMarshalBinary(asset)
	store.Set(GetAssetKey(asset.ID), bz)
}

func (k Keeper) setAssetByAccountIndex(ctx sdk.Context, asset Asset) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(asset.ID)
	store.Set(GetAccountAssetKey(asset.Owner, asset.ID), bz)
}

func (k Keeper) removeAssetByAccountIndex(ctx sdk.Context, asset Asset) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetAccountAssetKey(asset.Owner, asset.ID))
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
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}

	if len(asset.Parent) != 0 {
		return nil, ErrInvalidAssetRoot(asset.ID)
	}

	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add", msg.Sender))
	}
	asset.Quantity += msg.Quantity
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"sender", []byte(msg.Sender.String()),
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
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}

	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to sbutract", msg.Sender))
	}

	if asset.Quantity < msg.Quantity {
		return nil, ErrInvalidAssetQuantity(asset.ID)
	}
	asset.Quantity -= msg.Quantity
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"sender", []byte(msg.Sender.String()),
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
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}

	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to finalize", msg.Sender))
	}

	asset.Final = true
	k.removeAssetByAccountIndex(ctx, asset)
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(msg.AssetID),
		"sender", []byte(msg.Sender.String()),
	)
	return tags, nil
}

// Transfer transfer asset
func (k Keeper) Transfer(ctx sdk.Context, msg MsgTransfer) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costTransfer, "transferAsset")
	assets := []Asset{}
	for _, a := range msg.Assets {
		asset, found := k.GetAsset(ctx, a)
		if !found {
			return nil, ErrAssetNotFound(a)
		}
		if asset.Final {
			return nil, ErrAssetAlreadyFinal(asset.ID)
		}
		if !asset.IsOwner(msg.Sender) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to send", msg.Sender))
		}
		assets = append(assets, asset)
	}
	tags := sdk.NewTags(
		"sender", []byte(msg.Sender.String()),
		"recipient", []byte(msg.Recipient.String()),
	)
	for _, asset := range assets {
		k.removeAssetByAccountIndex(ctx, asset)
		// change ownership
		asset.Owner = msg.Recipient
		// clear reporter
		asset.Reporters = nil
		k.setAsset(ctx, asset)
		k.setAssetByAccountIndex(ctx, asset)
		// add asset tag
		tags = tags.AppendTag("asset_id", []byte(asset.ID))
	}
	return tags, nil
}
