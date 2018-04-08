package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/commands"
	"github.com/icheckteam/ichain/commands"
	"github.com/icheckteam/ichain/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tmlibs/cli"
)

const (
	defaultConfigBaseDir = ".ichainctl"
)

var (
	ichainctlCmd = &cobra.Command{
		Use:   "ichainctl",
		Short: "Ichain light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := types.MakeCodec()
	rpc.AddCommands(ichainctlCmd)
	ichainctlCmd.AddCommand(client.LineBreak)
	tx.AddCommands(ichainctlCmd, cdc)
	ichainctlCmd.AddCommand(client.LineBreak)

	// add ichain-specific commands
	ichainctlCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetAccountDecoder(cdc)),
		)...)
	ichainctlCmd.AddCommand(
		client.PostCommands(
			commands.GetCreateAdminTxCmd(cdc),
			commands.GetCreateOperatorTxCmd(cdc),
			commands.GetCreateAssetAccountTxCmd(cdc),
		)...)
	ichainctlCmd.AddCommand(commands.GetExportPubCmd(cdc))
	//clearchainctlCmd.AddCommand(commands.GetImportPubCmd(cdc))

	// add proxy, version and key info
	ichainctlCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		commands.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(ichainctlCmd, "IC", os.ExpandEnv(defaultConfigBaseDir))
	executor.Execute()
}
