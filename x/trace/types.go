package trace

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "trace"

// MsgRecordCreate A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type MsgRecordCreate struct {
	Sender     sdk.Address
	RecordID   string
	RecordName string
}

// NewMsgRecordCreate new record create msg
func NewMsgRecordCreate(sender sdk.Address, recordID, recordName string) MsgRecordCreate {
	return MsgRecordCreate{
		Sender:     sender,
		RecordID:   recordID,
		RecordName: recordName,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = MsgRecordCreate{}

// nolint
func (msg MsgRecordCreate) Type() string                            { return msgType }
func (msg MsgRecordCreate) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgRecordCreate) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgRecordCreate) String() string {
	return fmt.Sprintf("MsgRecordCreate{Sender: %v, RecordID: %s, RecordName: %s}", msg.Sender, msg.RecordID, msg.RecordName)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgRecordCreate) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgRecordCreate) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// MsgChangeRecordOwner A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type MsgChangeRecordOwner struct {
	Sender   sdk.Address
	To       sdk.Address
	RecordID string
}

// nolint
func (msg MsgChangeRecordOwner) Type() string                            { return msgType }
func (msg MsgChangeRecordOwner) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgChangeRecordOwner) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg MsgChangeRecordOwner) String() string {
	return fmt.Sprintf("MsgChangeRecordOwner{Sender: %v, To: %s, RecordID: %s}", msg.Sender, msg.To, msg.RecordID)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgChangeRecordOwner) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.To) == 0 {
		return sdk.ErrUnknownAddress(msg.To.String()).Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgChangeRecordOwner) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
