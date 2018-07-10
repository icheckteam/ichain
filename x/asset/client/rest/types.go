package rest

import (
	"errors"

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
