package restb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

type createAssetBody struct {
	Msg        asset.RegisterMsg   `json:"asset"`
	Fee        auth.StdFee         `json:"fee"`
	Signatures []auth.StdSignature `json:"signatures"`
}

// CreateAssetHandlerFn ...
func CreateAssetHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m createAssetBody
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

// BuildAndBroadcastTx ...
func BuildAndBroadcastTx(ctx context.CoreContext, cdc *wire.Codec, msg sdk.Msg, fee auth.StdFee, sigs []auth.StdSignature) ([]byte, error) {
	// marshal bytes
	tx := auth.NewStdTx(msg, fee, sigs)

	b, err := cdc.MarshalBinary(tx)

	if err != nil {
		return nil, err
	}

	// send
	res, err := ctx.BroadcastTx(b)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(res, "", "  ")
}
