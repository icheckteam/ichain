package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyNextIdentityID = []byte{0x01}
	IdentitiesKey     = []byte{0x02}
)

// Key for getting a identity from the store
func KeyIdentity(identityID int64) []byte {
	return append(IdentitiesKey, []byte(fmt.Sprintf("%d", identityID))...)
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
func KeyTrusted(trustor sdk.Address) []byte {
	return []byte(fmt.Sprintf("trust:%s", trustor.String()))
}

// Key for getting a cert from the store
func KeyCert(identityID int64, property string, certifier sdk.Address) []byte {
	return []byte(fmt.Sprintf("identity:%d:%s:%s", identityID, property, certifier.String()))
}

// Key for getting all certs from the store
func KeyCerts(identityID int64) []byte {
	return []byte(fmt.Sprintf("identity:%d:cert:", identityID))
}
