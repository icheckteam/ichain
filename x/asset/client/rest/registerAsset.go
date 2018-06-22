package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

type createAssetBody struct {
	baseBody
	Asset asset.MsgCreateAsset `json:"asset"`
}

// Create asset REST handler
func CreateAssetHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m createAssetBody
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

		if m.Asset.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("asset.name is required"))
			return
		}

		if m.Asset.Quantity == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("asset.quantity is required"))
			return
		}

		if m.Asset.AssetID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("asset.id is required"))
			return
		}
		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		m.Asset.Sender = info.PubKey.Address()
		msg := m.Asset
		if err != nil { // XXX rechecking same error ?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		ctx = ctx.WithGas(m.Gas)
		ctx = ctx.WithAccountNumber(m.AccountNumber)
		ctx = ctx.WithSequence(m.Sequence)

		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, msg, cdc)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// send
		res, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
