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
	RestAccount    = "address"
)

// resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {

	// trusts
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/trusts", RestAccount), trustsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/trusts", RestAccount), SetTrustHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/claimed", RestAccount), claimedIdentHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/identities", RestAccount), identsByAccountHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/accounts/{%s}/certs", RestAccount), queryAccountCertsHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc("/identities", CreateIdentityHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/identities/{identityID}/certs", SetCertsHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/identities/{identityID}/certs", certsHandlerFn(ctx, cdc)).Methods("GET")
}
