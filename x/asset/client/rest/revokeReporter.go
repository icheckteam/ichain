package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
)

type revokeReporterBody struct {
	BaseReq baseBody `json:"base_req"`
}

func RevokeReporterHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m revokeReporterBody
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

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			return err
		}

		// build message

		msg := asset.MsgRevokeReporter{
			Sender:   sdk.AccAddress(info.GetPubKey().Address()),
			Reporter: address,
			AssetID:  vars["id"],
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
