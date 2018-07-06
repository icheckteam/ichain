package identity

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler ...
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateClaim:
			return handleCreate(ctx, k, msg)
		case MsgRevokeClaim:
			return handleRevokeMsg(ctx, k, msg)
		case MsgAnswerClaim:
			return handleAnswerMsg(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCreate(ctx sdk.Context, k Keeper, msg MsgCreateClaim) sdk.Result {
	tags, err := k.CreateClaim(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Tags: tags}
}

func handleRevokeMsg(ctx sdk.Context, k Keeper, msg MsgRevokeClaim) sdk.Result {
	tags, err := k.RevokeClaim(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleAnswerMsg(ctx sdk.Context, k Keeper, msg MsgAnswerClaim) sdk.Result {
	tags, err := k.AnswerClaim(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
