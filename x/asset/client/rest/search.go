package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/icheckteam/ichain/client/tx"
	"github.com/icheckteam/ichain/x/asset"
)

func queryAssetTxs(ctx context.CoreContext, assetID string, cdc *wire.Codec, height int64) ([]tx.TxInfo, error) {
	record, err := getRecord(ctx, assetID, cdc)
	if err != nil {
		return nil, err
	}
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("asset_id='%s'", record.ID)
	page := 0
	perPage := 500
	prove := false
	res, err := node.TxSearch(query, prove, page, perPage)
	if err != nil {
		return nil, err
	}
	info, err := tx.FormatTxResults(cdc, res.Txs)
	if err != nil {
		return nil, err
	}

	if height > 0 {
		info = filterByHeight(info, height)
	}

	// load tx from parents ....
	if record.Parent != "" {
		txs, err := queryAssetTxs(ctx, record.Parent, cdc, record.Height)
		if err != nil {
			return nil, err
		}
		info = append(info, txs...)
	}

	// get block time
	for index, inf := range info {
		block, err := tx.GetBlock(ctx, &inf.Height)
		if err != nil {
			return nil, err
		}
		info[index].Time = block.Block.Header.Time.Unix()
	}
	return info, nil
}

func filterByHeight(infos []tx.TxInfo, height int64) []tx.TxInfo {
	newInfos := make([]tx.TxInfo, len(infos))
	var index = 0
	for _, info := range infos {
		if info.Height > height {
			continue
		}
		newInfos[index] = info
		index++
	}
	return newInfos
}

func newHistoryUpdateProperties(sender sdk.AccAddress, memo string, time int64, props asset.Properties, name string) []asset.HistoryUpdateProperty {
	history := make([]asset.HistoryUpdateProperty, len(props))
	var i = 0
	for index, p := range props {
		if name != "" && name != p.Name {
			continue
		}
		i++
		history[index] = asset.HistoryUpdateProperty{
			Reporter: sender,
			Type:     asset.PropertyTypeToString(p.Type),
			Name:     p.Name,
			Value:    p.GetValue(),
			Time:     time,
			Memo:     memo,
		}
	}
	return history[:i]
}

func filterTxUpdateProperties(infos []tx.TxInfo, name string) []asset.HistoryUpdateProperty {
	history := []asset.HistoryUpdateProperty{}
	for _, info := range infos {
		tx, _ := info.Tx.(auth.StdTx)
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgCreateAsset:
				history = append(history,
					newHistoryUpdateProperties(msg.Sender, tx.Memo, info.Time, msg.Properties, name)...)
				break
			case asset.MsgUpdateProperties:
				history = append(history,
					newHistoryUpdateProperties(msg.Sender, tx.Memo, info.Time, msg.Properties, name)...)
				break
			default:
				break
			}
		}
	}
	return history
}
func filterTxChangeOwner(infos []tx.TxInfo) []asset.HistoryTransferOutput {
	history := []asset.HistoryTransferOutput{}
	for _, info := range infos {
		tx, _ := info.Tx.(auth.StdTx)
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgCreateAsset:
				history = append(history, asset.HistoryTransferOutput{
					Time:  info.Time,
					Memo:  tx.Memo,
					Owner: msg.Sender,
				})
				break
			case asset.MsgAnswerProposal:
				if msg.Role == asset.RoleOwner {
					history = append(history, asset.HistoryTransferOutput{
						Time:  info.Time,
						Memo:  tx.Memo,
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
