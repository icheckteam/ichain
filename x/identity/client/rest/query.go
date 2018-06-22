package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
)

///////////////////////////
// REST

// get key REST handler
func QueryClaimRequestHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		claim, err := queryClaim(storeName, ctx, cdc, vars["id"])
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

// QueryAccountClaims
func QueryAccountClaims(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]
		bz, err := hex.DecodeString(addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		key := identity.GetAccountClaimsKey(bz)
		kvs, err := ctx.QuerySubspace(cdc, key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query account claims. Error: %s", err.Error())))
			return
		}
		claims := []identity.Claim{}
		for _, kv := range kvs {
			var claimID string
			err = cdc.UnmarshalBinary(kv.Value, &claimID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode claim. Error: %s", err.Error())))
				return
			}

			claim, err := queryClaim(storeName, ctx, cdc, claimID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode claim. Error: %s", err.Error())))
				return
			}
			claims = append(claims, *claim)
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

func queryClaim(storeName string, ctx context.CoreContext, cdc *wire.Codec, claimID string) (*identity.Claim, error) {
	key := identity.GetClaimKey(claimID)
	res, err := ctx.Query(key, storeName)
	if err != nil {
		return nil, err
	}
	var claim *identity.Claim
	err = cdc.UnmarshalBinary(res, claim)
	if err != nil {
		return nil, err
	}

	return claim, nil
}
