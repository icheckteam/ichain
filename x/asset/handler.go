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
		case MsgTransfer:
			return handleTransfer(ctx, k, msg)
		case MsgCreateAsset:
			return handleCreateAsset(ctx, k, msg)
		case MsgSubtractQuantity:
			return handleSubtractQuantity(ctx, k, msg)
		case MsgAddQuantity:
			return handleAddQuantity(ctx, k, msg)
		case MsgUpdateProperties:
			return handleUpdateProperties(ctx, k, msg)
		case MsgAddMaterials:
			return handleAddMaterials(ctx, k, msg)
		case MsgFinalize:
			return handleFinalize(ctx, k, msg)
		case MsgRevokeReporter:
			return handleRevokeReporter(ctx, k, msg)
		case MsgCreateReporter:
			return handleCreateReporter(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleTransfer(ctx sdk.Context, k Keeper, msg MsgTransfer) sdk.Result {
	tags, err := k.Transfer(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleCreateAsset(ctx sdk.Context, k Keeper, msg MsgCreateAsset) sdk.Result {
	tags, err := k.CreateAsset(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleUpdateProperties(ctx sdk.Context, k Keeper, msg MsgUpdateProperties) sdk.Result {
	tags, err := k.UpdateProperties(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleAddQuantity(ctx sdk.Context, k Keeper, msg MsgAddQuantity) sdk.Result {
	tags, err := k.AddQuantity(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleAddMaterials(ctx sdk.Context, k Keeper, msg MsgAddMaterials) sdk.Result {
	tags, err := k.AddMaterials(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleFinalize(ctx sdk.Context, k Keeper, msg MsgFinalize) sdk.Result {
	tags, err := k.Finalize(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleSubtractQuantity(ctx sdk.Context, k Keeper, msg MsgSubtractQuantity) sdk.Result {
	tags, err := k.SubtractQuantity(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleCreateReporter(ctx sdk.Context, k Keeper, msg MsgCreateReporter) sdk.Result {
	tags, err := k.CreateReporter(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleRevokeReporter(ctx sdk.Context, k Keeper, msg MsgRevokeReporter) sdk.Result {
	tags, err := k.RevokeReporter(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
