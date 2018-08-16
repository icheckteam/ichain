package epcis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendMsg send data message
type SendMsg struct {
	Sender    sdk.AccAddress   `json:"sender"`
	Receiver  sdk.AccAddress   `json:"receiver"`
	Actors    []sdk.AccAddress `json:"actors"`
	Products  []string         `json:"products"`
	Events    []Event          `json:"events"`
	Batches   []string         `json:"batches"`
	Locations []string         `json:"locations"`
}
