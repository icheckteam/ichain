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
	Memo          string `json:"memo"`
}

func (b baseBody) Validate() error {
	if b.Name == "" {
		return errors.New("name required but not specified")
	}
	if b.Password == "" {
		return errors.New("password required but not specified")
	}
	if b.Gas == 0 {
		return errors.New("gas required but not specified")
	}
	if len(b.ChainID) == 0 {
		return errors.New("chain_id required but not specified")
	}
	if b.AccountNumber < 0 {
		return errors.New("account_number required but not specified")
	}

	if b.Sequence < 0 {
		return errors.New("sequence required but not specified")
	}
	return nil
}

type msgCreateCreateProposalBody struct {
	BaseReq baseBody `json:"base_req"`

	Recipient  string             `json:"recipient"`
	Properties []string           `json:"properties"`
	Role       asset.ProposalRole `json:"role"`
}

func (b msgCreateCreateProposalBody) ValidateBasic() error {
	err := b.BaseReq.Validate()
	if err != nil {
		return err
	}
	if b.Recipient == "" {
		return errors.New("recipient is required")
	}

	switch b.Role {
	case asset.RoleOwner, asset.RoleReporter:
		break
	default:
		return errors.New("invalid role")
	}

	if b.Role == asset.RoleReporter && len(b.Properties) == 0 {
		return errors.New("properties is required")
	}
	return nil
}

type msgAnswerProposalBody struct {
	BaseReq baseBody `json:"base_req"`

	Response asset.ProposalStatus `json:"response"`
	AssetID  string               `json:"asset_id"`
	Role     asset.ProposalRole   `json:"role"`
}

func (b msgAnswerProposalBody) ValidateBasic() error {
	err := b.BaseReq.Validate()
	if err != nil {
		return err
	}
	switch b.Response {
	case asset.StatusAccepted, asset.StatusCancel, asset.StatusRejected:
		break
	default:
		return errors.New("invalid response")
	}

	return nil
}
