package identity

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSetCerts{}, "identity/SetCerts", nil)
	cdc.RegisterConcrete(MsgSetTrust{}, "identity/SetTrust", nil)
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
