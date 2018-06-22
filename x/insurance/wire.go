package insurance

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateClaim{}, "insurance/MsgCreateClaim", nil)
	cdc.RegisterConcrete(MsgCreateContract{}, "insurance/MsgCreateContract", nil)
	cdc.RegisterConcrete(MsgProcessClaim{}, "insurance/MsgProcessClaim", nil)
}
