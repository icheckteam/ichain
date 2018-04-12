package chain

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	RecordKey = []byte{0x00} // prefix for each key to a candidate
)

// GetRecordKey get the key for the record with address
func GetRecordKey(uid []byte) []byte {
	return append(RecordKey, uid...)
}
