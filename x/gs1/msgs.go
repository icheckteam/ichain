package gs1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "gs1"

var _ sdk.Msg = &MsgSend{}

// MsgSend send data message
type MsgSend struct {
	Sender    sdk.AccAddress `json:"sender"`
	Receiver  sdk.AccAddress `json:"receiver"`
	Actors    []Actor        `json:"actors"`
	Products  []Product      `json:"products"`
	Events    []Event        `json:"events"`
	Batches   []Batch        `json:"batches"`
	Locations []Location     `json:"locations"`
}

// Type ...
func (msg MsgSend) Type() string { return msgType }

// GetSigners ...
func (msg MsgSend) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Sender} }

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgSend) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg MsgSend) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
