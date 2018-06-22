package insurance

import (
	"testing"
	"time"

	"github.com/icheckteam/ichain/x/asset"
	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, keeper := createTestInput(t, false, 0)

	msg := MsgCreateContract{
		ID:        "1",
		Issuer:    addrs[0],
		Recipient: addrs[1],
		Expires:   time.Now(),
		Serial:    "100495",
		AssetID:   "demi",
	}

	// invalid asset
	_, err := keeper.CreateContract(ctx, msg)
	assert.True(t, err != nil)

	keeper.assetKeeper.CreateAsset(ctx, asset.MsgCreateAsset{
		AssetID:  "demi",
		Quantity: 1,
		Sender:   addrs[0],
	})

	// Test create contract
	_, err = keeper.CreateContract(ctx, msg)
	newContract := keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.ID == msg.ID)
	assert.True(t, newContract.Issuer.String() == msg.Issuer.String())
	assert.True(t, newContract.Recipient.String() == msg.Recipient.String())
	assert.True(t, newContract.Expires.Unix() == msg.Expires.Unix())
	assert.True(t, newContract.Serial == msg.Serial)
	assert.True(t, newContract.AssetID == msg.AssetID)
	assert.True(t, newContract.Claim == nil)

	// Test create claim
	msgCreateClaim := MsgCreateClaim{ContractID: newContract.ID, Issuer: addrs[1], Recipient: addrs[3]}
	keeper.CreateClaim(ctx, msgCreateClaim)
	newContract = keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.Claim.Status == ClaimStatusPending)
	assert.True(t, newContract.Claim.Recipient.String() == addrs[3].String())

	// claim is processing
	err = keeper.CreateClaim(ctx, msgCreateClaim)
	assert.True(t, err != nil)

	// invalid issuer
	msgCreateClaim = MsgCreateClaim{ContractID: newContract.ID, Issuer: addrs[2], Recipient: addrs[3]}
	err = keeper.CreateClaim(ctx, msgCreateClaim)
	assert.True(t, err != nil)

	// invalid contract ID
	msgCreateClaim = MsgCreateClaim{ContractID: "1212", Issuer: addrs[1], Recipient: addrs[3]}
	err = keeper.CreateClaim(ctx, msgCreateClaim)
	assert.True(t, err != nil)

	// Test process claim
	msgProcessClaim := MsgProcessClaim{ContractID: msg.ID, Issuer: addrs[3], Status: ClaimStatusClaimRepair}
	err = keeper.ProcessClaim(ctx, msgProcessClaim)
	newContract = keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.Claim.Status == ClaimStatusClaimRepair)

	msgCreateClaim = MsgCreateClaim{ContractID: newContract.ID, Issuer: addrs[1], Recipient: addrs[3]}
	keeper.CreateClaim(ctx, msgCreateClaim)
	newContract = keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.Claim.Status == ClaimStatusPending)
	assert.True(t, newContract.Claim.Recipient.String() == addrs[3].String())

	msgProcessClaim = MsgProcessClaim{ContractID: msg.ID, Issuer: addrs[3], Status: ClaimStatusPending}
	err = keeper.ProcessClaim(ctx, msgProcessClaim)
	assert.True(t, err != nil)

	msgProcessClaim = MsgProcessClaim{ContractID: msg.ID, Issuer: addrs[3], Status: ClaimStatusRejected}
	keeper.ProcessClaim(ctx, msgProcessClaim)
	newContract = keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.Claim.Status == ClaimStatusRejected)

	msgCreateClaim = MsgCreateClaim{ContractID: newContract.ID, Issuer: addrs[1], Recipient: addrs[3]}
	keeper.CreateClaim(ctx, msgCreateClaim)
	newContract = keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.Claim.Status == ClaimStatusPending)
	assert.True(t, newContract.Claim.Recipient.String() == addrs[3].String())

	msgProcessClaim = MsgProcessClaim{ContractID: msg.ID, Issuer: addrs[3], Status: ClaimStatusTheftConfirmed}
	keeper.ProcessClaim(ctx, msgProcessClaim)
	newContract = keeper.GetContract(ctx, msg.ID)
	assert.True(t, newContract.Claim.Status == ClaimStatusTheftConfirmed)

	msgCreateClaim = MsgCreateClaim{ContractID: newContract.ID, Issuer: addrs[1], Recipient: addrs[3]}
	err = keeper.CreateClaim(ctx, msgCreateClaim)
	assert.True(t, err != nil)

}
