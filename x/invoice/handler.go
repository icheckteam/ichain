package invoice

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MakeHandle(k InvoiceKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		var err sdk.Error

		switch msg.(type) {
		case MsgCreate:
			err = k.CreateInvoice(ctx, msg.(MsgCreate))
		default:
			err = sdk.ErrUnknownRequest(fmt.Sprintf("Unrecognized trace Msg type: %v", reflect.TypeOf(msg).Name()))
		}

		if err != nil {
			return err.Result()
		}

		return sdk.Result{}
	}
}
