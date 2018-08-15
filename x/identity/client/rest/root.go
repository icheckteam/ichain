package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

const (
	// RestIdentityID ...
	RestIdentityID = "identityID"
	// RestAccount ...
	RestAccount = "address"
)

// RegisterRoutes REST routes
func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/trusts", RestAccount), trustsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/trusts", RestAccount), SetTrustHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/certs", RestAccount), queryCertsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/certs", RestAccount), SetCertsHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/register", RestAccount), registerHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/owners", RestAccount), getOwnersHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/owners", RestAccount), addOwnerHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/idents/{%s}/owners/{owner}", RestAccount), delOwnerHandlerFn(ctx, cdc, kb)).Methods("DELETE")
}
