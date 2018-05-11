package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/bank"
)

// Keeper ...
type Keeper struct {
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc      *wire.Codec
	bank     bank.CoinKeeper
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, bank bank.CoinKeeper) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
		bank:     bank,
	}
}

// Register register new asset
func (k Keeper) RegisterAsset(ctx sdk.Context, asset Asset) sdk.Error {
	if asset.ID == "icc" {
		return InvalidTransaction("Asset already exists")
	}

	if k.Has(ctx, asset.ID) {
		return InvalidTransaction("Asset already exists")
	}
	// update asset info
	k.setAsset(ctx, asset)

	// add coin ...
	_, err := k.bank.AddCoins(ctx, asset.Issuer, sdk.Coins{
		sdk.Coin{Denom: asset.ID, Amount: asset.Quantity},
	})

	return err
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
		return nil
	}
	return asset
}

// UpdateAttribute ...
func (k Keeper) UpdateAttribute(ctx sdk.Context, msg UpdateAttrMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if !asset.IsOwner(msg.Issuer) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}

	asset.Attributes[msg.Name] = msg.Value
	k.setAsset(ctx, *asset)
	return nil
}

// AddQuantity ...
func (k Keeper) AddQuantity(ctx sdk.Context, msg AddQuantityMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if !asset.IsOwner(msg.Issuer) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}

	asset.Quantity += msg.Quantity
	k.setAsset(ctx, *asset)
	// add coin ...
	_, err := k.bank.AddCoins(ctx, asset.Issuer, sdk.Coins{
		sdk.Coin{Denom: asset.ID, Amount: msg.Quantity},
	})
	return err
}

// SubtractQuantity ...
func (k Keeper) SubtractQuantity(ctx sdk.Context, msg SubtractQuantityMsg) sdk.Error {
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return ErrUnknownAsset("Asset not found")
	}
	if !asset.IsOwner(msg.Issuer) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}

	// add coin ...
	_, err := k.bank.SubtractCoins(ctx, asset.Issuer, sdk.Coins{
		sdk.Coin{Denom: asset.ID, Amount: msg.Quantity},
	})

	if err != nil {
		return err
	}
	asset.Quantity -= msg.Quantity
	k.setAsset(ctx, *asset)
	return err
}

// ------------------------------------------
// AddQuantity Tests
