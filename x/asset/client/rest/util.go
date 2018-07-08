package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func signAndBuild(ctx context.CoreContext, cdc *wire.Codec, w http.ResponseWriter, m baseBody, msg sdk.Msg) {
	ctx = ctx.WithGas(m.Gas)
	ctx = ctx.WithAccountNumber(m.AccountNumber)
	ctx = ctx.WithSequence(m.Sequence)
	ctx = ctx.WithChainID(m.ChainID)
	txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, []sdk.Msg{msg}, cdc)
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
