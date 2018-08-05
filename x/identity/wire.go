package identity

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// RegisterWire Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSetCerts{}, "identity/SetCerts", nil)
	cdc.RegisterConcrete(MsgSetTrust{}, "identity/SetTrust", nil)
}

// MsgCdc generic sealed codec to be used throughout sdk
var MsgCdc *wire.Codec

func init() {
	cdc := wire.NewCodec()
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	MsgCdc = cdc
}
