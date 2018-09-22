package gs1

import "github.com/cosmos/cosmos-sdk/wire"

var msgCdc = wire.NewCodec()

// RegisterWire Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "gs1/MsgSend", nil)
}

func init() {
	RegisterWire(msgCdc)
}
