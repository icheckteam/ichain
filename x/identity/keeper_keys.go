package identity

var (
	// ClaimRecordKeyPrefix for store prefixes
	ClaimRecordKeyPrefix = []byte{0x00} // prefix for each key to a candidate
)

// GetClaimRecordKey ...
func GetClaimRecordKey(uuid string) []byte {
	return append(ClaimRecordKeyPrefix, []byte(uuid)...)
}
