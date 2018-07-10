package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
)

const storeName = "identity"

func identsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kvs, err := ctx.QuerySubspace(cdc, identity.IdentitiesKey, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query idents. Error: %s", err.Error())))
			return
		}

		idents := make([]identity.Identity, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			ident := identity.Identity{}
			err = cdc.UnmarshalBinary(addr, &ident)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
				return
			}

			idents[i] = ident
		}

		output, err := cdc.MarshalJSON(idents)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func certsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		identID, err := strconv.Atoi(vars[RestIdentityID])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode ident_id. Error: %s", err.Error())))
			return
		}
		kvs, err := ctx.QuerySubspace(cdc, identity.KeyCerts(int64(identID), vars["property"]), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query idents. Error: %s", err.Error())))
			return
		}

		certs := make([]identity.Cert, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			cert := identity.Cert{}
			err = cdc.UnmarshalBinary(addr, &cert)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
				return
			}
			certs[i] = cert
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
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars[RestTrusting])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode address. Error: %s", err.Error())))
			return
		}
		kvs, err := ctx.QuerySubspace(cdc, identity.KeyTrusts(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query trusts. Error: %s", err.Error())))
			return
		}

		trusts := make([]identity.Trust, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			trust := identity.Trust{}
			err = cdc.UnmarshalBinary(addr, &trust)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode trust. Error: %s", err.Error())))
				return
			}

			trusts[i] = trust
		}

		output, err := cdc.MarshalJSON(trusts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
