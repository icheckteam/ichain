package rest

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

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

		output, err := cdc.MarshalJSON(a)
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
		output, err := cdc.MarshalJSON(items)
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

func QueryInventoryHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
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
			var assetId string

			err = cdc.UnmarshalBinary(kv.Value, &assetId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(asset.GetInventoryKey(address, assetId), storeName)
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

func QueryReporterAssetsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
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
			var assetId string

			err = cdc.UnmarshalBinary(kv.Value, &assetId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(asset.GetReporterAssetKey(address, assetId), storeName)
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

func QueryProposalsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
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

func QueryAccountProposalsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
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
			var assetId string

			err = cdc.UnmarshalBinary(kv.Value, &assetId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode proposal. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(asset.GetProposalAccountKey(address, assetId), storeName)
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

func HistortyHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		recordId := vars["id"]
		history, err := queryHistory(ctx, storeName, cdc, recordId, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
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

func queryHistory(ctx context.CoreContext, storeName string, cdc *wire.Codec, recordID string, fromHeight int64) (historyOutput, error) {
	historyOutput := historyOutput{}
	record, err := queryAsset(ctx, storeName, cdc, recordID)
	if err != nil {
		return historyOutput, err
	}
	history, err := searchTxs(ctx, cdc, recordID, fromHeight)
	if err != nil {
		return historyOutput, err
	}

	if record.Parent != "" {
		otherHistory, err := queryHistory(ctx, storeName, cdc, record.Parent, record.Height)
		if err != nil {
			return historyOutput, err
		}
		history.Properties = append(history.Properties, otherHistory.Properties...)
		history.Quantity = append(history.Quantity, otherHistory.Quantity...)
		history.Transfers = append(history.Transfers, otherHistory.Transfers...)
	}
	return history, nil
}

func searchTxs(ctx context.CoreContext, cdc *wire.Codec, assetID string, fromHeight int64) (historyOutput, error) {
	history := historyOutput{}
	// get the node
	node, err := ctx.GetNode()
	if err != nil {
		return history, err
	}

	prove := !viper.GetBool(client.FlagTrustNode)
	// TODO: take these as args
	page := 0
	perPage := 300
	res, err := node.TxSearch(fmt.Sprintf("asset_id='%s'", assetID), prove, page, perPage)
	if err != nil {
		return history, err
	}

	infos, err := formatTxResults(cdc, res.Txs)
	if err != nil {
		return history, err
	}

	proposals := map[string]asset.MsgCreateProposal{}
	recordAmount := map[string]sdk.Int{}

	for index, info := range infos {
		if info.Height < fromHeight {
			continue
		}
		block, err := getBlock(ctx, &info.Height)
		if err != nil {
			return history, err
		}
		infos[index].Time = block.Block.Time.Unix()

	}

	sort.SliceStable(infos, func(i, j int) bool { return infos[i].Time < infos[j].Time })
	for _, info := range infos {
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgCreateAsset:
				recordAmount[msg.AssetID] = msg.Quantity
				actionType := "add"
				if len(msg.Parent) > 0 && msg.Parent == assetID {
					actionType = "subtract"
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

	return history, nil
}

// historyOutput
type historyOutput struct {
	Transfers  []historyTransferOutput       `json:"transfers"`
	Quantity   []historyChangeQuantityOutput `json:"quantity"`
	Properties []historyUpdateProperty       `json:"properties"`
	Materials  []historyAddMaterial          `json:"materials"`
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
