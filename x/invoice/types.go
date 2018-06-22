package invoice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Item struct {
	AssetID string `json:"asset_id"`
}

type Invoice struct {
	ID         string      `json:"id"`
	Issuer     sdk.Address `json:"issuer"`
	Receiver   sdk.Address `json:"receiver"`
	Items      []Item      `json:"items"`
	CreateTime int64       `json:"create_time"`
}
