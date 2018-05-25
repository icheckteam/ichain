package identity

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateClaim{}, "ichain/MsgCreateClaim", nil)
	cdc.RegisterConcrete(MsgRevokeClaim{}, "ichain/MsgRevokeClaim", nil)

}
