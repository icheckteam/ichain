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
func KeyIdentityByOwnerIndex(owner sdk.AccAddress, identityID int64) []byte {
	return []byte(fmt.Sprintf("account:%s:%d", owner.String(), identityID))
}

func KeyIdentitiesByOwnerIndex(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("account:%s", owner.String()))
}

func KeyClaimedIdentity(address sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("claim:%s", address.String()))
}

// Key for getting all trusting from the store
func KeyTrust(trustor, trusting sdk.AccAddress) []byte {
	return append(KeyTrusts(trustor), []byte(fmt.Sprintf("trust:%s:%s", trustor.String(), trusting.String()))...)
}

func KeyTrusts(trustor sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("trust:%s", trustor.String()))
}

// Key for getting a cert from the store
func KeyCert(identityID int64, property string, certifier sdk.AccAddress) []byte {
	return append(KeyCerts(identityID, property), []byte(fmt.Sprintf(":%s", certifier.String()))...)
}

// Key for getting all certs from the store
func KeyCerts(identityID int64, property string) []byte {
	return []byte(fmt.Sprintf("identity:%d:%s", identityID, property))
}
