package bank

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tmlibs/common"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(ck CoinKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SendMsg:
			return handleSendMsg(ctx, ck, msg)
		case IssueMsg:
			return handleIssueMsg(ctx, ck, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SendMsg.
func handleSendMsg(ctx sdk.Context, ck CoinKeeper, msg SendMsg) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked

	err := ck.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	tags := []common.KVPair{}
	for _, in := range msg.Inputs {
		tags = append(tags, common.KVPair{
			Key:   []byte(in.Address.String()),
			Value: in.Address.Bytes(),
		})
	}
	for _, out := range msg.Outputs {
		tags = append(tags, common.KVPair{
			Key:   []byte("owner"),
			Value: []byte(out.Address.String()),
		})
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Handle IssueMsg.
func handleIssueMsg(ctx sdk.Context, ck CoinKeeper, msg IssueMsg) sdk.Result {
	panic("not implemented yet")
}
