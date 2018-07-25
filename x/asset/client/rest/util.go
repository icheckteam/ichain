package rest

import (
	"net/http"

	"github.com/icheckteam/ichain/client/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func signAndBuild(ctx context.CoreContext, cdc *wire.Codec, w http.ResponseWriter, m baseBody, msg sdk.Msg) {
	ctx = ctx.WithGas(m.Gas)
	ctx = ctx.WithAccountNumber(m.AccountNumber)
	ctx = ctx.WithSequence(m.Sequence)
	ctx = ctx.WithChainID(m.ChainID)

	if len(m.Memo) > 0 {
		ctx = ctx.WithMemo(m.Memo)
	}

	txBytes, err := ctx.SignAndBuild(m.Name, m.Password, []sdk.Msg{msg}, cdc)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	// send
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BroadcastTx:" + err.Error()))
		return
	}

	output, err := wire.MarshalJSONIndent(cdc, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(output)
}

func withErrHandler(fn func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			errors.WriteError(w, err)
			return
		}
	}
}
