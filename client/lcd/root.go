package lcd

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/log"

	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
	cmn "github.com/tendermint/tmlibs/common"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	rpc "github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/wire"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	tx "github.com/icheckteam/ichain/client/tx"
	version "github.com/icheckteam/ichain/version"
	asset "github.com/icheckteam/ichain/x/asset/client/rest"
	bank "github.com/icheckteam/ichain/x/bank/client/rest"
	ibc "github.com/icheckteam/ichain/x/ibc/client/rest"
	identity "github.com/icheckteam/ichain/x/identity/client/rest"
	stake "github.com/icheckteam/ichain/x/stake/client/rest"
	warranty "github.com/icheckteam/ichain/x/warranty/client/rest"
)

const (
	flagListenAddr = "laddr"
	flagCORS       = "cors"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE:  startRESTServerFn(cdc),
	}
	cmd.Flags().StringP(flagListenAddr, "a", "tcp://localhost:1317", "Address for server to listen on")
	cmd.Flags().String(flagCORS, "", "Set to domains that can make CORS requests (* for all)")
	cmd.Flags().StringP(client.FlagChainID, "c", "", "ID of chain we connect to")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:46657", "Node to connect to")
	return cmd
}

func startRESTServerFn(cdc *wire.Codec) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		listenAddr := viper.GetString(flagListenAddr)
		handler := createHandler(cdc)
		logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
			With("module", "rest-server")
		listener, err := tmserver.StartHTTPServer(listenAddr, handler, logger)
		if err != nil {
			return err
		}

		// Wait forever and cleanup
		cmn.TrapSignal(func() {
			err := listener.Close()
			logger.Error("Error closing listener", "err", err)
		})
		return nil
	}
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", version.RequestHandler).Methods("GET")

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	ctx := context.NewCoreContextFromViper()

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(ctx, r)
	tx.RegisterRoutes(r, cdc)
	auth.RegisterRoutes(ctx, r, cdc, "acc")
	bank.RegisterRoutes(ctx, r, cdc, kb)
	ibc.RegisterRoutes(ctx, r, cdc, kb)
	asset.RegisterRoutes(ctx, r, cdc, kb, "asset")
	identity.RegisterRoutes(ctx, r, cdc, kb, "identity")
	stake.RegisterRoutes(ctx, r, cdc, kb)
	warranty.RegisterRoutes(ctx, r, cdc, kb, "warranty")
	return r
}
