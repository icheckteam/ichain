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
		{"basic good", addr1, addr2, []CertValue{CertValue{Property: "owner", Confidence: true, Owner: addr3}}, true},
		{"empty certifier", nil, addr2, []CertValue{CertValue{Property: "owner", Confidence: true}}, false},
		{"empty identity id", addr1, nil, []CertValue{CertValue{Property: "owner", Confidence: true}}, false},
		{"empty property address", addr1, addr2, []CertValue{CertValue{Confidence: true}}, false},
	}

	for _, tc := range tests {
		msg := NewMsgSetCerts(tc.certifier, tc.certifier, tc.values)
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
	signers := MsgSetCerts{Sender: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgSetCertsGetSignBytes(t *testing.T) {
	signBytes := MsgSetCerts{
		Sender: addr1,
		Issuer: addr2,
		Values: []CertValue{CertValue{Property: "owner", Confidence: true}},
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), "{\"type\":\"identity/SetCerts\",\"value\":{\"issuer\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"sender\":\"cosmosaccaddr1v9jxgu333rmgrm\",\"values\":[{\"confidence\":true,\"data\":null,\"owner\":\"cosmosaccaddr16y6p2v\",\"property\":\"owner\"}]}}")
}

// MsgReg
// ------------------------------------------
func TestMsgRegType(t *testing.T) {
	msg := MsgReg{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgRegGetSigner(t *testing.T) {
	signers := MsgReg{Sender: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgRegGetSignBytes(t *testing.T) {
	signBytes := MsgReg{
		Sender:  addr1,
		Address: addr2,
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), "{\"address\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"sender\":\"cosmosaccaddr1v9jxgu333rmgrm\"}")
}

func TestMsgRegValidation(t *testing.T) {
	tests := []struct {
		name       string
		sender     sdk.AccAddress
		address    sdk.AccAddress
		expectPass bool
	}{
		{"basic good", addr1, addr2, true},
		{"empty sender", nil, addr2, false},
		{"empty address", addr1, nil, false},
	}

	for _, tc := range tests {
		msg := MsgReg{tc.sender, tc.address}
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// MsgAddOwner
// ------------------------------------------
func TestMsgAddOwnerType(t *testing.T) {
	msg := MsgAddOwner{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgAddOwnerGetSigner(t *testing.T) {
	signers := MsgAddOwner{Sender: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgAddOwnerGetSignBytes(t *testing.T) {
	signBytes := MsgAddOwner{
		Sender:  addr1,
		Address: addr2,
		Owner:   addr3,
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), "{\"address\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"owner\":\"cosmosaccaddr1v9jxgu3nzx6tsk\",\"sender\":\"cosmosaccaddr1v9jxgu333rmgrm\"}")
}

func TestMsgAddOwnerValidation(t *testing.T) {
	tests := []struct {
		name       string
		sender     sdk.AccAddress
		address    sdk.AccAddress
		owner      sdk.AccAddress
		expectPass bool
	}{
		{"basic good", addr1, addr2, addr3, true},
		{"empty sender", nil, addr2, addr3, false},
		{"empty address", addr1, nil, addr3, false},
		{"empty owner", addr1, addr2, nil, false},
	}

	for _, tc := range tests {
		msg := MsgAddOwner{tc.sender, tc.address, tc.owner}
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// MsgDelOwner
// ------------------------------------------
func TestMsgDelOwnerType(t *testing.T) {
	msg := MsgDelOwner{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgDelOwnerGetSigner(t *testing.T) {
	signers := MsgDelOwner{Sender: addr1}.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", signers), `[6164647231]`)
}

func TestMsgDelOwnerGetSignBytes(t *testing.T) {
	signBytes := MsgDelOwner{
		Sender:  addr1,
		Address: addr2,
		Owner:   addr3,
	}.GetSignBytes()
	assert.Equal(t, string(signBytes), "{\"address\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"owner\":\"cosmosaccaddr1v9jxgu3nzx6tsk\",\"sender\":\"cosmosaccaddr1v9jxgu333rmgrm\"}")
}

func TestMsgDelOwnerValidation(t *testing.T) {
	tests := []struct {
		name       string
		sender     sdk.AccAddress
		address    sdk.AccAddress
		owner      sdk.AccAddress
		expectPass bool
	}{
		{"basic good", addr1, addr2, addr3, true},
		{"empty sender", nil, addr2, addr3, false},
		{"empty address", addr1, nil, addr3, false},
		{"empty owner", addr1, addr2, nil, false},
	}

	for _, tc := range tests {
		msg := MsgDelOwner{tc.sender, tc.address, tc.owner}
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}
