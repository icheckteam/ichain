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
		Quantity: 100,
		Properties: Properties{Property{
			Name:        "unit",
			Type:        PropertyTypeString,
			StringValue: "kg",
		}},
	}

	asset2 = MsgCreateAsset{
		AssetID:  "asset2",
		Sender:   addr,
		Name:     "asset 2",
		Quantity: 100,
	}

	asset3 = MsgCreateAsset{
		AssetID:  "asset3",
		Sender:   addr,
		Name:     "asset 3",
		Quantity: 100,
	}

	assetChild = MsgCreateAsset{
		AssetID:  "asset4",
		Sender:   addr,
		Name:     "asset 3",
		Quantity: 100,
		Parent:   "asset3",
	}

	assetChild1 = MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: 100,
		Parent:   "asset4",
	}

	assetParentNotfound = MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
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
	newAsset, _ := keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.ID == asset.AssetID)
	assert.True(t, newAsset.Owner.String() == asset.Sender.String())
	assert.True(t, newAsset.Name == asset.Name)
	assert.True(t, newAsset.Quantity == asset.Quantity)
	assert.True(t, newAsset.Unit == "kg")

	keeper.CreateAsset(ctx, asset2)
	keeper.CreateAsset(ctx, asset3)

	// asset already exists
	_, err := keeper.CreateAsset(ctx, asset)
	assert.True(t, err != nil)

	// create asset child
	keeper.CreateAsset(ctx, assetChild)
	newAsset, _ = keeper.GetAsset(ctx, assetChild.AssetID)
	assert.True(t, newAsset.Parent == asset3.AssetID)
	assert.True(t, newAsset.Root == asset3.AssetID)

	// invalid asset quantity
	assetChild.Quantity += 1
	_, err = keeper.CreateAsset(ctx, assetChild)
	assert.True(t, err != nil)

	keeper.CreateAsset(ctx, assetChild1)
	newAsset, _ = keeper.GetAsset(ctx, assetChild1.AssetID)
	assert.True(t, newAsset.Parent == assetChild.AssetID)
	assert.True(t, newAsset.Root == asset3.AssetID)

	// invalid parent
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: 100,
		Parent:   "asset45",
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

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
	newAsset, _ = keeper.GetAsset(ctx, msgAddMaterials.AssetID)
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

	// invalid sender
	msgFinalize := MsgFinalize{
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

	// valid
	msgFinalize = MsgFinalize{
		Sender:  addr,
		AssetID: assetChild1.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)
	newAsset, _ = keeper.GetAsset(ctx, msgFinalize.AssetID)
	assert.True(t, newAsset.Final == true)

	// create asset invalid parent
	msgCreateAsset = MsgCreateAsset{
		AssetID:  "asset5",
		Sender:   addr,
		Name:     "asset 5",
		Quantity: 100,
		Parent:   assetChild1.AssetID,
	}
	_, err = keeper.CreateAsset(ctx, msgCreateAsset)
	assert.True(t, err != nil)

	//-------------------------------------------------
	// Test Add Quantity

	// add quantity err
	_, err = keeper.AddQuantity(ctx, MsgAddQuantity{AssetID: assetChild1.AssetID, Sender: addr, Quantity: 50})
	assert.True(t, err != nil)

	// Test add quantity
	keeper.AddQuantity(ctx, MsgAddQuantity{AssetID: asset.AssetID, Sender: addr, Quantity: 50})
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Quantity == 150)

	// Test subtract quantity
	keeper.SubtractQuantity(ctx, MsgSubtractQuantity{AssetID: asset.AssetID, Sender: addr, Quantity: 50})
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, newAsset.Quantity == 100)

	// Test subtract quantity error
	_, err = keeper.SubtractQuantity(ctx, MsgSubtractQuantity{AssetID: asset.AssetID, Sender: addr, Quantity: 102})
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

	// invalid issuer
	_, err = keeper.UpdateProperties(ctx, MsgUpdateProperties{AssetID: asset.AssetID, Sender: addr2, Properties: props})
	assert.True(t, err != nil)

	// Test CreateReporter
	msgCreateReporter := MsgCreateReporter{
		AssetID:    asset.AssetID,
		Sender:     addr,
		Reporter:   addr2,
		Properties: []string{"size"},
	}
	keeper.CreateReporter(ctx, msgCreateReporter)
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, bytes.Equal(newAsset.Reporters[0].Addr, msgCreateReporter.Reporter))
	assert.True(t, newAsset.Reporters[0].Properties[0] == msgCreateReporter.Properties[0])

	keeper.CreateReporter(ctx, msgCreateReporter)
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, bytes.Equal(newAsset.Reporters[0].Addr, msgCreateReporter.Reporter))
	assert.True(t, newAsset.Reporters[0].Properties[0] == msgCreateReporter.Properties[0])

	//Test invalid sender
	msgCreateReporter.Sender = addr4
	_, err = keeper.CreateReporter(ctx, msgCreateReporter)
	assert.True(t, err != nil)

	// Tét Transfer asset

	// Test invalid sender
	msgTransfer := MsgTransfer{
		Assets:    []string{asset.AssetID},
		Sender:    addr3,
		Recipient: addr3,
	}
	_, err = keeper.Transfer(ctx, msgTransfer)
	assert.True(t, err != nil)

	msgTransfer = MsgTransfer{
		Assets:    []string{"adasd"},
		Sender:    addr,
		Recipient: addr3,
	}
	_, err = keeper.Transfer(ctx, msgTransfer)
	assert.True(t, err != nil)

	msgTransfer = MsgTransfer{
		Assets:    []string{asset.AssetID},
		Sender:    addr,
		Recipient: addr3,
	}
	keeper.Transfer(ctx, msgTransfer)
	newAsset, _ = keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, bytes.Equal(newAsset.Owner, msgTransfer.Recipient))
	assert.True(t, newAsset.Reporters == nil)

}

