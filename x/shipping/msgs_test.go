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
		name  string
		valid bool
		tx    CreateOrderMsg
	}{
		{"only 1", false, CreateOrderMsg{ID: "1"}},
		{"only 2", false, CreateOrderMsg{TransportedAssets: transportedAssets}},
		{"only 3", false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0]}},
		{"only 4", false, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[1]}},
		{"only 9", false, CreateOrderMsg{ID: "1", TransportedAssets: []TransportedAsset{}, Issuer: addrs[0], Carrier: addrs[1], Receiver: addrs[2]}},
		{"only 10", true, CreateOrderMsg{ID: "1", TransportedAssets: transportedAssets, Issuer: addrs[0], Carrier: addrs[1], Receiver: addrs[2]}},
	}

	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%s : %d: %+v", tc.name, i, err)
		} else {
			assert.NotNil(t, err, "%s : %d", tc.name, i)
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
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6100]`)
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
	assert.Equal(t, string(res), "{\"type\":\"shipping/CreateOrderMsg\",\"value\":{\"carrier\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgpg7hpw0\",\"id\":\"1\",\"issuer\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq4gr5na\",\"receiver\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzxdzhqs\",\"transported_assets\":[{\"id\":\"tomato\"}]}}")
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
		{false, ConfirmOrderMsg{OrderID: "1", Carrier: sdk.AccAddress([]byte(""))}},
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

func TestConfirmOrderMsgGetSigners(t *testing.T) {
	msg := ConfirmOrderMsg{
		Carrier: addrs[0],
	}
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6100]`)
}

func TestConfirmOrderMsgGetSignBytes(t *testing.T) {
	msg := ConfirmOrderMsg{
		OrderID: "1",
		Carrier: addrs[3],
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"type\":\"shipping/ConfirmOrderMsg\",\"value\":{\"carrier_id\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrmmkzaz\",\"order_id\":\"1\"}}")
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
		{false, CompleteOrderMsg{OrderID: "1", Receiver: sdk.AccAddress([]byte(""))}},
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
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6100]`)
}

func TestCompleteOrderMsgGetSignBytes(t *testing.T) {
	msg := CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addrs[3],
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"type\":\"shipping/CompleteOrderMsg\",\"value\":{\"order_id\":\"1\",\"receiver_id\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrmmkzaz\"}}")
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
		{false, CancelOrderMsg{OrderID: "1", Issuer: sdk.AccAddress([]byte(""))}},
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
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6100]`)
}

func TestCancelOrderMsgGetSignBytes(t *testing.T) {
	msg := CancelOrderMsg{
		OrderID: "2",
		Issuer:  addrs[3],
	}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"type\":\"shipping/CancelOrderMsg\",\"value\":{\"issuer_id\":\"cosmosaccaddr15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrmmkzaz\",\"order_id\":\"2\"}}")
}
