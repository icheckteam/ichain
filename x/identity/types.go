package identity

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim struct
type Claim struct {
	ID       string    `json:"id"`
	Context  string    `json:"context"`
	Content  Content   `json:"content"`
	Metadata Metadata  `json:"metadata"`
	Fee      sdk.Coins `json:"fee"`
	Paid     bool      `json:"paid"`
}

// Metadata the claim metadata
type Metadata struct {
	CreateTime time.Time   `json:"create_time"`
	Issuer     sdk.Address `json:"issuer"`
	Recipient  sdk.Address `json:"recipient"`
	Expires    time.Time   `json:"expires"`
	Revocation string      `json:"revocation"`
}

func (c Claim) IsOwner(addr sdk.Address) bool {
	return bytes.Equal(c.Metadata.Issuer, addr)
}

type Content []byte
