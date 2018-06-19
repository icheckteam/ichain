package rest

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client/context"
)

type baseBody struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	Sequence         int64  `json:"sequence"`
	AccountNumber    int64  `json:"account_number"`
	Gas              int64  `json:"gas"`
}

func (b baseBody) Validate() error {
	if b.LocalAccountName == "" {
		return errors.New("account_name is required")
	}
	if b.Password == "" {
		return errors.New("password is required")
	}
	if b.Gas == 0 {
		return errors.New("gas is required")
	}
	return nil
}

func (b baseBody) WithContext(ctx context.CoreContext) context.CoreContext {
	ctx = ctx.WithGas(b.Gas)
	ctx = ctx.WithAccountNumber(b.AccountNumber)
	ctx = ctx.WithSequence(b.Sequence)
	return ctx
}
