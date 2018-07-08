package asset

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	attr = Property{
		Name:        "weight",
		Type:        3,
		NumberValue: 100,
	}

	props = Properties{
		attr,
	}
)

// ----------------------------------------
// Register Tests

func TestRegisterMsg(t *testing.T) {
	msg := NewMsgCreateAsset(addr, "1", "1", sdk.NewInt(10), "1")
	assert.Equal(t, msg.Sender, addr)
	assert.Equal(t, msg.AssetID, "1")
	assert.Equal(t, msg.Name, "1")
	assert.Equal(t, msg.Quantity, sdk.NewInt(10))
	assert.Equal(t, msg.Parent, "1")

}

func TestCreateAssetMsgType(t *testing.T) {
	addr := sdk.Address([]byte("input"))

	var msg = MsgCreateAsset{
		AssetID:  "1",
		Name:     "asset name",
		Quantity: sdk.NewInt(100),
		Sender:   addr,
	}

	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestCreateAssetMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgCreateAsset
	}{
		{false, MsgCreateAsset{}},                                                                              // no asset info
		{false, MsgCreateAsset{Sender: addr1, Quantity: sdk.NewInt(0), Name: "name", AssetID: "1"}},            // missing quantity
		{false, MsgCreateAsset{Sender: addr1, Quantity: sdk.NewInt(1), Name: "name"}},                          // missing id
		{false, MsgCreateAsset{Sender: addr1, Quantity: sdk.NewInt(1), Name: "name", AssetID: "1"}},            //
		{true, MsgCreateAsset{Sender: addr1, Quantity: sdk.NewInt(1), Name: "name", AssetID: "1", Unit: "kg"}}, //
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

func TestRegisterMsgGet(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateAsset{
		Sender: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestRegisterGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateAsset{
		Sender:   addr1,
		AssetID:  "1212",
		Name:     "name",
		Quantity: sdk.NewInt(1),
		Unit:     "kg",
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"asset_id":"1212","asset_type":"","name":"name","parent":"","quantity":"1","sender":"cosmosaccaddr1d9h8qat5e4ehc5","unit":"kg"}`)
}

func TestRegisterGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateAsset{
		Sender:   addr1,
		AssetID:  "1212",
		Name:     "name",
		Quantity: sdk.NewInt(1),
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// Update Attribute Tests
func TestUpdateAttrMsgMsg(t *testing.T) {

}

func TestUpdateAttrMsgType(t *testing.T) {
	addr := sdk.Address([]byte("input"))

	var msg = MsgUpdateProperties{
		Sender: addr,
		Properties: []Property{
			attr,
		},
	}

	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")

}

func TestUpdateAttrMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgUpdateProperties
	}{
		{false, MsgUpdateProperties{}},              // no asset info
		{false, MsgUpdateProperties{Sender: addr1}}, // only set owner
		{false, MsgUpdateProperties{Sender: addr1, Properties: []Property{
			attr,
		}}}, // missing id
		{true, MsgUpdateProperties{Sender: addr1, Properties: []Property{
			attr,
		}, AssetID: "1212"}},
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

func TestUpdateAttrMsgGet(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgUpdateProperties{
		Sender: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestUpdateAttrMsgGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgUpdateProperties{
		Sender:     addr1,
		AssetID:    "1",
		Properties: props,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"sender":"696E707574","asset_id":"1","properties":[{"name":"weight","type":3,"number_value":100,"location_value":{"latitude":0,"longitude":0}}]}`)
}

func TestUpdateAttrGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgUpdateProperties{
		Sender: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// AddQuantity Tests
func TestAddQuantityMsg(t *testing.T) {

}

func TestAddQuantityType(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgAddQuantity{
		Sender:   addr,
		AssetID:  "!",
		Quantity: sdk.NewInt(1),
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestAddQuantityMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgAddQuantity
	}{
		{false, MsgAddQuantity{}},                                                     // no asset info
		{false, MsgAddQuantity{Sender: addr1, Quantity: sdk.NewInt(0), AssetID: "1"}}, // missing quantity
		{false, MsgAddQuantity{Sender: addr1, Quantity: sdk.NewInt(1)}},               // missing id
		{true, MsgAddQuantity{Sender: addr1, Quantity: sdk.NewInt(1), AssetID: "1"}},  //
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

func TestAddQuantityMsgGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgAddQuantity{
		Sender:   addr1,
		AssetID:  "1",
		Quantity: sdk.NewInt(1),
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"sender\":\"696E707574\",\"asset_id\":\"1\",\"quantity\":\"1\"}")
}

func TestAddQuantityGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgAddQuantity{
		Sender: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// AddQuantity Tests
func TestSubtractQuantityMsg(t *testing.T) {

}

func TestSubtractQuantityMsgType(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgSubtractQuantity{
		Sender:   addr,
		AssetID:  "!",
		Quantity: sdk.NewInt(1),
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestSubtractQuantityMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgSubtractQuantity
	}{
		{false, MsgSubtractQuantity{}},                                                     // no asset info
		{false, MsgSubtractQuantity{Sender: addr1, Quantity: sdk.NewInt(0), AssetID: "1"}}, // missing quantity
		{false, MsgSubtractQuantity{Sender: addr1, Quantity: sdk.NewInt(1)}},               // missing id
		{true, MsgSubtractQuantity{Sender: addr1, Quantity: sdk.NewInt(1), AssetID: "1"}},  //
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

func TestSubtractQuantityMsgGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgSubtractQuantity{
		Sender:   addr1,
		AssetID:  "1",
		Quantity: sdk.NewInt(1),
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"sender":"696E707574","asset_id":"1","quantity":"1"}`)
}

func TestSubtractQuantityMsgGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgAddQuantity{
		Sender: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// TestAddMaterialsMsg Tests
func TestAddMaterialsMsg(t *testing.T) {

}

func TestMsgAddMaterialsType(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgAddMaterials{
		Sender: addr,
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestMsgAddMaterialsValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgAddMaterials
	}{
		{false, MsgAddMaterials{}},                                                                                                   // no asset info
		{false, MsgAddMaterials{Sender: addr1, AssetID: "1"}},                                                                        // missing quantity
		{false, MsgAddMaterials{Sender: addr1}},                                                                                      // missing id
		{false, MsgAddMaterials{Sender: addr1, AssetID: "1", Materials: Materials{Material{AssetID: "1", Quantity: sdk.NewInt(0)}}}}, //
		{true, MsgAddMaterials{Sender: addr1, AssetID: "1", Materials: Materials{Material{AssetID: "1", Quantity: sdk.NewInt(1)}}}},  //
		{false, MsgAddMaterials{Sender: addr1, AssetID: "1", Materials: Materials{Material{Quantity: sdk.NewInt(1)}}}},               //
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

func TestMsgAddMaterialsGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgAddMaterials{
		Sender:    addr1,
		AssetID:   "1",
		Materials: Materials{Material{AssetID: "1", Quantity: sdk.NewInt(1)}},
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"asset_id":"1","sender":"696E707574","materials":[{"asset_id":"1","quantity":"1"}]}`)
}

func TestMsgGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgAddMaterials{
		Sender: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// TestMsgFinalize Tests
func TestMsgFinalize(t *testing.T) {

}

func TestMsgFinalizeType(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgFinalize{
		Sender:  addr,
		AssetID: "121",
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestMsgMsgFinalizeValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgFinalize
	}{
		{false, MsgFinalize{}},              // no asset info
		{false, MsgFinalize{Sender: addr1}}, // missing id
		{true, MsgFinalize{Sender: addr1, AssetID: "23213"}},
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

func TestMsgFinalizeGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgFinalize{
		Sender:  addr1,
		AssetID: "3434",
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"sender":"696E707574","asset_id":"3434"}`)
}

func TestMsgFinalizeGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgFinalize{
		Sender: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// TestMsgRevokeReporter Tests
func TestMsgRevokeReporter(t *testing.T) {

}

func TestMsgMsgRevokeReporterType(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	var msg = MsgRevokeReporter{
		Sender:  addr,
		AssetID: "121",
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestMsgRevokeReporterValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    MsgRevokeReporter
	}{
		{false, MsgRevokeReporter{}},              // no asset info
		{false, MsgRevokeReporter{Sender: addr1}}, // missing id
		{false, MsgRevokeReporter{Sender: addr1, Reporter: addr3}},
		{true, MsgRevokeReporter{Sender: addr1, Reporter: addr3, AssetID: "1"}},
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

func TestMsgRevokeReporterGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgRevokeReporter{
		Sender:   addr1,
		Reporter: addr2,
		AssetID:  "3434",
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"sender":"696E707574","reporter":"6164647232","asset_id":"3434"}`)
}

func TestMsgRevokeReporterGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgRevokeReporter{
		Sender: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}
