package asset

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	AssetKey = []byte{0x00} // prefix for each key to a candidate
)

// GetAssetKey get the key for the record with address
func GetAssetKey(uid []byte) []byte {
	return append(AssetKey, uid...)
}
