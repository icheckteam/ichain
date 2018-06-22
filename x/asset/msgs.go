package asset

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

// MsgCreateAsset A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type MsgCreateAsset struct {
	Sender     sdk.Address `json:"sender"`
	AssetID    string      `json:"asset_id"`
	Name       string      `json:"name"`
	Quantity   int64       `json:"quantity"`
	Parent     string      `json:"parent"` // the id of the  parent asset
	Materials  Materials   `json:"materials"`
	Properties Properties  `json:"properties"`
	Precision  int         `json:"precision"`
}

// NewMsgCreateAsset new record create msg
func NewMsgCreateAsset(sender sdk.Address, id, name string, quantity int64, parent string) MsgCreateAsset {
	return MsgCreateAsset{
		Sender:   sender,
		AssetID:  id,
		Quantity: quantity,
		Name:     name,
		Parent:   parent,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = MsgCreateAsset{}

// nolint ...
func (msg MsgCreateAsset) Type() string                            { return msgType }
func (msg MsgCreateAsset) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateAsset) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgCreateAsset) String() string {
	return fmt.Sprintf("MsgCreateAsset{%s}", msg.Name)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgCreateAsset) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if len(msg.AssetID) == 0 {
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
func (msg MsgCreateAsset) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgUpdateProperties ...
// ---------------------------------------------------------------
type MsgUpdateProperties struct {
	Sender     sdk.Address `json:"sender"`
	AssetID    string      `json:"asset_id"`
	Properties Properties  `json:"properties"`
}

func (msg MsgUpdateProperties) Type() string                            { return msgType }
func (msg MsgUpdateProperties) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgUpdateProperties) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgUpdateProperties) String() string {
	return fmt.Sprintf("MsgUpdateProperties{%s->%v}", msg.AssetID, msg.Properties)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgUpdateProperties) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Properties) == 0 {
		return ErrMissingField("name")
	}
	for _, attr := range msg.Properties {
		switch attr.Type {
		case PropertyTypeBoolean,
			PropertyTypeBytes,
			PropertyTypeEnum,
			PropertyTypeLocation,
			PropertyTypeNumber,
			PropertyTypeString:
			break
		default:
			return ErrInvalidField("properties")
		}
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgUpdateProperties) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgAddQuantity ...
// ---------------------------------------------------------------
type MsgAddQuantity struct {
	Sender   sdk.Address `json:"sender"`
	AssetID  string      `json:"asset_id"`
	Quantity int64       `json:"quantity"`
}

func (msg MsgAddQuantity) Type() string                            { return msgType }
func (msg MsgAddQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAddQuantity) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgAddQuantity) String() string {
	return fmt.Sprintf("MsgAddQuantity{Sender: %v, quantity: %v}", msg.Sender, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAddQuantity) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if msg.Quantity <= 0 {
		return ErrMissingField("quantity")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgAddQuantity) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgSubtractQuantity ...
// ---------------------------------------------------------------
type MsgSubtractQuantity struct {
	Sender   sdk.Address `json:"sender"`
	AssetID  string      `json:"asset_id"`
	Quantity int64       `json:"quantity"`
}

func (msg MsgSubtractQuantity) Type() string                            { return msgType }
func (msg MsgSubtractQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgSubtractQuantity) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgSubtractQuantity) String() string {
	return fmt.Sprintf("MsgSubtractQuantity{Issuer: %v, quantity: %v}", msg.Sender, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgSubtractQuantity) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if msg.Quantity <= 0 {
		return ErrMissingField("quantity")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgSubtractQuantity) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgAddMaterials ...
type MsgAddMaterials struct {
	AssetID   string      `json:"asset_id"`
	Sender    sdk.Address `json:"sender"`
	Materials Materials   `json:"materials"`
}

func (msg MsgAddMaterials) Type() string                            { return msgType }
func (msg MsgAddMaterials) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAddMaterials) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgAddMaterials) String() string {
	return fmt.Sprintf(`MsgAddMaterials{%v->%s->%v}`, msg.Sender, msg.AssetID, msg.Materials)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAddMaterials) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Materials) == 0 {
		return ErrMissingField("asset_id")
	}
	for i, material := range msg.Materials {
		if len(material.AssetID) == 0 {
			return ErrMissingField(fmt.Sprintf("materials[%d].asset_id is required", i))
		}
		if material.Quantity == 0 {
			return ErrMissingField(fmt.Sprintf("materials[%d].quantity is required", i))
		}
	}

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgAddMaterials) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgTransfer ...
type MsgTransfer struct {
	Sender    sdk.Address `json:"sender"`
	Recipient sdk.Address `json:"recipient"`
	Assets    []string    `json:"assets"`
}

func (msg MsgTransfer) Type() string                            { return msgType }
func (msg MsgTransfer) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgTransfer) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgTransfer) String() string {
	return fmt.Sprintf(`MsgTransfer{%s->%s->%v}`, msg.Sender, msg.Recipient, msg.Assets)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.Recipient) == 0 {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}
	if len(msg.Assets) == 0 {
		return ErrMissingField("assets")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgTransfer) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgSend ...
type MsgFinalize struct {
	Sender  sdk.Address `json:"sender"`
	AssetID string      `json:"asset_id"`
}

func (msg MsgFinalize) Type() string                            { return msgType }
func (msg MsgFinalize) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgFinalize) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgFinalize) String() string {
	return fmt.Sprintf(`MsgFinalize{%s->%s}`, msg.Sender, msg.AssetID)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgFinalize) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgFinalize) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgSend ...
type MsgCreateReporter struct {
	Sender     sdk.Address `json:"sender"`
	Reporter   sdk.Address `json:"reporter"`
	AssetID    string      `json:"asset_id"`
	Properties []string    `json:"properties"`
}

func (msg MsgCreateReporter) Type() string                            { return msgType }
func (msg MsgCreateReporter) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateReporter) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgCreateReporter) String() string {
	return fmt.Sprintf(`MsgCreateReporter{%s->%s}`, msg.Sender, msg.Reporter)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgCreateReporter) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.Reporter) == 0 {
		return sdk.ErrInvalidAddress(msg.Reporter.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Properties) == 0 {
		return ErrMissingField("properties")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgCreateReporter) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgRevokeReporter ...
type MsgRevokeReporter struct {
	Sender   sdk.Address `json:"sender"`
	Reporter sdk.Address `json:"reporter"`
	AssetID  string      `json:"asset_id"`
}

func (msg MsgRevokeReporter) Type() string                            { return msgType }
func (msg MsgRevokeReporter) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgRevokeReporter) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgRevokeReporter) String() string {
	return fmt.Sprintf(`MsgRevokeReporter{%s->%s}`, msg.Sender, msg.Reporter)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgRevokeReporter) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.Reporter) == 0 {
		return sdk.ErrInvalidAddress(msg.Reporter.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgRevokeReporter) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
