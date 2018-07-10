package invoice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Item struct {
	AssetID string `json:"asset_id"`
}

type Invoice struct {
	ID         string         `json:"id"`
	Issuer     sdk.AccAddress `json:"issuer"`
	Receiver   sdk.AccAddress `json:"receiver"`
	Items      []Item         `json:"items"`
	CreateTime int64          `json:"create_time"`
}
