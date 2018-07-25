package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/client/errors"
	"github.com/icheckteam/ichain/x/asset"
)

type updateAttributeBody struct {
	BaseReq    baseBody         `json:"base_req"`
	Properties asset.Properties `json:"properties"`
}

func UpdateAttributeHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m updateAttributeBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			return err
		}

		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}

		if len(m.Properties) == 0 {
			return errors.New("properties is required")
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message

		msg := asset.MsgUpdateProperties{
			AssetID:    vars["id"],
			Properties: m.Properties,
			Sender:     sdk.AccAddress(info.GetPubKey().Address()),
		}

		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
