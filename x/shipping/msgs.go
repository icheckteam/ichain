package shipping

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/types"
)

const msgType = "shipping"

// ------------------------------------ Create Order

// CreateOrderMsg is sent by the issuer to construct a new order,
// all fields are required
type CreateOrderMsg struct {
	ID                string             `json:"id"`                 // ID of the order, provided by the client
	TransportedAssets []TransportedAsset `json:"transported_assets"` // The assets to be transported
	Issuer            sdk.Address        `json:"issuer"`             // The issuer of the order, must also be the owner of the asset (coin)
	Carrier           sdk.Address        `json:"carrier"`            // The carrier
	Receiver          sdk.Address        `json:"receiver"`           // The receiver, often a buyer
}

// enforce the Msg type at compile time
var _ sdk.Msg = CreateOrderMsg{}

// nolint ...
func (msg CreateOrderMsg) Type() string                            { return msgType }
func (msg CreateOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CreateOrderMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg CreateOrderMsg) String() string {
	return fmt.Sprintf(`
		CreateOrderMsg{
			ID: %s, 
			TransportedAssets: %v,
			Issuer: %s,
			Carrier %s,
			Receiver: %s
		}	
	`,
		msg.ID, msg.TransportedAssets,
		msg.Issuer, msg.Carrier, msg.Receiver,
	)
}

// GetSignBytes ...
func (msg CreateOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic ...
func (msg CreateOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}

	if len(msg.Carrier) == 0 {
		return sdk.ErrUnknownAddress(msg.Carrier.String()).Trace("")
	}

	if len(msg.Receiver) == 0 {
		return sdk.ErrUnknownAddress(msg.Receiver.String()).Trace("")
	}

	if len(msg.ID) == 0 {
		return types.ErrMissingField("id")
	}

	if len(msg.TransportedAssets) == 0 {
		return types.ErrMissingField("transported_assets")
	}

	return nil
}

// ------------------------------------ Confirm Order

// ReceiveOrderMsg is sent by the carrier to confirm
// that the carrier has received the asset from the issuer
type ConfirmOrderMsg struct {
	OrderID string      `json:"order_id"`   // ID of the order received
	Carrier sdk.Address `json:"carrier_id"` // the carrier, who confirms the order
}

// enforce the Msg type at compile time
var _ sdk.Msg = ConfirmOrderMsg{}

// nolint ...
func (msg ConfirmOrderMsg) Type() string                            { return msgType }
func (msg ConfirmOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg ConfirmOrderMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Carrier} }
func (msg ConfirmOrderMsg) String() string {
	return fmt.Sprintf(`
		ConfirmOrderMsg{
			OrderID: %s, 
			Carrier: %s,
		}	
	`,
		msg.OrderID, msg.Carrier,
	)
}

// GetSignBytes ...
func (msg ConfirmOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic ...
func (msg ConfirmOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Carrier) == 0 {
		return sdk.ErrUnknownAddress(msg.Carrier.String()).Trace("")
	}

	if len(msg.OrderID) == 0 {
		return types.ErrMissingField("order_id")
	}

	return nil
}

// ------------------------------------ Complete Order

// CompleteOrderMsg is sent by the receiver to confirm
// that the receiver has received the asset from the carrier
type CompleteOrderMsg struct {
	OrderID  string      `json:"order_id"`    // ID of the order completed
	Receiver sdk.Address `json:"receiver_id"` // the receiver, who confirm the completion of the order
}

// enforce the Msg type at compile time
var _ sdk.Msg = CompleteOrderMsg{}

// nolint ...
func (msg CompleteOrderMsg) Type() string                            { return msgType }
func (msg CompleteOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CompleteOrderMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Receiver} }
func (msg CompleteOrderMsg) String() string {
	return fmt.Sprintf(`
		ReceiveOrderMsg{
			OrderID: %s, 
			Receiver: %s,
		}	
	`,
		msg.OrderID, msg.Receiver,
	)
}

// GetSignBytes ...
func (msg CompleteOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic ...
func (msg CompleteOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Receiver) == 0 {
		return sdk.ErrUnknownAddress(msg.Receiver.String()).Trace("")
	}

	if len(msg.OrderID) == 0 {
		return types.ErrMissingField("order_id")
	}

	return nil
}

// ------------------------------------ Cancel Order

// CompleteOrderMsg is sent by the receiver to confirm
// that the receiver has received the asset from the carrier
type CancelOrderMsg struct {
	OrderID string      `json:"order_id"`    // ID of the order to be cancelled
	Issuer  sdk.Address `json:"receiver_id"` // the issuer
}

// enforce the Msg type at compile time
var _ sdk.Msg = CancelOrderMsg{}

// nolint ...
func (msg CancelOrderMsg) Type() string                            { return msgType }
func (msg CancelOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CancelOrderMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg CancelOrderMsg) String() string {
	return fmt.Sprintf(`
		CancelOrderMsg{
			OrderID: %s, 
			Issuer: %s,
		}	
	`,
		msg.OrderID, msg.Issuer,
	)
}

// GetSignBytes ...
func (msg CancelOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic ...
func (msg CancelOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}

	if len(msg.OrderID) == 0 {
		return types.ErrMissingField("order_id")
	}

	return nil
}
