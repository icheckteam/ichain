package app

import (
	"encoding/json"
	"errors"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/icheckteam/ichain/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/go-crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"

	// bonded tokens given to genesis validators/accounts
	freeFermionVal  = int64(100)
	freeFermionsAcc = int64(50)
)

// get app init parameters for server init command
func GaiaAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         GaiaAppGenTx,
		AppGenState:      GaiaAppGenStateJSON,
	}
}

// simple genesis tx
type GaiaGenTx struct {
	Name    string        `json:"name"`
	Address sdk.Address   `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
}

// Generate a gaia genesis transaction with flags
func GaiaAppGenTx(cdc *wire.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	clientRoot := viper.GetString(flagClientHome)
	overwrite := viper.GetBool(flagOWK)
	name := viper.GetString(flagName)
	if name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}

	var addr sdk.Address
	var secret string
	addr, secret, err = server.GenerateSaveCoinKey(clientRoot, name, "1234567890", overwrite)
	if err != nil {
		return
	}
	mm := map[string]string{"secret": secret}
	var bz []byte
	bz, err = cdc.MarshalJSON(mm)
	if err != nil {
		return
	}
	cliPrint = json.RawMessage(bz)
	appGenTx, _, validator, err = GaiaAppGenTxNF(cdc, pk, addr, name, overwrite)
	return
}

// Generate a gaia genesis transaction without flags
func GaiaAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.Address, name string, overwrite bool) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var bz []byte
	gaiaGenTx := GaiaGenTx{
		Name:    name,
		Address: addr,
		PubKey:  pk,
	}
	bz, err = wire.MarshalJSONIndent(cdc, gaiaGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  freeFermionVal,
	}
	return
}

// Create the core parameters for genesis initialization for gaia
// note that the pubkey input is this machines pubkey
func GaiaAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState types.GenesisState, err error) {

	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// start with the default staking genesis state
	stakeData := stake.DefaultGenesisState()

	// get genesis flag account information
	genaccs := make([]types.GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {

		var genTx GaiaGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		// create the genesis account, give'm few steaks and a buncha token with there name
		accAuth := types.AppAccount{
			BaseAccount: auth.NewBaseAccountWithAddress(genTx.Address),
		}
		accAuth.Coins = sdk.Coins{
			{genTx.Name + "Token", 1000},
			{"steak", freeFermionsAcc},
		}
		acc := types.NewGenesisAccount(&accAuth)
		genaccs[i] = acc
		stakeData.Pool.LooseUnbondedTokens += freeFermionsAcc // increase the supply

		// add the validator
		if len(genTx.Name) > 0 {
			desc := stake.NewDescription(genTx.Name, "", "", "")
			validator := stake.NewValidator(genTx.Address, genTx.PubKey, desc)
			validator.PoolShares = stake.NewBondedShares(sdk.NewRat(freeFermionVal))
			stakeData.Validators = append(stakeData.Validators, validator)

			// pool logic
			stakeData.Pool.BondedTokens += freeFermionVal
			stakeData.Pool.BondedShares = sdk.NewRat(stakeData.Pool.BondedTokens)
		}
	}

	// create the final app state
	genesisState = types.GenesisState{
		Accounts:  genaccs,
		StakeData: stakeData,
	}
	return
}

// GaiaAppGenState but with JSON
func GaiaAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := GaiaAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
