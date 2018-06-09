package restb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

type updateAttributesBody struct {
	Msg        asset.UpdateAttrMsg `json:"update_attribute"`
	Fee        auth.StdFee         `json:"fee"`
	Signatures []auth.StdSignature `json:"signatures"`
}

// UpdateAttributesHandlerFn ...
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
