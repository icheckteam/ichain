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
	Unit     string         `json:"unit"`
	Created  int64          `json:"created"`
	Height   int64          `json:"height"`
}

// RecordOutput ...
type RecordOutput struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Owner      sdk.AccAddress `json:"owner"`
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
