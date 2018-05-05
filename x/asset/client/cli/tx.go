package cli

import (
	"github.com/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

// CreateAssetCmd will create asset
func CreateAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-asset",
		Short: "create new asset",
	}
	return cmd
}
