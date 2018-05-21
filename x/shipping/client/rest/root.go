package rest

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// RegisterRoutes resgisters REST routes
func RegisterRoutes(r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/shipping", CreateOrderHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/shipping/{id}", QueryOrderRequestHandlerFn(storeName, cdc)).Methods("GET")
	r.HandleFunc("/shipping/{id}/confirm", ConfirmOrderHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/shipping/{id}/complete", CompleteOrderHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/shipping/{id}/cancel", CancelOrderHandlerFn(cdc, kb)).Methods("POST")
}
