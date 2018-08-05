package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
)

const storeName = "identity"

func queryCertsByOwner(ctx context.CoreContext, cdc *wire.Codec, owner sdk.AccAddress, vars map[string]string) (identity.Certs, error) {
	return nil, nil
}

func queryCertsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		certs, err := queryCertsByOwner(ctx, cdc, address, vars)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query certs. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(certs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func trustsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not thing here"))
	}
}

func hasTrust(ctx context.CoreContext, cdc *wire.Codec, trustor, trusting sdk.AccAddress) bool {
	res, err := ctx.QueryStore(identity.KeyTrust(trustor, trusting), "identity")
	if err != nil {
		panic(err)
	}

	if len(res) > 0 {
		return true
	}

	return false
}
