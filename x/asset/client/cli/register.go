package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/spf13/cobra"
)

// CreateAssetCmd will create asset
func CreateAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := Commander{cdc}
	cmd := &cobra.Command{
		Use:   "create-asset",
		Short: "create new asset",
		RunE:  cmdr.registerCmd,
	}
	return cmd
}

type Commander struct {
	Cdc *wire.Codec
}

func (c Commander) registerCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()
	// get the from address
	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	// build message
	msg := buildCreateAssetMsg(from)
	// default to next sequence number if none provided
	ctx, err = context.EnsureSequence(ctx)
	if err != nil {
		return err
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, err := ctx.SignBuildBroadcast(ctx.FromAddressName, msg, c.Cdc)
	if err != nil {
		return err
	}
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func buildCreateAssetMsg(creator sdk.Address) sdk.Msg {
	return asset.NewRegisterMsg(
		creator,
		"1",
		"1",
		1,
		"1",
		"1",
	)
}
