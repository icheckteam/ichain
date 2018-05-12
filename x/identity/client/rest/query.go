package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
)

///////////////////////////
// REST

// get key REST handler
func QueryClaimRequestHandlerFn(storeName string, cdc *wire.Codec) http.HandlerFunc {
	ctx := context.NewCoreContextFromViper()
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assetID := vars["id"]
		key := identity.GetClaimRecordKey(assetID)
		res, err := ctx.Query(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query asset. Error: %s", err.Error())))
			return
		}
		var claim identity.Claim
		err = cdc.UnmarshalBinary(res, &claim)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode claim. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(claim)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// QueryClaimsOwner
func QueryClaimsOwner(storeName string, cdc *wire.Codec) http.HandlerFunc {
	ctx := context.NewCoreContextFromViper()
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]

		bz, err := hex.DecodeString(addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		key := identity.GetClaimsOwnerKey(sdk.Address(bz))
		res, err := ctx.Query(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query claim. Error: %s", err.Error())))
			return
		}
		var claims []identity.Claim
		err = cdc.UnmarshalBinary(res, &claims)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode claim. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(claims)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// QueryClaimsOwner
func QueryClaimsAccount(storeName string, cdc *wire.Codec) http.HandlerFunc {
	ctx := context.NewCoreContextFromViper()
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]
		bz, err := hex.DecodeString(addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		key := identity.GetClaimsAccountKey(bz)
		res, err := ctx.Query(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query claim. Error: %s", err.Error())))
			return
		}
		var claims []identity.Claim
		err = cdc.UnmarshalBinary(res, &claims)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode claim. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(claims)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
