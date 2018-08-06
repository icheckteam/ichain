package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/client/errors"
	"github.com/icheckteam/ichain/x/asset"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type addAssetQuantityBody struct {
	BaseReq  baseBody `json:"base_req"`
	Quantity sdk.Int  `json:"quantity"`
}

func (b addAssetQuantityBody) ValidateBasic() error {
	err := b.BaseReq.Validate()
	if err != nil {
		return err
	}
	if b.Quantity.IsZero() {
		return errors.New("Quantity is required")
	}
	return nil
}

func addAssetQuantityHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m addAssetQuantityBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message
		msg := asset.MsgAddQuantity{
			Sender:   sdk.AccAddress(info.GetPubKey().Address()),
			AssetID:  vars["id"],
			Quantity: m.Quantity,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
