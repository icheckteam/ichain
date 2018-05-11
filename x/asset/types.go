package asset

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Assets ...
type Asset struct {
	ID         string
	Name       string
	Issuer     sdk.Address
	Quantity   int64
	Attributes map[string]interface{}
	Company    string
	Email      string
}

// IsOwner ....
func (a Asset) IsOwner(addr sdk.Address) bool {
	return hex.EncodeToString(a.Issuer) == hex.EncodeToString(addr)
}
