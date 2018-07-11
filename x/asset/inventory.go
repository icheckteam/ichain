package asset

import sdk "github.com/cosmos/cosmos-sdk/types"

// Inventory
// ......................................
func (k Keeper) setInventory(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(amount)
	store.Set(GetInventoryKey(addr, amount.Denom), bz)
}

func (k Keeper) getInventory(ctx sdk.Context, addr sdk.AccAddress, assetID string) (amount sdk.Coin, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetInventoryKey(addr, assetID))
	if b == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(b, &amount)
	found = true
	return
}

func (k Keeper) deleteInventory(ctx sdk.Context, addr sdk.AccAddress, assetID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetInventoryKey(addr, assetID))
}

func (k Keeper) addInventory(ctx sdk.Context, addr sdk.AccAddress, addAmount sdk.Coin) {
	amount, found := k.getInventory(ctx, addr, addAmount.Denom)

	if !found {
		amount = sdk.NewCoin(addAmount.Denom, 0)
	}

	amount = amount.Plus(addAmount)
	k.setInventory(ctx, addr, amount)
}

func (k Keeper) subtractInventory(ctx sdk.Context, addr sdk.AccAddress, subAmount sdk.Coin) {
	amount, found := k.getInventory(ctx, addr, subAmount.Denom)
	if !found {
		return
	}

	amount = amount.Minus(subAmount)

	if amount.IsZero() {
		k.deleteInventory(ctx, addr, subAmount.Denom)
		return
	}

	k.setInventory(ctx, addr, amount)
}
