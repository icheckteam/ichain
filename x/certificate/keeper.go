package trace

import (
	"github.com/cosmos-sdk/x/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper ...
type Keeper struct {
	ck bank.CoinKeeper

	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper) Keeper {
	return Keeper{bankKeeper, key}
}

func (k Keeper) createRecord(ctx sdk.Context, record Record) {

}
