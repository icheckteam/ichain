package invoice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgCreate_Type(t *testing.T) {
	msg := MsgCreate{}
	assert.Equal(t, msg.Type(), "invoice")
}

func TestMsgCreate_GetSignBytes(t *testing.T) {
	msg := NewMsgCreate("1", addrs[0], addrs[1], items)
	res := msg.GetSignBytes()
	assert.Equal(t, string(res), "{\"id\":\"1\",\"issuer\":\"A58856F0FD53BF058B4909A21AEC019107BA6160\",\"receiver\":\"A58856F0FD53BF058B4909A21AEC019107BA6161\",\"items\":[{\"asset_id\":\"jav\",\"quantity\":1}]}")
}

func TestMsgCreate_ValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgCreate
	}{
		{false, MsgCreate{ID: "1"}},
		{false, MsgCreate{ID: "1", Issuer: addrs[0]}},
		{false, MsgCreate{ID: "1", Issuer: addrs[0], Receiver: addrs[1]}},
		{true, MsgCreate{ID: "1", Issuer: addrs[0], Receiver: addrs[1], Items: items}},
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

func TestMsgCreate_GetSigners(t *testing.T) {
	msg := NewMsgCreate("1", addrs[0], addrs[1], items)
	res := msg.GetSigners()
	assert.Equal(t, fmt.Sprintf("%v", res), `[A58856F0FD53BF058B4909A21AEC019107BA6160]`)
}
