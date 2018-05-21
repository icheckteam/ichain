package warranty

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	ContractKeyPrefix = []byte{0x00} // prefix for each key to a contract
)

func GetContractKey(contractID string) []byte {
	return append(ContractKeyPrefix, []byte(contractID)...)
}
