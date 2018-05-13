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
			return handleCreate(ctx, k, msg)
		case RevokeMsg:
			return handleRevokeMsg(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCreate(ctx sdk.Context, k Keeper, msg CreateMsg) sdk.Result {
	tags, err := k.Create(ctx, Claim{
		ID:       msg.ID,
		Context:  msg.Context,
		Content:  msg.Content,
		Metadata: msg.Metadata,
	})
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Tags: tags}
}

func handleRevokeMsg(ctx sdk.Context, k Keeper, msg RevokeMsg) sdk.Result {
	tags, err := k.Revoke(ctx, msg.Owner, msg.ID, msg.Revocation)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
