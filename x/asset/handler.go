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
		case AssetCreateMsg:
			return handleCreateAsset(ctx, k, msg)
		case TransferMsg:
			return handleTrasfer(ctx, k, msg)
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

func handleCreateAsset(ctx sdk.Context, k Keeper, msg AssetCreateMsg) sdk.Result {
	if k.Has(ctx, msg.AssetID) {
		return InvalidTransaction("Asset already exists").Result()
	}
	k.setAsset(ctx, Asset{
		ID:     msg.AssetID,
		Issuer: msg.Sender,
		Name:   msg.AssetName,
	})
	return sdk.Result{}
}

func handleTrasfer(ctx sdk.Context, k Keeper, msg TransferMsg) sdk.Result {
	if err := k.Transfer(ctx, msg); err != nil {
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
