package rest

import (
	"fmt"
	"net/http"

	"github.com/icheckteam/ichain/client/errors"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/icheckteam/ichain/x/asset"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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
			w.WriteHeader(http.StatusNotFound)
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
		info, err := queryAssetTxs(ctx, vars["id"], cdc)
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
		info, err := queryAssetTxs(ctx, vars["id"], cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := cdc.MarshalJSON(filterTxUpdateProperties(info, vars["name"]))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
	}
}

func queryAssetTxs(ctx context.CoreContext, assetID string, cdc *wire.Codec) (txInfos, error) {

	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("asset_id='%s'", assetID)
	page := 0
	perPage := 500
	prove := false
	res, err := node.TxSearch(query, prove, page, perPage)
	if err != nil {
		return nil, err
	}
	info, err := formatTxResults(cdc, res.Txs)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func filterTxUpdateProperties(infos txInfos, name string) []historyUpdateProperty {
	history := []historyUpdateProperty{}
	for _, info := range infos {
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgUpdateProperties:
				for _, p := range msg.Properties {
					if name != "" && name != p.Name {
						continue
					}
					history = append(history, historyUpdateProperty{
						Type:  asset.PropertyTypeToString(p.Type),
						Name:  p.Name,
						Value: p.GetValue(),
						Time:  info.Time,
						Memo:  info.Tx.Memo,
					})
				}
				break
			default:
				break
			}
		}
	}
	return history
}

func filterTxChangeOwner(infos txInfos) []historyTransferOutput {
	history := []historyTransferOutput{}
	for _, info := range infos {
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgAnswerProposal:
				if msg.Role == asset.RoleOwner {
					history = append(history, historyTransferOutput{
						Time:  info.Time,
						Memo:  info.Tx.Memo,
						Owner: msg.Sender,
					})
				}
				break
			default:
				break
			}
		}
	}
	return history
}

func queryAssetChildrensHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		items, err := queryAssetChildrens(ctx, storeName, cdc, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
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

func queryAsset(ctx context.CoreContext, storeName string, cdc *wire.Codec, assetID string) (*asset.Asset, error) {
	key := asset.GetAssetKey(assetID)
	res, err := ctx.QueryStore(key, storeName)

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

func queryAssetsByIds(ctx context.CoreContext, storeName string, cdc *wire.Codec, assetID []string) (map[string]asset.Asset, error) {
	records := map[string]asset.Asset{}

	for _, id := range assetID {
		record, err := queryAsset(ctx, storeName, cdc, id)
		if err != nil {
			return nil, err
		}
		records[id] = *record
	}
	return records, nil
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
		output, err := cdc.MarshalJSON(records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
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

		proposals := make([]asset.Proposal, len(kvs))
		for index, kv := range kvs {
			proposal := asset.Proposal{}
			err = cdc.UnmarshalBinary(kv.Value, &proposal)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}
			proposals[index] = proposal
		}
		output, err := cdc.MarshalJSON(proposals)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
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
			err = cdc.UnmarshalBinary(res, &proposal)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			proposals[index] = ToProposalOutput(proposal, assetID)
		}
		output, err := cdc.MarshalJSON(proposals)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// historyOutput
type historyOutput struct {
	Transfers  []historyTransferOutput       `json:"transfers"`
	Quantity   []historyChangeQuantityOutput `json:"quantity"`
	Properties []historyUpdateProperty       `json:"properties"`
	Materials  []historyAddMaterial          `json:"materials"`
	AssetByID  map[string]asset.Asset        `json:"asset_by_id"`
}

type historyTransferOutput struct {
	Owner sdk.AccAddress `json:"recipient"`
	Time  int64          `json:"time"`
	Memo  string         `json:"memo"`
}

type historyChangeQuantityOutput struct {
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Int        `json:"amount"`
	Type   string         `json:"type"`
	Time   int64          `json:"time"`
	Memo   string         `json:"memo"`
}

type historyUpdateProperty struct {
	Reporter sdk.AccAddress `json:"reporter"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Value    interface{}    `json:"value"`
	Time     int64          `json:"time"`
	Memo     string         `json:"memo"`
}

type historyAddMaterial struct {
	Sender  sdk.AccAddress `json:"sender"`
	Amount  sdk.Int        `json:"amount"`
	AssetID string         `json:"asset_id"`
	Time    int64          `json:"time"`
	Memo    string         `json:"memo"`
}

func formatTxResults(cdc *wire.Codec, res []*ctypes.ResultTx) (txInfos, error) {
	var err error
	out := make(txInfos, len(res))
	for i := range res {
		out[i], err = formatTxResult(cdc, res[i])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func getBlock(ctx context.CoreContext, height *int64) (*ctypes.ResultBlock, error) {
	// get the node
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	// TODO: actually honor the --select flag!
	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	return node.Block(height)
}

func formatTxResult(cdc *wire.Codec, res *ctypes.ResultTx) (txInfo, error) {
	// TODO: verify the proof if requested
	tx, err := parseTx(cdc, res.Tx)
	if err != nil {
		return txInfo{}, err
	}

	info := txInfo{
		Hash:   res.Hash,
		Height: res.Height,
		Tx:     tx,
		Result: res.TxResult,
	}
	return info, nil
}

// txInfo is used to prepare info to display
type txInfo struct {
	Hash   common.HexBytes        `json:"hash"`
	Height int64                  `json:"height"`
	Tx     auth.StdTx             `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
	Time   int64                  `json:"time"`
}

type txInfos []txInfo

func parseTx(cdc *wire.Codec, txBytes []byte) (auth.StdTx, error) {
	var tx auth.StdTx
	err := cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return tx, err
	}
	return tx, nil
}
