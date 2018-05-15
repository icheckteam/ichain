package app

import (
	"encoding/json"

	"github.com/icheckteam/ichain/x/identity"

	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/icheckteam/ichain/x/bank"
	"github.com/icheckteam/ichain/x/ibc"
	"github.com/icheckteam/ichain/x/stake"

	"github.com/icheckteam/ichain/types"
	"github.com/icheckteam/ichain/x/asset"
)

const (
	appName = "IchainApp"
)

// IchainApp Extended ABCI application
type IchainApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	keyMain     *sdk.KVStoreKey
	keyAccount  *sdk.KVStoreKey
	keyIBC      *sdk.KVStoreKey
	keyIdentity *sdk.KVStoreKey
	keyStake    *sdk.KVStoreKey
	keyAsset    *sdk.KVStoreKey

	// Manage getting and setting accounts
	// Manage getting and setting accounts
	accountMapper  sdk.AccountMapper
	coinKeeper     bank.Keeper
	ibcMapper      ibc.Mapper
	stakeKeeper    stake.Keeper
	assetKeeper    asset.Keeper
	identityKeeper identity.Keeper

	// Handle fees
	feeHandler sdk.FeeHandler
}

// NewIchainApp  new ichain application
func NewIchainApp(logger log.Logger, db dbm.DB) *IchainApp {
	// Create app-level codec for txs and accounts.
	var cdc = MakeCodec()
	// create your application object
	var app = &IchainApp{
		BaseApp:     bam.NewBaseApp(appName, cdc, logger, db),
		cdc:         cdc,
		keyMain:     sdk.NewKVStoreKey("main"),
		keyAccount:  sdk.NewKVStoreKey("acc"),
		keyIBC:      sdk.NewKVStoreKey("ibc"),
		keyIdentity: sdk.NewKVStoreKey("identity"),
		keyAsset:    sdk.NewKVStoreKey("asset"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyMain,         // target store
		&types.AppAccount{}, // prototype
	)

	// add handlers
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.assetKeeper = asset.NewKeeper(app.keyAsset, cdc, app.coinKeeper)
	app.identityKeeper = identity.NewKeeper(app.keyIdentity, cdc)
	app.ibcMapper = ibc.NewMapper(cdc, app.keyIBC, ibc.DefaultCodespace)
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("asset", asset.NewHandler(app.assetKeeper)).
		AddRoute("identity", identity.NewHandler(app.identityKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper))

	// Define the feeHandler.
	app.feeHandler = auth.BurnFeeHandler

	// initialize BaseApp
	app.SetTxDecoder(app.txDecoder)
	app.SetInitChainer(app.initChainer)

	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeHandler))
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyAsset, app.keyIdentity)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// MakeCodec Custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	wire.RegisterCrypto(cdc) // Register crypto.
	sdk.RegisterWire(cdc)    // Register Msgs
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	asset.RegisterWire(cdc)

	// register custom AppAccount
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterConcrete(&types.AppAccount{}, "ichain/Account", nil)
	return cdc

	return cdc
}

// custom logic for transaction decoding
func (app *IchainApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	if len(txBytes) == 0 {
		return nil, sdk.ErrTxDecode("txBytes are empty")
	}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec in bank.RegisterAmino.
	err := app.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").Trace(err.Error())
	}
	return tx, nil
}

// custom logic for basecoin initialization
func (app *IchainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		app.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}

// Custom logic for state export
func (app *IchainApp) ExportAppStateJSON() (appState json.RawMessage, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []*types.GenesisAccount{}
	appendAccount := func(acc sdk.Account) (stop bool) {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := types.GenesisState{
		Accounts: accounts,
	}
	return wire.MarshalJSONIndent(app.cdc, genState)
}
