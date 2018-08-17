package gs1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
