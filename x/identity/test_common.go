package identity

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tmlibs/log"

	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func makeTestCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)

	return cdc
}

// hogpodge of all sorts of input required for testing
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	db := dbm.NewMemDB()
	keyIdentity := sdk.NewKVStoreKey("identity")
	keyAcc := sdk.NewKVStoreKey("acc")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyIdentity, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, nil, log.NewNopLogger())
	cdc := makeTestCodec()

	accountMapper := auth.NewAccountMapper(
		cdc,                 // amino codec
		keyAcc,              // target store
		&auth.BaseAccount{}, // prototype
	)
	ck := bank.NewKeeper(accountMapper)
	keeper := NewKeeper(keyIdentity, cdc, ck)
	return ctx, keeper
}
