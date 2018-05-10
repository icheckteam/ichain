package asset

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

// AssetCreateMsg A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type AssetCreateMsg struct {
	Sender    sdk.Address
	AssetID   string
	AssetName string
	Quantity  int64

	// Company info
	Company string
	Email   string
}

// NewAssetCreateMsg new record create msg
func NewAssetCreateMsg(sender sdk.Address, assetID, assetName string, quantity int64, company, email string) AssetCreateMsg {
	return AssetCreateMsg{
		Sender:    sender,
		AssetID:   assetID,
		Quantity:  quantity,
		AssetName: assetID,
		Company:   company,
		Email:     email,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = AssetCreateMsg{}

// nolint ...
func (msg AssetCreateMsg) Type() string                            { return msgType }
func (msg AssetCreateMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AssetCreateMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg AssetCreateMsg) String() string {
	return fmt.Sprintf("AssetCreateMsg{Sender: %v, RecordID: %s, RecordName: %s}", msg.Sender, msg.AssetID, msg.AssetName)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AssetCreateMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg AssetCreateMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// UpdateAttrMsg ...
// ---------------------------------------------------------------
type UpdateAttrMsg struct {
	Sender         sdk.Address
	AssetID        string
	AttributeName  string
	AttributeValue interface{}
}

func (msg UpdateAttrMsg) Type() string                            { return msgType }
func (msg UpdateAttrMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg UpdateAttrMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg UpdateAttrMsg) String() string {
	return fmt.Sprintf("UpdateAttrMsg{Sender: %v}", msg.Sender)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg UpdateAttrMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return sdk.ErrUnknownAddress(msg.AssetID).Trace("")
	}
	if len(msg.AttributeName) == 0 {
		return sdk.ErrInternal("Attribute name is required").Trace("")
	}
	if msg.AttributeValue == 0 {
		return sdk.ErrInternal("Attribute value is required").Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg UpdateAttrMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// AddQuantityMsg ...
// ---------------------------------------------------------------
type AddQuantityMsg struct {
	Sender   sdk.Address
	AssetID  string
	Quantity int64
}

func (msg AddQuantityMsg) Type() string                            { return msgType }
func (msg AddQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AddQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg AddQuantityMsg) String() string {
	return fmt.Sprintf("AddQuantityMsg{Sender: %v, quantity: %v}", msg.Sender, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AddQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return sdk.ErrUnknownAddress(msg.AssetID).Trace("")
	}
	if msg.Quantity <= 0 {
		return sdk.ErrInternal("Quantity is required").Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg AddQuantityMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// AddQuantityMsg ...
// ---------------------------------------------------------------
type SubtractQuantityMsg struct {
	Sender   sdk.Address
	AssetID  string
	Quantity int64
}

func (msg SubtractQuantityMsg) Type() string                            { return msgType }
func (msg SubtractQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg SubtractQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg SubtractQuantityMsg) String() string {
	return fmt.Sprintf("SubtractQuantityMsg{Sender: %v, quantity: %v}", msg.Sender, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg SubtractQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return sdk.ErrUnknownAddress(msg.AssetID).Trace("")
	}
	if msg.Quantity <= 0 {
		return sdk.ErrInternal("Quantity is required").Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg SubtractQuantityMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
