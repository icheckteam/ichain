package identity

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func makeTestCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(CreateMsg{}, "test/asset/CreateMsg", nil)
	cdc.RegisterConcrete(RevokeMsg{}, "test/asset/RevokeMsg", nil)
	wire.RegisterCrypto(cdc)

	return cdc
}

// hogpodge of all sorts of input required for testing
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	db := dbm.NewMemDB()
	keyStake := sdk.NewKVStoreKey("identity")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyStake, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, nil)
	cdc := makeTestCodec()
	keeper := NewKeeper(keyStake, cdc)
	return ctx, keeper
}
