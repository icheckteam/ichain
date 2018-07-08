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
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	}

	asset2 = MsgCreateAsset{
		AssetID:  "asset2",
		Sender:   addr,
		Name:     "asset 2",
		Quantity: sdk.NewInt(100),
	}

	asset3 = MsgCreateAsset{
		AssetID:  "asset3",
		Sender:   addr,
		Name:     "asset 3",
		Quantity: sdk.NewInt(100),
	}

	asset10 = MsgCreateAsset{
		AssetID:  "asset10",
		Sender:   addr,
		Name:     "asset10",
		Quantity: sdk.NewInt(100),
	}

	asset11 = MsgCreateAsset{
		AssetID:  "asset11",
		Sender:   addr2,
		Name:     "asset11",
		Quantity: sdk.NewInt(100),
	}

	assetChild = MsgCreateAsset{
		AssetID:  "asset4",
		Sender:   addr,
		Name:     "asset 3",
		Quantity: sdk.NewInt(100),
		Parent:   "asset3",
	}

	assetChild1 = MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100),
		Parent:   "asset4",
	}

	assetParentNotfound = MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100),
		Parent:   "asset4",
	}
)

func TestKeeper(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	// ----------------------------------
	// Test Assets

	// Test register asset
	keeper.CreateAsset(ctx, asset)
	newAsset, _ := keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.ID == asset.AssetID)
	assert.True(t, newAsset.Owner.String() == asset.Sender.String())
	assert.True(t, newAsset.Name == asset.Name)
	assert.True(t, newAsset.Quantity.Equal(asset.Quantity))
	assert.True(t, newAsset.Unit == "kg")

	keeper.CreateAsset(ctx, asset2)
	keeper.CreateAsset(ctx, asset3)
	keeper.CreateAsset(ctx, asset10)
	keeper.CreateAsset(ctx, asset11)

	// asset already exists
	_, err := keeper.CreateAsset(ctx, asset)
	assert.True(t, err != nil)

	// create asset child
	keeper.CreateAsset(ctx, assetChild)
	newAsset, _ = keeper.GetAsset(ctx, assetChild.AssetID)
	assert.True(t, newAsset.Parent == asset3.AssetID)
	assert.True(t, newAsset.Root == asset3.AssetID)

	// invalid asset quantity
	assetChild.Quantity = assetChild.Quantity.Add(sdk.NewInt(1))
	_, err = keeper.CreateAsset(ctx, assetChild)
	assert.True(t, err != nil)

	keeper.CreateAsset(ctx, assetChild1)
	newAsset, _ = keeper.GetAsset(ctx, assetChild1.AssetID)
	assert.True(t, newAsset.Parent == assetChild.AssetID)
	assert.True(t, newAsset.Root == asset3.AssetID)

	// valid
	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: asset10.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)
	newAsset, _ = keeper.GetAsset(ctx, msgFinalize.AssetID)
	assert.True(t, newAsset.Final == true)

	// invalid parent
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset575765",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100),
		Parent:   "asset10",
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

	// invalid owner
	msgCreateAsset = MsgCreateAsset{
		AssetID:  "asset575765",
		Sender:   addr3,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100),
		Parent:   "asset1",
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

	// invalid quantity
	msgCreateAsset = MsgCreateAsset{
		AssetID:  "6456546",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100000),
		Parent:   "asset2",
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

	msgCreateAsset = MsgCreateAsset{
		AssetID:  "asset575765",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100),
		Parent:   "asset6456456",
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

	// -----------------------------------------
	// Test Add Materials
	// -----------------------------------------------

	// test valid
	msgAddMaterials := MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr,
		Materials: Materials{
			Material{AssetID: assetChild1.AssetID, Quantity: sdk.NewInt(1)},
			Material{AssetID: asset2.AssetID, Quantity: sdk.NewInt(1)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	newAsset, _ = keeper.GetAsset(ctx, msgAddMaterials.AssetID)
	msgAddMaterials.Materials = msgAddMaterials.Materials.Sort()
	assert.True(t, newAsset.Materials[0].AssetID == msgAddMaterials.Materials[0].AssetID)
	assert.True(t, newAsset.Materials[0].Quantity.Equal(msgAddMaterials.Materials[0].Quantity))
	assert.True(t, newAsset.Materials[1].AssetID == msgAddMaterials.Materials[1].AssetID)
	assert.True(t, newAsset.Materials[1].Quantity.Equal(msgAddMaterials.Materials[1].Quantity))

	// add materials error
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr4,
		Materials: Materials{
			Material{AssetID: assetChild1.AssetID, Quantity: sdk.NewInt(1)},
			Material{AssetID: asset2.AssetID, Quantity: sdk.NewInt(1)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// invalid sender
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr4,
		Materials: Materials{
			Material{AssetID: assetChild1.AssetID, Quantity: sdk.NewInt(1000)},
			Material{AssetID: asset2.AssetID, Quantity: sdk.NewInt(1)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	//  invalid asset id
	msgAddMaterials = MsgAddMaterials{
		AssetID: "445",
		Sender:  addr4,
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// invalid asset final
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset10.AssetID,
		Sender:  addr4,
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// invalid asset
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr,
		Materials: Materials{
			Material{AssetID: "12121", Quantity: sdk.NewInt(1)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// invalid asset
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr,
		Materials: Materials{
			Material{AssetID: asset10.AssetID, Quantity: sdk.NewInt(1)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// invalid owner
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr,
		Materials: Materials{
			Material{AssetID: asset11.AssetID, Quantity: sdk.NewInt(1)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// invalid quantity
	msgAddMaterials = MsgAddMaterials{
		AssetID: asset3.AssetID,
		Sender:  addr,
		Materials: Materials{
			Material{AssetID: asset2.AssetID, Quantity: sdk.NewInt(100000)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	//-------------------------------------------
	// Test Finalize

	// invalid sender
	msgFinalize = MsgFinalize{
		Sender:  addrs[0],
		AssetID: assetChild1.AssetID,
	}
	_, err = keeper.Finalize(ctx, msgFinalize)
	assert.True(t, err != nil)

	// invalid owner
	msgFinalize = MsgFinalize{
		Sender:  addrs[1],
		AssetID: assetChild1.AssetID,
	}
	_, err = keeper.Finalize(ctx, msgFinalize)
	assert.True(t, err != nil)

	// invalid asset id
	msgFinalize = MsgFinalize{
		Sender:  addrs[1],
		AssetID: "1",
	}
	_, err = keeper.Finalize(ctx, msgFinalize)
	assert.True(t, err != nil)

	// create asset invalid parent
	msgCreateAsset = MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: sdk.NewInt(100),
		Parent:   assetChild1.AssetID,
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

	//-------------------------------------------------
	// Test Add Quantity

	// add quantity err
	_, err = keeper.AddQuantity(ctx, MsgAddQuantity{AssetID: assetChild1.AssetID, Sender: addr, Quantity: sdk.NewInt(50)})
	assert.True(t, err != nil)

	// Test add quantity
	keeper.AddQuantity(ctx, MsgAddQuantity{AssetID: asset.AssetID, Sender: addr, Quantity: sdk.NewInt(50)})
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Quantity.Equal(sdk.NewInt(150)))

	// Test subtract quantity
	keeper.SubtractQuantity(ctx, MsgSubtractQuantity{AssetID: asset.AssetID, Sender: addr, Quantity: sdk.NewInt(50)})
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Quantity.Equal(sdk.NewInt(100)))

	// Test subtract quantity error
	_, err = keeper.SubtractQuantity(ctx, MsgSubtractQuantity{AssetID: asset.AssetID, Sender: addr, Quantity: sdk.NewInt(102)})
	assert.True(t, err != nil)

	// Test Update Properties
	props := Properties{Property{Name: "weight", NumberValue: 100}, Property{Name: "size", NumberValue: 2}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr, Properties: props})
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	props = props.Sort()
	assert.True(t, newAsset.Properties[0].Name == props[0].Name)
	assert.True(t, newAsset.Properties[0].NumberValue == props[0].NumberValue)
	assert.True(t, newAsset.Properties[1].Name == props[1].Name)
	assert.True(t, newAsset.Properties[1].NumberValue == props[1].NumberValue)

	props2 := Properties{Property{Name: "weight", NumberValue: 150}, Property{Name: "shock", NumberValue: 2}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr, Properties: props2})
	props2 = props2.Sort()
	props = props.Adds(props2...)
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)

	assert.True(t, newAsset.Properties[0].Name == props[0].Name)
	assert.True(t, newAsset.Properties[0].NumberValue == props[0].NumberValue)
	assert.True(t, newAsset.Properties[1].Name == props[1].Name)
	assert.True(t, newAsset.Properties[1].NumberValue == props[1].NumberValue)
	assert.True(t, newAsset.Properties[2].Name == props[2].Name)
	assert.True(t, newAsset.Properties[2].NumberValue == props[2].NumberValue)

	// Invalid property type
	props = Properties{Property{Name: "weight", NumberValue: 100, Type: 10}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr, Properties: props})

	// invalid asset
	props = Properties{Property{Name: "weight", NumberValue: 100, Type: 10}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset10.AssetID, Sender: addr, Properties: props})

	// invalid asset
	props = Properties{Property{Name: "weight", NumberValue: 100, Type: 10}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: "adasdas", Sender: addr, Properties: props})

	// invalid issuer
	_, err = keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr2, Properties: props})
	assert.True(t, err != nil)

	// CreateProposal Tests
	// -------------------------------------------------------
	msgCreateProposal := MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"size"},
		Recipient:  addr2,
		Role:       RoleReporter,
	}
	keeper.AddProposal(ctx, msgCreateProposal)
	proposal, found := keeper.GetProposal(ctx, msgCreateProposal.AssetID, msgCreateProposal.Recipient)
	assert.True(t, found == true)
	assert.True(t, bytes.Equal(proposal.Issuer, msgCreateProposal.Sender))
	assert.True(t, bytes.Equal(proposal.Recipient, msgCreateProposal.Recipient))
	assert.True(t, proposal.Properties[0] == msgCreateProposal.Properties[0])
	assert.True(t, proposal.Role == msgCreateProposal.Role)

	msgCreateProposal = MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"size"},
		Recipient:  addr3,
		Role:       RoleReporter,
	}
	keeper.AddProposal(ctx, msgCreateProposal)

	// invalid sender
	msgCreateProposal = MsgCreateProposal{
		Sender:  addr2,
		AssetID: asset.AssetID,
	}
	_, err = keeper.AddProposal(ctx, msgCreateProposal)
	assert.True(t, err != nil)

	// invalid asset
	msgCreateProposal = MsgCreateProposal{
		Sender:  addr2,
		AssetID: "adad",
	}
	_, err = keeper.AddProposal(ctx, msgCreateProposal)
	assert.True(t, err != nil)

	// AnswerProposal Tests
	// -------------------------------------------------------
	msgAnswerProposal := MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Recipient: addr2,
		Response:  StatusAccepted,
	}
	keeper.AnswerProposal(ctx, msgAnswerProposal)
	newAsset, _ = keeper.GetAsset(ctx, msgAnswerProposal.AssetID)
	assert.True(t, bytes.Equal(newAsset.Reporters[0].Addr, msgAnswerProposal.Recipient))
	assert.True(t, newAsset.Reporters[0].Properties[0] == "size")

	msgCreateProposal = MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"size"},
		Recipient:  addr2,
		Role:       RoleReporter,
	}
	keeper.AddProposal(ctx, msgCreateProposal)

	// valid
	keeper.AnswerProposal(ctx, msgAnswerProposal)
	keeper.AnswerProposal(ctx, msgAnswerProposal)
	newAsset, _ = keeper.GetAsset(ctx, msgAnswerProposal.AssetID)
	assert.True(t, bytes.Equal(newAsset.Reporters[0].Addr, msgAnswerProposal.Recipient))
	assert.True(t, newAsset.Reporters[0].Properties[0] == "size")

	// invalid proposal
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Recipient: addr2,
		Response:  1,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err != nil)

	msgCreateProposal = MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"size"},
		Recipient:  addr2,
		Role:       RoleOwner,
	}
	keeper.AddProposal(ctx, msgCreateProposal)
	msgCreateProposal = MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"size"},
		Recipient:  addr3,
		Role:       RoleReporter,
	}
	keeper.AddProposal(ctx, msgCreateProposal)
	keeper.AnswerProposal(ctx, msgAnswerProposal)
	newAsset, _ = keeper.GetAsset(ctx, msgAnswerProposal.AssetID)
	assert.True(t, bytes.Equal(newAsset.Owner, msgAnswerProposal.Recipient))

	keeper.AnswerProposal(ctx, MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Recipient: addr3,
		Response:  1,
	})
}

