package insurance

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
			return handleCreateClaim(ctx, k, msg)
		case MsgCreateContract:
			return handleCreateContract(ctx, k, msg)
		case MsgProcessClaim:
			return handleProcessClaim(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCreateClaim(ctx sdk.Context, k Keeper, msg MsgCreateClaim) sdk.Result {
	if err := k.CreateClaim(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleCreateContract(ctx sdk.Context, k Keeper, msg MsgCreateContract) sdk.Result {
	tags, err := k.CreateContract(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleProcessClaim(ctx sdk.Context, k Keeper, msg MsgProcessClaim) sdk.Result {
	if err := k.ProcessClaim(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
