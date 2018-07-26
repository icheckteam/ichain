package insurance

import (
	"bytes"
	"encoding/hex"
	"os"
	"strconv"
	"testing"

	"github.com/icheckteam/ichain/types"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// dummy addresses used for testing
var (
	addrs       = createTestAddrs(100)
	pks         = createTestPubKeys(100)
	emptyAddr   sdk.AccAddress
	emptyPubkey crypto.PubKey
)

// XXX reference the common declaration of this function
func subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}

func makeTestCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	bank.RegisterWire(cdc)
	RegisterWire(cdc)

	// Register AppAccount
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&types.AppAccount{}, "test/asset/Account", nil)
	wire.RegisterCrypto(cdc)

	return cdc
}

// hogpodge of all sorts of input required for testing
func createTestInput(t *testing.T, isCheckTx bool, initCoins int64) (sdk.Context, Keeper) {
	db := dbm.NewMemDB()
	keyInsurance := sdk.NewKVStoreKey("insurance")
	keyAsset := sdk.NewKVStoreKey("asset")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyInsurance, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAsset, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	cdc := makeTestCodec()
	assetKeeper := asset.NewKeeper(keyAsset, cdc)

	keeper := NewKeeper(keyInsurance, cdc, assetKeeper)
	return ctx, keeper
}
func NewPubKey(pk string) (res ed25519.PubKeyEd25519) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

// for incode address generation
func TestAddr(addr string, bech string) sdk.AccAddress {

	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(bechres, res) != 0 {
		panic("Bech decode and hex decode don't match")
	}

	return res
}

// nolint: unparam
func createTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, TestAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

// nolint: unparam
func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

//_____________________________________________________________________________________

// does a certain by-power index record exist
func ValidatorByPowerIndexExists(ctx sdk.Context, keeper Keeper, power []byte) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Get(power) != nil
}
