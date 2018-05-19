package rest

import (
	"net/http"

	"github.com/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto/keys"
)

type createContractBody struct {
	AccountName string `json:"name"`
}

// CreateContractHandlderFn
func CreateContractHandlderFn(cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
