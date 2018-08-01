package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

var _, _, _ sdk.Msg = &MsgCreateAsset{}, &MsgAddMaterials{}, &MsgAddQuantity{}
var _, _, _ sdk.Msg = &MsgCreateProposal{}, &MsgRevokeReporter{}, &MsgFinalize{}
var _, _, _ sdk.Msg = &MsgSubtractQuantity{}, &MsgAnswerProposal{}, &MsgUpdateProperties{}

// MsgCreateAsset A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type MsgCreateAsset struct {
	Sender     sdk.AccAddress `json:"sender"`
	AssetID    string         `json:"asset_id"`
	Name       string         `json:"name"`
	Quantity   sdk.Int        `json:"quantity"`
	Parent     string         `json:"parent"` // the id of the  parent asset
	Unit       string         `json:"unit"`
	Properties Properties     `json:"properties"`
}

// NewMsgCreateAsset new record create msg
func NewMsgCreateAsset(sender sdk.AccAddress, id, name string, quantity sdk.Int, parent string) MsgCreateAsset {
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
func (msg MsgCreateAsset) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgUpdateProperties ...
// ---------------------------------------------------------------
type MsgUpdateProperties struct {
	Sender     sdk.AccAddress `json:"sender"`
	AssetID    string         `json:"asset_id"`
	Properties Properties     `json:"properties"`
}

func (msg MsgUpdateProperties) Type() string                            { return msgType }
func (msg MsgUpdateProperties) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgUpdateProperties) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgUpdateProperties) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Properties) == 0 {
		return ErrMissingField("properties")
	}

	if err := msg.Properties.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgUpdateProperties) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgAddQuantity ...
// ---------------------------------------------------------------
type MsgAddQuantity struct {
	Sender   sdk.AccAddress `json:"sender"`
	AssetID  string         `json:"asset_id"`
	Quantity sdk.Int        `json:"quantity"`
}

func (msg MsgAddQuantity) Type() string                            { return msgType }
func (msg MsgAddQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAddQuantity) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgSubtractQuantity ...
// ---------------------------------------------------------------
type MsgSubtractQuantity struct {
	Sender   sdk.AccAddress `json:"sender"`
	AssetID  string         `json:"asset_id"`
	Quantity sdk.Int        `json:"quantity"`
}

func (msg MsgSubtractQuantity) Type() string                            { return msgType }
func (msg MsgSubtractQuantity) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgSubtractQuantity) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgAddMaterials ...
type MsgAddMaterials struct {
	AssetID string         `json:"asset_id"`
	Sender  sdk.AccAddress `json:"sender"`
	Amount  Materials      `json:"amount"`
}

func (msg MsgAddMaterials) Type() string                            { return msgType }
func (msg MsgAddMaterials) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAddMaterials) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAddMaterials) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Amount) == 0 {
		return ErrMissingField("amount")
	}
	if err := msg.Amount.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgAddMaterials) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgSend ...
type MsgFinalize struct {
	Sender  sdk.AccAddress `json:"sender"`
	AssetID string         `json:"asset_id"`
}

func (msg MsgFinalize) Type() string                            { return msgType }
func (msg MsgFinalize) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgFinalize) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgRevokeReporter ...
type MsgRevokeReporter struct {
	Sender   sdk.AccAddress `json:"sender"`
	Reporter sdk.AccAddress `json:"reporter"`
	AssetID  string         `json:"asset_id"`
}

func (msg MsgRevokeReporter) Type() string                            { return msgType }
func (msg MsgRevokeReporter) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgRevokeReporter) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// CreateProposalMsg ...
type MsgCreateProposal struct {
	AssetID    string         `json:"asset_id"`
	Sender     sdk.AccAddress `json:"sender"`
	Recipient  sdk.AccAddress `json:"recipient"`
	Properties []string       `json:"properties"`
	Role       ProposalRole   `json:"role"`
}

func (msg MsgCreateProposal) Type() string                            { return msgType }
func (msg MsgCreateProposal) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateProposal) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Sender} }

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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgAnswerProposal ...
type MsgAnswerProposal struct {
	AssetID   string         `json:"asset_id"`
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	Response  ProposalStatus `json:"response"`
	Role      ProposalRole   `json:"role"`
}

func (msg MsgAnswerProposal) Type() string                            { return msgType }
func (msg MsgAnswerProposal) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAnswerProposal) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Recipient} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAnswerProposal) ValidateBasic() sdk.Error {
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	if len(msg.Sender) == 0 {
		return ErrMissingField("sender")
	}
	if msg.Role == 0 {
		return ErrMissingField("role")
	}
	switch msg.Response {
	case StatusAccepted, StatusCancel, StatusRejected:
		break
	default:
		return ErrInvalidField("response")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgAnswerProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
