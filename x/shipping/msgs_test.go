package shipping

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	transportedAssets = []TransportedAsset{TransportedAsset{"tomato"}}
)

// ----------------------------------------
// Create Order Tests

func TestCreateOrderMsgType(t *testing.T) {
	msg := CreateOrderMsg{}
	assert.Equal(t, msg.Type(), "shipping")
}

func TestCreateOrderMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    CreateOrderMsg
	}{
		{false, CreateOrderMsg{ID: "1"}},
		{false, CreateOrderMsg{TransportedAssets: transportedAssets}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0]}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[1]}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[0], Receiver: addrs[1]}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[1], Receiver: addrs[0]}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[1], Carrier: addrs[0], Receiver: addrs[1]}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[0], Receiver: addrs[0]}},
		{false, CreateOrderMsg{ID: "1", TransportedAssets: []TransportedAsset{}, Issuer: addrs[0], Carrier: addrs[1], Receiver: addrs[2]}},
		{true, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[1], Receiver: addrs[2]}},
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

func TestCreateOrderMsgGet(t *testing.T) {
	msg := CreateOrderMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCreateOrderMsgGetSigners(t *testing.T) {
	msg := CreateOrderMsg{
		Issuer: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestCreateOrderMsgGetSignBytes(t *testing.T) {
	msg := CreateOrderMsg{
		ID:                "1",
		Issuer:            addrs[0],
		Carrier:           addrs[1],
		Receiver:          addrs[2],
		TransportedAssets: transportedAssets,
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"id\":\"1\",\"transported_assets\":[{\"id\":\"tomato\"}],\"issuer\":\"A58856F0FD53BF058B4909A21AEC019107BA6160\",\"carrier\":\"A58856F0FD53BF058B4909A21AEC019107BA6161\",\"receiver\":\"A58856F0FD53BF058B4909A21AEC019107BA6162\"}")
}

// ----------------------------------------
// Confirm Order Tests

func TestConfirmOrderMsgType(t *testing.T) {
	msg := ConfirmOrderMsg{}
	assert.Equal(t, msg.Type(), "shipping")
}

func TestConfirmOrderMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    ConfirmOrderMsg
	}{
		{false, ConfirmOrderMsg{OrderID: "1"}},
		{false, ConfirmOrderMsg{Carrier: addrs[3]}},
		{false, ConfirmOrderMsg{OrderID: "1", Carrier: sdk.Address([]byte(""))}},
		{true, ConfirmOrderMsg{OrderID: "1", Carrier: addrs[3]}},
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

func TestConfirmOrderMsgGet(t *testing.T) {
	msg := ConfirmOrderMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestConfirmOrderMsgGetSigners(t *testing.T) {
	msg := ConfirmOrderMsg{
		Carrier: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestConfirmOrderMsgGetSignBytes(t *testing.T) {
	msg := ConfirmOrderMsg{
		OrderID: "1",
		Carrier: addrs[3],
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"order_id\":\"1\",\"carrier_id\":\"A58856F0FD53BF058B4909A21AEC019107BA6163\"}")
}

// ----------------------------------------
// Complete Order Tests

func TestCompleteOrderMsgType(t *testing.T) {
	msg := CompleteOrderMsg{}
	assert.Equal(t, msg.Type(), "shipping")
}

func TestCompleteOrderMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    CompleteOrderMsg
	}{
		{false, CompleteOrderMsg{OrderID: "1"}},
		{false, CompleteOrderMsg{Receiver: addrs[3]}},
		{false, CompleteOrderMsg{OrderID: "1", Receiver: sdk.Address([]byte(""))}},
		{true, CompleteOrderMsg{OrderID: "1", Receiver: addrs[3]}},
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

func TestCompleteOrderMsgGet(t *testing.T) {
	msg := CompleteOrderMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCompleteOrderMsgGetSigners(t *testing.T) {
	msg := CompleteOrderMsg{
		Receiver: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestCompleteOrderMsgGetSignBytes(t *testing.T) {
	msg := CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addrs[3],
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"order_id\":\"1\",\"receiver_id\":\"A58856F0FD53BF058B4909A21AEC019107BA6163\"}")
}

// ----------------------------------------
// Cancel Order Tests

func TestCancelOrderMsgType(t *testing.T) {
	msg := CancelOrderMsg{}
	assert.Equal(t, msg.Type(), "shipping")
}

func TestCancelOrderMsgValidation(t *testing.T) {
	cases := []struct {
		valid bool
		tx    CancelOrderMsg
	}{
		{false, CancelOrderMsg{OrderID: "1"}},
		{false, CancelOrderMsg{Issuer: addrs[3]}},
		{false, CancelOrderMsg{OrderID: "1", Issuer: sdk.Address([]byte(""))}},
		{true, CancelOrderMsg{OrderID: "1", Issuer: addrs[3]}},
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

func TestCancelOrderMsgGet(t *testing.T) {
	msg := CancelOrderMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCancelOrderMsgGetSigners(t *testing.T) {
	msg := CancelOrderMsg{
		Issuer: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}

func TestCancelOrderMsgGetSignBytes(t *testing.T) {
	msg := CancelOrderMsg{
		OrderID: "2",
		Issuer:  addrs[3],
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"order_id\":\"2\",\"issuer_id\":\"A58856F0FD53BF058B4909A21AEC019107BA6163\"}")
}
