package identity

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	addr1 = sdk.AccAddress([]byte("addr1"))
	addr2 = sdk.AccAddress([]byte("addr2"))
	addr3 = sdk.AccAddress([]byte("addr3"))
)

func TestMsgSetTrust(t *testing.T) {
	tests := []struct {
		name       string
		trustor    sdk.AccAddress
		trusting   sdk.AccAddress
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
	assert.Equal(t, string(signBytes), "{\"type\":\"identity/SetTrust\",\"value\":{\"trust\":true,\"trusting\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"trustor\":\"cosmosaccaddr1v9jxgu333rmgrm\"}}")
}

func TestMsgSetTrustGetSigner(t *testing.T) {
	signers := MsgSetTrust{Trustor: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgSetCerts(t *testing.T) {
	tests := []struct {
		name       string
		certifier  sdk.AccAddress
		recipient  sdk.AccAddress
		values     []CertValue
		expectPass bool
	}{
		{"basic good", addr1, addr2, []CertValue{CertValue{Property: "owner", Confidence: true}}, true},
		{"empty certifier", nil, addr2, []CertValue{CertValue{Property: "owner", Confidence: true}}, false},
		{"empty identity id", addr1, nil, []CertValue{CertValue{Property: "owner", Confidence: true}}, false},
		{"empty property address", addr1, addr2, []CertValue{CertValue{Confidence: true}}, false},
	}

	for _, tc := range tests {
		msg := NewMsgSetCerts(tc.certifier, tc.recipient, tc.values)
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
		Certifier: addr1,
		Recipient: addr2,
		Values:    []CertValue{CertValue{Property: "owner", Confidence: true}},
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), "{\"type\":\"identity/SetCerts\",\"value\":{\"certifier\":\"cosmosaccaddr1v9jxgu333rmgrm\",\"recipient\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"values\":[{\"confidence\":true,\"context\":\"\",\"data\":null,\"expires\":\"0\",\"id\":\"\",\"property\":\"owner\",\"revocation\":{\"id\":\"\",\"type\":\"\"}}]}}")
}
