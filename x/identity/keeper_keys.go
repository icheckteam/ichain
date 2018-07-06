package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyNextIdentityID = []byte("keyNextIdentityID")
)

// Key for getting a identity from the store
func KeyIdentity(identityID int64) []byte {
	return []byte(fmt.Sprintf("identities:%d", identityID))
}

// Key for getting all identities from the store
func KeyIdentities() []byte {
	return []byte(fmt.Sprintf("identities:"))
}

// Key for getting a identity id  of the account from the store
func KeyIdentityByOwnerIndex(owner sdk.Address, identityID int64) []byte {
	return []byte(fmt.Sprintf("accounts:%s:%d", owner.String(), identityID))
}

// Key for getting all identity id  of the account from the store
func KeyIdentitiesByOwnerIndex(owner sdk.Address, identityID int64) []byte {
	return []byte(fmt.Sprintf("accounts:%s:", owner.String()))
}

// Key for getting all trusting from the store
func KeyTrust(trustor, trusting sdk.Address) []byte {
	return []byte(fmt.Sprintf("trusts:%s:%s", trustor.String(), trusting.String()))
}

// Key for getting all trusting from the store
func KeyTrusting(trustor, trusting sdk.Address) []byte {
	return []byte(fmt.Sprintf("trustings:%s:%s", trusting.String(), trustor.String()))
}

// Key for getting a cert from the store
func KeyCert(identityID int64, certifier sdk.Address) []byte {
	return []byte(fmt.Sprintf("certs:%d:%s", identityID, certifier.String()))
}

// Key for getting all certs from the store
func KeyCerts(identityID int64) []byte {
	return []byte(fmt.Sprintf("certs:%d", identityID))
}
