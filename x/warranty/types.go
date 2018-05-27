package warranty

import (
	"bytes"
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

	if c.Claim != nil {
		switch c.Claim.Status {
		case ClaimStatusPending, ClaimStatusTheftConfirmed:
			return
		default:
		}
	}

	if !bytes.Equal(addr, c.Recipient) {
		return
	}
	valid = true
	return
}

// ValidateClaimProcess ...
func (c Contract) ValidateClaimProcess(addr sdk.Address, status ClaimStatus) (valid bool) {
	valid = false
	if c.Claim == nil {
		return
	}

	if !bytes.Equal(addr, c.Claim.Recipient) {
		return
	}

	if c.Claim.Status != ClaimStatusPending {
		return
	}

	switch status {
	case ClaimStatusPending:
		return
	default:
		valid = true
		return
	}
}

// Claim the claim of the contract
type Claim struct {
	Status    ClaimStatus
	Recipient sdk.Address // warranty address
}

// ClaimStatus status of a claim
type ClaimStatus int

const (
	// The claim is pending
	ClaimStatusPending ClaimStatus = iota
	// The claim has been rejected
	ClaimStatusRejected
	// The item is up for repair or had been repaired
	ClaimStatusClaimRepair
	// The customer should be reimbursed
	ClaimStatusReimbursement
	// The theft of the item has been confirmed by authorities
	ClaimStatusTheftConfirmed
)
