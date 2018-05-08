package asset

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Asset struct {
	ID         string
	Name       string
	Issuer     sdk.Address
	Quantity   int64
	Attributes map[string]interface{}
}

// IsOwner ....
func (a Asset) IsOwner(addr sdk.Address) bool {
	return hex.EncodeToString(a.Issuer) == hex.EncodeToString(addr)
}
