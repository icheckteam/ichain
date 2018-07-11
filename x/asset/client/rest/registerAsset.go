package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
)

type createAssetBody struct {
	BaseReq    baseBody         `json:"base_req"`
	AssetID    string           `json:"asset_id"`
	Name       string           `json:"name"`
	Quantity   sdk.Int          `json:"quantity"`
	Parent     string           `json:"parent"`
	Unit       string           `json:"unit"`
	Properties asset.Properties `json:"properties"`
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

		err = m.BaseReq.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("name is required"))
			return
		}

		if m.Quantity.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("quantity is required"))
			return
		}

		if m.AssetID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("asset.id is required"))
			return
		}
		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := asset.MsgCreateAsset{
			AssetID:    m.AssetID,
			Name:       m.Name,
			Parent:     m.Parent,
			Properties: m.Properties,
			Sender:     sdk.AccAddress(info.GetPubKey().Address()),
			Quantity:   m.Quantity,
			Unit:       m.Unit,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}
