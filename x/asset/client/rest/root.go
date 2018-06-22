package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/assets", CreateAssetHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}", QueryAssetRequestHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/assets/{id}/add", AddAssetQuantityHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/subtract", SubtractQuantityBodyHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/properties", UpdateAttributeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/materials", AddMaterialsHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/finalize", FinalizeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/reporters", CreateReporterHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/reporters/{address}/revoke", RevokeReporterHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/accounts/{address}/transfer-asset", TrasferHandlerFn(ctx, cdc, kb)).Methods("POST")
}
