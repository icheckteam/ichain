package asset

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(RegisterMsg{}, "ichain/RegisterMsg", nil)
	cdc.RegisterConcrete(AddQuantityMsg{}, "ichain/AddQuantityMsg", nil)
	cdc.RegisterConcrete(SubtractQuantityMsg{}, "ichain/SubtractQuantityMsg", nil)
	cdc.RegisterConcrete(UpdateAttrMsg{}, "ichain/UpdateAttrMsg", nil)

}
