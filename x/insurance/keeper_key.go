package insurance

import sdk "github.com/cosmos/cosmos-sdk/types"

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	ContractKeyPrefix  = []byte{0x00} // prefix for each key to a contract
	AccountContractKey = []byte{0x01}
)

func GetContractKey(contractID string) []byte {
	return append(ContractKeyPrefix, []byte(contractID)...)
}

func GetAccountContractKey(addr sdk.Address, contractID string) []byte {
	return append(GetAccountContractsKey(addr), []byte(contractID)...)
}

func GetAccountContractsKey(addr sdk.Address) []byte {
	return append(AccountContractKey, []byte(addr.String())...)
}
