package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateUpdateProperties ...
func (k Keeper) ValidateUpdateProperties(ctx sdk.Context, record Asset, sender sdk.AccAddress, properties Properties) sdk.Error {
	if record.Final {
		return ErrAssetAlreadyFinal(record.ID)
	}
	if record.IsOwner(sender) {
		return nil
	}
	reporter, found := k.GetReporter(ctx, record.ID, sender)
	if !found {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized", sender))
	}

	// check role permissions
	for _, attr := range properties {
		authorized := k.CheckUpdateAttributeAuthorization(ctx, record, reporter, attr)
		if !authorized {
			return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to update", sender))
		}
	}

	return nil
}

// CheckUpdateAttributeAuthorization returns whether the address is authorized to update the attribute
func (k Keeper) CheckUpdateAttributeAuthorization(ctx sdk.Context, record Asset, reporter Reporter, prop Property) bool {
	attributeName := prop.Name

	// Check if the address exist in the asset's reporters
	// then check if the reporter's properties includes the attribute name
	for _, property := range reporter.Properties {
		if property == attributeName {
			return true
		}
	}
	return false
}
