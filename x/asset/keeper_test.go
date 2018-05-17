package asset

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	addr := sdk.Address([]byte("addr1"))
	addr2 := sdk.Address([]byte("addr2"))
	addr3 := sdk.Address([]byte("addr3"))
	addr4 := sdk.Address([]byte("addr4"))

	asset := Asset{
		ID:       "asset1",
		Issuer:   addr,
		Name:     "asset 1",
		Quantity: 100,
	}

	// Test register asset
	keeper.RegisterAsset(ctx, asset)
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))
	keeper.RegisterAsset(ctx, asset)
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))

	// Test add quantity
	keeper.AddQuantity(ctx, AddQuantityMsg{ID: asset.ID, Issuer: addr, Quantity: 50})
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 150}}))

	// Test subtract quantity
	keeper.SubtractQuantity(ctx, SubtractQuantityMsg{ID: asset.ID, Issuer: addr, Quantity: 50})
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))
	keeper.SubtractQuantity(ctx, SubtractQuantityMsg{ID: asset.ID, Issuer: addr, Quantity: 102})
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))

	// Test update attributes
	attrs := []Attribute{Attribute{Name: "weight", NumberValue: 100}, Attribute{Name: "size", NumberValue: 2}}
	keeper.UpdateAttribute(ctx, UpdateAttrMsg{ID: asset.ID, Issuer: addr, Attributes: attrs})
	newAsset := keeper.GetAsset(ctx, asset.ID)
	assert.True(t, newAsset.Attributes[0].Name == "weight")
	assert.True(t, newAsset.Attributes[0].NumberValue == 100)
	assert.True(t, newAsset.Attributes[1].Name == "size")
	assert.True(t, newAsset.Attributes[1].NumberValue == 2)

	//-------------- Test create proposal
	createProposalMsg := CreateProposalMsg{asset.ID, addr, addr2, []string{"weight"}, RoleReporter}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	assert.True(t, newAsset.Proposals[0].Role == RoleReporter)
	assert.True(t, newAsset.Proposals[0].Issuer.String() == addr.String())
	assert.True(t, newAsset.Proposals[0].Recipient.String() == addr2.String())
	assert.True(t, newAsset.Proposals[0].Status == StatusPending)
	assert.True(t, len(newAsset.Proposals[0].Properties) == 1)
	assert.True(t, newAsset.Proposals[0].Properties[0] == "weight")

	createProposalMsg = CreateProposalMsg{asset.ID, addr, addr3, []string{"size", "weight"}, RoleOwner}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	assert.True(t, newAsset.Proposals[1].Role == RoleOwner)
	assert.True(t, newAsset.Proposals[1].Issuer.String() == addr.String())
	assert.True(t, newAsset.Proposals[1].Recipient.String() == addr3.String())
	assert.True(t, newAsset.Proposals[1].Status == StatusPending)
	assert.True(t, len(newAsset.Proposals[1].Properties) == 2)
	assert.True(t, newAsset.Proposals[1].Properties[0] == "size")
	assert.True(t, newAsset.Proposals[1].Properties[1] == "weight")

	// Invalid role
	createProposalMsg = CreateProposalMsg{asset.ID, addr, addr3, []string{"size", "weight"}, 123}
	_, err := keeper.CreateProposal(ctx, createProposalMsg)
	assert.True(t, err != nil)

	//-------------- Test update properties
	createProposalMsg = CreateProposalMsg{asset.ID, addr, addr2, []string{"weight", "size"}, RoleOwner}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	proposal := newAsset.Proposals[0]
	assert.True(t, proposal.Role == RoleOwner)
	assert.True(t, proposal.Issuer.String() == addr.String())
	assert.True(t, proposal.Recipient.String() == addr2.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "weight")
	assert.True(t, proposal.Properties[1] == "size")

	//-------------- Test answer proposal
	// Valid answer
	answerProposalMsg := AnswerProposalMsg{asset.ID, addr2, StatusAccepted}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	proposal = newAsset.Proposals[0]
	assert.True(t, err == nil)
	assert.True(t, proposal.Role == RoleOwner)
	assert.True(t, proposal.Issuer.String() == addr.String())
	assert.True(t, proposal.Recipient.String() == addr2.String())
	assert.True(t, proposal.Status == StatusAccepted)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "weight")
	assert.True(t, proposal.Properties[1] == "size")

	// Answer an already answered proposal
	answerProposalMsg = AnswerProposalMsg{asset.ID, addr2, StatusAccepted}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	assert.True(t, err != nil)

	// Answer with invalid response
	answerProposalMsg = AnswerProposalMsg{asset.ID, addr3, StatusPending}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	assert.True(t, err != nil)

	// Answer with invalid response
	answerProposalMsg = AnswerProposalMsg{asset.ID, addr2, 123}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	assert.True(t, err != nil)

	// Valid answer
	answerProposalMsg = AnswerProposalMsg{asset.ID, addr3, StatusRefused}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	proposal = newAsset.Proposals[1]
	assert.True(t, err == nil)
	assert.True(t, proposal.Role == RoleOwner)
	assert.True(t, proposal.Issuer.String() == addr.String())
	assert.True(t, proposal.Recipient.String() == addr3.String())
	assert.True(t, proposal.Status == StatusRefused)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "size")
	assert.True(t, proposal.Properties[1] == "weight")

	// Refused recipient is not authorized
	createProposalMsg = CreateProposalMsg{asset.ID, addr3, addr4, []string{"weight", "size"}, RoleOwner}
	_, err = keeper.CreateProposal(ctx, createProposalMsg)
	assert.True(t, err != nil)

	// Accepted recipient with role owner is
	createProposalMsg = CreateProposalMsg{asset.ID, addr2, addr4, []string{"weight", "size"}, RoleOwner}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	proposal = newAsset.Proposals[2]
	assert.True(t, proposal.Role == RoleOwner)
	assert.True(t, proposal.Issuer.String() == addr2.String())
	assert.True(t, proposal.Recipient.String() == addr4.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "weight")
	assert.True(t, proposal.Properties[1] == "size")

	// Asset issuer is always authorized
	createProposalMsg = CreateProposalMsg{asset.ID, addr, addr4, []string{}, RoleReporter}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	proposal = newAsset.Proposals[2]
	assert.True(t, proposal.Role == RoleReporter)
	assert.True(t, proposal.Issuer.String() == addr2.String())
	assert.True(t, proposal.Recipient.String() == addr4.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "weight")
	assert.True(t, proposal.Properties[1] == "size")

	// Test update attributes
	attrs = []Attribute{Attribute{Name: "weight", NumberValue: 250}, Attribute{Name: "size", NumberValue: 3}}
	_, err = keeper.UpdateAttribute(ctx, UpdateAttrMsg{ID: asset.ID, Issuer: addr2, Attributes: attrs})
	newAsset = keeper.GetAsset(ctx, asset.ID)
	fmt.Println(newAsset.Attributes)
	assert.True(t, err == nil)
	assert.True(t, newAsset.Attributes[0].Name == "weight")
	assert.True(t, newAsset.Attributes[0].NumberValue == 250)
	assert.True(t, newAsset.Attributes[1].Name == "size")
	assert.True(t, newAsset.Attributes[1].NumberValue == 3)

	//-------------- Test revoke proposal
	// Refused recipient is not authorized
	revokeProposalMsg := RevokeProposalMsg{asset.ID, addr3, addr4, []string{"weight"}}
	_, err = keeper.RevokeProposal(ctx, revokeProposalMsg)
	assert.True(t, err != nil)

	// addr2 is authorized
	revokeProposalMsg = RevokeProposalMsg{asset.ID, addr2, addr4, []string{"weight"}}
	_, err = keeper.RevokeProposal(ctx, revokeProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	proposal = newAsset.Proposals[2]
	assert.True(t, err == nil)
	assert.True(t, proposal.Role == RoleReporter)
	assert.True(t, proposal.Issuer.String() == addr2.String())
	assert.True(t, proposal.Recipient.String() == addr4.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 1)
	assert.True(t, proposal.Properties[0] == "size")

	// proposal is deleted when there is no more property
	revokeProposalMsg = RevokeProposalMsg{asset.ID, addr, addr4, []string{"size"}}
	_, err = keeper.RevokeProposal(ctx, revokeProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.ID)
	assert.True(t, len(newAsset.Proposals) == 2)
}
