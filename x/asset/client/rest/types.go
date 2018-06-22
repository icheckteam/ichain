package rest

import (
	"errors"
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
