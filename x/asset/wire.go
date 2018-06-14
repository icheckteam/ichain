package asset

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateAsset{}, "ichain/MsgCreateAsset", nil)
	cdc.RegisterConcrete(AddQuantityMsg{}, "ichain/AddQuantityMsg", nil)
	cdc.RegisterConcrete(SubtractQuantityMsg{}, "ichain/SubtractQuantityMsg", nil)
	cdc.RegisterConcrete(MsgUpdatePropertipes{}, "ichain/MsgUpdatePropertipes", nil)
	cdc.RegisterConcrete(MsgSend{}, "ichain/SendAsset", nil)
}
