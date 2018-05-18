package shipping

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	coin "github.com/icheckteam/ichain/x/bank"
)

// Keeper manages shipping orders
type Keeper struct {
	storeKey   sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc        *wire.Codec
	coinKeeper coin.Keeper
}
