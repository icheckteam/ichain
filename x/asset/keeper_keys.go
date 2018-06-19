package asset

import sdk "github.com/cosmos/cosmos-sdk/types"

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	AssetKey        = []byte{0x00} // prefix for each key to an asset
	AccountAssetKey = []byte{0x01} // prefix for each key to an account
)

// GetAssetKey get the key for the record with address
func GetAssetKey(assetID string) []byte {
	return append(AssetKey, []byte(assetID)...)
}

// GetAccountAssetKey get the key for an account for an asset
func GetAccountAssetKey(assetID string, addr sdk.Address) []byte {
	return append(GetAccountAssetsKey(addr), []byte(assetID)...)
}

// GetAccountAssetsKey get the key for an account for all assets
func GetAccountAssetsKey(addr sdk.Address) []byte {
	return append(AccountAssetKey, []byte(addr.String())...)
}
