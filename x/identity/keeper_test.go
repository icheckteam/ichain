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

	// Add Trust
	// ----------------------------------------------------

	// invalid trustor/trusting
	msgSetTrust := MsgSetTrust{Trustor: addr1, Trusting: addr2, Trust: true}
	keeper.AddTrust(ctx, msgSetTrust)
	found := keeper.hasTrust(ctx, msgSetTrust.Trustor, msgSetTrust.Trusting)
	assert.True(t, found)

	// get trust not found
	found = keeper.hasTrust(ctx, addr1, addr3)
	assert.True(t, !found)

	// trust = false
	msgSetTrust = MsgSetTrust{Trustor: addr1, Trusting: addr2, Trust: false}
	keeper.AddTrust(ctx, msgSetTrust)
	found = keeper.hasTrust(ctx, msgSetTrust.Trustor, msgSetTrust.Trusting)
	assert.True(t, !found)

	// Add Certs
	// ----------------------------------------------------
	msgSetCerts := MsgSetCerts{Certifier: addr1, Recipient: addr1, Values: []CertValue{CertValue{Property: "owner", Confidence: true}}}
	keeper.AddCerts(ctx, msgSetCerts)
	cert, found := keeper.GetCert(ctx, addr1, "owner", addr1)
	assert.True(t, found)
	assert.True(t, bytes.Equal(cert.Certifier, addr1))

}

func TestKeeperCreate(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	keeper.Register(ctx, MsgReg{
		Address: addr1,
		Sender:  addr2,
	})
}
