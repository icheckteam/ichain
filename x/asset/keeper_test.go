package asset

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	addr  = sdk.AccAddress([]byte("addr1"))
	addr2 = sdk.AccAddress([]byte("addr2"))
	addr3 = sdk.AccAddress([]byte("addr3"))
	addr4 = sdk.AccAddress([]byte("addr4"))

	asset = MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
		Properties: Properties{
			Property{Name: "size", StringValue: "size"},
			Property{Name: "barcode", StringValue: "barcode"},
			Property{Name: "type", StringValue: "type"},
			Property{Name: "subtype", StringValue: "subtype"},
		},
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
	properties := keeper.GetProperties(ctx, newAsset.ID)
	assert.True(t, len(properties) == 4)
	assert.True(t, properties[0].Name == "barcode")
	assert.True(t, properties[0].StringValue == "barcode")

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

	// answer role reporter

	// Cancel/invalid owner
	msgAnswerProposal := MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr2,
		Response:  StatusCancel,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err != nil)

	// StatusRejected/invalid recipient
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr2,
		Response:  StatusRejected,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err != nil)

	// Accepted/invalid recipient
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr2,
		Response:  StatusAccepted,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err != nil)

	// invalid response
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr2,
		Response:  5,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err != nil)

	// invalid proposal
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr4,
		Response:  StatusAccepted,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err != nil)

	// valid
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr2,
		Recipient: addr2,
		Response:  StatusAccepted,
		Role:      RoleReporter,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	reporters := keeper.GetReporters(ctx, msgAnswerProposal.AssetID)
	assert.True(t, len(reporters) == 1)

	// test validate update property authorization
	props = Properties{Property{Name: "size", NumberValue: 100, Type: 10}}
	_, err = keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr2, Properties: props})
	assert.True(t, err == nil)

	// update reporter
	msgCreateProposal = MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"weight"},
		Recipient:  addr2,
		Role:       RoleReporter,
	}
	keeper.AddProposal(ctx, msgCreateProposal)
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr2,
		Recipient: addr2,
		Response:  StatusAccepted,
		Role:      RoleReporter,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err == nil)

	// create proposal change owner
	msgCreateProposal = MsgCreateProposal{
		Sender:     addr,
		AssetID:    asset.AssetID,
		Properties: []string{"size"},
		Recipient:  addr2,
		Role:       RoleOwner,
	}
	keeper.AddProposal(ctx, msgCreateProposal)

	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr2,
		Recipient: addr2,
		Response:  StatusAccepted,
		Role:      RoleOwner,
	}
	keeper.AnswerProposal(ctx, msgAnswerProposal)
	newAsset, _ = keeper.GetAsset(ctx, msgAnswerProposal.AssetID)
	assert.True(t, bytes.Equal(msgAnswerProposal.Recipient, newAsset.Owner))

	// delete proposal
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr3,
		Response:  StatusAccepted,
		Role:      RoleReporter,
	}
	_, err = keeper.AnswerProposal(ctx, msgAnswerProposal)
	assert.True(t, err == nil)
	_, found = keeper.GetProposal(ctx, asset.AssetID, addr3)
	assert.True(t, !found)

	// RevokeReporter Test
	//--------------------------------------------------------------

	// invalid asset id
	msgRevokeReporter := MsgRevokeReporter{
		Sender:   addr2,
		Reporter: addr3,
		AssetID:  "adasdas",
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	// invalid asset final
	msgRevokeReporter = MsgRevokeReporter{
		Sender:   addr2,
		Reporter: addr3,
		AssetID:  asset10.AssetID,
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	// invalid owner
	msgRevokeReporter = MsgRevokeReporter{
		Sender:   addr,
		Reporter: addr3,
		AssetID:  asset.AssetID,
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	// invalid reporter
	msgRevokeReporter = MsgRevokeReporter{
		Sender:   addr2,
		Reporter: addr,
		AssetID:  asset.AssetID,
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	// create reporter test revoke
	msgCreateProposal = MsgCreateProposal{
		Sender:     addr2,
		AssetID:    asset.AssetID,
		Properties: []string{"weight"},
		Recipient:  addr3,
		Role:       RoleReporter,
	}
	keeper.AddProposal(ctx, msgCreateProposal)
	msgAnswerProposal = MsgAnswerProposal{
		AssetID:   asset.AssetID,
		Sender:    addr3,
		Recipient: addr3,
		Response:  StatusAccepted,
	}
	keeper.AnswerProposal(ctx, msgAnswerProposal)

	// valid reporter
	msgRevokeReporter = MsgRevokeReporter{
		Sender:   addr2,
		Reporter: addr3,
		AssetID:  asset.AssetID,
	}
	keeper.RevokeReporter(ctx, msgRevokeReporter)
	reporters = keeper.GetReporters(ctx, msgAnswerProposal.AssetID)
	assert.True(t, len(reporters) == 0)
}

func TestKeeperUpdateProperties(t *testing.T) {
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

	// Test Update Properties
	props := Properties{Property{Name: "weight", NumberValue: 100}, Property{Name: "size", NumberValue: 2}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr, Properties: props})
	properties := keeper.GetProperties(ctx, msgCreateAsset.AssetID)
	assert.True(t, len(properties) == 2)

	props2 := Properties{Property{Name: "weight", NumberValue: 150}, Property{Name: "shock", NumberValue: 2}}
	keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr, Properties: props2})
	properties = keeper.GetProperties(ctx, msgCreateAsset.AssetID)
	assert.True(t, len(properties) == 3)

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
	_, err := keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr2, Properties: props})
	assert.True(t, err != nil)
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
	createRecordTest(ctx, keeper)

	// Invalid Record Owner
	msg := MsgSubtractQuantity{
		AssetID:  "asseta",
		Quantity: sdk.NewInt(10),
		Sender:   addr2,
	}
	_, err := keeper.SubtractQuantity(ctx, msg)
	assert.True(t, err != nil)

	// Invalid Record ID (ID does not exists)
	msg = MsgSubtractQuantity{
		AssetID:  "asset6",
		Quantity: sdk.NewInt(10),
		Sender:   addr2,
	}
	_, err = keeper.SubtractQuantity(ctx, msg)
	assert.True(t, err != nil)

	// Invalid Record = Record has final
	msg = MsgSubtractQuantity{
		AssetID:  "asset4",
		Quantity: sdk.NewInt(10),
		Sender:   addr,
	}
	_, err = keeper.SubtractQuantity(ctx, msg)
	assert.True(t, err != nil)

	// Valid
	msg = MsgSubtractQuantity{
		AssetID:  "asseta",
		Quantity: sdk.NewInt(10),
		Sender:   addr,
	}
	keeper.SubtractQuantity(ctx, msg)
	record, _ := keeper.GetAsset(ctx, "asseta")
	assert.True(t, record.Quantity.Equal(sdk.NewInt(90)))
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

