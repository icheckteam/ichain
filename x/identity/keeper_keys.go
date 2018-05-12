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
func GetClaimsAccountKey(addr []byte) []byte {
	return append(ClaimsAccountKeyPrefix, addr...)
}

// GetClaimRecordKey ...
func GetClaimsOwnerKey(addr []byte) []byte {
	return append(ClaimsOwnerKeyPrefix, addr...)
}
