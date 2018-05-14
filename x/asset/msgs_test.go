package asset

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	attr = Attribute{
		Name:        "weight",
		Type:        3,
		NumberValue: 100,
	}
)

// ----------------------------------------
// Register Tests

func TestRegisterMsg(t *testing.T) {}

func TestCreateAssetMsgType(t *testing.T) {
	addr := sdk.Address([]byte("input"))

	var msg = RegisterMsg{
		ID:       "1",
		Name:     "asset name",
		Company:  "company name",
		Email:    "email",
		Quantity: 100,
		Issuer:   addr,
	}

	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestCreateAssetMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    RegisterMsg
	}{
		{false, RegisterMsg{}},                                                  // no asset info
		{false, RegisterMsg{Issuer: addr1, Quantity: 0, Name: "name", ID: "1"}}, // missing quantity
		{false, RegisterMsg{Issuer: addr1, Quantity: 1, Name: "name"}},          // missing id
		{true, RegisterMsg{Issuer: addr1, Quantity: 1, Name: "name", ID: "1"}},  //
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
	var msg = RegisterMsg{
		Issuer: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestRegisterGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = RegisterMsg{
		Issuer:   addr1,
		ID:       "1212",
		Name:     "name",
		Quantity: 1,
		Company:  "1",
		Email:    "1@gmail.com",
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"issuer":"696E707574","id":"1212","name":"name","quantity":1,"company":"1","email":"1@gmail.com"}`)
}

func TestRegisterGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = RegisterMsg{
		Issuer:   addr1,
		ID:       "1212",
		Name:     "name",
		Quantity: 1,
		Company:  "1",
		Email:    "1@gmail.com",
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

	var msg = UpdateAttrMsg{
		Issuer: addr,
		Attributes: []Attribute{
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
		tx    UpdateAttrMsg
	}{
		{false, UpdateAttrMsg{}},              // no asset info
		{false, UpdateAttrMsg{Issuer: addr1}}, // only set owner
		{false, UpdateAttrMsg{Issuer: addr1, Attributes: []Attribute{
			attr,
		}}}, // missing id
		{true, UpdateAttrMsg{Issuer: addr1, Attributes: []Attribute{
			attr,
		}, ID: "1212"}},
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
	var msg = UpdateAttrMsg{
		Issuer: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestUpdateAttrMsgGetSignBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = UpdateAttrMsg{
		Issuer: addr1,
		ID:     "1",
		Attributes: []Attribute{
			attr,
		},
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"issuer\":\"696E707574\",\"id\":\"1\",\"attribute\":[{\"name\":\"weight\",\"type\":3,\"bytes_value\":null,\"string_value\":\"\",\"boolean_value\":false,\"number_value\":100,\"enum_value\":null,\"location_value\":{\"latitude\":0,\"longitude\":0}}]}")
}

func TestUpdateAttrGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = UpdateAttrMsg{
		Issuer: addr1,
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
	var msg = AddQuantityMsg{
		Issuer:   addr,
		ID:       "!",
		Quantity: 1,
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestAddQuantityMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    AddQuantityMsg
	}{
		{false, AddQuantityMsg{}},                                    // no asset info
		{false, AddQuantityMsg{Issuer: addr1, Quantity: 0, ID: "1"}}, // missing quantity
		{false, AddQuantityMsg{Issuer: addr1, Quantity: 1}},          // missing id
		{true, AddQuantityMsg{Issuer: addr1, Quantity: 1, ID: "1"}},  //
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
	var msg = AddQuantityMsg{
		Issuer:   addr1,
		ID:       "1",
		Quantity: 1,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"issuer":"696E707574","id":"1","quantity":1}`)
}

func TestAddQuantityGetGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = AddQuantityMsg{
		Issuer: addr1,
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
	var msg = SubtractQuantityMsg{
		Issuer:   addr,
		ID:       "!",
		Quantity: 1,
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "asset")
}

