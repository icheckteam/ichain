package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto/keys"
)

type addAssetQuantityBody struct {
	baseBody
	Quantity int64 `json:"quantity"`
}

func AddAssetQuantityHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m addAssetQuantityBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.LocalAccountName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("account_name is required"))
			return
		}

		if m.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("password is required"))
			return
		}

		if m.Quantity == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Quantity is required"))
			return
		}

		if m.Quantity == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Quantity is required"))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		// build message
		msg := buildAdAssetQuantityMsg(info.PubKey.Address(), vars["id"], m.Quantity)
		signAndBuild(ctx, cdc, w, m.baseBody, msg)
	}
}

func buildAdAssetQuantityMsg(creator sdk.Address, assetID string, qty int64) sdk.Msg {
	return asset.MsgAddQuantity{
		Sender:   creator,
		AssetID:  assetID,
		Quantity: qty,
	}
}
