package asset

import sdk "github.com/cosmos/cosmos-sdk/types"

//--------------------------------------------------

// Proposal is an invitation to manage an asset
type Proposal struct {
	Role       ProposalRole   `json:"role"`       // The role assigned to the recipient
	Status     ProposalStatus `json:"status"`     // The response of the recipient
	Properties []string       `json:"properties"` // The asset's attributes name that the recipient is authorized to update
	Issuer     sdk.Address    `json:"issuer"`     // The proposal issuer
	Recipient  sdk.Address    `json:"recipient"`  // The recipient of the proposal
}

// Proposals is a sclice of Proposal
type Proposals []Proposal

// ProposalRole defines the authority of the proposal's recipient
type ProposalRole int

const (
	// RoleReporter is authorized to update the asset's attributes
	// whose name is included in the proposal's properties field
	RoleReporter ProposalRole = iota + 1

	// RoleOwner has the same authorization as RoleReporter
	// but also authorized to make proposal to other recipient
	RoleOwner
)

// ProposalStatus define the status of the proposal
type ProposalStatus int

// All available status of the proposal
const (
	StatusPending  ProposalStatus = iota // The recipient has not answered
	StatusAccepted                       // The recipient accepted the proposal
	StatusRefused                        // The recipient refused the proposal
)
