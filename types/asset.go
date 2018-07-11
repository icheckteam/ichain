package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AssetAmount
type AssetAmount struct {
	AssetID string `json:"asset_id"`
	Amount  sdk.Int
}
