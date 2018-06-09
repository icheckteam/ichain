package restb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type updateAttributesBody struct {
	Msg        asset.UpdateAttrMsg `json:"update_attribute"`
	Fee        sdk.StdFee          `json:"fee"`
	Signatures []sdk.StdSignature  `json:"signatures"`
}

// Create asset REST handler
func UpdateAttributesHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m updateAttributesBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := BuildAndBroadcastTx(ctx, cdc, m.Msg, m.Fee, m.Signatures)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
