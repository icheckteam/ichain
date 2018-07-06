package insurance

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	addr  = sdk.Address([]byte("input"))
	addr1 = sdk.Address([]byte("input"))
)

// ------------------------------------------------------------
// CreateContract Tests
func TestMsgCreateContractType(t *testing.T) {

	var msg = MsgCreateContract{
		ID:        "1",
		Issuer:    addr,
		Recipient: addr1,
		Serial:    "1",
		Expires:   time.Now(),
	}

	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "insurance")
}

func TestMsgCreateContractValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgCreateContract
	}{
		{false, MsgCreateContract{}},
		{false, MsgCreateContract{ID: "1"}},
		{false, MsgCreateContract{ID: "1", AssetID: "1"}},
		{false, MsgCreateContract{ID: "1", AssetID: "1", Issuer: addr}},
		{false, MsgCreateContract{ID: "1", AssetID: "1", Issuer: addr, Recipient: addr1}},
		{false, MsgCreateContract{ID: "1", AssetID: "1", Issuer: addr, Recipient: addr1, Serial: "1"}},
		{true, MsgCreateContract{ID: "1", AssetID: "1", Issuer: addr, Recipient: addr1, Serial: "1", Expires: time.Now()}},
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

func TestCreateContractMsgGet(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateContract{
		Issuer: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestMsgCreateContractnBytes(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	expiration, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	var msg = MsgCreateContract{
		ID:        "1",
		Issuer:    addr,
		Recipient: addr1,
		Serial:    "1",
		Expires:   expiration,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"id\":\"1\",\"issuer\":\"696E707574\",\"recipient\":\"696E707574\",\"expires\":\"2018-05-11T16:28:45.78807557+07:00\",\"serial\":\"1\",\"asset_id\":\"\"}")
}

func TestMsgCreateContractSigners(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateContract{
		ID:        "1",
		Issuer:    addr,
		Recipient: addr1,
		Serial:    "1",
		Expires:   time.Now(),
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// CreateClaim Tests

func TestMsgCreateClaimType(t *testing.T) {

	var msg = MsgCreateClaim{
		ContractID: "1",
		Issuer:     addr,
	}

	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "insurance")
}

func TestMsgCreateClaimValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgCreateClaim
	}{
		{false, MsgCreateClaim{}},
		{false, MsgCreateClaim{ContractID: "1"}},
		{false, MsgCreateClaim{Issuer: addr}},
		{false, MsgCreateClaim{Issuer: addr, Recipient: addr1}},
		{true, MsgCreateClaim{ContractID: "1", Issuer: addr, Recipient: addr1}},
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

func TestMsgCreateClaimGet(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgCreateClaim{
		Issuer: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestMsgCreateClaimGetBytes(t *testing.T) {
	var msg = MsgCreateClaim{
		ContractID: "1",
		Issuer:     addr,
		Recipient:  addr1,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"contract_id\":\"1\",\"issuer\":\"696E707574\",\"recipient\":\"696E707574\"}")
}

func TestMsgCreateClaimSigners(t *testing.T) {
	var msg = MsgCreateClaim{
		Issuer: addr,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}

// ------------------------------------------------------------
// ProcessClaim Tests
func TestMsgProcessClaimType(t *testing.T) {

	var msg = MsgProcessClaim{
		ContractID: "1",
		Issuer:     addr,
	}

	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "insurance")
}

func TestMsgProcessClaimValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgProcessClaim
	}{
		{false, MsgProcessClaim{}},
		{false, MsgProcessClaim{ContractID: "1"}},
		{false, MsgProcessClaim{Issuer: addr}},
		{true, MsgProcessClaim{ContractID: "1", Issuer: addr}},
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

func TestMsgCompleteClaimGet(t *testing.T) {
	addr1 := sdk.Address([]byte("input"))
	var msg = MsgProcessClaim{
		Issuer: addr1,
	}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestMsgCompleteClaimGetBytes(t *testing.T) {
	var msg = MsgProcessClaim{
		ContractID: "1",
		Issuer:     addr,
	}
	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"contract_id\":\"1\",\"issuer\":\"696E707574\",\"status\":0}")
}

func TestMsgCompleteClaimSigners(t *testing.T) {
	var msg = MsgCreateClaim{
		Issuer: addr,
	}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}
