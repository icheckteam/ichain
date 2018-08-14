package epcis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendMsg send data message
type SendMsg struct {
	Sender     sdk.AccAddress `json:"sender"`     // the sender address
	Receiver   sdk.AccAddress `json:"receiver"`   // the receiver address
	Properties []string       `json:"properties"` // list all properties
}
