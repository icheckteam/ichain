package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/client/errors"
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
func createAssetHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		var m createAssetBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			return err
		}

		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}

		if m.Name == "" {
			return errors.New("name is required")
		}

		if m.Quantity.IsZero() {
			return errors.New("quantity is required")
		}

		if m.AssetID == "" {
			return errors.New("asset.id is required")
		}
		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return errors.New("asset.id is required")
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
		return nil
	})
}
