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

type addMaterialsBody struct {
	BaseReq baseBody  `json:"base_req"`
	Amount  sdk.Coins `json:"amount"`
}

// AddMaterialsHandlerFn  REST handler
func AddMaterialsHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m addMaterialsBody
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

		if len(m.Amount) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("amount is required"))
			return
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		msg := asset.MsgAddMaterials{
			AssetID: vars["id"],
			Sender:  sdk.AccAddress(info.GetPubKey().Address()),
			Amount:  m.Amount,
		}

		// sign
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}
