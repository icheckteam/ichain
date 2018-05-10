package asset

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler ...
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case RegisterMsg:
			return handleRegisterAsset(ctx, k, msg)
		case SubtractQuantityMsg:
			return handleSubtractQuantity(ctx, k, msg)
		case AddQuantityMsg:
			return handleAddQuantity(ctx, k, msg)
		case UpdateAttrMsg:
			return handleUpdateAttr(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleRegisterAsset(ctx sdk.Context, k Keeper, msg RegisterMsg) sdk.Result {
	asset := Asset{
		ID:       msg.ID,
		Name:     msg.Name,
		Issuer:   msg.Issuer,
		Quantity: msg.Quantity,
	}
	if err := k.RegisterAsset(ctx, asset); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleUpdateAttr(ctx sdk.Context, k Keeper, msg UpdateAttrMsg) sdk.Result {
	if err := k.UpdateAttribute(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleAddQuantity(ctx sdk.Context, k Keeper, msg AddQuantityMsg) sdk.Result {
	if err := k.AddQuantity(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleSubtractQuantity(ctx sdk.Context, k Keeper, msg SubtractQuantityMsg) sdk.Result {
	if err := k.SubtractQuantity(ctx, msg); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
