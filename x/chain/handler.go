package chain

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler ...
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRecordCreate:
			return handleCreateRecord(ctx, k, msg)
		case MsgChangeRecordOwner:
			return handlerChangeRecordOwner(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCreateRecord(ctx sdk.Context, k Keeper, msg MsgRecordCreate) sdk.Result {
	record := k.GetRecord(ctx, msg.RecordID)
	if record.ID != "" {
		return InvalidTransaction("Record already exists").Result()
	}
	k.createRecord(ctx, Record{
		ID:    msg.RecordID,
		Owner: msg.Sender,
		Name:  msg.RecordName,
	})
	return sdk.Result{}
}

func handlerChangeRecordOwner(ctx sdk.Context, k Keeper, msg MsgChangeRecordOwner) sdk.Result {
	if err := k.Transfer(ctx, msg.Sender, msg.To, msg.RecordID); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
