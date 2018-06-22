package identity

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "identity"

// MsgCreateClaim ...
type MsgCreateClaim struct {
	ID       string    `json:"id"`
	Context  string    `json:"context"`
	Content  Content   `json:"content"`
	Metadata Metadata  `json:"metadata"`
	Fee      sdk.Coins `json:"fee"`
}

// nolint ...
func (msg MsgCreateClaim) Type() string                            { return msgType }
func (msg MsgCreateClaim) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateClaim) GetSigners() []sdk.Address               { return []sdk.Address{msg.Metadata.Issuer} }
func (msg MsgCreateClaim) String() string {
	return fmt.Sprintf("MsgCreateClaim{Issuer: %v, Recipient: %s, Expires: %s}", msg.Metadata.Issuer, msg.Metadata.Recipient, msg.Metadata.Expires)
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgCreateClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgCreateClaim) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return sdk.ErrTxDecode("id is requried")
	}
	if len(msg.Context) == 0 {
		return sdk.ErrTxDecode("Context is requried")
	}

	if msg.Content == nil {
		return sdk.ErrTxDecode("Content is requried")
	}
	if len(msg.Metadata.Issuer) == 0 {
		return sdk.ErrTxDecode("Metadata.Issuer is requried")
	}
	if msg.Metadata.CreateTime.IsZero() {
		return sdk.ErrTxDecode("Metadata.CreateTime is requried")
	}

	if msg.Metadata.Expires.IsZero() {
		return sdk.ErrTxDecode("Metadata.expires is requried")
	}

	if len(msg.Metadata.Recipient) == 0 {
		return sdk.ErrInvalidAddress("Metadata.Recipient is requried")
	}
	return nil
}

// RevokeMsg ...
type MsgRevokeClaim struct {
	ClaimID    string      `json:"claim_id"`
	Sender     sdk.Address `json:"sender"`
	Revocation string      `json:"revocation"`
}

// nolint ...
func (msg MsgRevokeClaim) Type() string                            { return msgType }
func (msg MsgRevokeClaim) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgRevokeClaim) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgRevokeClaim) String() string {
	return fmt.Sprintf("MsgRevokeClaim{ID: %v, Sender: %s}", msg.ClaimID, msg.Sender)
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgRevokeClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgRevokeClaim) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.ClaimID) == 0 {
		return sdk.ErrTxDecode("ClaimID is requried")
	}
	if len(msg.Revocation) == 0 {
		return sdk.ErrTxDecode("Revocation is requried")
	}
	return nil
}

// RevokeMsg ...
type MsgAnswerClaim struct {
	ClaimID  string      `json:"claim_id"`
	Sender   sdk.Address `json:"sender"`
	Response int         `json:"response"`
}

// nolint ...
func (msg MsgAnswerClaim) Type() string                            { return msgType }
func (msg MsgAnswerClaim) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgAnswerClaim) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgAnswerClaim) String() string {
	return fmt.Sprintf("MsgAnswerClaim{ID: %v, Sender: %s}", msg.ClaimID, msg.Sender)
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgAnswerClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgAnswerClaim) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.ClaimID) == 0 {
		return sdk.ErrTxDecode("ClaimID is requried")
	}
	return nil
}
