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
		case MsgReg:
			return handleRegister(ctx, k, msg)
		case MsgAddOwner:
			return handleAddOwner(ctx, k, msg)
		case MsgDelOwner:
			return handleDelOwner(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func mapKeeperToHandler(mapFn func() (sdk.Tags, sdk.Error)) sdk.Result {
	tags, err := mapFn()
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleSetTrust(ctx sdk.Context, k Keeper, msg MsgSetTrust) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		return k.AddTrust(ctx, msg)
	})
}

func handleSetCerts(ctx sdk.Context, k Keeper, msg MsgSetCerts) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		return k.AddCerts(ctx, msg)
	})
}

func handleRegister(ctx sdk.Context, k Keeper, msg MsgReg) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		return k.Register(ctx, msg)
	})
}

func handleAddOwner(ctx sdk.Context, k Keeper, msg MsgAddOwner) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		return k.AddOwner(ctx, msg)
	})
}

func handleDelOwner(ctx sdk.Context, k Keeper, msg MsgDelOwner) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		return k.DeleteOwner(ctx, msg)
	})
}
