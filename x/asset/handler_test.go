package asset

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

type msgInvalid struct{}

// Type ...
func (msg msgInvalid) Type() string { return msgType }

// GetSigners ...
func (msg msgInvalid) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg msgInvalid) ValidateBasic() sdk.Error {

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg msgInvalid) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
func TestHandler(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	result := NewHandler(keeper)(ctx, MsgCreateAsset{
		AssetID:  "asseta",
		Sender:   addr,
		Name:     "asset 1",
		Quantity: sdk.NewInt(100),
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgCreateAsset{
		AssetID:  "1",
		Sender:   addr,
		Name:     "asset 1",
		Quantity: sdk.NewInt(100),
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgAddQuantity{
		AssetID:  "asseta",
		Sender:   addr,
		Quantity: sdk.NewInt(100),
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgSubtractQuantity{
		AssetID:  "asseta",
		Sender:   addr,
		Quantity: sdk.NewInt(100),
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgAddMaterials{
		AssetID: "asseta",
		Sender:  addr,
		Amount: Materials{
			Material{RecordID: "1", Amount: sdk.NewInt(1)},
		},
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgUpdateProperties{
		AssetID: "asseta",
		Sender:  addr,
		Properties: Properties{
			Property{Name: "demo", Type: 1, StringValue: "1"},
		},
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgCreateProposal{
		AssetID:    "asseta",
		Sender:     addr,
		Recipient:  addr2,
		Role:       RoleReporter,
		Properties: []string{"size"},
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgAnswerProposal{
		AssetID:   "asseta",
		Sender:    addr2,
		Recipient: addr2,
		Role:      RoleReporter,
		Response:  StatusAccepted,
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgRevokeReporter{
		AssetID:  "asseta",
		Sender:   addr,
		Reporter: addr2,
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, MsgFinalize{
		AssetID: "asseta",
		Sender:  addr,
	})
	assert.True(t, result.IsOK() == true)

	result = NewHandler(keeper)(ctx, msgInvalid{})
	assert.True(t, result.IsOK() != true)

}
