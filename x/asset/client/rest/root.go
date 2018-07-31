package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

// RegisterRoutes resgister REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, storeName string) {
	r.HandleFunc("/assets", createAssetHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/history", HistortyHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/assets/{id}", queryAssetRequestHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/assets/{id}/txs", assetTxsHandlerFn(ctx, storeName, cdc)).Methods("GET")
	r.HandleFunc("/assets/{id}/children", queryAssetChildrensHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/assets/{id}/add", addAssetQuantityHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/subtract", subtractQuantityBodyHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/properties", updateAttributeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/properties/{name}/history", queryHistoryUpdatePropertiesHandlerFn(ctx, cdc)).Methods("GET")
	r.HandleFunc("/assets/{id}/materials", addMaterialsHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/finalize", finalizeHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/reporters/{address}/revoke", revokeReporterHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/proposals", createProposalHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/proposals", queryProposalsHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/assets/{id}/proposals/{recipient}/answer", answerProposalHandlerFn(ctx, cdc, kb)).Methods("POST")
	r.HandleFunc("/accounts/{address}/assets", queryAccountAssetsHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/accounts/{address}/inventory", queryInventoryHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/accounts/{address}/proposals", queryAccountProposalsHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
	r.HandleFunc("/accounts/{address}/report-assets", queryReporterAssetsHandlerFn(ctx, storeName, cdc, kb)).Methods("GET")
}
