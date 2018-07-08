package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
)

type createReporterBody struct {
	baseBody

	Reporter   string   `json:"reporter"`
	Properties []string `json:"properties"`
}

func CreateReporterHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m createReporterBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = m.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.Reporter == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Reporter is required"))
			return
		}

		if len(m.Properties) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("properties is required"))
			return
		}
		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		address, err := sdk.GetAccAddressBech32(m.Reporter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// build message

		msg := asset.MsgCreateReporter{
			Sender:     info.GetPubKey().Address(),
			Reporter:   address,
			Properties: m.Properties,
			AssetID:    vars["id"],
		}

		signAndBuild(ctx, cdc, w, m.baseBody, msg)
	}
}
