package app

import (
	"encoding/json"
	"os"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"

	"github.com/icheckteam/ichain/types"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/icheckteam/ichain/x/identity"
	"github.com/icheckteam/ichain/x/insurance"
	"github.com/icheckteam/ichain/x/invoice"
	"github.com/icheckteam/ichain/x/shipping"
)

const (
	appName = "IchainApp"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.ichaincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.ichaind")
)

// IchainApp Extended ABCI application
type IchainApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	keyMain      *sdk.KVStoreKey
	keyAccount   *sdk.KVStoreKey
	keyIBC       *sdk.KVStoreKey
	keyIdentity  *sdk.KVStoreKey
	keyStake     *sdk.KVStoreKey
	keyAsset     *sdk.KVStoreKey
	keyInsurance *sdk.KVStoreKey
	keyShipping  *sdk.KVStoreKey
	keyInvoice   *sdk.KVStoreKey
	keySlashing  *sdk.KVStoreKey
	keyGov       *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountMapper
	bankKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	govKeeper           gov.Keeper
	keyFeeCollection    *sdk.KVStoreKey

	assetKeeper     asset.Keeper
	identityKeeper  identity.Keeper
	insuranceKeeper insurance.Keeper
	shippingKeeper  shipping.Keeper
	invoiceKeeper   invoice.InvoiceKeeper
}

// NewIchainApp  new ichain application
func NewIchainApp(logger log.Logger, db dbm.DB) *IchainApp {
	// Create app-level codec for txs and accounts.
	var cdc = MakeCodec()
	// create your application object
	var app = &IchainApp{
		BaseApp:          bam.NewBaseApp(appName, cdc, logger, db),
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyIBC:           sdk.NewKVStoreKey("ibc"),
		keyIdentity:      sdk.NewKVStoreKey("identity"),
		keyAsset:         sdk.NewKVStoreKey("asset"),
		keyInsurance:     sdk.NewKVStoreKey("insurance"),
		keyShipping:      sdk.NewKVStoreKey("shipping"),
		keyInvoice:       sdk.NewKVStoreKey("invoice"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyStake:         sdk.NewKVStoreKey("stake"),
		keyGov:           sdk.NewKVStoreKey("gov"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,        // target store
		types.ProtoAppAccount, // prototype
	)

	// add handlers
	app.bankKeeper = bank.NewKeeper(app.accountMapper)
	app.assetKeeper = asset.NewKeeper(app.keyAsset, cdc)
	app.identityKeeper = identity.NewKeeper(app.keyIdentity, cdc)
	app.ibcMapper = ibc.NewMapper(cdc, app.keyIBC, ibc.DefaultCodespace)
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.bankKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.insuranceKeeper = insurance.NewKeeper(app.keyInsurance, cdc, app.assetKeeper)
	app.shippingKeeper = shipping.NewKeeper(app.keyShipping, cdc, app.assetKeeper)
	app.invoiceKeeper = invoice.NewInvoiceKeeper(app.keyInvoice, cdc, app.assetKeeper)
	app.slashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, app.stakeKeeper, app.RegisterCodespace(slashing.DefaultCodespace))
	app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.bankKeeper, app.stakeKeeper, app.RegisterCodespace(gov.DefaultCodespace))
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, app.keyFeeCollection)
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute("asset", asset.NewHandler(app.assetKeeper)).
		AddRoute("identity", identity.NewHandler(app.identityKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("insurance", insurance.NewHandler(app.insuranceKeeper)).
		AddRoute("shipping", shipping.NewHandler(app.shippingKeeper)).
		AddRoute("invoice", invoice.MakeHandle(app.invoiceKeeper)).
		AddRoute("slashing", slashing.NewHandler(app.slashingKeeper)).
		AddRoute("gov", gov.NewHandler(app.govKeeper))

	// initialize Ichain App
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.MountStoresIAVL(
		app.keyMain,
		app.keyAccount,
		app.keyIBC,
		app.keyStake,
		app.keySlashing,
		app.keyGov,
		app.keyFeeCollection,

		app.keyAsset,
		app.keyIdentity,
		app.keyShipping,
		app.keyInsurance,
		app.keyInvoice,
	)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// MakeCodec Custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	ibc.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	slashing.RegisterWire(cdc)
	gov.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	sdk.RegisterWire(cdc)
	wire.RegisterCrypto(cdc)

	asset.RegisterWire(cdc)
	insurance.RegisterWire(cdc)
	shipping.RegisterWire(cdc)
	invoice.RegisterWire(cdc)
	identity.RegisterWire(cdc)
	// register custom AppAccount
	cdc.RegisterConcrete(&types.AppAccount{}, "ichain/Account", nil)
	return cdc
}

// application updates every end block
func (app *IchainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
func (app *IchainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)

	tags, _ := gov.EndBlocker(ctx, app.govKeeper)

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// Custom logic for basecoin initialization
func (app *IchainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.BaseAccount.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	// load the initial stake information
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	gov.InitGenesis(ctx, app.govKeeper, gov.DefaultGenesisState())
	return abci.ResponseInitChain{}
}

// Custom logic for state export
func (app *IchainApp) ExportAppStateJSON() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []types.GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := types.NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := types.GenesisState{
		Accounts:  accounts,
		StakeData: stake.WriteGenesis(ctx, app.stakeKeeper),
	}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return
}
