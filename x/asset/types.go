package asset

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Asset asset infomation
type Asset struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	Parent   string         `json:"parent"` // the id of the asset parent
	Root     string         `json:"root"`   // the id of the asset root
	Final    bool           `json:"final"`
	Quantity sdk.Int        `json:"quantity"`
	Created  int64          `json:"created"`
	Height   int64          `json:"height"`
}

// RecordOutput ...
type RecordOutput struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Owner      sdk.AccAddress `json:"owner"`
	Type       string         `json:"type"`
	SubType    string         `json:"subtype"`
	Barcode    string         `json:"barcode"`
	Parent     string         `json:"parent"` // the id of the asset parent
	Root       string         `json:"root"`   // the id of the asset root
	Final      bool           `json:"final"`
	Quantity   sdk.Int        `json:"quantity"`
	Unit       string         `json:"unit"`
	Created    int64          `json:"created"`
	Height     int64          `json:"height"`
	Materials  []Material     `json:"materials"`
	Reporters  []Reporter     `json:"reporters"`
	Properties Properties     `json:"properties"`
}

// HistoryTransferOutput ...
type HistoryTransferOutput struct {
	Owner sdk.AccAddress `json:"recipient"`
	Time  int64          `json:"time"`
	Memo  string         `json:"memo"`
}

// HistoryChangeQuantityOutput ...
type HistoryChangeQuantityOutput struct {
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Int        `json:"amount"`
	Type   string         `json:"type"`
	Time   int64          `json:"time"`
	Memo   string         `json:"memo"`
}

// HistoryUpdateProperty ...
type HistoryUpdateProperty struct {
	Reporter sdk.AccAddress `json:"reporter"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Value    interface{}    `json:"value"`
	Time     int64          `json:"time"`
	Memo     string         `json:"memo"`
}

// HistoryTransferMaterial ...
type HistoryTransferMaterial struct {
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Int        `json:"amount"`
	From   string         `json:"from"`
	To     string         `json:"to"`
	Time   int64          `json:"time"`
	Memo   string         `json:"memo"`
}

// ProposalOutput ...
type ProposalOutput struct {
	Role       ProposalRole   `json:"role"`       // The role assigned to the recipient
	Status     ProposalStatus `json:"status"`     // The response of the recipient
	Properties []string       `json:"properties"` // The asset's attributes name that the recipient is authorized to update
	Issuer     sdk.AccAddress `json:"issuer"`     // The proposal issuer
	Recipient  sdk.AccAddress `json:"recipient"`  // The recipient of the proposal
	AssetID    string         `json:"asset_id"`   // The id of the asset
}

// ToProposalOutput ...
func ToProposalOutput(proposal Proposal, assetID string) ProposalOutput {
	return ProposalOutput{
		Role:       proposal.Role,
		Status:     proposal.Status,
		Properties: proposal.Properties,
		Issuer:     proposal.Issuer,
		Recipient:  proposal.Recipient,
		AssetID:    assetID,
	}
}

// IsOwner check is owner of the asset
func (a Asset) IsOwner(addr sdk.AccAddress) bool {
	return bytes.Equal(a.Owner, addr)
}

// UnmarshalReporter ...
func UnmarshalReporter(cdc *wire.Codec, value []byte) (reporter Reporter, err error) {
	err = cdc.UnmarshalBinary(value, &reporter)
	return
}

// UnmarshalProperty ...
func UnmarshalProperty(cdc *wire.Codec, value []byte) (property Property, err error) {
	err = cdc.UnmarshalBinary(value, &property)
	return
}

// UnmarshalMaterial ...
func UnmarshalMaterial(cdc *wire.Codec, value []byte) (material Material, err error) {
	err = cdc.UnmarshalBinary(value, &material)
	return
}

// UnmarshalRecord ...
func UnmarshalRecord(cdc *wire.Codec, value []byte) (record Asset, err error) {
	err = cdc.UnmarshalBinary(value, &record)
	return
}

// UnmarshalProposal ...
func UnmarshalProposal(cdc *wire.Codec, value []byte) (proposal Proposal, err error) {
	err = cdc.UnmarshalBinary(value, &proposal)
	return
}
