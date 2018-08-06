package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
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

// CreateAsset create new an asset
func (k Keeper) CreateAsset(ctx sdk.Context, msg MsgCreateAsset) (sdk.Tags, sdk.Error) {
	if k.has(ctx, msg.AssetID) {
		return nil, ErrInvalidTransaction(fmt.Sprintf("Asset {%s} already exists", msg.AssetID))
	}

	var parent Asset
	var found bool
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
	}

	if len(msg.Parent) > 0 {
		// get asset to check quantity and check authorized
		parent, found = k.GetAsset(ctx, msg.Parent)
		if !found {
			return nil, ErrAssetNotFound(msg.Parent)
		}
		if parent.Final {
			return nil, ErrAssetAlreadyFinal(parent.ID)
		}

		if !parent.IsOwner(msg.Sender) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to revoke", msg.Sender))
		}

		if parent.Quantity.LT(msg.Quantity) {
			return nil, ErrInvalidAssetQuantity(parent.ID)
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
	}

	if msg.Parent != "" {
		// clone data
		k.setAsset(ctx, parent)
	}

	if len(msg.Properties) > 0 {
		k.SetProperties(ctx, msg.AssetID, msg.Properties)
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
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(asset)
	store.Set(GetAssetKey(asset.ID), bz)
}

func (k Keeper) setAssetByAccountIndex(ctx sdk.Context, assetID string, recipient sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetAccountAssetKey(recipient, assetID), []byte{})
}

func (k Keeper) setAssetByParentIndex(ctx sdk.Context, asset Asset) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(asset.ID)
	store.Set(GetAssetChildrenKey(asset.Parent, asset.ID), bz)
}

// SetAsset set the main record holding asset details
func (k Keeper) SetAsset(ctx sdk.Context, asset Asset) {
	k.setAsset(ctx, asset)
}

// Has asset
func (k Keeper) has(ctx sdk.Context, assetID string) bool {
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey(assetID)
	return store.Has(assetKey)
}

// GetAsset get asset by IDS
func (k Keeper) GetAsset(ctx sdk.Context, assetID string) (asset Asset, found bool) {
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
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}

	if asset.Root != "" || !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add", msg.Sender))
	}

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
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}

	if asset.Quantity.LT(msg.Quantity) {
		return nil, ErrInvalidAssetQuantity(asset.ID)
	}

	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}

	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to revoke", msg.Sender))
	}

	asset.Quantity = asset.Quantity.Sub(msg.Quantity)
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}

// Finalize ...
func (k Keeper) Finalize(ctx sdk.Context, msg MsgFinalize) (sdk.Tags, sdk.Error) {
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}
	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to revoke", msg.Sender))
	}
	asset.Final = true
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagAsset, []byte(msg.AssetID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}
