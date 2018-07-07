package identity

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// Add Identity
	// ----------------------------------------------------

	// valid
	msgCreateIdentity := MsgCreateIdentity{
		Sender: addr1,
	}
	keeper.AddIdentity(ctx, msgCreateIdentity)
	identity, found := keeper.GetIdentity(ctx, 1)
	assert.True(t, found)
	assert.True(t, identity.ID == 1)
	assert.True(t, bytes.Equal(identity.Owner, msgCreateIdentity.Sender))

	msgCreateIdentity = MsgCreateIdentity{
		Sender: addr1,
	}
	keeper.AddIdentity(ctx, msgCreateIdentity)
	identity, found = keeper.GetIdentity(ctx, 2)
	assert.True(t, found)
	assert.True(t, identity.ID == 2)
	assert.True(t, bytes.Equal(identity.Owner, msgCreateIdentity.Sender))

	// get identity not found
	_, found = keeper.GetIdentity(ctx, 3)
	assert.True(t, !found)

	// Add Trust
	// ----------------------------------------------------

	// invalid trustor/trusting
	msgSetTrust := MsgSetTrust{Trustor: addr1, Trusting: addr2, Trust: true}
	keeper.AddTrust(ctx, msgSetTrust)
	trust, found := keeper.GetTrust(ctx, msgSetTrust.Trustor, msgSetTrust.Trusting)
	assert.True(t, found)
	assert.True(t, bytes.Equal(trust.Trusting, msgSetTrust.Trusting))

	// get trust not found
	_, found = keeper.GetTrust(ctx, addr1, addr3)
	assert.True(t, !found)

	// trust = false
	msgSetTrust = MsgSetTrust{Trustor: addr1, Trusting: addr2, Trust: false}
	keeper.AddTrust(ctx, msgSetTrust)
	_, found = keeper.GetTrust(ctx, msgSetTrust.Trustor, msgSetTrust.Trusting)
	assert.True(t, !found)

	// Add Certs
	// ----------------------------------------------------
	msgSetCerts := MsgSetCerts{Certifier: addr1, IdentityID: 1, Values: []CertValue{CertValue{Property: addr2, Confidence: true}}}
	keeper.AddCerts(ctx, msgSetCerts)
	cert, found := keeper.GetCert(ctx, 1, addr2, addr1)
	assert.True(t, found)
	assert.True(t, bytes.Equal(cert.Property, addr2))

	// Claim identity
	msgSetCerts = MsgSetCerts{Certifier: addr2, IdentityID: 1, Values: []CertValue{CertValue{Property: addr2, Confidence: true}}}
	keeper.AddCerts(ctx, msgSetCerts)
	found = keeper.HasClaimedIdentity(ctx, addr2, 1)
	assert.True(t, found)

	msgSetCerts = MsgSetCerts{Certifier: addr1, IdentityID: 1, Values: []CertValue{CertValue{Property: addr2, Confidence: false}}}
	keeper.AddCerts(ctx, msgSetCerts)
	_, found = keeper.GetCert(ctx, 1, addr2, addr1)
	assert.True(t, !found)

	// unClaim identity
	msgSetCerts = MsgSetCerts{Certifier: addr2, IdentityID: 1, Values: []CertValue{CertValue{Property: addr2, Confidence: false}}}
	keeper.AddCerts(ctx, msgSetCerts)
	found = keeper.HasClaimedIdentity(ctx, addr2, 1)
	assert.True(t, !found)

	// invalid identity id
	msgSetCerts = MsgSetCerts{Certifier: addr2, IdentityID: 1000, Values: []CertValue{CertValue{Property: addr2, Confidence: false}}}
	err := keeper.AddCerts(ctx, msgSetCerts)
	assert.True(t, err != nil)

}
