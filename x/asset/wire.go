package asset

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

var msgCdc = wire.NewCodec()

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateAsset{}, "asset/CreateAsset", nil)
	cdc.RegisterConcrete(MsgAddQuantity{}, "asset/AddQuantity", nil)
	cdc.RegisterConcrete(MsgSubtractQuantity{}, "asset/SubtractQuantity", nil)
	cdc.RegisterConcrete(MsgUpdateProperties{}, "asset/UpdateProperties", nil)
	cdc.RegisterConcrete(MsgFinalize{}, "asset/FinalizeAsset", nil)
	cdc.RegisterConcrete(MsgAddMaterials{}, "asset/AddMaterials", nil)
	cdc.RegisterConcrete(MsgCreateProposal{}, "asset/CreateProposal", nil)
	cdc.RegisterConcrete(MsgAnswerProposal{}, "asset/AnswerProposal", nil)
	cdc.RegisterConcrete(MsgRevokeReporter{}, "asset/RevokeReporter", nil)
}

func init() {
	RegisterWire(msgCdc)
}
