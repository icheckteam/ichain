package asset

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset asset infomation
type Asset struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Subtype    string      `json:"subtype"`
	Name       string      `json:"name"`
	Owner      sdk.Address `json:"owner"`
	Reporters  Reporters   `json:"reporters"`
	Parent     string      `json:"parent"` // the id of the asset parent
	Root       string      `json:"root"`   // the id of the asset root
	Final      bool        `json:"final"`
	Properties Properties  `json:"properties"`
	Materials  Materials   `json:"materials"`
	Quantity   sdk.Int     `json:"quantity"`
	Unit       string      `json:"unit"`
	Created    int64       `json:"created"`
	Height     int64       `json:"height"`
}

// IsOwner check is owner of the asset
func (a Asset) IsOwner(addr sdk.Address) bool {
	return bytes.Equal(a.Owner, addr)
}

// CheckUpdateAttributeAuthorization returns whether the address is authorized to update the attribute
func (a Asset) CheckUpdateAttributeAuthorization(address sdk.Address, prop Property) bool {
	if a.IsOwner(address) {
		return true
	}

	attributeName := prop.Name

	// Check if the address exist in the asset's reporters
	// then check if the reporter's properties includes the attribute name
	for _, reporter := range a.Reporters {
		if bytes.Equal(reporter.Addr, address) {
			for _, property := range reporter.Properties {
				if property == attributeName {
					return true
				}
			}
		}
	}
	return false
}

// CheckUpdateAttributeAuthorization returns whether the address is authorized to update the attribute
func (a Asset) GetReporter(address sdk.Address) (*Reporter, int) {
	// Check if the address exist in the asset's reporters
	// then check if the reporter's properties includes the attribute name
	for index, reporter := range a.Reporters {
		if bytes.Equal(reporter.Addr, address) {
			return &reporter, index
		}
	}
	return nil, -1
}
