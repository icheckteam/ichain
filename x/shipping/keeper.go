package shipping

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/types"
	"github.com/icheckteam/ichain/x/asset"
)

// Keeper manages shipping orders
type Keeper struct {
	storeKey    sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc         *wire.Codec
	assetKeeper asset.Keeper
}

// NewKeeper constructs a new keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *wire.Codec, assetKeeper asset.Keeper) Keeper {
	return Keeper{storeKey, cdc, assetKeeper}
}

// HasOrder checks if an order with the provided ID exists
func (k Keeper) hasOrder(ctx sdk.Context, uid string) bool {
	store := ctx.KVStore(k.storeKey)
	orderKey := GetOrderKey([]byte(uid))
	return store.Has(orderKey)
}

// GetOrder get order by id
func (k Keeper) getOrder(ctx sdk.Context, uid string) *Order {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetOrderKey([]byte(uid)))
	order := &Order{}

	// marshal the record and add to the state
	if err := k.cdc.UnmarshalBinary(b, order); err != nil {
		return nil
	}
	return order
}

// setOrder persists the order
func (k Keeper) setOrder(ctx sdk.Context, order Order) {
	store := ctx.KVStore(k.storeKey)
	orderKey := GetOrderKey([]byte(order.ID))

	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(order)
	if err != nil {
		panic(err)
	}

	store.Set(orderKey, bz)
}

// CreateOrder validates and creates a new order
func (k Keeper) CreateOrder(ctx sdk.Context, msg CreateOrderMsg) (sdk.Tags, sdk.Error) {
	if k.hasOrder(ctx, msg.ID) {
		return nil, ErrDuplicateOrder(msg.ID)
	}
	if len(msg.TransportedAssets) == 0 {
		return nil, ErrInvalidAsset()
	}

	assetItems := []asset.Asset{}
	for _, item := range msg.TransportedAssets {
		aitem, found := k.assetKeeper.GetAsset(ctx, item.ID)
		if !found {
			return nil, asset.ErrAssetNotFound(item.ID)
		}
		if aitem.Final {
			return nil, asset.ErrAssetAlreadyFinal(aitem.ID)
		}

		if !aitem.IsOwner(msg.Issuer) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Issuer))
		}
	}

	order := Order{
		ID:                msg.ID,
		TransportedAssets: msg.TransportedAssets,
		Issuer:            msg.Issuer,
		Carrier:           msg.Carrier,
		Receiver:          msg.Receiver,
		Status:            OrderStatusPending,
	}

	k.setOrder(ctx, order)

	allTags := sdk.EmptyTags()
	allTags.AppendTag("sender", types.AddrToBytes(msg.Issuer))
	allTags.AppendTag("recipient", types.AddrToBytes(msg.Carrier))
	allTags.AppendTag("order_id", []byte(msg.ID))

	for _, item := range assetItems {
		item.Final = true
		k.assetKeeper.SetAsset(ctx, item)
		allTags = allTags.AppendTag("asset_id", []byte(item.ID))
	}

	return allTags, nil
}

// ConfirmOrder validate the message and
// set the status of the target order to Completed
//
// Only the carrier can confirm the order, and only when the order is pending
func (k Keeper) ConfirmOrder(ctx sdk.Context, msg ConfirmOrderMsg) (sdk.Tags, sdk.Error) {
	order := k.getOrder(ctx, msg.OrderID)
	if order == nil {
		return nil, ErrUnknownOrder(msg.OrderID)
	}

	if msg.Carrier.String() != order.Carrier.String() {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to confirm", msg.Carrier))
	}

	if order.Status != OrderStatusPending {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("order id %s cannot be confirmed", msg.OrderID))
	}

	order.Status = OrderStatusConfirmed
	k.setOrder(ctx, *order)
	return nil, nil
}

// CompleteOrder vakudates the message and
// set the status of the target order to Completed
//
// Only the receiver can complete the order, and only when the order is confirmed
func (k Keeper) CompleteOrder(ctx sdk.Context, msg CompleteOrderMsg) (sdk.Tags, sdk.Error) {
	order := k.getOrder(ctx, msg.OrderID)
	if order == nil {
		return nil, ErrUnknownOrder(msg.OrderID)
	}

	if msg.Receiver.String() != order.Receiver.String() {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to set the order status to complete", msg.Receiver))
	}

	if order.Status != OrderStatusConfirmed {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("order id %s cannot be completed", msg.OrderID))
	}

	order.Status = OrderStatusCompleted
	k.setOrder(ctx, *order)
	return nil, nil
}

// CancelOrder vakudates the message and
// set the status of the target order to Cancelled
//
// Only the issuer can cancel the order
func (k Keeper) CancelOrder(ctx sdk.Context, msg CancelOrderMsg) (sdk.Tags, sdk.Error) {
	order := k.getOrder(ctx, msg.OrderID)
	if order == nil {
		return nil, ErrUnknownOrder(msg.OrderID)
	}

	if msg.Issuer.String() != order.Issuer.String() {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to cancel the order", msg.Issuer))
	}

	if order.Status == OrderStatusCancelled || order.Status == OrderStatusCompleted {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("order id %s cannot be cancelled", msg.OrderID))
	}

	for _, asset := range order.TransportedAssets {
		a, _ := k.assetKeeper.GetAsset(ctx, asset.ID)
		a.Final = false
		k.assetKeeper.SetAsset(ctx, a)
	}

	order.Status = OrderStatusCancelled
	k.setOrder(ctx, *order)
	return nil, nil
}
