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

func (k Keeper) setAsset(ctx sdk.Context, asset Asset) {
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey([]byte(asset.ID))

	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(asset)
	if err != nil {
		panic(err)
	}

	store.Set(assetKey, bz)
}

// Has asset
func (k Keeper) Has(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey([]byte(id))
	return store.Has(assetKey)
}

// GetAsset get asset by IDS
func (k Keeper) GetAsset(ctx sdk.Context, uid string) *Asset {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetAssetKey([]byte(uid)))
	asset := &Asset{}

	// marshal the record and add to the state
	if err := k.cdc.UnmarshalBinary(b, asset); err != nil {
		panic(err)
	}
	return asset
}

func (k Keeper) Transfer(ctx sdk.Context, msg TransferMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if asset.IsOwner(msg.Sender) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Sender))
	}
	return nil
}

// UpdateAttribute ...
func (k Keeper) UpdateAttribute(ctx sdk.Context, msg UpdateAttrMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if asset.IsOwner(msg.Sender) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Sender))
	}

	asset.Attributes[msg.AttributeName] = msg.AttributeValue
	k.setAsset(ctx, *asset)
	return nil
}

// AddQuantity ...
func (k Keeper) AddQuantity(ctx sdk.Context, msg AddQuantityMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if asset.IsOwner(msg.Sender) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Sender))
	}

	asset.Quantity += msg.Quantity
	k.setAsset(ctx, *asset)
	return nil
}

// SubtractQuantity ...
func (k Keeper) SubtractQuantity(ctx sdk.Context, msg SubtractQuantityMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if asset.IsOwner(msg.Sender) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Sender))
	}

	asset.Quantity -= msg.Quantity
	k.setAsset(ctx, *asset)
	return nil
}