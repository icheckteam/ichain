package identity

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// MsgCreateClaim tests
// ------------------------------------
func TestMsgCreateClaimType(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr2 := sdk.Address([]byte("input2"))
	var msg = MsgCreateClaim{
		Context:   "claim:context",
		Content:   []byte(`{"demo": 1}`),
		Expires:   time.Now().Add(time.Hour * 100000).Unix(),
		Issuer:    addr,
		Recipient: addr2,
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgCreateClaimValidation(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr2 := sdk.Address([]byte("input2"))

	cases := []struct {
		valid bool
		tx    MsgCreateClaim
	}{
		{false, MsgCreateClaim{}}, // no id
		{false, MsgCreateClaim{
			ClaimID: "1",
		}}, // no context
		{false, MsgCreateClaim{
			ClaimID: "1",
			Context: "1212",
		}}, // no content
		{false, MsgCreateClaim{
			ClaimID: "1",
			Context: "1212",
			Content: []byte(`{"demo": 1}`),
		}}, // no meta

		{false, MsgCreateClaim{
			ClaimID:   "1",
			Context:   "1",
			Content:   []byte(`{"demo": 1}`),
			Recipient: addr}},
		{false, MsgCreateClaim{
			ClaimID: "1",
			Context: "1",
			Content: []byte(`{"demo": 1}`),
			Issuer:  addr2}},
		{false, MsgCreateClaim{
			ClaimID:   "1",
			Context:   "1",
			Content:   []byte(`{"demo": 1}`),
			Issuer:    addr2,
			Recipient: addr}},

		{true, MsgCreateClaim{
			ClaimID:   "1",
			Context:   "1",
			Content:   []byte(`{"demo": 1}`),
			Recipient: addr,
			Issuer:    addr2,
			Expires:   time.Now().Add(time.Hour * 100000).Unix(),
		}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

func TestCreateMsgGet(t *testing.T) {
	var msg = MsgCreateClaim{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCreateMsgGetSignBytes(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr1 := sdk.Address([]byte("input1"))

	expires, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")

	var msg = MsgCreateClaim{
		ClaimID:   "1",
		Context:   "1",
		Content:   []byte(`{"demo": 1}`),
		Recipient: addr,
		Issuer:    addr1,
		Expires:   expires.Unix(),
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"claim_id\":\"1\",\"issuer\":\"696E70757431\",\"recipient\":\"696E707574\",\"context\":\"1\",\"content\":\"eyJkZW1vIjogMX0=\",\"fee\":null,\"expires\":1526030925}")
}

func TestMsgCreateClaimGetSigners(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr1 := sdk.Address([]byte("input1"))

	expiration, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")

	var msg = MsgCreateClaim{
		ClaimID:   "1",
		Context:   "1",
		Content:   []byte(`{"demo": 1}`),
		Recipient: addr1,
		Issuer:    addr,
		Expires:   expiration.Unix(),
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// Revoke tests
// ----------------------------------------------
func TestMsgRevokeClaimType(t *testing.T) {
	msg := MsgRevokeClaim{}
	assert.Equal(t, msg.Type(), "identity")
}

func TestMsgRevokeClaimValidation(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	cases := []struct {
		valid bool
		tx    MsgRevokeClaim
	}{
		{false, MsgRevokeClaim{}},
		{false, MsgRevokeClaim{ClaimID: "1"}},               // only id
		{false, MsgRevokeClaim{ClaimID: "1", Sender: addr}}, // no revocation
		{true, MsgRevokeClaim{ClaimID: "1", Sender: addr, Revocation: "2323"}},
	}

	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgRevokeClaimGet(t *testing.T) {
	var msg = MsgRevokeClaim{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestMsgRevokeClaimGetSignBytes(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgRevokeClaim{
		Sender:     addr,
		Revocation: "demo",
		ClaimID:    "1",
	}
	res := msg.GetSignBytes()
	assert.Equal(t, string(res), `{"claim_id":"1","sender":"696E707574","revocation":"demo"}`)
}

func TestMsgRevokeClaimGetSigners(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgRevokeClaim{
		Sender:     addr,
		Revocation: "demo",
		ClaimID:    "1",
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), "[696E707574]")
}
