package identity

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "identity"

// ClaimIssueMsg ...
type CreateMsg struct {
	ID       string                 `json:"id"`
	Context  string                 `json:"context"`
	Content  map[string]interface{} `json:"content"`
	Metadata ClaimMetadata          `json:"metadata"`
}

// nolint ...
func (msg CreateMsg) Type() string                            { return msgType }
func (msg CreateMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CreateMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Metadata.Issuer} }
func (msg CreateMsg) String() string {
	return fmt.Sprintf("CreateMsg{Issuer: %v, Recipient: %s, ExpirationTime: %s}", msg.Metadata.Issuer, msg.Metadata.Recipient, msg.Metadata.ExpirationTime)
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg CreateMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg CreateMsg) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return sdk.ErrTxDecode("id is requried")
	}
	if len(msg.Context) == 0 {
		return sdk.ErrTxDecode("Context is requried")
	}

	if len(msg.Content) == 0 {
		return sdk.ErrTxDecode("Content is requried")
	}
	if len(msg.Metadata.Issuer) == 0 {
		return sdk.ErrTxDecode("Metadata.Issuer is requried")
	}
	if msg.Metadata.CreateTime.IsZero() {
		return sdk.ErrTxDecode("Metadata.CreateTime is requried")
	}

	if msg.Metadata.ExpirationTime.IsZero() {
		return sdk.ErrTxDecode("Metadata.ExpirationTime is requried")
	}

	if len(msg.Metadata.Recipient) == 0 {
		return sdk.ErrInvalidAddress("Metadata.Recipient is requried")
	}
	return nil
}

// RevokeMsg ...
type RevokeMsg struct {
	ClaimID string
	Sender  sdk.Address
}

// nolint ...
func (msg RevokeMsg) Type() string                            { return msgType }
func (msg RevokeMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg RevokeMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg RevokeMsg) String() string {
	return fmt.Sprintf("RevokeMsg{ClaimID: %v, Sender: %s}", msg.ClaimID, msg.Sender)
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg RevokeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg RevokeMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.ClaimID) == 0 {
		return sdk.ErrTxDecode("ClaimID is requried")
	}

	return nil
}
