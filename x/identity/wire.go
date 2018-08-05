package identity

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// RegisterWire Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSetCerts{}, "identity/SetCerts", nil)
	cdc.RegisterConcrete(MsgSetTrust{}, "identity/SetTrust", nil)
	cdc.RegisterConcrete(MsgAddOwner{}, "identity/MsgAddOwner", nil)
	cdc.RegisterConcrete(MsgReg{}, "identity/MsgReg", nil)
	cdc.RegisterConcrete(MsgDelOwner{}, "identity/MsgDelOwner", nil)
}

// MsgCdc generic sealed codec to be used throughout sdk
var MsgCdc *wire.Codec

func init() {
	cdc := wire.NewCodec()
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	MsgCdc = cdc
}
