package warranty

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Contract
type Contract struct {
	ID        string      `json:"id"`
	Issuer    sdk.Address `json:"issuer"`
	Recipient sdk.Address `json:"recipient"`
	AssetID   string      `json:"asset_id"` // the id of asset
	Serial    string      `json:"serial"`   // the serial of asset
	Expires   time.Time   `json:"expires"`
	Claim     *Claim      // the claim of contract
}

// CanCreateClaim
func (c Contract) ValidateCreateClaim(addr sdk.Address) (valid bool) {
	valid = false

	if c.Claim != nil && c.Claim.Status == ClaimStatusPending {
		return
	}

	addrsAccepted := map[string]bool{
		c.Issuer.String():    true,
		c.Recipient.String(): true,
	}
	if addrsAccepted[addr.String()] {
		valid = true
		return
	}
	return
}

// ValidateClaimProcess ...
func (c Contract) ValidateClaimProcess(addr sdk.Address, status ClaimStatus) (valid bool) {
	valid = false
	if c.Claim == nil {
		return
	}

	addrsAccepted := map[string]bool{
		c.Issuer.String():          true,
		c.Claim.Recipient.String(): true,
	}

	if addrsAccepted[addr.String()] == false {
		return
	}

	if c.Claim.Status != ClaimStatusPending {
		return
	}

	switch status {
	case ClaimStatusClaimRepair, ClaimStatusRejected:
		valid = true
		return
	default:
		return
	}
}

// Claim the claim of the contract
type Claim struct {
	Status    ClaimStatus
	Recipient sdk.Address
}

// ClaimStatus status of a claim
type ClaimStatus int

const (
	// The claim is pending
	ClaimStatusPending ClaimStatus = iota
	// The claim has been rejected
	ClaimStatusRejected
	// The item is up for repair
	ClaimStatusClaimRepair
)
