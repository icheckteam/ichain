package asset

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	addr  = sdk.Address([]byte("addr1"))
	addr2 = sdk.Address([]byte("addr2"))
	addr3 = sdk.Address([]byte("addr3"))
	addr4 = sdk.Address([]byte("addr4"))

	asset = MsgCreateAsset{
		AssetID:  "asset1",
		Issuer:   addr,
		Name:     "asset 1",
		Quantity: 100,
	}

	asset2 = MsgCreateAsset{
		AssetID:  "asset2",
		Issuer:   addr,
		Name:     "asset 2",
		Quantity: 100,
	}

	asset3 = MsgCreateAsset{
		AssetID:  "asset3",
		Issuer:   addr,
		Name:     "asset 3",
		Quantity: 100,
	}

	assetChild = MsgCreateAsset{
		AssetID:  "asset4",
		Issuer:   addr,
		Name:     "asset 3",
		Quantity: 100,
		Parent:   "asset3",
	}

	assetChild1 = MsgCreateAsset{
		AssetID:  "asset5",
		Issuer:   addr,
		Name:     "asset 5",
		Quantity: 100,
		Parent:   "asset4",
	}
)

func TestKeeper(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	// ----------------------------------
	// Test Assets

	// Test register asset
	keeper.CreateAsset(ctx, asset)
	newAsset := keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.ID == asset.AssetID)
	assert.True(t, newAsset.Issuer.String() == asset.Issuer.String())
	assert.True(t, newAsset.Name == asset.Name)
	assert.True(t, newAsset.Quantity == asset.Quantity)

	keeper.CreateAsset(ctx, asset2)
	keeper.CreateAsset(ctx, asset3)

	// asset already exists
	_, err := keeper.CreateAsset(ctx, asset)
	assert.True(t, err != nil)

	// create asset child
	keeper.CreateAsset(ctx, assetChild)
	newAsset = keeper.GetAsset(ctx, assetChild.AssetID)
	assert.True(t, newAsset.Parent == asset3.AssetID)
	assert.True(t, newAsset.Root == asset3.AssetID)

	// invalid asset quantity
	assetChild.Quantity += 1
	_, err = keeper.CreateAsset(ctx, assetChild)
	assert.True(t, err != nil)

	keeper.CreateAsset(ctx, assetChild1)
	newAsset = keeper.GetAsset(ctx, assetChild1.AssetID)
	assert.True(t, newAsset.Parent == assetChild.AssetID)
	assert.True(t, newAsset.Root == asset3.AssetID)

	// -----------------------------------------
	// Test Add Materials
	msgAddMaterials := MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr,
		Materials: Materials{
			Material{AssetID: assetChild1.AssetID, Quantity: 1},
			Material{AssetID: asset2.AssetID, Quantity: 1},
		},
	}
	keeper.AddMaterials(ctx, msgAddMaterials)
	newAsset = keeper.GetAsset(ctx, msgAddMaterials.AssetID)
	msgAddMaterials.Materials = msgAddMaterials.Materials.Sort()
	assert.True(t, newAsset.Materials[0].AssetID == msgAddMaterials.Materials[0].AssetID)
	assert.True(t, newAsset.Materials[0].Quantity == msgAddMaterials.Materials[0].Quantity)
	assert.True(t, newAsset.Materials[1].AssetID == msgAddMaterials.Materials[1].AssetID)
	assert.True(t, newAsset.Materials[1].Quantity == msgAddMaterials.Materials[1].Quantity)

	// add materials error
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr4,
		Materials: Materials{
			Material{AssetID: assetChild1.AssetID, Quantity: 1},
			Material{AssetID: asset2.AssetID, Quantity: 1},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr4,
		Materials: Materials{
			Material{AssetID: assetChild1.AssetID, Quantity: 1000},
			Material{AssetID: asset2.AssetID, Quantity: 1},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	//-------------------------------------------
	// Test Finalize
	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: assetChild1.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)
	newAsset = keeper.GetAsset(ctx, msgFinalize.AssetID)
	assert.True(t, newAsset.Final == true)

	msgFinalize = MsgFinalize{
		Sender:  addrs[0],
		AssetID: assetChild1.AssetID,
	}
	_, err = keeper.Finalize(ctx, msgFinalize)
	assert.True(t, err != nil)

	//-------------------------------------------------
	// Test Send
	msgSend := MsgSend{Assets: []string{asset2.AssetID}, Sender: addr, Recipient: addrs[1]}
	keeper.Send(ctx, msgSend)
	newAsset = keeper.GetAsset(ctx, asset2.AssetID)
	assert.True(t, bytes.Equal(newAsset.Owner, msgSend.Recipient))

	msgSend = MsgSend{Assets: []string{asset2.AssetID}, Sender: addr, Recipient: addrs[1]}
	_, err = keeper.Send(ctx, msgSend)
	assert.True(t, err != nil)

	//-------------------------------------------------
	// Test Add Quantity

	// add quantity err
	_, err = keeper.AddQuantity(ctx, AddQuantityMsg{AssetID: assetChild1.AssetID, Issuer: addr, Quantity: 50})
	assert.True(t, err != nil)

	// Test add quantity
	keeper.AddQuantity(ctx, AddQuantityMsg{AssetID: asset.AssetID, Issuer: addr, Quantity: 50})
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Quantity == 150)

	// Test subtract quantity
	keeper.SubtractQuantity(ctx, MsgSubtractQuantity{AssetID: asset.AssetID, Issuer: addr, Quantity: 50})
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Quantity == 100)

	// Test subtract quantity error
	_, err = keeper.SubtractQuantity(ctx, MsgSubtractQuantity{AssetID: asset.AssetID, Issuer: addr, Quantity: 102})
	assert.True(t, err != nil)

	// Test Update Properties
	props := Properties{Property{Name: "weight", NumberValue: 100}, Property{Name: "size", NumberValue: 2}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Issuer: addr, Properties: props})
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	props = props.Sort()
	assert.True(t, newAsset.Properties[0].Name == props[0].Name)
	assert.True(t, newAsset.Properties[0].NumberValue == props[0].NumberValue)
	assert.True(t, newAsset.Properties[1].Name == props[1].Name)
	assert.True(t, newAsset.Properties[1].NumberValue == props[1].NumberValue)

	props2 := Properties{Property{Name: "weight", NumberValue: 150}, Property{Name: "shock", NumberValue: 2}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Issuer: addr, Properties: props2})
	props2 = props2.Sort()
	props = props.Adds(props2...)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)

	assert.True(t, newAsset.Properties[0].Name == props[0].Name)
	assert.True(t, newAsset.Properties[0].NumberValue == props[0].NumberValue)
	assert.True(t, newAsset.Properties[1].Name == props[1].Name)
	assert.True(t, newAsset.Properties[1].NumberValue == props[1].NumberValue)
	assert.True(t, newAsset.Properties[2].Name == props[2].Name)
	assert.True(t, newAsset.Properties[2].NumberValue == props[2].NumberValue)

	// Invalid property type
	props = Properties{Property{Name: "weight", NumberValue: 100, Type: 10}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Issuer: addr, Properties: props})

	// invalid issuer
	_, err = keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Issuer: addr2, Properties: props})
	assert.True(t, err != nil)

	//-------------- Test create proposal
	createProposalMsg := CreateProposalMsg{asset.AssetID, addr, addr2, []string{"weight"}, RoleReporter}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Proposals[0].Role == RoleReporter)
	assert.True(t, newAsset.Proposals[0].Issuer.String() == addr.String())
	assert.True(t, newAsset.Proposals[0].Recipient.String() == addr2.String())
	assert.True(t, newAsset.Proposals[0].Status == StatusPending)
	assert.True(t, len(newAsset.Proposals[0].Properties) == 1)
	assert.True(t, newAsset.Proposals[0].Properties[0] == "weight")

	createProposalMsg = CreateProposalMsg{asset.AssetID, addr, addr3, []string{"size", "weight"}, RoleOwner}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Proposals[1].Role == RoleOwner)
	assert.True(t, newAsset.Proposals[1].Issuer.String() == addr.String())
	assert.True(t, newAsset.Proposals[1].Recipient.String() == addr3.String())
	assert.True(t, newAsset.Proposals[1].Status == StatusPending)
	assert.True(t, len(newAsset.Proposals[1].Properties) == 2)
	assert.True(t, newAsset.Proposals[1].Properties[0] == "size")
	assert.True(t, newAsset.Proposals[1].Properties[1] == "weight")

	// Invalid role
	createProposalMsg = CreateProposalMsg{asset.AssetID, addr, addr3, []string{"size", "weight"}, 123}
	_, err = keeper.CreateProposal(ctx, createProposalMsg)
	assert.True(t, err != nil)

	//-------------- Test update properties
	createProposalMsg = CreateProposalMsg{asset.AssetID, addr, addr2, []string{"weight", "size"}, RoleOwner}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
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
	answerProposalMsg := AnswerProposalMsg{asset.AssetID, addr2, StatusAccepted}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
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
	answerProposalMsg = AnswerProposalMsg{asset.AssetID, addr2, StatusAccepted}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	assert.True(t, err != nil)

	// Answer with invalid response
	answerProposalMsg = AnswerProposalMsg{asset.AssetID, addr3, StatusPending}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	assert.True(t, err != nil)

	// Answer with invalid response
	answerProposalMsg = AnswerProposalMsg{asset.AssetID, addr2, 123}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	assert.True(t, err != nil)

	// Valid answer
	answerProposalMsg = AnswerProposalMsg{asset.AssetID, addr3, StatusRefused}
	_, err = keeper.AnswerProposal(ctx, answerProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
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
	createProposalMsg = CreateProposalMsg{asset.AssetID, addr3, addr4, []string{"weight", "size"}, RoleOwner}
	_, err = keeper.CreateProposal(ctx, createProposalMsg)
	assert.True(t, err != nil)

	// Accepted recipient with role owner is
	createProposalMsg = CreateProposalMsg{asset.AssetID, addr2, addr4, []string{"weight", "size"}, RoleOwner}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	proposal = newAsset.Proposals[2]
	assert.True(t, proposal.Role == RoleOwner)
	assert.True(t, proposal.Issuer.String() == addr2.String())
	assert.True(t, proposal.Recipient.String() == addr4.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "weight")
	assert.True(t, proposal.Properties[1] == "size")

	// Asset issuer is always authorized
	createProposalMsg = CreateProposalMsg{asset.AssetID, addr, addr4, []string{}, RoleReporter}
	keeper.CreateProposal(ctx, createProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	proposal = newAsset.Proposals[2]
	assert.True(t, proposal.Role == RoleReporter)
	assert.True(t, proposal.Issuer.String() == addr2.String())
	assert.True(t, proposal.Recipient.String() == addr4.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 2)
	assert.True(t, proposal.Properties[0] == "weight")
	assert.True(t, proposal.Properties[1] == "size")

	// Test UpdateProperties
	props = Properties{Property{Name: "weight", NumberValue: 250}}
	_, err = keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Issuer: addr2, Properties: props})
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, err == nil)

	//-------------- Test revoke proposal
	// Refused recipient is not authorized
	revokeProposalMsg := RevokeProposalMsg{asset.AssetID, addr3, addr4, []string{"weight"}}
	_, err = keeper.RevokeProposal(ctx, revokeProposalMsg)
	assert.True(t, err != nil)

	// addr2 is authorized
	revokeProposalMsg = RevokeProposalMsg{asset.AssetID, addr2, addr4, []string{"weight"}}
	_, err = keeper.RevokeProposal(ctx, revokeProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	proposal = newAsset.Proposals[2]
	assert.True(t, err == nil)
	assert.True(t, proposal.Role == RoleReporter)
	assert.True(t, proposal.Issuer.String() == addr2.String())
	assert.True(t, proposal.Recipient.String() == addr4.String())
	assert.True(t, proposal.Status == StatusPending)
	assert.True(t, len(proposal.Properties) == 1)
	assert.True(t, proposal.Properties[0] == "size")

	// proposal is deleted when there is no more property
	revokeProposalMsg = RevokeProposalMsg{asset.AssetID, addr, addr4, []string{"size"}}
	_, err = keeper.RevokeProposal(ctx, revokeProposalMsg)
	newAsset = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, len(newAsset.Proposals) == 2)
}
