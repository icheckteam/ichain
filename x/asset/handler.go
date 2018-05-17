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
		case CreateProposalMsg:
			return handleCreateProposal(ctx, k, msg)
		case AnswerProposalMsg:
			return handleAnswerProposal(ctx, k, msg)
		case RevokeProposalMsg:
			return handleRevokeProposal(ctx, k, msg)
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
		Company:  msg.Company,
		Email:    msg.Email,
	}

	_, tags, err := k.RegisterAsset(ctx, asset)

	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

func handleUpdateAttr(ctx sdk.Context, k Keeper, msg UpdateAttrMsg) sdk.Result {
	tags, err := k.UpdateAttribute(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleAddQuantity(ctx sdk.Context, k Keeper, msg AddQuantityMsg) sdk.Result {
	_, tags, err := k.AddQuantity(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleSubtractQuantity(ctx sdk.Context, k Keeper, msg SubtractQuantityMsg) sdk.Result {
	_, tags, err := k.SubtractQuantity(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleCreateProposal(ctx sdk.Context, k Keeper, msg CreateProposalMsg) sdk.Result {
	_, err := k.CreateProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleRevokeProposal(ctx sdk.Context, k Keeper, msg RevokeProposalMsg) sdk.Result {
	_, err := k.RevokeProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleAnswerProposal(ctx sdk.Context, k Keeper, msg AnswerProposalMsg) sdk.Result {
	_, err := k.AnswerProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
