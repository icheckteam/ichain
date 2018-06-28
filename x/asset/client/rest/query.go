package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

///////////////////////////
// REST

// get key REST handler
func QueryAssetRequestHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		a, err := queryAsset(ctx, storeName, cdc, vars["id"])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't decode asset. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(ToAssetOutput(*a))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

func QueryAccountAssetsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAccountAssets(ctx, storeName, cdc, vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get assets. Error: %s", err.Error())))
			return
		}
		output, err := cdc.MarshalJSON(ToAssetsOutput(items))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

func QueryAssetChildrensHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAssetChildrens(ctx, storeName, cdc, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get assets. Error: %s", err.Error())))
			return
		}
		output, err := cdc.MarshalJSON(ToAssetsOutput(items))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

func queryAsset(ctx context.CoreContext, storeName string, cdc *wire.Codec, assetID string) (*asset.Asset, error) {
	key := asset.GetAssetKey(assetID)
	res, err := ctx.Query(key, storeName)

	if res == nil {
		return nil, errors.New("asset not found")
	}

	if err != nil {
		return nil, err
	}

	var a asset.Asset

	err = cdc.UnmarshalBinary(res, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func queryAccountAssets(ctx context.CoreContext, storeName string, cdc *wire.Codec, account string) ([]asset.Asset, error) {
	address, err := sdk.GetAccAddressBech32(account)
	if err != nil {
		return nil, err
	}

	items := []asset.Asset{}
	kvs, err := ctx.QuerySubspace(cdc, asset.GetAccountAssetsKey(address), storeName)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		var itemID string
		err = cdc.UnmarshalBinary(kv.Value, &itemID)
		if err != nil {
			return nil, err
		}
		item, err := queryAsset(ctx, storeName, cdc, itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

func queryAssetChildrens(ctx context.CoreContext, storeName string, cdc *wire.Codec, assetID string) ([]asset.Asset, error) {
	items := []asset.Asset{}
	kvs, err := ctx.QuerySubspace(cdc, asset.GetAssetChildrensKey(assetID), storeName)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		var itemID string
		err = cdc.UnmarshalBinary(kv.Value, &itemID)
		if err != nil {
			return nil, err
		}
		item, err := queryAsset(ctx, storeName, cdc, itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}
