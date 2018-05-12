package bank

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(ck Keeper) sdk.Handler {
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
func handleSendMsg(ctx sdk.Context, ck Keeper, msg SendMsg) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	tags, err := ck.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Handle IssueMsg.
func handleIssueMsg(ctx sdk.Context, ck Keeper, msg IssueMsg) sdk.Result {
	panic("not implemented yet")
}
