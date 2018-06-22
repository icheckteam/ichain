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
	Issuer     sdk.Address `json:"issuer"`
	AssetID    string      `json:"asset_id"`
	Name       string      `json:"name"`
	Quantity   int64       `json:"quantity"`
	Parent     string      `json:"parent"` // the id of the  parent asset
	Materials  Materials   `json:"materials"`
	Properties Properties  `json:"properties"`
	Precision  int         `json:"precision"`
}

// NewMsgCreateAsset new record create msg
func NewMsgCreateAsset(issuer sdk.Address, id, name string, quantity int64, parent string) MsgCreateAsset {
	return MsgCreateAsset{
		Issuer:   issuer,
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
func (msg MsgCreateAsset) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg MsgCreateAsset) String() string {
	return fmt.Sprintf("MsgCreateAsset{%s}", msg.Name)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgCreateAsset) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
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
	Issuer     sdk.Address `json:"issuer"`
	AssetID    string      `json:"asset_id"`
	Properties Properties  `json:"properties"`
}

func (msg MsgUpdateProperties) Type() string                            { return msgType }
func (msg MsgUpdateProperties) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgUpdateProperties) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg MsgUpdateProperties) String() string {
	return fmt.Sprintf("MsgUpdateProperties{%s->%v}", msg.AssetID, msg.Properties)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgUpdateProperties) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrInvalidAddress(msg.Issuer.String()).Trace("")
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

// AddQuantityMsg ...
// ---------------------------------------------------------------
type AddQuantityMsg struct {
	Issuer   sdk.Address `json:"issuer"`
	AssetID  string      `json:"asset_id"`
	Quantity int64       `json:"quantity"`
}

func (msg AddQuantityMsg) Type() string                            { return msgType }
func (msg AddQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AddQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg AddQuantityMsg) String() string {
	return fmt.Sprintf("AddQuantityMsg{Sender: %v, quantity: %v}", msg.Issuer, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AddQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrInvalidAddress(msg.Issuer.String())
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
func (msg AddQuantityMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgSubtractQuantity ...
// ---------------------------------------------------------------
type MsgSubtractQuantity struct {
	Issuer   sdk.Address `json:"issuer"`
	AssetID  string      `json:"asset_id"`
	Quantity int64       `json:"quantity"`
}

func (msg MsgSubtractQuantity) Type() string                            { return msgType }
func (msg MsgSubtractQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgSubtractQuantity) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg MsgSubtractQuantity) String() string {
	return fmt.Sprintf("MsgSubtractQuantity{Issuer: %v, quantity: %v}", msg.Issuer, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgSubtractQuantity) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrInvalidAddress(msg.Issuer.String())
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

// CreateProposalMsg ...
type CreateProposalMsg struct {
	AssetID    string       `json:"asset_id"`
	Issuer     sdk.Address  `json:"issuer"`
	Recipient  sdk.Address  `json:"recipient"`
	Properties []string     `json:"properties"`
	Role       ProposalRole `json:"role"`
}

func (msg CreateProposalMsg) Type() string                            { return msgType }
func (msg CreateProposalMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CreateProposalMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg CreateProposalMsg) String() string {
	return fmt.Sprintf(`
		CreateProposalMsg{
			AssetID: %s, 
			Issuer: %v,
			Recipient:%v,
			Properties:%v
		}	
	`, msg.AssetID, msg.Issuer, msg.Recipient, msg.Properties)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg CreateProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrInvalidAddress(msg.Issuer.String())
	}
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	if len(msg.Properties) == 0 {
		return ErrMissingField("properties")
	}
	switch msg.Role {
	case 0, 1:
		break
	default:
		return ErrInvalidField("role")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg CreateProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// AnswerProposalMsg ...
type AnswerProposalMsg struct {
	AssetID   string         `json:"asset_id"`
	Recipient sdk.Address    `json:"recipient"`
	Response  ProposalStatus `json:"response"`
}

func (msg AnswerProposalMsg) Type() string                            { return msgType }
func (msg AnswerProposalMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AnswerProposalMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Recipient} }
func (msg AnswerProposalMsg) String() string {
	return fmt.Sprintf(`
		AnswerProposalMsg{
			AssetID: %s, 
			Recipient: %v,
			Response:%v,
		}	
	`, msg.AssetID, msg.Recipient, msg.Response)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AnswerProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Recipient) == 0 {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}
	if msg.Response > 2 {
		return ErrMissingField("response")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg AnswerProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// RevokeProposalMsg ...
type RevokeProposalMsg struct {
	AssetID    string      `json:"asset_id"`
	Issuer     sdk.Address `json:"issuer"`
	Recipient  sdk.Address `json:"recipient"`
	Properties []string    `json:"properties"`
}

func (msg RevokeProposalMsg) Type() string                            { return msgType }
func (msg RevokeProposalMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg RevokeProposalMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg RevokeProposalMsg) String() string {
	return fmt.Sprintf(`
		RevokeProposalMsg{
			AssetID: %s, 
			Issuer: %s,
			Recipient: %v,
			Properties:%v,
		}	
	`, msg.AssetID, msg.Issuer, msg.Recipient, msg.Properties)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg RevokeProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrInvalidAddress(msg.Issuer.String())
	}
	if len(msg.Properties) == 0 {
		return ErrMissingField("properties")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg RevokeProposalMsg) GetSignBytes() []byte {
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

// MsgSend ...
type MsgSend struct {
	Sender    sdk.Address `json:"sender"`
	Recipient sdk.Address `json:"recipient"`
	Assets    []string    `json:"assets"`
}

func (msg MsgSend) Type() string                            { return msgType }
func (msg MsgSend) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgSend) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgSend) String() string {
	return fmt.Sprintf(`MsgSend{%s->%s->%v}`, msg.Sender, msg.Recipient, msg.Assets)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgSend) ValidateBasic() sdk.Error {
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
func (msg MsgSend) GetSignBytes() []byte {
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
