package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
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
		m.Asset.Sender = info.GetPubKey().Address()
		msg := m.Asset
		signAndBuild(ctx, cdc, w, m.baseBody, msg)
	}
}
