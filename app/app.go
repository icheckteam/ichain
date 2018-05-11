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
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"

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
	capKeyMainStore     *sdk.KVStoreKey
	capKeyAccountStore  *sdk.KVStoreKey
	capKeyIBCStore      *sdk.KVStoreKey
	capKeyStakingStore  *sdk.KVStoreKey
	capKeyIdentityStore *sdk.KVStoreKey
	capKeyAssetStore    *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper sdk.AccountMapper

	// Handle fees
	feeHandler sdk.FeeHandler
}

// NewIchainApp  new ichain application
func NewIchainApp(logger log.Logger, dbs map[string]dbm.DB) *IchainApp {
	// Create app-level codec for txs and accounts.
	var cdc = MakeCodec()
	// create your application object
	var app = &IchainApp{
		BaseApp:             bam.NewBaseApp(appName, logger, dbs["main"]),
		cdc:                 cdc,
		capKeyMainStore:     sdk.NewKVStoreKey("main"),
		capKeyAccountStore:  sdk.NewKVStoreKey("acc"),
		capKeyIBCStore:      sdk.NewKVStoreKey("ibc"),
		capKeyStakingStore:  sdk.NewKVStoreKey("stake"),
		capKeyIdentityStore: sdk.NewKVStoreKey("identity"),
		capKeyAssetStore:    sdk.NewKVStoreKey("asset"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	).Seal()

	// add handlers
	coinKeeper := bank.NewCoinKeeper(app.accountMapper)
	assetKeeper := asset.NewKeeper(app.capKeyAssetStore, cdc, coinKeeper)
	identityKeeper := identity.NewKeeper(app.capKeyIdentityStore, cdc, app.accountMapper)
	ibcMapper := ibc.NewIBCMapper(cdc, app.capKeyIBCStore)
	app.Router().
		AddRoute("bank", bank.NewHandler(coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(ibcMapper, coinKeeper)).
		AddRoute("asset", asset.NewHandler(assetKeeper)).
		AddRoute("identity", identity.NewHandler(identityKeeper))

	// Define the feeHandler.
	app.feeHandler = auth.BurnFeeHandler

	// initialize BaseApp
	app.SetTxDecoder(app.txDecoder)
	app.SetInitChainer(app.initChainer)
	app.MountStoreWithDB(app.capKeyMainStore, sdk.StoreTypeIAVL, dbs["main"])
	app.MountStoreWithDB(app.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	app.MountStoreWithDB(app.capKeyIBCStore, sdk.StoreTypeIAVL, dbs["ibc"])
	app.MountStoreWithDB(app.capKeyStakingStore, sdk.StoreTypeIAVL, dbs["staking"])
	app.MountStoreWithDB(app.capKeyIdentityStore, sdk.StoreTypeIAVL, dbs["identity"])
	app.MountStoreWithDB(app.capKeyAssetStore, sdk.StoreTypeIAVL, dbs["asset"])

	// NOTE: Broken until #532 lands
	//app.MountStoresIAVL(app.capKeyMainStore, app.capKeyIBCStore, app.capKeyStakingStore)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeHandler))
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// MakeCodec Custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.SendMsg{}, "ichain/Send", nil)
	cdc.RegisterConcrete(bank.IssueMsg{}, "ichain/Issue", nil)
	cdc.RegisterConcrete(ibc.IBCTransferMsg{}, "ichain/IBCTransferMsg", nil)
	cdc.RegisterConcrete(ibc.IBCReceiveMsg{}, "ichain/IBCReceiveMsg", nil)

	cdc.RegisterConcrete(asset.RegisterMsg{}, "ichain/RegisterMsg", nil)
	cdc.RegisterConcrete(asset.AddQuantityMsg{}, "ichain/AddQuantityMsg", nil)
	cdc.RegisterConcrete(asset.SubtractQuantityMsg{}, "ichain/SubtractQuantityMsg", nil)
	cdc.RegisterConcrete(asset.UpdateAttrMsg{}, "ichain/UpdateAttrMsg", nil)

	cdc.RegisterConcrete(identity.CreateMsg{}, "ichain/ClaimIssueMsg", nil)
	cdc.RegisterConcrete(identity.RevokeMsg{}, "ichain/RevokeMsg", nil)

	// Register AppAccount
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterConcrete(&types.AppAccount{}, "ichain/Account", nil)

	// Register crypto.
	wire.RegisterCrypto(cdc)

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
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for basecoin initialization
func (app *IchainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := json.Unmarshal(stateJSON, genesisState)
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
