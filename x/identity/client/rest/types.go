package rest

import (
	"errors"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/identity"
)

type baseBody struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainID       string `json:"chain_id"`
	Sequence      int64  `json:"sequence"`
	AccountNumber int64  `json:"account_number"`
	Gas           int64  `json:"gas"`
}

func (b baseBody) Validate() error {
	if b.Name == "" {
		return errors.New("name required but not specified")
	}
	if b.Password == "" {
		return errors.New("password required but not specified")
	}
	if b.Gas == 0 {
		return errors.New("gas required but not specified")
	}
	if len(b.ChainID) == 0 {
		return errors.New("chain_id required but not specified")
	}
	if b.AccountNumber < 0 {
		return errors.New("account_number required but not specified")
	}

	if b.Sequence < 0 {
		return errors.New("sequence required but not specified")
	}
	return nil
}

func signAndBuild(ctx context.CoreContext, cdc *wire.Codec, w http.ResponseWriter, m baseBody, msg sdk.Msg) {
	ctx = ctx.WithGas(m.Gas)
	ctx = ctx.WithAccountNumber(m.AccountNumber)
	ctx = ctx.WithSequence(m.Sequence)
	ctx = ctx.WithChainID(m.ChainID)
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

type msgSetTrustBody struct {
	BaseReq  baseBody `json:"base_req"`
	Trusting string   `json:"trusting"`
	Trust    bool     `json:"trust"`
}

type msgCreateIdentityBody struct {
	BaseReq baseBody `json:"base_req"`
}

type msgSetCertsBody struct {
	BaseReq baseBody             `json:"base_req"`
	Values  []identity.CertValue `json:"values"`
}
