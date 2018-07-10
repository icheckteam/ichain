package shipping

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/types"
)

const msgType = "shipping"

// enforce the Msg type at compile time
var _, _, _ sdk.Msg = ConfirmOrderMsg{}, CreateOrderMsg{}, CancelOrderMsg{}

// ------------------------------------ Create Order

// CreateOrderMsg is sent by the issuer to construct a new order,
// all fields are required
type CreateOrderMsg struct {
	ID                string             `json:"id"`                 // ID of the order, provided by the client
	TransportedAssets []TransportedAsset `json:"transported_assets"` // The assets to be transported
	Issuer            sdk.AccAddress     `json:"issuer"`             // The issuer of the order, must also be the owner of the asset (coin)
	Carrier           sdk.AccAddress     `json:"carrier"`            // The carrier
	Receiver          sdk.AccAddress     `json:"receiver"`           // The receiver, often a buyer
}

// enforce the Msg type at compile time
var _ sdk.Msg = CreateOrderMsg{}

// nolint ...
func (msg CreateOrderMsg) Type() string                            { return msgType }
func (msg CreateOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CreateOrderMsg) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Issuer} }
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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic ...
func (msg CreateOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil issuer address")
	}
	if len(msg.Carrier) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil carrier address")
	}

	if len(msg.Receiver) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil receiver address")
	}

	if len(msg.ID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "id")
	}

	if len(msg.TransportedAssets) == 0 {
		return types.ErrMissingField(DefaultCodespace, "transported_assets")
	}

	return nil
}

// ------------------------------------ Confirm Order

// ReceiveOrderMsg is sent by the carrier to confirm
// that the carrier has received the asset from the issuer
type ConfirmOrderMsg struct {
	OrderID string         `json:"order_id"`   // ID of the order received
	Carrier sdk.AccAddress `json:"carrier_id"` // the carrier, who confirms the order
}

// nolint ...
func (msg ConfirmOrderMsg) Type() string                 { return msgType }
func (msg ConfirmOrderMsg) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Carrier} }

// GetSignBytes ...
func (msg ConfirmOrderMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic ...
func (msg ConfirmOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Carrier) == 0 {
		return sdk.ErrUnknownAddress(msg.Carrier.String())
	}

	if len(msg.OrderID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "order_id")
	}

	return nil
}

// ------------------------------------ Complete Order

// CompleteOrderMsg is sent by the receiver to confirm
// that the receiver has received the asset from the carrier
type CompleteOrderMsg struct {
	OrderID  string         `json:"order_id"`    // ID of the order completed
	Receiver sdk.AccAddress `json:"receiver_id"` // the receiver, who confirm the completion of the order
}

// enforce the Msg type at compile time
var _ sdk.Msg = CompleteOrderMsg{}

// nolint ...
func (msg CompleteOrderMsg) Type() string                            { return msgType }
func (msg CompleteOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CompleteOrderMsg) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Receiver} }
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
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic ...
func (msg CompleteOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Receiver) == 0 {
		return sdk.ErrUnknownAddress(msg.Receiver.String())
	}

	if len(msg.OrderID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "order_id")
	}

	return nil
}

// ------------------------------------ Cancel Order

// CompleteOrderMsg is sent by the receiver to confirm
// that the receiver has received the asset from the carrier
type CancelOrderMsg struct {
	OrderID string         `json:"order_id"`  // ID of the order to be cancelled
	Issuer  sdk.AccAddress `json:"issuer_id"` // the issuer
}

// enforce the Msg type at compile time
var _ sdk.Msg = CancelOrderMsg{}

// nolint ...
func (msg CancelOrderMsg) Type() string                            { return msgType }
func (msg CancelOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CancelOrderMsg) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Issuer} }

// GetSignBytes ...
func (msg CancelOrderMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic ...
func (msg CancelOrderMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String())
	}

	if len(msg.OrderID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "order_id")
	}

	return nil
}
