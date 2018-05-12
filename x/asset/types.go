package asset

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Assets ...
type Asset struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Issuer     sdk.Address            `json:"issuer"`
	Quantity   int64                  `json:"quantity"`
	Attributes map[string]interface{} `json:"attributes"`
	Company    string                 `json:"company"`
	Email      string                 `json:"email"`
}

// IsOwner ....
func (a Asset) IsOwner(addr sdk.Address) bool {
	return hex.EncodeToString(a.Issuer) == hex.EncodeToString(addr)
}
