package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

///////////////////////////
// REST

// get key REST handler
func QueryClaimRequestHandlerFn(storeName string, cdc *wire.Codec) http.HandlerFunc {
	ctx := context.NewCoreContextFromViper()
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assetID := vars["id"]

		hash, err := hex.DecodeString(assetID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := ctx.Query(hash, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query asset. Error: %s", err.Error())))
			return
		}
		w.Write(res)
	}
}
