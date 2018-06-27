package identity

import (
	"bytes"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr1 := sdk.Address([]byte("input1"))
	expires, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	ctx, keeper := createTestInput(t, false)

	msgCreateClaim := MsgCreateClaim{
		ClaimID:   "claimID",
		Context:   "claim:identity",
		Content:   []byte(`{"demo": 1}`),
		Expires:   expires.Unix(),
		Issuer:    addr,
		Recipient: addr1,
	}

	// Test create claim
	// -----------------------------------------
	keeper.CreateClaim(ctx, msgCreateClaim)
	newClaim := keeper.GetClaim(ctx, msgCreateClaim.ClaimID)
	assert.True(t, newClaim.ID == msgCreateClaim.ClaimID)
	assert.True(t, newClaim.Context == msgCreateClaim.Context)
	assert.True(t, bytes.Equal(newClaim.Content, msgCreateClaim.Content))
	assert.True(t, newClaim.Expires == msgCreateClaim.Expires)
	assert.True(t, bytes.Equal(newClaim.Issuer, msgCreateClaim.Issuer))
	assert.True(t, bytes.Equal(newClaim.Recipient, msgCreateClaim.Recipient))

	// test update claims
	msgCreateClaim.Content = []byte(`{"demo": 2}`)
	keeper.CreateClaim(ctx, msgCreateClaim)
	newClaim = keeper.GetClaim(ctx, msgCreateClaim.ClaimID)
	assert.True(t, newClaim.ID == msgCreateClaim.ClaimID)
	assert.True(t, newClaim.Context == msgCreateClaim.Context)
	assert.True(t, bytes.Equal(newClaim.Content, msgCreateClaim.Content))
	assert.True(t, newClaim.Expires == msgCreateClaim.Expires)
	assert.True(t, bytes.Equal(newClaim.Issuer, msgCreateClaim.Issuer))
	assert.True(t, bytes.Equal(newClaim.Recipient, msgCreateClaim.Recipient))

	// id already exists
	msgCreateClaim2 := MsgCreateClaim{
		ClaimID:   "claimID",
		Context:   "claim:identity",
		Content:   []byte(`{"demo": 1}`),
		Expires:   expires.Unix(),
		Issuer:    addr1,
		Recipient: addr,
	}
	_, err := keeper.CreateClaim(ctx, msgCreateClaim2)
	assert.True(t, err != nil)

	msgAnswerClaim := MsgAnswerClaim{ClaimID: msgCreateClaim.ClaimID, Sender: addr}
	keeper.AnswerClaim(ctx, msgAnswerClaim)
	newClaim = keeper.GetClaim(ctx, msgAnswerClaim.ClaimID)
	assert.True(t, newClaim.Paid == true)

	// Test Revoke Claim
	// ------------------------------------------------------------
	msgRevokeClaim := MsgRevokeClaim{ClaimID: msgCreateClaim.ClaimID, Sender: addr1, Revocation: "323232"}
	keeper.RevokeClaim(ctx, msgRevokeClaim)
	newClaim = keeper.GetClaim(ctx, msgCreateClaim.ClaimID)
	assert.True(t, newClaim.Revocation == "323232")

	// asset not found
	msgRevokeClaim = MsgRevokeClaim{ClaimID: "12121", Sender: addr, Revocation: "323232"}
	_, err = keeper.RevokeClaim(ctx, msgRevokeClaim)
	assert.True(t, err != nil)

	// invalid Sender
	msgRevokeClaim = MsgRevokeClaim{ClaimID: "12121", Sender: addr1, Revocation: "323232"}
	_, err = keeper.RevokeClaim(ctx, msgRevokeClaim)
	assert.True(t, err != nil)
}
