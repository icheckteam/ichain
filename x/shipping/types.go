package shipping

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Order is a shipping order
type Order struct {
	ID                string             `json:"id"`                 // ID of the order, provided by the client
	TransportedAssets []TransportedAsset `json:"transported_assets"` // The assets to be transported
	Issuer            sdk.AccAddress     `json:"issuer"`             // The issuer of the order, must also be the owner of the asset (coin)
	Carrier           sdk.AccAddress     `json:"carrier"`            // The carrier
	Receiver          sdk.AccAddress     `json:"receiver"`           // The receiver, often a buyer
	Status            OrderStatus        `json:"status"`             // The status of the order
}

// OrderStatus represents the order's status
type OrderStatus int

// Valid status of an order
const (
	OrderStatusPending   OrderStatus = iota // Order has been created
	OrderStatusConfirmed                    // The carrier received the asset from the issuer
	OrderStatusCompleted                    // The receiver received the asset from the carrier
	OrderStatusCancelled                    // Order is cancelled
)

// TransportedAsset contains the id of the asset
// and the quantity to be transported
type TransportedAsset struct {
	ID string `json:"id"` // ID of the asset (coin) to be transported
}
