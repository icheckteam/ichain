package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Assets ...
type Asset struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Issuer     sdk.Address `json:"issuer"`
	Quantity   int64       `json:"quantity"`
	Company    string      `json:"company"`
	Email      string      `json:"email"`
	Attributes []Attribute `json:"attributes"`
	Proposals  Proposals   `json:"proposals"`
}

// IsOwner ....
func (a Asset) IsOwner(addr sdk.Address) bool {
	return a.Issuer.String() == addr.String()
}

// CheckUpdateAttributeAuthorization returns whether the address is authorized to update the attribute
func (a Asset) CheckUpdateAttributeAuthorization(address sdk.Address, attr Attribute) bool {
	if a.IsOwner(address) {
		return true
	}

	attributeName := attr.Name

	// Check if the address exist in the asset's proposals
	// then check if the proposal's properties includes the attribute name
	for _, proposal := range a.Proposals {
		if proposal.IsAccepted() && proposal.Recipient.String() == address.String() {
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

	if issuer.String() == recipient.String() {
		return nil, -1, false
	}

	if a.IsOwner(issuer) {
		authorized = true
	}

	for index, p := range a.Proposals {
		// Check if recipient already exists in the proposals list
		if p.Recipient.String() == recipient.String() {
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
		if p.Role == RoleOwner && p.IsAccepted() && p.Recipient.String() == issuer.String() {
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
		if p.Recipient.String() == recipient.String() {

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

// Attribute ...
type Attribute struct {
	Name         string        `json:"name"`
	Type         AttributeType `json:"type"`
	BytesValue   []byte        `json:"bytes_value"`
	StringValue  string        `json:"string_value"`
	BooleanValue bool          `json:"boolean_value"`
	NumberValue  int64         `json:"number_value"`
	EnumValue    []string      `json:"enum_value"`
	Location     Location      `json:"location_value"`
}

type Location struct {
	Latitude  float64 `json:"latitude" amino:"unsafe"`
	Longitude float64 `json:"longitude" amino:"unsafe"`
}

//--------------------------------------------------

// Proposal is an invitation to manage an asset
type Proposal struct {
	Role       ProposalRole   `json:"role"`       // The role assigned to the recipient
	Status     ProposalStatus `json:"status"`     // The response of the recipient
	Properties []string       `json:"properties"` // The asset's attributes name that the recipient is authorized to update
	Issuer     sdk.Address    `json:"issuer"`     // The proposal issuer
	Recipient  sdk.Address    `json:"recipient"`  // The recipient of the proposal
}

// IsAccepted returns true if the proposal was accepted
func (p *Proposal) IsAccepted() bool {
	return p.Status == StatusAccepted
}

// AddProperties add new properties to the proposal, filtering existing value
func (p *Proposal) AddProperties(properties []string) {
OuterLoop:
	for _, addedProperty := range properties {
		for _, currentProperty := range p.Properties {
			if addedProperty == currentProperty {
				continue OuterLoop
			}
		}
		p.Properties = append(p.Properties, addedProperty)
	}
}

// RemoveProperties removes the listed properties from the proposal
func (p *Proposal) RemoveProperties(removedProperties []string) {
	properties := []string{}

OuterLoop:
	for _, currentProperty := range p.Properties {
		for _, removedProperty := range removedProperties {
			if removedProperty == currentProperty {
				continue OuterLoop
			}
		}
		properties = append(properties, currentProperty)
	}
	p.Properties = properties
}

// Proposals is a sclice of Proposal
type Proposals []Proposal

// ProposalRole defines the authority of the proposal's recipient
type ProposalRole int

const (
	// RoleReporter is authorized to update the asset's attributes
	// whose name is included in the proposal's properties field
	RoleReporter ProposalRole = iota

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

// AttributeType define the type ò the attribute
type AttributeType int

// All avaliable type ò the attribute
const (
	AttributeTypeBytes AttributeType = iota
	AttributeTypeString
	AttributeTypeBoolean
	AttributeTypeNumber
	AttributeTypeEnum
	AttributeTypeLocation
)
