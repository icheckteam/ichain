package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/warranties", CreateContractHandlerFn(ctx, cdc, kb))
	r.HandleFunc("/warranties", CreateClaimHandlerFn(ctx, cdc, kb))
	r.HandleFunc("/warranties/{id}/process", ProcessClaimHandlerFn(ctx, cdc, kb))
	r.HandleFunc("/warranties/{id}", QueryContractHandlerFn(ctx, storeName, cdc))
}
