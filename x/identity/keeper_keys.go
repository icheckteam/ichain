package identity

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	// AccountClaimsKey for store prefixes
	AccountClaimsKey = []byte{0x00}

	ClaimKey = []byte{0x01}
)

// GetAccoGetClaimKeyuntClaimKey get the key for an account for a claim
func GetClaimKey(claimID string) []byte {
	return append(ClaimKey, []byte(claimID)...)
}

// GetAccountClaimKey get the key for an account for a claim
func GetAccountClaimKey(addr sdk.Address, claimID string) []byte {
	return append(GetAccountClaimsKey(addr), []byte(claimID)...)
}

// GetAccountClaimsKey get the key for an account for all claims
func GetAccountClaimsKey(addr sdk.Address) []byte {
	return append(AccountClaimsKey, []byte(addr.String())...)
}
