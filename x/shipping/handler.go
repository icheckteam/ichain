package shipping

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler ...
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateOrderMsg:
			return handleCreateOrder(ctx, k, msg)
		case ConfirmOrderMsg:
			return handleConfirmOrder(ctx, k, msg)
		case CompleteOrderMsg:
			return handleCompleteOrder(ctx, k, msg)
		case CancelOrderMsg:
			return handleCancelOrder(ctx, k, msg)
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

func handleCreateOrder(ctx sdk.Context, k Keeper, msg CreateOrderMsg) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		tags, err := k.CreateOrder(ctx, msg)
		return tags, err
	})
}

func handleConfirmOrder(ctx sdk.Context, k Keeper, msg ConfirmOrderMsg) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		tags, err := k.ConfirmOrder(ctx, msg)
		return tags, err
	})
}

func handleCompleteOrder(ctx sdk.Context, k Keeper, msg CompleteOrderMsg) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		tags, err := k.CompleteOrder(ctx, msg)
		return tags, err
	})
}

func handleCancelOrder(ctx sdk.Context, k Keeper, msg CancelOrderMsg) sdk.Result {
	return mapKeeperToHandler(func() (sdk.Tags, sdk.Error) {
		tags, err := k.CancelOrder(ctx, msg)
		return tags, err
	})
}
