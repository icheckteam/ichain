package gs1

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler ...
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleSend(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {
	tags, err := k.CreateRecord(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
