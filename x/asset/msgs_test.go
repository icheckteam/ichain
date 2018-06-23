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

func TestRegisterMsg(t *testing.T) {}

func TestCreateAssetMsgType(t *testing.T) {
	addr := sdk.Address([]byte("input"))

	var msg = MsgCreateAsset{
		AssetID:  "1",
		Name:     "asset name",
		Quantity: 100,
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
		{false, MsgCreateAsset{}},                                                       // no asset info
		{false, MsgCreateAsset{Sender: addr1, Quantity: 0, Name: "name", AssetID: "1"}}, // missing quantity
		{false, MsgCreateAsset{Sender: addr1, Quantity: 1, Name: "name"}},               // missing id
		{true, MsgCreateAsset{Sender: addr1, Quantity: 1, Name: "name", AssetID: "1"}},  //
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
		Quantity: 1,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"sender\":\"696E707574\",\"asset_id\":\"1212\",\"asset_type\":\"\",\"name\":\"name\",\"quantity\":1,\"parent\":\"\",\"materials\":null,\"properties\":null,\"precision\":0}")
}

func TestRegisterGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateAsset{
		Sender:   addr1,
		AssetID:  "1212",
		Name:     "name",
		Quantity: 1,
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
	assert.Equal(t, string(res), "{\"sender\":\"696E707574\",\"asset_id\":\"1\",\"properties\":[{\"name\":\"weight\",\"type\":3,\"bytes_value\":null,\"string_value\":\"\",\"boolean_value\":false,\"number_value\":100,\"enum_value\":null,\"location_value\":{\"latitude\":\"\",\"longitude\":\"\"},\"precision\":0}]}")
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
		Quantity: 1,
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
		{false, MsgAddQuantity{}},                                         // no asset info
		{false, MsgAddQuantity{Sender: addr1, Quantity: 0, AssetID: "1"}}, // missing quantity
		{false, MsgAddQuantity{Sender: addr1, Quantity: 1}},               // missing id
		{true, MsgAddQuantity{Sender: addr1, Quantity: 1, AssetID: "1"}},  //
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
		Quantity: 1,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"sender\":\"696E707574\",\"asset_id\":\"1\",\"quantity\":1}")
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
		Quantity: 1,
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
		{false, MsgSubtractQuantity{}},                                         // no asset info
		{false, MsgSubtractQuantity{Sender: addr1, Quantity: 0, AssetID: "1"}}, // missing quantity
		{false, MsgSubtractQuantity{Sender: addr1, Quantity: 1}},               // missing id
		{true, MsgSubtractQuantity{Sender: addr1, Quantity: 1, AssetID: "1"}},  //
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
		Quantity: 1,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"sender":"696E707574","asset_id":"1","quantity":1}`)
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
