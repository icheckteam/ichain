package asset

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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
		{true, RegisterMsg{Issuer: addr1, Quantity: 1, Name: "name", ID: "1"}},  // missing quantity
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
