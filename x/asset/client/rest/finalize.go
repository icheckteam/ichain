package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type finalizeBody struct {
	BaseReq baseBody `json:"base_req"`
}

// FinalizeHandlerFn ...
func finalizeHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		var m finalizeBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			return err
		}

		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message
		msg := asset.MsgFinalize{
			Sender:  sdk.AccAddress(info.GetPubKey().Address()),
			AssetID: vars["id"],
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
