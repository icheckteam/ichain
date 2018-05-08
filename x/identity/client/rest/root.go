package rest

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// resgister REST routes
func RegisterRoutes(r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/claims", CreateClaimHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/claims/{id}/revoke", RevokeHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/claims/{id}", QueryClaimRequestHandlerFn("identity", cdc)).Methods("GET")
}
