package identity

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	addr1 = sdk.Address([]byte("addr1"))
	addr2 = sdk.Address([]byte("addr2"))
	addr3 = sdk.Address([]byte("addr3"))
)

func TestMsgCreateIdentity(t *testing.T) {
	tests := []struct {
		name       string
		sender     sdk.Address
		expectPass bool
	}{
		{"basic good", addr1, true},
		{"empty sender", nil, false},
	}

	for _, tc := range tests {
		msg := NewMsgCreateIdentity(tc.sender)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgCreateIdentityType(t *testing.T) {
	msg := MsgCreateIdentity{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgCreateIdentityGetSignBytes(t *testing.T) {
	signBytes := MsgCreateIdentity{
		Sender: addr1,
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), `{"sender":"cosmosaccaddr1v9jxgu333rmgrm"}`)
}

func TestMsgCreateIdentityGetSigner(t *testing.T) {
	signers := MsgCreateIdentity{Sender: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgSetTrust(t *testing.T) {
	tests := []struct {
		name       string
		trustor    sdk.Address
		trusting   sdk.Address
		trust      bool
		expectPass bool
	}{
		{"basic good", addr1, addr2, true, true},
		{"empty trustor", nil, addr2, true, false},
		{"empty trusting", addr1, nil, true, false},
	}

	for _, tc := range tests {
		msg := NewMsgSetTrust(tc.trustor, tc.trusting, tc.trust)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgSetTrustType(t *testing.T) {
	msg := MsgSetTrust{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgSetTrustGetSignBytes(t *testing.T) {
	signBytes := MsgSetTrust{
		Trustor:  addr1,
		Trusting: addr2,
		Trust:    true,
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), `{"trustor":"cosmosaccaddr1v9jxgu333rmgrm","trusting":"cosmosaccaddr1v9jxgu3jlsw7dy","trust":true}`)
}

func TestMsgSetTrustGetSigner(t *testing.T) {
	signers := MsgSetTrust{Trustor: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgSetCerts(t *testing.T) {
	tests := []struct {
		name       string
		certifier  sdk.Address
		identity   int64
		values     []CertValue
		expectPass bool
	}{
		{"basic good", addr1, 1, []CertValue{CertValue{Property: "owner", Type: "realname", Confidence: true}}, true},
		{"empty certifier", nil, 1, []CertValue{CertValue{Property: "owner", Type: "realname", Confidence: true}}, false},
		{"empty identity id", addr1, 0, []CertValue{CertValue{Property: "owner", Type: "realname", Confidence: true}}, false},
		{"empty property address", addr1, 1, []CertValue{CertValue{Type: "realname", Confidence: true}}, false},
	}

	for _, tc := range tests {
		msg := NewMsgSetCerts(tc.certifier, tc.identity, tc.values)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgSetCertsType(t *testing.T) {
	msg := MsgSetCerts{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgSetCertsGetSigner(t *testing.T) {
	signers := MsgSetCerts{Certifier: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgSetCertsGetSignBytes(t *testing.T) {
	signBytes := MsgSetCerts{
		Certifier:  addr1,
		IdentityID: 1,
		Values:     []CertValue{CertValue{Property: "owner", Type: "realname", Confidence: true}},
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), `{"certifier":"cosmosaccaddr1v9jxgu333rmgrm","identity_id":"1","values":[{"property":"owner","type":"realname","data":null,"confidence":true}]}`)
}
