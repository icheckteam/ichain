package asset

import (
	"bytes"
	"sort"
	"strings"

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

// Property property of the asset
type Property struct {
	Name         string       `json:"name"`
	Type         PropertyType `json:"type"`
	BytesValue   []byte       `json:"bytes_value"`
	StringValue  string       `json:"string_value"`
	BooleanValue bool         `json:"boolean_value"`
	NumberValue  int64        `json:"number_value"`
	EnumValue    []string     `json:"enum_value"`
	Location     Location     `json:"location_value"`
}

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// list all properties
type Properties []Property

func (propertiesA Properties) Adds(othersB ...Property) Properties {
	sum := ([]Property)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(propertiesA), len(othersB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, othersB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, propertiesA[indexA:]...)
		}
		propertyA, propertyB := propertiesA[indexA], othersB[indexB]
		switch strings.Compare(propertyA.Name, propertyB.Name) {
		case -1:
			sum = append(sum, propertyA)
			indexA++
		case 0:
			sum = append(sum, propertyB)
			indexA++
			indexB++
		case 1:
			indexB++
			sum = append(sum, propertyB)
		}
	}
}

//----------------------------------------
// Sort interface

//nolint
func (properties Properties) Len() int           { return len(properties) }
func (properties Properties) Less(i, j int) bool { return properties[i].Name < properties[j].Name }
func (properties Properties) Swap(i, j int) {
	properties[i], properties[j] = properties[j], properties[i]
}

var _ sort.Interface = Properties{}

// Sort is a helper function to sort the set of materials inplace
func (properties Properties) Sort() Properties {
	sort.Sort(properties)
	return properties
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

// PropertyType define the type of the property
type PropertyType int

// All avaliable type ò the attribute
const (
	PropertyTypeBytes PropertyType = iota + 1
	PropertyTypeString
	PropertyTypeBoolean
	PropertyTypeNumber
	PropertyTypeEnum
	PropertyTypeLocation
)

// Material defines the total material of new asset
type Material struct {
	AssetID  string `json:"asset_id"`
	Quantity int64  `json:"quantity"`
}

// Materials - list of materials
type Materials []Material

// SameDenomAs returns true if the two assets are the same asset
func (material Material) SameAssetAs(other Material) bool {
	return (material.AssetID == other.AssetID)
}

// Adds quantities of two assets with same asset
func (material Material) Plus(materialB Material) Material {
	if !material.SameAssetAs(materialB) {
		return material
	}
	return Material{material.AssetID, material.Quantity + materialB.Quantity}
}

// Plus combines two sets of materials
// CONTRACT: Plus will never return materials where one Material has a 0 quantity.
func (materials Materials) Plus(materialsB Materials) Materials {
	sum := ([]Material)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(materials), len(materialsB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, materialsB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, materials[indexA:]...)
		}
		materialA, materialB := materials[indexA], materialsB[indexB]
		switch strings.Compare(materialA.AssetID, materialB.AssetID) {
		case -1:
			sum = append(sum, materialA)
			indexA++
		case 0:
			if materialA.Quantity+materialB.Quantity == 0 {
				// ignore 0 sum coin type
			} else {
				sum = append(sum, materialA.Plus(materialB))
			}
			indexA++
			indexB++
		case 1:
			sum = append(sum, materialB)
			indexB++
		}
	}
}

//----------------------------------------
// Sort interface

//nolint
func (materials Materials) Len() int           { return len(materials) }
func (materials Materials) Less(i, j int) bool { return materials[i].AssetID < materials[j].AssetID }
func (materials Materials) Swap(i, j int)      { materials[i], materials[j] = materials[j], materials[i] }

var _ sort.Interface = Materials{}

// Sort is a helper function to sort the set of materials inplace
func (materials Materials) Sort() Materials {
	sort.Sort(materials)
	return materials
}
