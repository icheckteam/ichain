package rest

import (
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

func (b createAssetBody) ValidateBasic() error {
	err := b.BaseReq.Validate()
	if err != nil {
		return err
	}
	if b.Name == "" {
		return errors.New("name is required")
	}

	if b.Quantity.IsZero() {
		return errors.New("quantity is required")
	}

	if b.AssetID == "" {
		return errors.New("asset.id is required")
	}
	return nil
}

// Create asset REST handler
func createAssetHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		var m createAssetBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
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
