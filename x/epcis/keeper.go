package epcis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Keeper ...
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
}

// RegisterActor register new an actor ...
func (k Keeper) RegisterActor() {

}
