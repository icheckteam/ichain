package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

const (
	RestIdentityID = "identityID"
	RestTrusting   = "address"
	RestTrustor    = "address"
)

// resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {

	// trusts
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/trusts", RestTrustor), trustsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/trusts", RestTrusting), SetTrustHandlerFn(ctx, cdc, kb)).Methods("POST")

	// identities
	r.HandleFunc("/identities", identsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/identities/{%s}", RestIdentityID), identsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc("/identities", CreateIdentityHandlerFn(ctx, cdc, kb)).Methods("POST")

	// certs
	r.HandleFunc(fmt.Sprintf("/identities/{%s}/certs", RestIdentityID), SetCertsHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/identities/{%s}/certs", RestIdentityID), certsHandlerFn(ctx, cdc)).Methods("GET")
}
