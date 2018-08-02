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
func queryAssetRequestHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
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

func queryAccountAssetsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAccountAssets(ctx, storeName, cdc, vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't get assets. Error: %s", err.Error())))
			return
		}
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
func assetTxsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
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

func queryHistoryUpdatePropertiesHandlerFn(ctx context.CoreContext, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
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

func queryHistoryOwnersHandlerFn(ctx context.CoreContext, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
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

func queryAssetChildrensHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAssetChildrens(ctx, storeName, cdc, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't get assets. Error: %s", err.Error())))
			return
		}
		WriteJSON(w, cdc, items)
	}
}

func queryAccountAssets(ctx context.CoreContext, storeName string, cdc *wire.Codec, account string) ([]*asset.RecordOutput, error) {
	address, err := sdk.AccAddressFromBech32(account)
	if err != nil {
		return nil, err
	}
	return getRecordsByAccount(ctx, address, cdc)
}

func queryAssetChildrens(ctx context.CoreContext, storeName string, cdc *wire.Codec, assetID string) ([]*asset.RecordOutput, error) {
	kvs, err := ctx.QuerySubspace(cdc, asset.GetAssetChildrensKey(assetID), storeName)
	if err != nil {
		return nil, err
	}

	return getRecordsByKvs(ctx, kvs, cdc)
}

func queryReporterAssetsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		kvs, err := ctx.QuerySubspace(cdc, asset.GetReporterAssetsKey(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		records, err := getRecordsByKvs(ctx, kvs, cdc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		WriteJSON(w, cdc, records)
	}
}

func queryProposalsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		kvs, err := ctx.QuerySubspace(cdc, asset.GetProposalsKey(vars["asset_id"]), storeName)
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

func queryAccountProposalsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		kvs, err := ctx.QuerySubspace(cdc, asset.GetProposalsAccountKey(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		proposals := make([]ProposalOutput, len(kvs))
		for index, kv := range kvs {
			proposal := asset.Proposal{}
			var assetID string

			err = cdc.UnmarshalBinary(kv.Value, &assetID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(asset.GetProposalAccountKey(address, assetID), storeName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}
			proposal, err = asset.UnmarshalProposal(cdc, res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			proposals[index] = ToProposalOutput(proposal, assetID)

		}
		WriteJSON(w, cdc, proposals)
	}
}

// HistoryTransferOutput ...
type HistoryTransferOutput struct {
	Owner sdk.AccAddress `json:"recipient"`
	Time  int64          `json:"time"`
	Memo  string         `json:"memo"`
}

// HistoryChangeQuantityOutput ...
type HistoryChangeQuantityOutput struct {
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Int        `json:"amount"`
	Type   string         `json:"type"`
	Time   int64          `json:"time"`
	Memo   string         `json:"memo"`
}

// HistoryUpdateProperty ...
type HistoryUpdateProperty struct {
	Reporter sdk.AccAddress `json:"reporter"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Value    interface{}    `json:"value"`
	Time     int64          `json:"time"`
	Memo     string         `json:"memo"`
}

// HistoryAddMaterial ...
type HistoryAddMaterial struct {
	Sender  sdk.AccAddress `json:"sender"`
	Amount  sdk.Int        `json:"amount"`
	AssetID string         `json:"asset_id"`
	Time    int64          `json:"time"`
	Memo    string         `json:"memo"`
}
