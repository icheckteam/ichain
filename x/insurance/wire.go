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

// generic sealed codec to be used throughout sdk
var MsgCdc *wire.Codec

func init() {
	cdc := wire.NewCodec()
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	MsgCdc = cdc
	//MsgCdc = cdc.Seal() //TODO use when upgraded to go-amino 0.9.10
}
