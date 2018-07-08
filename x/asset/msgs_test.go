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
		Sender:     addr1,
		AssetID:    "1212",
		Name:       "name",
		Quantity:   sdk.NewInt(1),
		Unit:       "kg",
		Properties: Properties{Property{Name: "size", StringValue: "50"}},
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"asset_id\":\"1212\",\"name\":\"name\",\"parent\":\"\",\"properties\":[{\"location_value\":{\"latitude\":\"0\",\"longitude\":\"0\"},\"name\":\"size\",\"string_value\":\"50\",\"type\":\"0\"}],\"quantity\":\"1\",\"sender\":\"cosmosaccaddr1d9h8qat5e4ehc5\",\"unit\":\"kg\"}")
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
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"properties\":[{\"location_value\":{\"latitude\":\"0\",\"longitude\":\"0\"},\"name\":\"weight\",\"number_value\":\"100\",\"type\":\"3\"}],\"sender\":\"cosmosaccaddr1d9h8qat5e4ehc5\"}")
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
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"quantity\":\"1\",\"sender\":\"cosmosaccaddr1d9h8qat5e4ehc5\"}")
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
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"materials\":[{\"asset_id\":\"1\",\"quantity\":\"1\"}],\"sender\":\"cosmosaccaddr1d9h8qat5e4ehc5\"}")
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
	assert.Equal(t, string(res), "{\"asset_id\":\"3434\",\"reporter\":\"cosmosaccaddr1v9jxgu3jlsw7dy\",\"sender\":\"cosmosaccaddr1d9h8qat5e4ehc5\"}")
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

// CreateProposal  Tests
// ------------------------------------------------------------
func TestCreateProposalMsgType(t *testing.T) {
	msg := MsgCreateProposal{}
	assert.Equal(t, msg.Type(), "asset")
}

func TestCreateProposalMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgCreateProposal
	}{
		{false, MsgCreateProposal{AssetID: "1"}},
		{false, MsgCreateProposal{AssetID: "1", Sender: addrs[0]}},
		{false, MsgCreateProposal{AssetID: "1", Recipient: addrs[0]}},
		{false, MsgCreateProposal{AssetID: "1", Properties: []string{"location"}}},
		{false, MsgCreateProposal{AssetID: "1", Sender: addrs[0], Recipient: addrs[1]}},
		{false, MsgCreateProposal{AssetID: "1", Sender: addrs[0], Recipient: addrs[1], Role: 0, Properties: []string{"location"}}},
		{true, MsgCreateProposal{AssetID: "1", Sender: addrs[0], Recipient: addrs[1], Role: 1, Properties: []string{"location"}}},
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

func TestCreateProposalMsgGet(t *testing.T) {
	msg := MsgCreateProposal{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCreateProposalMsgGetSigners(t *testing.T) {
	msg := MsgCreateProposal{
		Sender: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6100]`)
}

func TestCreateProposalMsgGetSignBytes(t *testing.T) {
	msg := MsgCreateProposal{
		AssetID:    "1",
		Properties: []string{"location"},
		Sender:     addrs[0],
		Recipient:  addrs[1],
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"properties\":[\"location\"],\"recipient\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgpg7hpw0\",\"role\":\"0\",\"sender\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq4gr5na\"}")
}

// AnswerProposal  Tests
// ------------------------------------------------------------

func TestAnswerProposalMsgType(t *testing.T) {
	msg := MsgCreateProposal{}
	assert.Equal(t, msg.Type(), "asset")
}

func TestAnswerProposalMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgAnswerProposal
	}{
		{false, MsgAnswerProposal{AssetID: "1"}},
		{false, MsgAnswerProposal{Recipient: addrs[0]}},
		{false, MsgAnswerProposal{AssetID: "1", Recipient: addrs[0], Response: 3}},
		{true, MsgAnswerProposal{AssetID: "1", Recipient: addrs[0], Response: 1}},
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

func TestAnswerProposalMsgGet(t *testing.T) {
	msg := MsgAnswerProposal{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestAnswerProposalGetSigners(t *testing.T) {
	msg := MsgAnswerProposal{
		Recipient: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6100]`)
}

func TestAnswerProposalMsgGetSignBytes(t *testing.T) {
	msg := MsgAnswerProposal{
		AssetID:   "1",
		Recipient: addrs[0],
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"recipient\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq4gr5na\",\"response\":\"0\"}")
}