func TestFinalize(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// create asset
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	}
	keeper.CreateAsset(ctx, msgCreateAsset)

	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: msgCreateAsset.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)

	// invalid owner
	msgFinalize = MsgFinalize{
		Sender:  addr,
		AssetID: msgCreateAsset.AssetID,
	}
	_, err := keeper.Finalize(ctx, msgFinalize)
	assert.True(t, err != nil)

}

func TestSubtractQuantity(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// create asset
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	}
	keeper.CreateAsset(ctx, msgCreateAsset)

	// invalid asset
	msgSubtractQuantity := MsgSubtractQuantity{
		AssetID:  "45345",
		Sender:   addr,
		Quantity: sdk.NewInt(102),
	}
	_, err := keeper.SubtractQuantity(ctx, msgSubtractQuantity)
	assert.True(t, err != nil)

	// invalid asset
	msgSubtractQuantity = MsgSubtractQuantity{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr2,
		Quantity: sdk.NewInt(102),
	}
	_, err = keeper.SubtractQuantity(ctx, msgSubtractQuantity)
	assert.True(t, err != nil)

	// invalid asset
	msgSubtractQuantity = MsgSubtractQuantity{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr2,
		Quantity: sdk.NewInt(102),
	}
	_, err = keeper.SubtractQuantity(ctx, msgSubtractQuantity)
	assert.True(t, err != nil)

	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: msgCreateAsset.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)

	_, err = keeper.SubtractQuantity(ctx, msgSubtractQuantity)
	assert.True(t, err != nil)

}

func TestAddQuantity(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// create asset
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	}
	keeper.CreateAsset(ctx, msgCreateAsset)

	// invalid asset
	msgAddQuantity := MsgAddQuantity{
		AssetID:  "45345",
		Sender:   addr,
		Quantity: sdk.NewInt(102),
	}
	_, err := keeper.AddQuantity(ctx, msgAddQuantity)
	assert.True(t, err != nil)

	// invalid asset
	msgAddQuantity = MsgAddQuantity{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr2,
		Quantity: sdk.NewInt(102),
	}
	_, err = keeper.AddQuantity(ctx, msgAddQuantity)
	assert.True(t, err != nil)

	// invalid asset
	msgAddQuantity = MsgAddQuantity{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr2,
		Quantity: sdk.NewInt(102),
	}
	_, err = keeper.AddQuantity(ctx, msgAddQuantity)
	assert.True(t, err != nil)

	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: msgCreateAsset.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)

	_, err = keeper.AddQuantity(ctx, msgAddQuantity)
	assert.True(t, err != nil)

}

func TestRevokeReporter(t *testing.T) {

}
