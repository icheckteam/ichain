package rest

import (
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

func (b revokeReporterBody) ValidateBasic() error {
	err := b.BaseReq.Validate()
	if err != nil {
		return err
	}

	return nil
}

func revokeReporterHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m revokeReporterBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
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
