package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/icheckteam/ichain/app"
)

// rootCmd is the entry point for this binary
var (
	context = server.NewDefaultContext()
	rootCmd = &cobra.Command{
		Use:   "ichaind",
		Short: "Ichaind Daemon (server)",
	}
)

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dbMain, err := dbm.NewGoLevelDB("ichain", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbAcc, err := dbm.NewGoLevelDB("ichain-acc", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbIBC, err := dbm.NewGoLevelDB("ichain-ibc", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbStaking, err := dbm.NewGoLevelDB("ichain-staking", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbs := map[string]dbm.DB{
		"main":    dbMain,
		"acc":     dbAcc,
		"ibc":     dbIBC,
		"staking": dbStaking,
	}
	bapp := app.NewIchainApp(logger, dbs)
	return bapp, nil
}

func main() {
	server.AddCommands(rootCmd, server.DefaultGenAppState, generateApp, context)
	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.ichaind")
	context.Config.RootDir = rootDir
	executor := cli.PrepareBaseCmd(rootCmd, "IC", rootDir)
	executor.Execute()
}
