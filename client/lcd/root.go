package lcd

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"

	cmn "github.com/tendermint/tendermint/libs/common"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	gov "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/client/rest"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/client/rest"
	stake "github.com/cosmos/cosmos-sdk/x/stake/client/rest"

	"github.com/icheckteam/ichain/client/rpc"
	"github.com/icheckteam/ichain/client/signature"
	"github.com/icheckteam/ichain/client/tx"
	asset "github.com/icheckteam/ichain/x/asset/client/rest"
	identity "github.com/icheckteam/ichain/x/identity/client/rest"
	insurance "github.com/icheckteam/ichain/x/insurance/client/rest"
	invoice "github.com/icheckteam/ichain/x/invoice/client/rest"
	shipping "github.com/icheckteam/ichain/x/shipping/client/rest"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"

	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			allowedOrigins := viper.GetString(flagCORS)
			listenAddr := viper.GetString(flagListenAddr)
			handler := createHandler(cdc)
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
			maxOpen := viper.GetInt(flagMaxOpenConnections)

			c := cors.New(cors.Options{
				AllowedOrigins:   []string{allowedOrigins},
				AllowCredentials: true,
				AllowedMethods:   []string{"PUT", "GET", "POST", "DELETE"},
				// Enable Debugging for testing, consider disabling in production
				Debug: true,
			})

			listener, err := tmserver.StartHTTPServer(
				listenAddr, c.Handler(handler), logger,
				tmserver.Config{MaxOpenConnections: maxOpen},
			)
			if err != nil {
				return err
			}

			logger.Info("REST server started")

			// wait forever and cleanup
			cmn.TrapSignal(func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})

			return nil
		},
	}
	cmd.Flags().String(flagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().String(flagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")
	return cmd
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	ctx := context.NewCoreContextFromViper()

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node_version", NodeVersionRequestHandler(cdc, ctx)).Methods("GET")

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(ctx, r, cdc)
	tx.RegisterRoutes(ctx, r, cdc)
	auth.RegisterRoutes(ctx, r, cdc, "acc")
	bank.RegisterRoutes(ctx, r, cdc, kb)
	ibc.RegisterRoutes(ctx, r, cdc, kb)
	stake.RegisterRoutes(ctx, r, cdc, kb)
	slashing.RegisterRoutes(ctx, r, cdc, kb)
	gov.RegisterRoutes(ctx, r, cdc)

	signature.RegisterRoutes(r)

	asset.RegisterRoutes(ctx, r, cdc, kb, "asset")
	identity.RegisterRoutes(ctx, r, cdc, kb, "identity")
	insurance.RegisterRoutes(ctx, r, cdc, kb, "insurance")
	shipping.RegisterRoutes(ctx, r, cdc, kb, "shipping")
	invoice.RegisterHTTPHandle(r, ctx, cdc, kb)

	return r
}
