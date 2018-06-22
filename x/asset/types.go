package asset

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Asset asset infomation
type Asset struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Issuer     sdk.Address `json:"issuer"`
	Owner      sdk.Address `json:"owner"`
	Parent     string      `json:"parent"` // the id of the asset parent
	Root       string      `json:"root"`   // the id of the asset root
	Quantity   int64       `json:"quantity"`
	Company    string      `json:"company"`
	Email      string      `json:"email"`
	Final      bool        `json:"final"`
	Properties Properties  `json:"properties"`
	Proposals  Proposals   `json:"proposals"`
	Materials  Materials   `json:"materials"`
	Precision  int         `json:"precision"`
}

// IsOwner check is owner of the asset
func (a Asset) IsOwner(addr sdk.Address) bool {
	return bytes.Equal(a.Owner, addr)
}

// IsIssuer check is issuer of the asset
func (a Asset) IsIssuer(addr sdk.Address) bool {
	return bytes.Equal(a.Issuer, addr)
}

// CheckUpdateAttributeAuthorization returns whether the address is authorized to update the attribute
func (a Asset) CheckUpdateAttributeAuthorization(address sdk.Address, prop Property) bool {
	if a.IsOwner(address) {
		return true
	}

	attributeName := prop.Name

	// Check if the address exist in the asset's proposals
	// then check if the proposal's properties includes the attribute name
	for _, proposal := range a.Proposals {
		if proposal.IsAccepted() && bytes.Equal(proposal.Recipient, address) {
			for _, property := range proposal.Properties {
				if property == attributeName {
					return true
				}
			}
		}
	}
	return false
}

// ValidatePropossal returns whether the address is authorized to create a new proposal,
// optionally return a proposal and its index if a proposal for the recipient already exists
func (a Asset) ValidatePropossal(issuer sdk.Address, recipient sdk.Address) (*Proposal, int, bool) {
	var (
		proposal      *Proposal
		proposalIndex = -1
		authorized    = false
	)

	if bytes.Equal(issuer, recipient) {
		return nil, -1, false
	}

	if a.IsOwner(issuer) {
		authorized = true
	}

	for index, p := range a.Proposals {
		// Check if recipient already exists in the proposals list
		if bytes.Equal(p.Recipient, recipient) {
			proposalIndex = index
			p := p
			proposal = &p
		}

		// Skip the check for role if already authorized
		if authorized {
			// Skip the loop if an existing proposal was also found
			if proposal != nil {
				break
			}
			continue
		}

		// Check if the issuer has the correct role
		if p.Role == RoleOwner && p.IsAccepted() && bytes.Equal(p.Recipient, issuer) {
			authorized = true
		}
	}
	return proposal, proposalIndex, authorized
}

// ValidateProposalAnswer checks whether the recipient is authorized to answer,
// if authorized then returns the existing proposal and its index
func (a Asset) ValidateProposalAnswer(recipient sdk.Address, answer ProposalStatus) (proposal *Proposal, proposalIndex int, authorized bool) {
	proposalIndex = -1
	authorized = false

	// Check for invalid answer
	switch answer {
	case StatusAccepted, StatusRefused:
		break
	default:
		return
	}

	for i, p := range a.Proposals {
		// Check for proposal with the same recipient
		if bytes.Equal(recipient, p.Recipient) {

			// Only proceed if the proposal's status is pending
			switch p.Status {
			case StatusPending:
				authorized = true
				proposalIndex = i
				proposal = &p
				return
			default:
				break
			}
		}
	}
	return
}
