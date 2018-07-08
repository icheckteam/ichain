package asset

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

var _, _, _ sdk.Msg = &MsgCreateAsset{}, &MsgAddMaterials{}, &MsgAddQuantity{}
var _, _, _ sdk.Msg = &MsgCreateProposal{}, &MsgRevokeReporter{}, &MsgFinalize{}
var _, _, _ sdk.Msg = &MsgSubtractQuantity{}, &MsgAnswerProposal{}, &MsgUpdateProperties{}

// MsgCreateAsset A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type MsgCreateAsset struct {
	Sender    sdk.Address `json:"sender"`
	AssetID   string      `json:"asset_id"`
	AssetType string      `json:"asset_type"`
	Name      string      `json:"name"`
	Quantity  sdk.Int     `json:"quantity"`
	Parent    string      `json:"parent"` // the id of the  parent asset
	Unit      string      `json:"unit"`
}

// NewMsgCreateAsset new record create msg
func NewMsgCreateAsset(sender sdk.Address, id, name string, quantity sdk.Int, parent string) MsgCreateAsset {
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
		return ErrMissingField("name")
	}

	if msg.Quantity.IsZero() {
		return ErrMissingField("quantity")
	}

	if msg.Unit == "" {
		return ErrMissingField("unit")
	}

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgCreateAsset) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(struct {
		Sender    string  `json:"sender"`
		AssetID   string  `json:"asset_id"`
		AssetType string  `json:"asset_type"`
		Name      string  `json:"name"`
		Quantity  sdk.Int `json:"quantity"`
		Parent    string  `json:"parent"`
		Unit      string  `json:"unit"`
	}{
		Sender:    sdk.MustBech32ifyAcc(msg.Sender),
		AssetID:   msg.AssetID,
		AssetType: msg.AssetType,
		Name:      msg.Name,
		Quantity:  msg.Quantity,
		Parent:    msg.Parent,
		Unit:      msg.Unit,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
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

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgUpdateProperties) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
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
	Quantity sdk.Int     `json:"quantity"`
}

func (msg MsgAddQuantity) Type() string                            { return msgType }
func (msg MsgAddQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAddQuantity) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAddQuantity) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if msg.Quantity.IsZero() {
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
	Quantity sdk.Int     `json:"quantity"`
}

func (msg MsgSubtractQuantity) Type() string                            { return msgType }
func (msg MsgSubtractQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgSubtractQuantity) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgSubtractQuantity) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if msg.Quantity.IsZero() {
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
		if material.Quantity.IsZero() {
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
type MsgFinalize struct {
	Sender  sdk.Address `json:"sender"`
	AssetID string      `json:"asset_id"`
}

func (msg MsgFinalize) Type() string                            { return msgType }
func (msg MsgFinalize) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgFinalize) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }

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

// MsgRevokeReporter ...
type MsgRevokeReporter struct {
	Sender   sdk.Address `json:"sender"`
	Reporter sdk.Address `json:"reporter"`
	AssetID  string      `json:"asset_id"`
}

func (msg MsgRevokeReporter) Type() string                            { return msgType }
func (msg MsgRevokeReporter) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgRevokeReporter) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }

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

// CreateProposalMsg ...
type MsgCreateProposal struct {
	AssetID    string       `json:"asset_id"`
	Sender     sdk.Address  `json:"sender"`
	Recipient  sdk.Address  `json:"recipient"`
	Properties []string     `json:"properties"`
	Role       ProposalRole `json:"role"`
}

func (msg MsgCreateProposal) Type() string                            { return msgType }
func (msg MsgCreateProposal) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateProposal) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgCreateProposal) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrMissingField("sender")
	}
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	switch msg.Role {
	case 1, 2:
		break
	default:
		return ErrInvalidField("role")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgCreateProposal) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgAnswerProposal ...
type MsgAnswerProposal struct {
	AssetID   string         `json:"asset_id"`
	Recipient sdk.Address    `json:"recipient"`
	Response  ProposalStatus `json:"response"`
}

func (msg MsgAnswerProposal) Type() string                            { return msgType }
func (msg MsgAnswerProposal) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAnswerProposal) GetSigners() []sdk.Address               { return []sdk.Address{msg.Recipient} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAnswerProposal) ValidateBasic() sdk.Error {
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	if msg.Response > 2 {
		return ErrMissingField("response")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgAnswerProposal) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
