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
		case CreateMsg:
			return handleClaimIssue(ctx, k, msg)
		case RevokeMsg:
			return handleRevokeMsg(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleClaimIssue(ctx sdk.Context, k Keeper, msg CreateMsg) sdk.Result {
	if err := k.ClaimIssue(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleRevokeMsg(ctx sdk.Context, k Keeper, msg RevokeMsg) sdk.Result {
	if err := k.Revoke(ctx, msg.ClaimID); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
