package rest

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/x/asset"
)

type baseBody struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainID       string `json:"chain_id"`
	Sequence      int64  `json:"sequence"`
	AccountNumber int64  `json:"account_number"`
	Gas           int64  `json:"gas"`
}

func (b baseBody) Validate() error {
	if b.Name == "" {
		return errors.New("account_name is required")
	}
	if b.Password == "" {
		return errors.New("password is required")
	}
	if b.Gas == 0 {
		return errors.New("gas is required")
	}
	if len(b.ChainID) == 0 {
		return errors.New("chain_id is required")
	}
	if b.AccountNumber == 0 {
		return errors.New("account_number is required")
	}
	return nil
}

type AssetOutput struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Height     int64            `json:"height"`
	Name       string           `json:"name"`
	Owner      string           `json:"owner"`
	Quantity   sdk.Int          `json:"quantity"`
	Reporters  []ReporterOutput `json:"reporters"`
	Parent     string           `json:"parent"` // the id of the asset parent
	Root       string           `json:"root"`   // the id of the asset root
	Final      bool             `json:"final"`
	Properties asset.Properties `json:"properties"`
	Materials  asset.Materials  `json:"materials"`
	Precision  int              `json:"precision"`
	Created    int64            `json:"created"`
}

type ReporterOutput struct {
	Addr       string   `json:"address"`
	Properties []string `json:"properties"`
	Created    int64    `json:"created"`
}

func ToAssetOutput(a asset.Asset) AssetOutput {
	reporters := []ReporterOutput{}
	for _, reporter := range a.Reporters {
		reporters = append(reporters, ReporterOutput{
			Addr:       sdk.MustBech32ifyAcc(reporter.Addr),
			Created:    reporter.Created,
			Properties: reporter.Properties,
		})
	}
	return AssetOutput{
		ID:         a.ID,
		Type:       a.Type,
		Height:     a.Height,
		Name:       a.Name,
		Quantity:   a.Quantity,
		Owner:      sdk.MustBech32ifyAcc(a.Owner),
		Reporters:  reporters,
		Parent:     a.Parent,
		Root:       a.Root,
		Final:      a.Final,
		Properties: a.Properties,
		Materials:  a.Materials,
		Created:    a.Created,
	}
}

func ToAssetsOutput(asa []asset.Asset) (asb []AssetOutput) {
	for _, a := range asa {
		asb = append(asb, ToAssetOutput(a))
	}
	return
}

type msgCreateCreateProposalBody struct {
	BaseReq baseBody `json:"base_req"`

	Recipient  string             `json:"recipient"`
	Properties []string           `json:"properties"`
	Role       asset.ProposalRole `json:"role"`
}

type msgAnswerProposalBody struct {
	BaseReq baseBody `json:"base_req"`

	Response asset.ProposalStatus `json:"response"`
	AssetID  string               `json:"asset_id"`
}

type ProposalOutput struct {
	Role       asset.ProposalRole   `json:"role"`       // The role assigned to the recipient
	Status     asset.ProposalStatus `json:"status"`     // The response of the recipient
	Properties []string             `json:"properties"` // The asset's attributes name that the recipient is authorized to update
	Issuer     string               `json:"issuer"`     // The proposal issuer
	Recipient  string               `json:"recipient"`  // The recipient of the proposal
}

func bech32ProposalOutput(proposal asset.Proposal) ProposalOutput {
	return ProposalOutput{
		Role:       proposal.Role,
		Issuer:     sdk.MustBech32ifyAcc(proposal.Issuer),
		Recipient:  sdk.MustBech32ifyAcc(proposal.Recipient),
		Status:     proposal.Status,
		Properties: proposal.Properties,
	}
}
