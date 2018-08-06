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

	msgRegister := MsgReg{
		Ident:  addrs[1],
		Sender: addrs[2],
	}
	keeper.Register(ctx, msgRegister)
	owners := keeper.GetOwners(ctx, addrs[1])
	assert.True(t, len(owners) == 1)

	msgRegister = MsgReg{
		Ident:  addrs[2],
		Sender: addrs[2],
	}
	keeper.Register(ctx, msgRegister)

	// Invalid id already exists
	_, err := keeper.Register(ctx, MsgReg{
		Ident:  addrs[1],
		Sender: addrs[2],
	})
	assert.True(t, err != nil)

	// add owner

	// invalid sender
	_, err = keeper.AddOwner(ctx, MsgAddOwner{
		Ident:  addrs[1],
		Owner:  addrs[3],
		Sender: addrs[4],
	})
	assert.True(t, err != nil)

	// valid
	keeper.AddOwner(ctx, MsgAddOwner{
		Ident:  addrs[1],
		Owner:  addrs[3],
		Sender: addrs[2],
	})
	owners = keeper.GetOwners(ctx, addrs[1])
	assert.True(t, len(owners) == 2)

	// delete owner invalid sender
	_, err = keeper.DeleteOwner(ctx, MsgDelOwner{
		Ident:  addrs[1],
		Owner:  addrs[3],
		Sender: addrs[4],
	})
	assert.True(t, err != nil)
	keeper.DeleteOwner(ctx, MsgDelOwner{
		Ident:  addrs[1],
		Owner:  addrs[3],
		Sender: addrs[2],
	})
	owners = keeper.GetOwners(ctx, addrs[1])
	assert.True(t, len(owners) == 1)

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
	msgSetCerts := MsgSetCerts{Issuer: addrs[2], Sender: addrs[2], Values: []CertValue{CertValue{Property: "owner", Owner: addrs[1], Confidence: true}}}
	keeper.AddCerts(ctx, msgSetCerts)
	cert, found := keeper.GetCert(ctx, addrs[1], "owner", addrs[2])
	assert.True(t, found)
	assert.True(t, bytes.Equal(cert.Certifier, addrs[2]))
	certs := keeper.GetCerts(ctx, addrs[1])
	assert.True(t, len(certs) == 1)

	msgSetCerts = MsgSetCerts{Issuer: addrs[2], Sender: addrs[2], Values: []CertValue{CertValue{Property: "owner", Owner: addrs[1], Confidence: true}}}
	_, err = keeper.AddCerts(ctx, msgSetCerts)
	assert.True(t, err == nil)

	msgSetCerts = MsgSetCerts{Issuer: addr1, Sender: addr2, Values: []CertValue{CertValue{Property: "owner", Owner: addr2, Confidence: true}}}
	_, err = keeper.AddCerts(ctx, msgSetCerts)
	assert.True(t, err != nil)

	msgSetCerts = MsgSetCerts{Issuer: addrs[2], Sender: addrs[2], Values: []CertValue{CertValue{Property: "owner", Owner: addrs[1], Confidence: false}}}
	_, err = keeper.AddCerts(ctx, msgSetCerts)
	certs = keeper.GetCerts(ctx, addr2)
	assert.True(t, len(certs) == 0)

}
