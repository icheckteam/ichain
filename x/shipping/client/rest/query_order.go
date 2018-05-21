package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/shipping"
)

///////////////////////////
// REST

// QueryOrderRequestHandlerFn gets key REST handler
func QueryOrderRequestHandlerFn(storeName string, cdc *wire.Codec) http.HandlerFunc {
	ctx := context.NewCoreContextFromViper()
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["id"]

		key := shipping.GetOrderKey([]byte(orderID))
		res, err := ctx.Query(key, storeName)
		var order shipping.Order
		err = cdc.UnmarshalBinary(res, &order)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode order. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(order)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
