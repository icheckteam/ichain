package asset

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

var msgCdc = wire.NewCodec()

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateAsset{}, "ichain/MsgCreateAsset", nil)
	cdc.RegisterConcrete(AddQuantityMsg{}, "ichain/AddQuantityMsg", nil)
	cdc.RegisterConcrete(MsgSubtractQuantity{}, "ichain/MsgSubtractQuantity", nil)
	cdc.RegisterConcrete(MsgUpdateProperties{}, "ichain/MsgUpdateProperties", nil)
	cdc.RegisterConcrete(MsgSend{}, "ichain/SendAsset", nil)
	cdc.RegisterConcrete(MsgFinalize{}, "ichain/FinalizeAsset", nil)
	cdc.RegisterConcrete(MsgAddMaterials{}, "ichain/AddMaterials", nil)
	cdc.RegisterConcrete(CreateProposalMsg{}, "ichain/CreateProposalMsg", nil)
	cdc.RegisterConcrete(AnswerProposalMsg{}, "ichain/AnswerProposalMsg", nil)
	cdc.RegisterConcrete(RevokeProposalMsg{}, "ichain/RevokeProposalMsg", nil)
}
