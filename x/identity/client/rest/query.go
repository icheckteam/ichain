package rest

import (
	"errors"
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
func QueryClaimRequestHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		claim, err := queryClaim(ctx, storeName, cdc, vars["id"])
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

// QueryClaimsByAccount
func QueryClaimsByAccount(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr, err := sdk.GetAccAddressBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		key := identity.GetAccountClaimsKey(addr)
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

			claim, err := queryClaim(ctx, storeName, cdc, claimID)
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

// QueryClaimsByIssuer
func QueryClaimsByIssuer(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr, err := sdk.GetAccAddressBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		key := identity.GetIssuerClaimsKey(addr)
		kvs, err := ctx.QuerySubspace(cdc, key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query issuer claims. Error: %s", err.Error())))
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

			claim, err := queryClaim(ctx, storeName, cdc, claimID)
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

func queryClaim(ctx context.CoreContext, storeName string, cdc *wire.Codec, assetID string) (*identity.Claim, error) {
	key := identity.GetClaimKey(assetID)
	res, err := ctx.Query(key, storeName)

	if res == nil {
		return nil, errors.New("asset not found")
	}

	if err != nil {
		return nil, err
	}

	var a identity.Claim

	err = cdc.UnmarshalBinary(res, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func queryClaimsByAccount(ctx context.CoreContext, storeName string, cdc *wire.Codec, account string) ([]identity.Claim, error) {
	address, err := sdk.GetAccAddressBech32(account)
	if err != nil {
		return nil, err
	}

	items := []identity.Claim{}
	kvs, err := ctx.QuerySubspace(cdc, identity.GetAccountClaimsKey(address), storeName)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		var itemID string
		err = cdc.UnmarshalBinary(kv.Value, &itemID)
		if err != nil {
			return nil, err
		}
		item, err := queryClaim(ctx, storeName, cdc, itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

func queryClaimsByIssuer(ctx context.CoreContext, storeName string, cdc *wire.Codec, account string) ([]identity.Claim, error) {
	address, err := sdk.GetAccAddressBech32(account)
	if err != nil {
		return nil, err
	}

	items := []identity.Claim{}
	kvs, err := ctx.QuerySubspace(cdc, identity.GetIssuerClaimsKey(address), storeName)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		var itemID string
		err = cdc.UnmarshalBinary(kv.Value, &itemID)
		if err != nil {
			return nil, err
		}
		item, err := queryClaim(ctx, storeName, cdc, itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}
