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
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCreateAsset(ctx sdk.Context, k Keeper, msg AssetCreateMsg) sdk.Result {
	asset := k.GetAsset(ctx, msg.RecordID)
	if asset.ID != "" {
		return InvalidTransaction("Record already exists").Result()
	}
	k.createAsset(ctx, Asset{
		ID:    msg.RecordID,
		Owner: msg.Sender,
		Name:  msg.RecordName,
	})
	return sdk.Result{}
}

func handleTrasfer(ctx sdk.Context, k Keeper, msg TransferMsg) sdk.Result {
	if err := k.Transfer(ctx, msg.Sender, msg.To, msg.RecordID); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
