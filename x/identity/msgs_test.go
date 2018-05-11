package identity

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// CreateMsg tests
func TestCreateMsg(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr2 := sdk.Address([]byte("input2"))
	var msg = CreateMsg{
		Context: "claim:context",
		Content: map[string]interface{}{
			"name": "claim name",
		},
		Metadata: ClaimMetadata{
			CreateTime:     time.Now(),
			ExpirationTime: time.Now().Add(time.Hour * 100000),
			Issuer:         addr,
			Recipient:      addr2,
		},
	}
	// TODO some failures for bad result
	assert.Equal(t, msg.Type(), "identity")
}

func TestCreateMsgValidation(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr2 := sdk.Address([]byte("input2"))

	cases := []struct {
		valid bool
		tx    CreateMsg
	}{
		{false, CreateMsg{}}, // no id
		{false, CreateMsg{
			ID: "1",
		}}, // no context
		{false, CreateMsg{
			ID:      "1",
			Context: "1212",
		}}, // no content
		{false, CreateMsg{
			ID:      "1",
			Context: "1212",
			Content: map[string]interface{}{"content": "1"},
		}}, // no meta

		{false, CreateMsg{
			ID:      "1",
			Context: "1",
			Content: map[string]interface{}{"content": "1"},
			Metadata: ClaimMetadata{
				Recipient: addr,
			}}}, // no issuer
		{false, CreateMsg{
			ID:      "1",
			Context: "1",
			Content: map[string]interface{}{"content": "1"},
			Metadata: ClaimMetadata{
				Issuer: addr2,
			}}}, // no recipient
		{false, CreateMsg{
			ID:      "1",
			Context: "1",
			Content: map[string]interface{}{"content": "1"},
			Metadata: ClaimMetadata{
				Issuer:     addr2,
				Recipient:  addr,
				CreateTime: time.Now(),
			}}}, // no expires
		{false, CreateMsg{
			ID:      "1",
			Context: "1",
			Content: map[string]interface{}{"content": "1"},
			Metadata: ClaimMetadata{
				Issuer:         addr2,
				Recipient:      addr,
				ExpirationTime: time.Now(),
			}}}, // no CreateTime

		{true, CreateMsg{
			ID:      "1",
			Context: "1",
			Content: map[string]interface{}{"content": "1"},
			Metadata: ClaimMetadata{
				Recipient:      addr,
				Issuer:         addr2,
				CreateTime:     time.Now(),
				ExpirationTime: time.Now().Add(time.Hour * 100000),
			}}},
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

func TestCreateMsgGet(t *testing.T) {
	var msg = CreateMsg{}
	res := msg.Get(nil)
	assert.Nil(t, res)
}

func TestCreateMsgGetSignBytes(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr1 := sdk.Address([]byte("input1"))

	creatTime, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	expiration, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")

	var msg = CreateMsg{
		ID:      "1",
		Context: "1",
		Content: map[string]interface{}{"content": "1"},
		Metadata: ClaimMetadata{
			Recipient:      addr,
			Issuer:         addr1,
			CreateTime:     creatTime,
			ExpirationTime: expiration,
		}}

	res := msg.GetSignBytes()
	// TODO bad results
	assert.Equal(t, string(res), "{\"id\":\"1\",\"context\":\"1\",\"content\":{\"content\":\"1\"},\"metadata\":{\"create_time\":\"2018-05-11T16:28:45.78807557+07:00\",\"issuer\":\"696E70757431\",\"recipient\":\"696E707574\",\"expiration_time\":\"2018-05-11T16:28:45.78807557+07:00\",\"revocation\":\"\"}}")
}

func TestCreateMsgGetSigners(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr1 := sdk.Address([]byte("input1"))

	creatTime, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	expiration, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")

	var msg = CreateMsg{
		ID:      "1",
		Context: "1",
		Content: map[string]interface{}{"content": "1"},
		Metadata: ClaimMetadata{
			Recipient:      addr1,
			Issuer:         addr,
			CreateTime:     creatTime,
			ExpirationTime: expiration,
		}}
	res := msg.GetSigners()
	// TODO bad results
	assert.Equal(t, fmt.Sprintf("%v", res), `[696E707574]`)
}
