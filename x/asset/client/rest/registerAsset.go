package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type createAssetBody struct {
	LocalAccountName string    `json:"account_name"`
	Password         string    `json:"password"`
	Asset            assetBody `json:"asset"`
	ChainID          string    `json:"chain_id"`
	Sequence         int64     `json:"sequence"`
	Gas              int64     `json:"gas"`
}

type assetBody struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	Quantity int64  `json:"quantity"`
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

		if m.Asset.ID == "" {
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
		msg := buildCreateAssetMsg(info.PubKey.Address(), m)
		if err != nil { // XXX rechecking same error ?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// sign
		ctx, err = withContext(ctx.WithFromAddressName(m.LocalAccountName), m.Gas)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("withContextErr:" + err.Error()))
			return
		}
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

func buildCreateAssetMsg(creator sdk.Address, body createAssetBody) sdk.Msg {
	return asset.NewRegisterMsg(
		creator,
		body.Asset.ID,
		body.Asset.Name,
		body.Asset.Quantity,
		body.Asset.Company,
		body.Asset.Email,
	)
}
