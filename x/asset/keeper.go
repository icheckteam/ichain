package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper ...
type Keeper struct {
	ck bank.CoinKeeper

	storeKey          sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc               *wire.Codec
	recordIndexNumber int
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey:          key,
		recordIndexNumber: 0,
		ck:                bankKeeper,
		cdc:               cdc,
	}
}

func (k Keeper) createAsset(ctx sdk.Context, asset Asset) {
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
func (k Keeper) GetAsset(ctx sdk.Context, uid string) Asset {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetAssetKey([]byte(uid)))
	asset := Asset{}

	// marshal the record and add to the state
	if err := k.cdc.UnmarshalBinary(b, &asset); err != nil {
		panic(err)
	}
	return asset
}

// Transfer change owner
func (k Keeper) Transfer(ctx sdk.Context, fromAddress sdk.Address, toAddress sdk.Address, uid string) sdk.Error {
	asset := k.GetAsset(ctx, uid)
	if asset.ID == "" {
		return ErrUnknownAsset("Asset not found")
	}
	return nil
}
