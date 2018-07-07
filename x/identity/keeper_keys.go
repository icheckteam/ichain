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
	return []byte(fmt.Sprintf("identity:%d", identityID))
}

// Key for getting all identities from the store
func KeyIdentities() []byte {
	return []byte(fmt.Sprintf("identity:"))
}

// Key for getting a identity id  of the account from the store
func KeyIdentityByOwnerIndex(owner sdk.Address, identityID int64) []byte {
	return []byte(fmt.Sprintf("account:%s:%d", owner.String(), identityID))
}

func KeyClaimedIdentity(address sdk.Address, identityID int64) []byte {
	return []byte(fmt.Sprintf("claimed:%s:%d", address.String(), identityID))
}

func KeyClaimedIdentities(address sdk.Address) []byte {
	return []byte(fmt.Sprintf("claimed:%s", address.String()))
}

// Key for getting all identity id  of the account from the store
func KeyIdentitiesByOwnerIndex(owner sdk.Address, identityID int64) []byte {
	return []byte(fmt.Sprintf("account:%s:", owner.String()))
}

// Key for getting all trusting from the store
func KeyTrust(trustor, trusting sdk.Address) []byte {
	return []byte(fmt.Sprintf("trust:%s:%s", trustor.String(), trusting.String()))
}

// Key for getting all trusting from the store
func KeyTrusting(trustor, trusting sdk.Address) []byte {
	return []byte(fmt.Sprintf("trusting:%s:%s", trusting.String(), trustor.String()))
}

// Key for getting a cert from the store
func KeyCert(identityID int64, certifier sdk.Address) []byte {
	return []byte(fmt.Sprintf("identity:%d:cert:%s", identityID, certifier.String()))
}

// Key for getting all certs from the store
func KeyCerts(identityID int64) []byte {
	return []byte(fmt.Sprintf("identity:%d:cert:", identityID))
}
