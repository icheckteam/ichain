package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/icheckteam/ichain/client/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
)

type addMaterialsBody struct {
	BaseReq baseBody         `json:"base_req"`
	Amount  []asset.Material `json:"amount"`
}

// AddMaterialsHandlerFn  REST handler
func addMaterialsHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m addMaterialsBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			return err
		}

		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}

		if len(m.Amount) == 0 {
			return errors.New("amount is required")
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}

		msg := asset.MsgAddMaterials{
			AssetID: vars["id"],
			Sender:  sdk.AccAddress(info.GetPubKey().Address()),
			Amount:  m.Amount,
		}

		// sign
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
