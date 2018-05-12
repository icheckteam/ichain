package identity

var (
	// ClaimRecordKeyPrefix for store prefixes
	ClaimRecordKeyPrefix   = []byte{0x00} // prefix for each key to a candidate
	ClaimsAccountKeyPrefix = []byte{0x01}
	ClaimsOwnerKeyPrefix   = []byte{0x02}
)

// GetClaimRecordKey ...
func GetClaimRecordKey(uuid string) []byte {
	return append(ClaimRecordKeyPrefix, []byte(uuid)...)
}

// GetClaimRecordKey ...
func GetClaimsAccountKey(uuid string) []byte {
	return append(ClaimsAccountKeyPrefix, []byte(uuid)...)
}

// GetClaimRecordKey ...
func GetClaimsOwnerKey(uuid string) []byte {
	return append(ClaimsOwnerKeyPrefix, []byte(uuid)...)
}
