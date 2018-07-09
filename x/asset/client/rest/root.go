package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

// resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/assets", CreateAssetHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}", QueryAssetRequestHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/assets/{id}/children", QueryAssetChildrensHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/assets/{id}/add", AddAssetQuantityHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/subtract", SubtractQuantityBodyHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/properties", UpdateAttributeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/materials", AddMaterialsHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/finalize", FinalizeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/reporters/{address}/revoke", RevokeReporterHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/proposals", CreateProposalHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/proposals", QueryProposalsHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/assets/{id}/proposals/{recipient}/answer", AnswerProposalHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/accounts/{address}/assets", QueryAccountAssetsHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
}
