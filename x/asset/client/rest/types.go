package rest

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/x/asset"
)

type baseBody struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	Sequence         int64  `json:"sequence"`
	AccountNumber    int64  `json:"account_number"`
	Gas              int64  `json:"gas"`
}

func (b baseBody) Validate() error {
	if b.LocalAccountName == "" {
		return errors.New("account_name is required")
	}
	if b.Password == "" {
		return errors.New("password is required")
	}
	if b.Gas == 0 {
		return errors.New("gas is required")
	}
	return nil
}

type AssetOutput struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Height     int64            `json:"height"`
	Name       string           `json:"name"`
	Owner      string           `json:"owner"`
	Quantity   int64            `json:"quantity"`
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
		Precision:  a.Precision,
		Created:    a.Created,
	}
}

func ToAssetsOutput(asa []asset.Asset) (asb []AssetOutput) {
	for _, a := range asa {
		asb = append(asb, ToAssetOutput(a))
	}
	return
}
