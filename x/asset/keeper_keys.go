package asset

import sdk "github.com/cosmos/cosmos-sdk/types"

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	AssetKey            = []byte{0x00} // prefix for each key to an asset
	AccountAssetKey     = []byte{0x01} // prefix for each key to an account
	ProposalsKey        = []byte{0x02} // prefix for each key to an account a proposal
	AssetChildrenKey    = []byte{0x03} // prefix for each key to an asset parent a an asset child
	AccountProposalsKey = []byte{0x04} // prefix for each key to an account a proposal
	ReportersKey        = []byte{0x05}
	PropertiesKey       = []byte{0x06}
	InventoryKey        = []byte{0x07}
	ReporterAssetsKey   = []byte{0x08}
	ProposalsAccountKey = []byte{0x09}
)

// GetAssetKey get the key for the record with address
func GetAssetKey(assetID string) []byte {
	return append(AssetKey, []byte(assetID)...)
}

// GetAccountAssetKey get the key for an account for an asset
func GetAccountAssetKey(addr sdk.AccAddress, assetID string) []byte {
	return append(GetAccountAssetsKey(addr), []byte(assetID)...)
}

// GetAccountAssetsKey get the key for an account for all assets
func GetAccountAssetsKey(addr sdk.AccAddress) []byte {
	return append(AccountAssetKey, []byte(addr.String())...)
}

// GetAssetChilrenKey get the key for an asset for an asset
func GetAssetChildrenKey(parent, children string) []byte {
	return append(GetAssetChildrensKey(parent), []byte(children)...)
}

func GetAssetChildrensKey(parent string) []byte {
	return append(AssetChildrenKey, []byte(parent)...)
}

func GetProposalKey(assetID string, recipient sdk.AccAddress) []byte {
	return append(GetProposalsKey(assetID), []byte(recipient.String())...)
}

func GetProposalsKey(assetID string) []byte {
	return append(ProposalsKey, []byte(assetID)...)
}

// GetInventoryKey ...
func GetInventoryKey(addr sdk.AccAddress, assetID string) []byte {
	return append(GetInventoryByAccountKey(addr), []byte(assetID)...)
}

// GetInventoryKey ...
func GetInventoryByAccountKey(addr sdk.AccAddress) []byte {
	return append(InventoryKey, []byte(addr.String())...)
}

// GetInventoryKey ...
func GetReporterAssetKey(addr sdk.AccAddress, assetID string) []byte {
	return append(GetReporterAssetsKey(addr), []byte(assetID)...)
}

// GetInventoryKey ...
func GetReporterAssetsKey(addr sdk.AccAddress) []byte {
	return append(ReporterAssetsKey, []byte(addr.String())...)
}

// GetProposalAccountKey ...
func GetProposalAccountKey(addr sdk.AccAddress, assetID string) []byte {
	return append(GetProposalsAccount(addr), []byte(assetID)...)
}

// GetProposalsAccount ...
func GetProposalsAccount(addr sdk.AccAddress) []byte {
	return append(ProposalsAccountKey, []byte(addr.String())...)
}
