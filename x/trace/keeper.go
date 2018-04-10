package trace

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper ...
type Keeper struct {
	ck bank.CoinKeeper

	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.

	recordIndexNumber uint
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper) Keeper {
	return Keeper{
		storeKey:          key,
		recordIndexNumber: 0,
		ck:                bankKeeper,
	}
}

func (k Keeper) createRecord(ctx sdk.Context, record Record) {

}
