package asset

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

// RegisterMsg A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type RegisterMsg struct {
	Issuer   sdk.Address `json:"issuer"`
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Quantity int64       `json:"quantity"`

	// Company info
	Company string `json:"company"`
	Email   string `json:"email"`
}

// NewRegisterMsg new record create msg
func NewRegisterMsg(issuer sdk.Address, id, name string, quantity int64, company, email string) RegisterMsg {
	return RegisterMsg{
		Issuer:   issuer,
		ID:       id,
		Quantity: quantity,
		Name:     name,
		Company:  company,
		Email:    email,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = RegisterMsg{}

// nolint ...
func (msg RegisterMsg) Type() string                            { return msgType }
func (msg RegisterMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg RegisterMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{%s->%s->%d}", msg.Issuer, msg.Name, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}

	if len(msg.ID) == 0 {
		return ErrMissingField("asset_id")
	}

	if len(msg.Name) == 0 {
		return ErrMissingField("asset_name")
	}

	if msg.Quantity == 0 {
		return ErrMissingField("asset_quantity")
	}

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg RegisterMsg) GetSignBytes() []byte {
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
