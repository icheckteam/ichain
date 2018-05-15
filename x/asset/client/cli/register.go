package cli

import (
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
