package idetify

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim struct
type Claim struct {
	ID       string                 `json:"id"`
	Context  string                 `json:"context"`
	Content  map[string]interface{} `json:"content"`
	Metadata Metadata               `json:"metadata"`
}

// Metadata the claim metadata
type Metadata struct {
	CreateTime     time.Time   `json:"create_time"`
	Issuer         sdk.Address `json:"issuer"`
	Recipient      sdk.Address `json:"recipient"`
	ExpirationTime time.Time   `json:"expiration_time"`
	Revocation     string      `json:"revocation"`
}
