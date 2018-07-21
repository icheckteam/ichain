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
		case MsgSetCerts:
			return handleSetCerts(ctx, k, msg)
		case MsgSetTrust:
			return handleSetTrust(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleSetTrust(ctx sdk.Context, k Keeper, msg MsgSetTrust) sdk.Result {
	err := k.AddTrust(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleSetCerts(ctx sdk.Context, k Keeper, msg MsgSetCerts) sdk.Result {
	err := k.AddCerts(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
