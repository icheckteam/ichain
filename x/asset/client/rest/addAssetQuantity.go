package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type addAssetQuantityBody struct {
	BaseReq  baseBody `json:"base_req"`
	Quantity sdk.Int  `json:"quantity"`
}

func AddAssetQuantityHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m addAssetQuantityBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = m.BaseReq.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.Quantity.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Quantity is required"))
			return
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		// build message
		msg := asset.MsgAddQuantity{
			Sender:   sdk.AccAddress(info.GetPubKey().Address()),
			AssetID:  vars["id"],
			Quantity: m.Quantity,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}
