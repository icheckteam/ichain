package rest

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

// resgister REST routes
func RegisterRoutes(r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/assets", CreateAssetHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}", QueryAssetRequestHandlerFn("asset", cdc)).Methods("POST")
	r.HandleFunc("/assets/{id}/add-quantity", AddAssetQuantityHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/subtract-quantity", SubtractQuantityBodyHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/update-attribute", UpdateAttributeHandlerFn(cdc, kb)).Methods("POST")
	r.HandleFunc("/assets/{id}/transfer", TransferAssetHandlerFn(cdc, kb)).Methods("POST")
}
