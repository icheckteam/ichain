package shipping

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

var msgCdc = wire.NewCodec()

// RegisterWire registers concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(CreateOrderMsg{}, "shipping/CreateOrderMsg", nil)
	cdc.RegisterConcrete(ConfirmOrderMsg{}, "shipping/ConfirmOrderMsg", nil)
	cdc.RegisterConcrete(CompleteOrderMsg{}, "shipping/CompleteOrderMsg", nil)
	cdc.RegisterConcrete(CancelOrderMsg{}, "shipping/CancelOrderMsg", nil)
}

func init() {
	RegisterWire(msgCdc)
}
