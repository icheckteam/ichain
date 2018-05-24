package invoice

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Item struct {
	AssetID  string `json:"asset_id"`
	Quantity int64  `json:"quantity"`
}

type Invoice struct {
	ID         string      `json:"id"`
	Issuer     sdk.Address `json:"issuer"`
	Receiver   sdk.Address `json:"receiver"`
	Items      []Item      `json:"items"`
	CreateTime time.Time   `json:"create_time"`
}
