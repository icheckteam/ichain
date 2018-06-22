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
	creatTime, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	expires, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	ctx, keeper := createTestInput(t, false)

	msgCreateClaim := MsgCreateClaim{
		ID:      "claimID",
		Context: "claim:identity",
		Content: []byte(`{"demo": 1}`),
		Metadata: Metadata{
			Expires:    expires,
			CreateTime: creatTime,
			Issuer:     addr,
			Recipient:  addr1,
		},
	}

	// Test create claim
	// -----------------------------------------
	keeper.CreateClaim(ctx, msgCreateClaim)
	newClaim := keeper.GetClaim(ctx, msgCreateClaim.ID)
	assert.True(t, newClaim.ID == msgCreateClaim.ID)
	assert.True(t, newClaim.Context == msgCreateClaim.Context)
	assert.True(t, bytes.Equal(newClaim.Content, msgCreateClaim.Content))
	assert.True(t, newClaim.Metadata.Expires.Unix() == msgCreateClaim.Metadata.Expires.Unix())
	assert.True(t, newClaim.Metadata.CreateTime.Unix() == msgCreateClaim.Metadata.CreateTime.Unix())
	assert.True(t, bytes.Equal(newClaim.Metadata.Issuer, msgCreateClaim.Metadata.Issuer))
	assert.True(t, bytes.Equal(newClaim.Metadata.Recipient, msgCreateClaim.Metadata.Recipient))

	// test update claims
	msgCreateClaim.Content = []byte(`{"demo": 2}`)
	keeper.CreateClaim(ctx, msgCreateClaim)
	newClaim = keeper.GetClaim(ctx, msgCreateClaim.ID)
	assert.True(t, newClaim.ID == msgCreateClaim.ID)
	assert.True(t, newClaim.Context == msgCreateClaim.Context)
	assert.True(t, bytes.Equal(newClaim.Content, msgCreateClaim.Content))
	assert.True(t, newClaim.Metadata.Expires.Unix() == msgCreateClaim.Metadata.Expires.Unix())
	assert.True(t, newClaim.Metadata.CreateTime.Unix() == msgCreateClaim.Metadata.CreateTime.Unix())
	assert.True(t, bytes.Equal(newClaim.Metadata.Issuer, msgCreateClaim.Metadata.Issuer))
	assert.True(t, bytes.Equal(newClaim.Metadata.Recipient, msgCreateClaim.Metadata.Recipient))

	// id already exists
	msgCreateClaim2 := MsgCreateClaim{
		ID:      "claimID",
		Context: "claim:identity",
		Content: []byte(`{"demo": 1}`),
		Metadata: Metadata{
			Expires:    expires,
			CreateTime: creatTime,
			Issuer:     addr1,
			Recipient:  addr,
		},
	}
	_, err := keeper.CreateClaim(ctx, msgCreateClaim2)
	assert.True(t, err != nil)

	msgAnswerClaim := MsgAnswerClaim{ClaimID: msgCreateClaim.ID, Sender: addr}
	keeper.AnswerClaim(ctx, msgAnswerClaim)
	newClaim = keeper.GetClaim(ctx, msgAnswerClaim.ClaimID)
	assert.True(t, newClaim.Paid == true)

	// Test Revoke Claim
	// ------------------------------------------------------------
	msgRevokeClaim := MsgRevokeClaim{ClaimID: msgCreateClaim.ID, Sender: addr1, Revocation: "323232"}
	keeper.RevokeClaim(ctx, msgRevokeClaim)
	newClaim = keeper.GetClaim(ctx, msgCreateClaim.ID)
	assert.True(t, newClaim.Metadata.Revocation == "323232")

	// asset not found
	msgRevokeClaim = MsgRevokeClaim{ClaimID: "12121", Sender: addr, Revocation: "323232"}
	_, err = keeper.RevokeClaim(ctx, msgRevokeClaim)
	assert.True(t, err != nil)

	// invalid Sender
	msgRevokeClaim = MsgRevokeClaim{ClaimID: "12121", Sender: addr1, Revocation: "323232"}
	_, err = keeper.RevokeClaim(ctx, msgRevokeClaim)
	assert.True(t, err != nil)
}
