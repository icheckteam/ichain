package asset

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset asset infomation
type Asset struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Subtype    string         `json:"subtype"`
	Name       string         `json:"name"`
	Owner      sdk.AccAddress `json:"owner"`
	Reporters  Reporters      `json:"reporters"`
	Parent     string         `json:"parent"` // the id of the asset parent
	Root       string         `json:"root"`   // the id of the asset root
	Final      bool           `json:"final"`
	Properties Properties     `json:"properties"`
	Materials  Materials      `json:"materials"`
	Quantity   sdk.Int        `json:"quantity"`
	Unit       string         `json:"unit"`
	Created    int64          `json:"created"`
	Height     int64          `json:"height"`
}

// IsOwner check is owner of the asset
func (a Asset) IsOwner(addr sdk.AccAddress) bool {
	return bytes.Equal(a.Owner, addr)
}

// CheckUpdateAttributeAuthorization returns whether the address is authorized to update the attribute
func (a Asset) CheckUpdateAttributeAuthorization(address sdk.AccAddress, prop Property) bool {
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
func (a Asset) GetReporter(address sdk.AccAddress) (*Reporter, int) {
	// Check if the address exist in the asset's reporters
	// then check if the reporter's properties includes the attribute name
	for index, reporter := range a.Reporters {
		if bytes.Equal(reporter.Addr, address) {
			return &reporter, index
		}
	}
	return nil, -1
}

func (a Asset) ValidateUpdateProperty(sender sdk.AccAddress, name string) sdk.Error {
	if a.Final {
		return ErrAssetAlreadyFinal(a.ID)
	}
	authorized := a.CheckUpdateAttributeAuthorization(sender, Property{Name: name})
	if !authorized {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized", sender))
	}
	return nil
}

// ValidateAddQuantity return error if invalid
func (a Asset) ValidateAddQuantity(sender sdk.AccAddress) sdk.Error {
	if len(a.Parent) != 0 {
		return ErrInvalidAssetRoot(a.ID)
	}
	return a.ValidateUpdateProperty(sender, "quantity")
}

func (a Asset) ValidateSubtractQuantity(sender sdk.AccAddress, quantity sdk.Int) sdk.Error {
	if a.Quantity.LT(quantity) {
		return ErrInvalidAssetQuantity(a.ID)
	}
	return a.ValidateUpdateProperty(sender, "quantity")
}

func (a Asset) ValidateFinalize(sender sdk.AccAddress) sdk.Error {
	return a.ValidateUpdateProperty(sender, "final")
}

func (a Asset) ValidateAddMaterial(sender sdk.AccAddress) sdk.Error {
	return a.ValidateUpdateProperty(sender, "materials")
}

func (a Asset) ValidateAddChildren(sender sdk.AccAddress, quantity sdk.Int) sdk.Error {
	if a.Final {
		return ErrAssetAlreadyFinal(a.ID)
	}
	if !a.IsOwner(sender) {
		return sdk.ErrUnauthorized(fmt.Sprintf("address {%v} not unauthorized to create asset", sender))
	}
	if a.Quantity.LT(quantity) {
		return ErrInvalidAssetQuantity(a.ID)
	}
	return nil
}

func (a Asset) ValidateUpdateProperties(sender sdk.AccAddress, properties Properties) sdk.Error {
	if a.Final {
		return ErrAssetAlreadyFinal(a.ID)
	}

	// check role permissions
	for _, attr := range properties {
		authorized := a.CheckUpdateAttributeAuthorization(sender, attr)
		if !authorized {
			return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to update", sender))
		}
	}

	return nil
}