func createRecordTest(ctx sdk.Context, keeper Keeper) {
	keeper.CreateAsset(ctx, MsgCreateAsset{
		AssetID:  "asseta",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	})

	keeper.CreateAsset(ctx, MsgCreateAsset{
		AssetID:  "assetb",
		Sender:   addr,
		Name:     "asset 2",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	})
	keeper.CreateAsset(ctx, MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 2",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	})
	keeper.CreateAsset(ctx, MsgCreateAsset{
		AssetID:  "asset2",
		Sender:   addr,
		Name:     "asset 2",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	})
	keeper.CreateAsset(ctx, MsgCreateAsset{
		AssetID:  "asset3",
		Sender:   addr2,
		Name:     "asset 3",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	})

	keeper.CreateAsset(ctx, MsgCreateAsset{
		AssetID:  "asset4",
		Sender:   addr2,
		Name:     "asset 4",
		Unit:     "kg",
		Quantity: sdk.NewInt(100),
	})

	keeper.Finalize(ctx, MsgFinalize{
		AssetID: "asset4",
		Sender:  addr2,
	})
}

func TestKeeperAddMaterials(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	createRecordTest(ctx, keeper)

	// Msg Add Materials
	msgAddMaterials := MsgAddMaterials{
		Sender:  addr,
		AssetID: "asseta",
		Amount: []Material{
			Material{RecordID: "assetb", Amount: sdk.NewInt(10)},
		},
	}

	// Test Valid
	keeper.AddMaterials(ctx, msgAddMaterials)
	materials := keeper.GetMaterials(ctx, "asseta")
	record, _ := keeper.GetAsset(ctx, "assetb")
	assert.True(t, record.Quantity.Equal(sdk.NewInt(90)), fmt.Sprintf("%s", record.Quantity))
	assert.True(t, materials[0].RecordID == "assetb")
	assert.True(t, materials[0].Amount.Equal(sdk.NewInt(10)))

	// Test Add More Material
	msgAddMaterials = MsgAddMaterials{
		Sender:  addr,
		AssetID: "asseta",
		Amount: []Material{
			Material{RecordID: "assetb", Amount: sdk.NewInt(10)},
		},
	}
	keeper.AddMaterials(ctx, msgAddMaterials)
	materials = keeper.GetMaterials(ctx, "asseta")
	assert.True(t, materials[0].RecordID == "assetb")
	assert.True(t, materials[0].Amount.Equal(sdk.NewInt(20)))

	// Test Add More Quantity
	msgAddMaterials = MsgAddMaterials{
		Sender:  addr,
		AssetID: "asseta",
		Amount: []Material{
			Material{RecordID: "asset1", Amount: sdk.NewInt(10)},
		},
	}
	keeper.AddMaterials(ctx, msgAddMaterials)
	materials = keeper.GetMaterials(ctx, "asseta")
	assert.True(t, materials[0].RecordID == "asset1")
	assert.True(t, materials[0].Amount.Equal(sdk.NewInt(10)))
	assert.True(t, materials[1].RecordID == "assetb")
	assert.True(t, materials[1].Amount.Equal(sdk.NewInt(20)))

	// Invalid Record Owner
	msgAddMaterials = MsgAddMaterials{
		Sender:  addr2,
		AssetID: "asseta",
		Amount: []Material{
			Material{RecordID: "asset1", Amount: sdk.NewInt(10)},
		},
	}
	_, err := keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// Invalid Material Record Owner
	msgAddMaterials = MsgAddMaterials{
		Sender:  addr,
		AssetID: "asseta",
		Amount: []Material{
			Material{RecordID: "asset3", Amount: sdk.NewInt(10)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// Invalid Material Record ID does not exists
	msgAddMaterials = MsgAddMaterials{
		Sender:  addr,
		AssetID: "asseta",
		Amount: []Material{
			Material{RecordID: "asset5", Amount: sdk.NewInt(10)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)

	// Invalid Record - Record Does not exists
	msgAddMaterials = MsgAddMaterials{
		Sender:  addr,
		AssetID: "asset6",
		Amount: []Material{
			Material{RecordID: "asset1", Amount: sdk.NewInt(10)},
		},
	}
	_, err = keeper.AddMaterials(ctx, msgAddMaterials)
	assert.True(t, err != nil)
}
