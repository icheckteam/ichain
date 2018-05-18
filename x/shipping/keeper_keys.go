package shipping

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	OrderKey = []byte{0x06} // prefix for each key to a candidate
)

// GetOrderKey get the key for the record with address
func GetOrderKey(uid []byte) []byte {
	return append(OrderKey, uid...)
}
