package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/client/errors"
	"github.com/icheckteam/ichain/x/identity"
)

// WriteJSON ...
func WriteJSON(w http.ResponseWriter, cdc *wire.Codec, data interface{}) {
	output, err := cdc.MarshalJSON(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

func withErr(fn func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			errors.WriteError(w, err)
			return
		}
	}
}

func getOwners(ctx context.CoreContext, ident sdk.AccAddress, cdc *wire.Codec) ([]sdk.AccAddress, error) {
	prefixKey := identity.KeyOwners(ident)
	kvs, err := ctx.QuerySubspace(cdc, prefixKey, storeName)
	if err != nil {
		return nil, err
	}
	owners := make([]sdk.AccAddress, len(kvs))
	var index int
	for _, kv := range kvs {
		owners[index] = sdk.AccAddress(kv.Key[1+sdk.AddrLen:])
		index++
	}
	return owners[:index], nil
}

func getCerts(ctx context.CoreContext, ident sdk.AccAddress, cdc *wire.Codec) (identity.Certs, error) {
	prefix := identity.KeyCerts(ident)
	kvs, err := ctx.QuerySubspace(cdc, prefix, storeName)
	if err != nil {
		return nil, err
	}
	certs := make(identity.Certs, len(kvs))
	for i, kv := range kvs {
		certs[i], err = identity.UnmarshalCert(cdc, kv.Value)
		if err != nil {
			return nil, err
		}
	}
	return certs, nil
}

func getTrusts(ctx context.CoreContext, ident sdk.AccAddress, cdc *wire.Codec) ([]sdk.AccAddress, error) {
	prefix := identity.KeyTrusts(ident)
	kvs, err := ctx.QuerySubspace(cdc, prefix, storeName)
	if err != nil {
		return nil, err
	}
	accounts := make([]sdk.AccAddress, len(kvs))
	var index = 0
	for _, kv := range kvs {
		accounts[index] = sdk.AccAddress(kv.Key[1+sdk.AddrLen:])
		index++
	}
	return accounts[:index], nil
}
