package rest

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/icheckteam/ichain/client/errors"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
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

		record, err := queryAsset(ctx, storeName, cdc, vars["id"])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't decode asset. Error: %s", err.Error())))
			return
		}

		assetOutput := AssetOutput{
			Asset:     *record,
			Materials: map[string]asset.Asset{},
		}
		for _, material := range record.Materials {
			record, err := queryAsset(ctx, storeName, cdc, material.Denom)
			assetOutput.Materials[record.ID] = *record
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("Couldn't decode asset. Error: %s", err.Error())))
				return
			}
		}

		output, err := cdc.MarshalJSON(assetOutput)
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
		var output []byte
		vars := mux.Vars(r)
		node, err := ctx.GetNode()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't get current Node information. Error: %s", err.Error())))
			return
		}

		query := fmt.Sprintf("asset_id='%s'", vars["id"])
		page := 0
		perPage := 500
		prove := false
		res, err := node.TxSearch(query, prove, page, perPage)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		info, err := formatTxResults(cdc, res.Txs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query txs. Error: %s", err.Error())))
			return
		}
		// success
		output, err = cdc.MarshalJSON(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output) // write
	}
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

func queryAccountAssets(ctx context.CoreContext, storeName string, cdc *wire.Codec, account string) ([]asset.Asset, error) {
	address, err := sdk.AccAddressFromBech32(account)
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

func queryInventoryHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		kvs, err := ctx.QuerySubspace(cdc, asset.GetInventoryByAccountKey(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Couldn't get proposals. Error: %s", err.Error())))
			return
		}

		assets := make([]asset.Asset, len(kvs))
		for index, kv := range kvs {
			a := asset.Asset{}
			var assetID string

			err = cdc.UnmarshalBinary(kv.Value, &assetID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(asset.GetInventoryKey(address, assetID), storeName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}
			err = cdc.UnmarshalBinary(res, &a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			assets[index] = a
		}
		output, err := cdc.MarshalJSON(assets)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
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

		assets := make([]asset.Asset, len(kvs))
		for index, kv := range kvs {
			a := asset.Asset{}
			var assetID string

			err = cdc.UnmarshalBinary(kv.Value, &assetID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(asset.GetReporterAssetKey(address, assetID), storeName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}
			err = cdc.UnmarshalBinary(res, &a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			assets[index] = a
		}
		output, err := cdc.MarshalJSON(assets)
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

		proposals := make([]asset.Proposal, len(kvs))
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

// HistortyHandlerFn ...
func HistortyHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		recordID := vars["id"]
		history, indexRecord, err := queryHistory(ctx, storeName, cdc, recordID, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		history.AssetByID = map[string]asset.Asset{}
		for _, id := range indexRecord {
			record, err := queryAsset(ctx, storeName, cdc, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			history.AssetByID[record.ID] = *record
		}

		output, err := cdc.MarshalJSON(history)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func queryHistory(ctx context.CoreContext, storeName string, cdc *wire.Codec, recordID string, fromHeight int64) (historyOutput, []string, error) {
	historyOutput := historyOutput{}
	record, err := queryAsset(ctx, storeName, cdc, recordID)
	if err != nil {
		return historyOutput, nil, err
	}
	history, indexRecord, err := searchTxs(ctx, cdc, recordID, fromHeight)
	if err != nil {
		return historyOutput, nil, err
	}

	if record.Parent != "" {
		otherHistory, otherIndexRecord, err := queryHistory(ctx, storeName, cdc, record.Parent, record.Height)
		if err != nil {
			return historyOutput, nil, err
		}
		history.Properties = append(history.Properties, otherHistory.Properties...)
		history.Quantity = append(history.Quantity, otherHistory.Quantity...)
		history.Transfers = append(history.Transfers, otherHistory.Transfers...)

		indexRecord = append(indexRecord, otherIndexRecord...)

	}
	return history, unique(indexRecord), nil
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func searchTxs(ctx context.CoreContext, cdc *wire.Codec, assetID string, fromHeight int64) (historyOutput, []string, error) {
	history := historyOutput{}
	// get the node
	node, err := ctx.GetNode()
	if err != nil {
		return history, nil, err
	}

	prove := !viper.GetBool(client.FlagTrustNode)
	// TODO: take these as args
	page := 0
	perPage := 300
	res, err := node.TxSearch(fmt.Sprintf("asset_id='%s'", assetID), prove, page, perPage)
	if err != nil {
		return history, nil, err
	}

	infos, err := formatTxResults(cdc, res.Txs)
	if err != nil {
		return history, nil, err
	}

	proposals := map[string]asset.MsgCreateProposal{}
	recordAmount := map[string]sdk.Int{}
	allRecords := []string{}
	allRecordsIndex := map[string]bool{}

	for index, info := range infos {
		if info.Height < fromHeight {
			continue
		}
		block, err := getBlock(ctx, &info.Height)
		if err != nil {
			return history, nil, err
		}
		infos[index].Time = block.Block.Time.Unix()
		allRecords = append(allRecords)
	}

	sort.SliceStable(infos, func(i, j int) bool { return infos[i].Time < infos[j].Time })

	allRecordsIndex[assetID] = true
	allRecords = append(allRecords, assetID)

	for _, info := range infos {
		if info.Height < fromHeight {
			continue
		}
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgCreateAsset:
				recordAmount[msg.AssetID] = msg.Quantity
				actionType := "add"
				if len(msg.Parent) > 0 && msg.Parent == assetID {
					actionType = "subtract"

					if allRecordsIndex[msg.Parent] == false {
						allRecordsIndex[msg.Parent] = true
						allRecords = append(allRecords, msg.Parent)
					}

				}
				history.Quantity = append(history.Quantity, historyChangeQuantityOutput{
					Sender: msg.Sender,
					Type:   actionType,
					Amount: msg.Quantity,
					Time:   info.Time,
					Memo:   info.Tx.Memo,
				})
				break
			case asset.MsgCreateProposal:
				proposals[msg.Recipient.String()] = msg
				break
			case asset.MsgAnswerProposal:
				proposal := proposals[msg.Recipient.String()]
				if msg.Response == asset.StatusAccepted {
					// index transfer asset
					if proposal.Role == asset.RoleOwner {
						history.Transfers = append(history.Transfers, historyTransferOutput{
							Sender:    proposal.Sender,
							Recipient: proposal.Recipient,
							Time:      info.Time,
							Memo:      info.Tx.Memo,
							Amount:    recordAmount[msg.AssetID],
						})
					}
				}
				break
			case asset.MsgUpdateProperties:
				for _, property := range msg.Properties {
					history.Properties = append(history.Properties, historyUpdateProperty{
						Sender: msg.Sender,
						Type:   asset.PropertyTypeToString(property.Type),
						Name:   property.Name,
						Value:  property.GetValue(),
						Time:   info.Time,
						Memo:   info.Tx.Memo,
					})
				}
				break
			case asset.MsgAddMaterials:

				for _, amount := range msg.Amount {
					actionType := "add"
					if msg.AssetID == assetID {
						actionType = "subtract"
					}

					if allRecordsIndex[amount.Denom] == false {
						allRecordsIndex[amount.Denom] = true
						allRecords = append(allRecords, amount.Denom)
					}

					history.Quantity = append(history.Quantity, historyChangeQuantityOutput{
						Sender: msg.Sender,
						Type:   actionType,
						Amount: amount.Amount,
						Time:   info.Time,
						Memo:   info.Tx.Memo,
					})

					history.Materials = append(history.Materials, historyAddMaterial{
						Amount:  amount.Amount,
						AssetID: amount.Denom,
						Sender:  msg.Sender,
						Time:    info.Time,
						Memo:    info.Tx.Memo,
					})
				}
				break
			default:
				break
			}

		}
	}
	return history, allRecords, nil
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
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	Time      int64          `json:"time"`
	Memo      string         `json:"memo"`
	Amount    sdk.Int        `json:"amount"`
}

type historyChangeQuantityOutput struct {
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Int        `json:"amount"`
	Type   string         `json:"type"`
	Time   int64          `json:"time"`
	Memo   string         `json:"memo"`
}

type historyUpdateProperty struct {
	Sender sdk.AccAddress `json:"sender"`
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Value  interface{}    `json:"value"`
	Time   int64          `json:"time"`
	Memo   string         `json:"memo"`
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
