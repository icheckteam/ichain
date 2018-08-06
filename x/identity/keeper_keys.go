package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// CertsKey ...
	CertsKey = []byte{0x01}
	// OwnersKey ...
	OwnersKey = []byte{0x02}
	//TrustsKey ...
	TrustsKey = []byte{0x03}
	// OwnerCountKey ...
	OwnerCountKey = []byte{0x04}
)

// KeyTrust Key for getting all trusting from the store
func KeyTrust(trustor, trusting sdk.AccAddress) []byte {
	return append(
		append(KeyTrusts(trustor), trustor.Bytes()...),
		trusting.Bytes()...,
	)
}

// KeyTrusts ...
func KeyTrusts(trustor sdk.AccAddress) []byte {
	return append(TrustsKey, trustor.Bytes()...)
}

// KeyCert Key for getting a cert from the store
func KeyCert(addr sdk.AccAddress, property string, certifier sdk.AccAddress) []byte {
	return append(
		append(KeyCerts(addr), []byte(property)...),
		certifier.Bytes()...,
	)
}

// KeyCerts Key for getting all certs from the store
func KeyCerts(addr sdk.AccAddress) []byte {
	return append(CertsKey, addr.Bytes()...)
}

// KeyOwners ...
func KeyOwners(id sdk.AccAddress) []byte {
	return append(OwnersKey, id.Bytes()...)
}

// KeyOwner ...
func KeyOwner(id, owner sdk.AccAddress) []byte {
	return append(KeyOwners(id), owner.Bytes()...)
}

// KeyOwnerCount ...
func KeyOwnerCount(id sdk.AccAddress) []byte {
	return append(OwnerCountKey, id.Bytes()...)
}
