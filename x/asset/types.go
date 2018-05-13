package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Assets ...
type Asset struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Issuer     sdk.Address `json:"issuer"`
	Quantity   int64       `json:"quantity"`
	Attributes []Attribute `json:"attributes"`
	Company    string      `json:"company"`
	Email      string      `json:"email"`
}

// IsOwner ....
func (a Asset) IsOwner(addr sdk.Address) bool {
	return a.Issuer.String() == addr.String()
}

// Attribute ...
type Attribute struct {
	Name         string   `json:"name"`
	Type         int      `json:"type"`
	BytesValue   []byte   `json:"bytes_value"`
	StringValue  string   `json:"string_value"`
	BooleanValue bool     `json:"boolean_value"`
	NumberValue  int64    `json:"number_value"`
	EnumValue    []string `json:"enum_value"`
	Location     Location `json:"location_value"`
}

type Location struct {
	Latitude  float64 `json:"latitude" amino:"unsafe"`
	Longitude float64 `json:"longitude" amino:"unsafe"`
}
