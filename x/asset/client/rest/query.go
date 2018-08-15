package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/icheckteam/ichain/x/asset"
)

// AssetOutput ..
type AssetOutput struct {
	Asset     asset.Asset            `json:"asset"`
	Materials map[string]asset.Asset `json:"material_by_id"`
}

///////////////////////////
// REST

// get key REST handler
func queryAssetRequestHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		record, err := getRecord(ctx, vars["id"], cdc)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't decode asset. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(record)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

func queryAccountAssetsHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAccountAssets(ctx, storeName, cdc, vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't get assets. Error: %s", err.Error())))
			return
		}
		items = items.Sort()
		output, err := cdc.MarshalJSON(items)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

// TxsHandlerFn ...
func assetTxsHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		info, err := queryAssetTxs(ctx, vars["id"], cdc, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		output, err := cdc.MarshalJSON(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
	}
}

func queryHistoryUpdatePropertiesHandlerFn(ctx context.CLIContext, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		info, err := queryAssetTxs(ctx, vars["id"], cdc, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		WriteJSON2(w, cdc, filterTxUpdateProperties(info, vars["name"]))
	}
}

func queryHistoryOwnersHandlerFn(ctx context.CLIContext, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		info, err := queryAssetTxs(ctx, vars["id"], cdc, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		WriteJSON2(w, cdc, filterTxChangeOwner(info))
	}
}
func queryHistoryTransferMaterialsHandlerFn(ctx context.CLIContext, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		info, err := queryAssetTxs(ctx, vars["id"], cdc, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		WriteJSON2(w, cdc, filterTxTransferMaterial(info))
	}
}

func queryAssetChildrensHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAssetChildrens(ctx, storeName, cdc, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't get assets. Error: %s", err.Error())))
			return
		}
		items = items.Sort()
		WriteJSON(w, cdc, items)
	}
}

func queryAccountAssets(ctx context.CLIContext, storeName string, cdc *wire.Codec, account string) (asset.RecordsOutput, error) {
	address, err := sdk.AccAddressFromBech32(account)
	if err != nil {
		return nil, err
	}
	return getRecordsByAccount(ctx, address, cdc)
}

func queryAssetChildrens(ctx context.CLIContext, storeName string, cdc *wire.Codec, assetID string) (asset.RecordsOutput, error) {
	kvs, err := ctx.QuerySubspace(asset.GetAssetChildrensKey(assetID), storeName)
	if err != nil {
		return nil, err
	}

	return getRecordsByKvs(ctx, kvs, cdc)
}

func queryReporterAssetsHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		kvs, err := ctx.QuerySubspace(asset.GetReporterAssetsKey(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		records, err := getRecordsByKvs(ctx, kvs, cdc)
		records = records.Sort()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		WriteJSON(w, cdc, records)
	}
}

func queryProposalsHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		kvs, err := ctx.QuerySubspace(asset.GetProposalsKey(vars["asset_id"]), storeName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		proposals, err := getProposals(ctx, kvs, cdc)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		WriteJSON(w, cdc, proposals)
	}
}

func queryAccountProposalsHandlerFn(ctx context.CLIContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		kvs, err := ctx.QuerySubspace(asset.GetProposalsAccountKey(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		proposals := make([]asset.ProposalOutput, len(kvs))
		for index, kv := range kvs {
			recordID := string(kv.Key[1+sdk.AddrLen:])
			proposal, err := getProposal(ctx, address, recordID, cdc)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("Couldn't get proposal. Error: %s", err.Error())))
				return
			}
			proposals[index] = asset.ToProposalOutput(proposal, recordID)

		}
		WriteJSON(w, cdc, proposals)
	}
}
