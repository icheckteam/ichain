package rest

import (
	"errors"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/icheckteam/ichain/x/identity"
)

type baseBody struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainID       string `json:"chain_id"`
	Sequence      int64  `json:"sequence"`
	AccountNumber int64  `json:"account_number"`
	Gas           int64  `json:"gas"`
	Memo          string `json:"memo"`
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

type msgRegBody struct {
	BaseReq baseBody `json:"base_req"`
}

type msgAddOwnerBody struct {
	BaseReq baseBody       `json:"base_req"`
	Owner   sdk.AccAddress `json:"owner"`
}

type msgDelOwnerBody struct {
	BaseReq baseBody       `json:"base_req"`
	Owner   sdk.AccAddress `json:"owner"`
}

func signAndBuild(ctx context.CLIContext, cdc *wire.Codec, w http.ResponseWriter, m baseBody, msg sdk.Msg) {

	txCtx := authctx.TxContext{
		Codec:         cdc,
		Gas:           m.Gas,
		ChainID:       m.ChainID,
		AccountNumber: m.AccountNumber,
		Sequence:      m.Sequence,
		Memo:          m.Memo,
	}

	txBytes, err := txCtx.BuildAndSign(m.Name, m.Password, []sdk.Msg{msg})
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
