package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// RegisterRoutes resgisters REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/shipping", CreateOrderHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/shipping/{id}", QueryOrderRequestHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/shipping/{id}/confirm", ConfirmOrderHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/shipping/{id}/complete", CompleteOrderHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/shipping/{id}/cancel", CancelOrderHandlerFn(ctx, cdc, kb)).Methods("POST")
}
