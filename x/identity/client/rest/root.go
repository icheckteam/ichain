package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/claims", CreateClaimHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/claims/{id}/revoke", RevokeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/claims/{id}/answer", AnswerHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/claims/{id}", QueryClaimRequestHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/accounts/{address}/claims", QueryClaimsByAccount(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/accounts/{address}/issuer/claims", QueryClaimsByIssuer(ctx, storeName, cdc)).Methods("GET")
}