func TestSubtractQuantityMsgValidation(t *testing.T) {
	addr1 := sdk.Address([]byte{1, 2})
	cases := []struct {
		valid bool
		tx    SubtractQuantityMsg
	}{
		{false, SubtractQuantityMsg{}},                                    // no asset info
		{false, SubtractQuantityMsg{Issuer: addr1, Quantity: 0, ID: "1"}}, // missing quantity
		{false, SubtractQuantityMsg{Issuer: addr1, Quantity: 1}},          // missing id
		{true, SubtractQuantityMsg{Issuer: addr1, Quantity: 1, ID: "1"}},  //
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
	var msg = SubtractQuantityMsg{
		Issuer:   addr1,
		ID:       "1",
		Quantity: 1,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"issuer":"696E707574","id":"1","quantity":1}`)
}

func TestSubtractQuantityMsgGetSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = AddQuantityMsg{
		Issuer: addr1,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// CreateProposal  Tests
// ------------------------------------------------------------
func TestCreateProposalMsgType(t *testing.T) {
	msg := CreateProposalMsg{}
	assert.Equal(t, msg.Type(), "asset")
}

func TestCreateProposalMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    CreateProposalMsg
	}{
		{false, CreateProposalMsg{AssetID: "1"}},
		{false, CreateProposalMsg{AssetID: "1", Issuer: addrs[0]}},
		{false, CreateProposalMsg{AssetID: "1", Recipient: addrs[0]}},
		{false, CreateProposalMsg{AssetID: "1", Propertipes: []string{"location"}}},
		{false, CreateProposalMsg{AssetID: "1", Issuer: addrs[0], Recipient: addrs[1]}},
		{true, CreateProposalMsg{AssetID: "1", Issuer: addrs[0], Recipient: addrs[1], Propertipes: []string{"location"}}},
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
	msg := CreateProposalMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCreateProposalMsgGetSigners(t *testing.T) {
	msg := CreateProposalMsg{
		Issuer: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestCreateProposalMsgGetSignBytes(t *testing.T) {
	msg := CreateProposalMsg{
		AssetID:     "1",
		Propertipes: []string{"location"},
		Issuer:      addrs[0],
		Recipient:   addrs[1],
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), `{"asset_id":"1","issuer":"A58856F0FD53BF058B4909A21AEC019107BA6160","recipient":"A58856F0FD53BF058B4909A21AEC019107BA6161","propertipes":["location"]}`)
}

// AnswerProposal  Tests
// ------------------------------------------------------------

func TestAnswerProposalMsgType(t *testing.T) {
	msg := AnswerProposalMsg{}
	assert.Equal(t, msg.Type(), "asset")
}

func TestAnswerProposalMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    AnswerProposalMsg
	}{
		{false, AnswerProposalMsg{AssetID: "1"}},
		{false, AnswerProposalMsg{Recipient: addrs[0]}},
		{false, AnswerProposalMsg{AssetID: "1", Recipient: addrs[0], Response: 3}},
		{true, AnswerProposalMsg{AssetID: "1", Recipient: addrs[0], Response: 1}},
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
	msg := AnswerProposalMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestAnswerProposalGetSigners(t *testing.T) {
	msg := AnswerProposalMsg{
		Recipient: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestAnswerProposalMsgGetSignBytes(t *testing.T) {
	msg := AnswerProposalMsg{
		AssetID:   "1",
		Recipient: addrs[0],
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"recipient\":\"A58856F0FD53BF058B4909A21AEC019107BA6160\",\"response\":0}")
}

// RevokeProposal Test
// -----------------------------------------------------
func TestRevokeProposalMsgType(t *testing.T) {
	msg := RevokeProposalMsg{}
	assert.Equal(t, msg.Type(), "asset")
}

func TestRevokeProposalMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    RevokeProposalMsg
	}{
		{false, RevokeProposalMsg{AssetID: "1"}},
		{false, RevokeProposalMsg{AssetID: "1", Recipient: addrs[1]}},
		{false, RevokeProposalMsg{AssetID: "1", Issuer: addrs[1]}},
		{false, RevokeProposalMsg{AssetID: "1", Issuer: addrs[1], Recipient: addrs[1]}},
		{false, RevokeProposalMsg{AssetID: "1", Issuer: addrs[1], Recipient: addrs[1], Propertipes: []string{}}},
		{true, RevokeProposalMsg{AssetID: "1", Recipient: addrs[0], Issuer: addrs[1], Propertipes: []string{"location"}}},
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

func TestRevokeProposalMsgGet(t *testing.T) {
	msg := RevokeProposalMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestRevokeProposalMsgGetSigners(t *testing.T) {
	msg := RevokeProposalMsg{
		Issuer: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestRevokeProposalMsgGetSignBytes(t *testing.T) {
	msg := RevokeProposalMsg{
		AssetID:     "1",
		Recipient:   addrs[0],
		Issuer:      addrs[1],
		Propertipes: []string{"location"},
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"asset_id\":\"1\",\"issuer\":\"A58856F0FD53BF058B4909A21AEC019107BA6161\",\"recipient\":\"A58856F0FD53BF058B4909A21AEC019107BA6160\",\"propertipes\":[\"location\"]}")
}
