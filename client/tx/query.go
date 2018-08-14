package tx

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// QueryTxCmd  Get the default command for a tx query
func QueryTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Matches this txhash over all committed blocks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// find the key to look up the account
			hashHexStr := args[0]
			trustNode := viper.GetBool(client.FlagTrustNode)

			output, err := QueryTx(cdc, context.NewCLIContext(), hashHexStr, trustNode)
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			return nil
		},
	}

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:46657", "Node to connect to")

	// TODO: change this to false when we can
	cmd.Flags().Bool(client.FlagTrustNode, true, "Don't verify proofs for responses")
	return cmd
}

// QueryTx ...
func QueryTx(cdc *wire.Codec, ctx context.CLIContext, hashHexStr string, trustNode bool) ([]byte, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	// get the node
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Tx(hash, !trustNode)
	if err != nil {

		return nil, err
	}
	info, err := FormatTxResult(cdc, res)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSON(info)
}

// FormatTxResult ...
func FormatTxResult(cdc *wire.Codec, res *ctypes.ResultTx) (TxInfo, error) {
	// TODO: verify the proof if requested
	tx, err := ParseTx(cdc, res.Tx)
	if err != nil {
		return TxInfo{}, err
	}

	info := TxInfo{
		Hash:   res.Hash,
		Height: res.Height,
		Tx:     tx,
		Result: res.TxResult,
	}
	return info, nil
}

// TxInfo  is used to prepare info to display
type TxInfo struct {
	Hash   common.HexBytes        `json:"hash"`
	Height int64                  `json:"height"`
	Tx     sdk.Tx                 `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
	Time   int64                  `json:"time"`
}

// ParseTx ...
func ParseTx(cdc *wire.Codec, txBytes []byte) (auth.StdTx, error) {
	var tx auth.StdTx
	err := cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return tx, err
	}
	return tx, nil
}

// REST

// QueryTxRequestHandlerFn  transaction query REST handler
func QueryTxRequestHandlerFn(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]
		trustNode, err := strconv.ParseBool(r.FormValue("trust_node"))
		// trustNode defaults to true
		if err != nil {
			trustNode = true
		}

		output, err := QueryTx(cdc, ctx, hashHexStr, trustNode)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
	}
}