func TestFinalize(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// create asset
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: 100,
	}
	keeper.CreateAsset(ctx, msgCreateAsset)

	// invalid asset
	msgCreateReporter := MsgCreateReporter{
		AssetID:    "dasdasd",
		Sender:     addr,
		Reporter:   addr2,
		Properties: []string{"size"},
	}
	_, err := keeper.CreateReporter(ctx, msgCreateReporter)
	assert.True(t, err != nil)

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
	_, err = keeper.Finalize(ctx, msgFinalize)
	assert.True(t, err != nil)

}

func TestCreateReporter(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// create asset
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: 100,
	}
	keeper.CreateAsset(ctx, msgCreateAsset)

	// invalid asset
	msgCreateReporter := MsgCreateReporter{
		AssetID:    "dasdasd",
		Sender:     addr,
		Reporter:   addr2,
		Properties: []string{"size"},
	}
	_, err := keeper.CreateReporter(ctx, msgCreateReporter)
	assert.True(t, err != nil)

	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: msgCreateAsset.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)

	msgCreateReporter = MsgCreateReporter{
		AssetID:    msgCreateAsset.AssetID,
		Sender:     addr,
		Reporter:   addr2,
		Properties: []string{"size"},
	}
	_, err = keeper.CreateReporter(ctx, msgCreateReporter)
	assert.True(t, err != nil)

}

func TestRevokeReporter(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// create asset
	msgCreateAsset := MsgCreateAsset{
		AssetID:  "asset1",
		Sender:   addr,
		Name:     "asset 1",
		Unit:     "kg",
		Quantity: 100,
	}
	keeper.CreateAsset(ctx, msgCreateAsset)

	// create reporter
	msgCreateReporter := MsgCreateReporter{
		AssetID:    msgCreateAsset.AssetID,
		Sender:     addr,
		Reporter:   addr2,
		Properties: []string{"size"},
	}
	keeper.CreateReporter(ctx, msgCreateReporter)

	// asset not found
	msgRevokeReporter := MsgRevokeReporter{
		AssetID:  "add3",
		Sender:   addr,
		Reporter: addr2,
	}
	_, err := keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	// invalid owner
	msgRevokeReporter = MsgRevokeReporter{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr2,
		Reporter: addr,
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	// invalid reporter
	msgRevokeReporter = MsgRevokeReporter{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr,
		Reporter: addr4,
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

	msgRevokeReporter = MsgRevokeReporter{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr,
		Reporter: addr2,
	}
	keeper.RevokeReporter(ctx, msgRevokeReporter)
	newAsset, _ := keeper.GetAsset(ctx, asset.AssetID)
	assert.True(t, len(newAsset.Reporters) == 0)

	// create reporter
	msgCreateReporter = MsgCreateReporter{
		AssetID:    msgCreateAsset.AssetID,
		Sender:     addr,
		Reporter:   addr2,
		Properties: []string{"size"},
	}
	keeper.CreateReporter(ctx, msgCreateReporter)

	msgFinalize := MsgFinalize{
		Sender:  addr,
		AssetID: msgCreateAsset.AssetID,
	}
	keeper.Finalize(ctx, msgFinalize)

	// invalid asset
	// invalid reporter
	msgRevokeReporter = MsgRevokeReporter{
		AssetID:  msgCreateAsset.AssetID,
		Sender:   addr,
		Reporter: addr2,
	}
	_, err = keeper.RevokeReporter(ctx, msgRevokeReporter)
	assert.True(t, err != nil)

}
