package shipping

import (
	"testing"

	"github.com/icheckteam/ichain/x/asset"
	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, keeper := createTestInput(t, false, 0)
	addr := addrs[0]
	addr2 := addrs[1]
	addr3 := addrs[2]
	addr4 := addrs[3]

	keeper.assetKeeper.CreateAsset(ctx, asset.MsgCreateAsset{
		AssetID:  "tomato",
		Quantity: 1,
		Sender:   addr,
	})
	keeper.assetKeeper.CreateAsset(ctx, asset.MsgCreateAsset{
		AssetID:  "eggs",
		Quantity: 1,
		Sender:   addr,
	})

	// -----------------------------------
	// MARK: - Create order
	//
	// Test not owning asset
	createOrderMsg := CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"tomato"}},
		Issuer:            addr2,
		Carrier:           addr,
		Receiver:          addr3,
	}
	_, err := keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// Valid order
	createOrderMsg = CreateOrderMsg{
		ID: "1",
		TransportedAssets: []TransportedAsset{
			{"tomato"},
			{"eggs"},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	order := keeper.getOrder(ctx, "1")
	assert.Nil(t, err)
	assert.True(t, order.ID == "1")
	assert.True(t, order.TransportedAssets[0].ID == "tomato")
	assert.True(t, order.TransportedAssets[1].ID == "eggs")
	assert.True(t, order.Issuer.String() == addr.String())
	assert.True(t, order.Carrier.String() == addr2.String())
	assert.True(t, order.Receiver.String() == addr3.String())

	// Invalid asset
	createOrderMsg = CreateOrderMsg{
		ID:                "2",
		TransportedAssets: []TransportedAsset{},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// -----------------------------------
	// MARK: - Confirm order
	//
	// Test invalid ID
	confirmOrderMsg := ConfirmOrderMsg{
		OrderID: "4",
		Carrier: addr2,
	}
	_, err = keeper.ConfirmOrder(ctx, confirmOrderMsg)
	assert.NotNil(t, err)

	// Test invalid carrier
	confirmOrderMsg = ConfirmOrderMsg{
		OrderID: "1",
		Carrier: addr,
	}
	_, err = keeper.ConfirmOrder(ctx, confirmOrderMsg)
	assert.NotNil(t, err)

	confirmOrderMsg = ConfirmOrderMsg{
		OrderID: "1",
		Carrier: addr3,
	}
	_, err = keeper.ConfirmOrder(ctx, confirmOrderMsg)
	assert.NotNil(t, err)

	// Valid confirm
	confirmOrderMsg = ConfirmOrderMsg{
		OrderID: "1",
		Carrier: addr2,
	}
	_, err = keeper.ConfirmOrder(ctx, confirmOrderMsg)
	order = keeper.getOrder(ctx, "1")
	assert.Nil(t, err)
	assert.True(t, order.Status == OrderStatusConfirmed)

	// Confirmed order cannot be confirmed again
	confirmOrderMsg = ConfirmOrderMsg{
		OrderID: "1",
		Carrier: addr2,
	}
	_, err = keeper.ConfirmOrder(ctx, confirmOrderMsg)
	assert.NotNil(t, err)

	// -----------------------------------
	// MARK: - Complete order
	//
	// Test invalid ID
	completeOrderMsg := CompleteOrderMsg{
		OrderID:  "4",
		Receiver: addr3,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	assert.NotNil(t, err)

	// Invalid receiver
	completeOrderMsg = CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addr2,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	assert.NotNil(t, err)

	completeOrderMsg = CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addr,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	assert.NotNil(t, err)

	completeOrderMsg = CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addr4,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	assert.NotNil(t, err)

	// Cannot complete order that is not confirmed
	completeOrderMsg = CompleteOrderMsg{
		OrderID:  "3",
		Receiver: addr3,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	assert.NotNil(t, err)

	// Valid completion
	completeOrderMsg = CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addr3,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	order = keeper.getOrder(ctx, "1")
	assert.Nil(t, err)
	assert.True(t, order.Status == OrderStatusCompleted)

	// Cannot complete a completed order
	completeOrderMsg = CompleteOrderMsg{
		OrderID:  "1",
		Receiver: addr3,
	}
	_, err = keeper.CompleteOrder(ctx, completeOrderMsg)
	assert.NotNil(t, err)

	// -----------------------------------
	// MARK: - Cancel order
	//
	// Test invalid ID
	cancelOrderMsg := CancelOrderMsg{
		OrderID: "4",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.NotNil(t, err)

	// Invalid issuer
	cancelOrderMsg = CancelOrderMsg{
		OrderID: "1",
		Issuer:  addr2,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.NotNil(t, err)

	// Cannot cancel a completed order
	cancelOrderMsg = CancelOrderMsg{
		OrderID: "1",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	order = keeper.getOrder(ctx, "1")
	assert.NotNil(t, err)
	assert.True(t, order.Status == OrderStatusCompleted)

}
