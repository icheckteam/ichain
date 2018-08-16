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

func getOwnersHandlerFn(ctx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ident, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		owners, err := getOwners(ctx, ident, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		WriteJSON(w, cdc, owners)
	}
}

func queryCertsHandlerFn(ctx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		certs, err := getCerts(ctx, address, cdc)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query certs. Error: %s", err.Error())))
			return
		}
		WriteJSON(w, cdc, certs)
	}
}

func trustsHandlerFn(ctx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		accs, err := getTrusts(ctx, address, cdc)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query accs. Error: %s", err.Error())))
			return
		}
		WriteJSON(w, cdc, accs)
	}
}

func hasTrust(ctx context.CLIContext, cdc *wire.Codec, trustor, trusting sdk.AccAddress) bool {
	res, err := ctx.QueryStore(identity.KeyTrust(trustor, trusting), "identity")
	if err != nil {
		panic(err)
	}

	if len(res) > 0 {
		return true
	}

	return false
}
