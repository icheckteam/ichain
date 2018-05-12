package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

type subtractAssetQuantityBody struct {
	LocalAccountName string `json:"account_name"`
	Password         string `json:"password"`
	Quantity         int64  `json:"quantity"`
	ChainID          string `json:"chain_id"`
	Sequence         int64  `json:"sequence"`
}

func SubtractQuantityBodyHandlerFn(cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	ctx := context.NewCoreContextFromViper()
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
		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		// build message
		msg := buildsubtractAssetQuantityMsg(info.PubKey.Address(), vars["id"], m)

		// sign
		ctx = ctx.WithSequence(m.Sequence).WithChainID(m.ChainID)
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

func buildsubtractAssetQuantityMsg(creator sdk.Address, assetID string, body addAssetQuantityBody) sdk.Msg {
	return asset.SubtractQuantityMsg{
		Issuer:   creator,
		ID:       assetID,
		Quantity: body.Quantity,
	}
}
