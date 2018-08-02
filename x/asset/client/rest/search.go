package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
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

	// load tx from parents ....
	if record.Parent != "" {
		txs, err := queryAssetTxs(ctx, record.Parent, cdc, record.Height)
		if err != nil {
			return nil, err
		}
		for _, tx := range txs {
			if tx.Height > record.Height {
				continue
			}
			info = append(info, tx)
		}
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

func filterTxUpdateProperties(infos []tx.TxInfo, name string) []HistoryUpdateProperty {
	history := []HistoryUpdateProperty{}
	for _, info := range infos {
		tx, _ := info.Tx.(auth.StdTx)
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgUpdateProperties:
				for _, p := range msg.Properties {
					if name != "" && name != p.Name {
						continue
					}
					history = append(history, HistoryUpdateProperty{
						Reporter: msg.Sender,
						Type:     asset.PropertyTypeToString(p.Type),
						Name:     p.Name,
						Value:    p.GetValue(),
						Time:     info.Time,
						Memo:     tx.Memo,
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
func filterTxChangeOwner(infos []tx.TxInfo) []HistoryTransferOutput {
	history := []HistoryTransferOutput{}
	for _, info := range infos {
		tx, _ := info.Tx.(auth.StdTx)
		for _, msg := range info.Tx.GetMsgs() {
			switch msg := msg.(type) {
			case asset.MsgAnswerProposal:
				if msg.Role == asset.RoleOwner {
					history = append(history, HistoryTransferOutput{
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
