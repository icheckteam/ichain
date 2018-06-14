package main

import (
	"encoding/json"
	"os"

	"github.com/icheckteam/ichain/app"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/server"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "ichaind",
		Short:             "Ichain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "ichain"),
		server.ConstructAppExporter(exportAppState, "ichain"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.ichaind")
	executor := cli.PrepareBaseCmd(rootCmd, "IC", rootDir)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewIchainApp(logger, db)
}

func exportAppState(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewIchainApp(logger, db)
	return bapp.ExportAppStateJSON()
}
