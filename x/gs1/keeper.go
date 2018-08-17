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

// CreateRecord create new record
func (k Keeper) CreateRecord(ctx sdk.Context, msg MsgSend) (sdk.Tags, sdk.Error) {
	return nil, nil
}
