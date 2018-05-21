package shipping

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// RegisterWire registers concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(CreateOrderMsg{}, "ichain/CreateOrderMsg", nil)
	cdc.RegisterConcrete(ConfirmOrderMsg{}, "ichain/ConfirmOrderMsg", nil)
	cdc.RegisterConcrete(CompleteOrderMsg{}, "ichain/CompleteOrderMsg", nil)
	cdc.RegisterConcrete(CancelOrderMsg{}, "ichain/CancelOrderMsg", nil)
}
