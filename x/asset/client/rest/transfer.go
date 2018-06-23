package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

type transferBody struct {
	baseBody
	Assets []string
}

func TrasferHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m transferBody
		vars := mux.Vars(r)
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = m.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if len(m.Assets) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("assets is required"))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		address, err := sdk.GetAccAddressBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// build message

		msg := asset.MsgTransfer{
			Sender:    info.PubKey.Address(),
			Recipient: address,
			Assets:    m.Assets,
		}

		ctx = ctx.WithGas(m.Gas)
		ctx = ctx.WithAccountNumber(m.AccountNumber)
		ctx = ctx.WithSequence(m.Sequence)
		ctx = ctx.WithChainID(m.ChainID)
		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, msg, cdc)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// send
		res, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
