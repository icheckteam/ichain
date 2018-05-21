package warranty

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateClaim{}, "ichain-warranty/MsgCreateClaim", nil)
	cdc.RegisterConcrete(MsgCreateContract{}, "ichain-warranty/MsgCreateContract", nil)
	cdc.RegisterConcrete(MsgProcessClaim{}, "ichain-warranty/MsgProcessClaim", nil)
}
