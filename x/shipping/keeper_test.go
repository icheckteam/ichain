package shipping

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	addr := addrs[0]
	addr2 := addrs[1]
	addr3 := addrs[2]
	addr4 := addrs[3]

	coins := sdk.Coins{
		{Denom: "tomato", Amount: 100},
		{Denom: "eggs", Amount: 200},
	}
	coins = coins.Sort()

	keeper.coinKeeper.AddCoins(ctx, addr, coins)

	// -----------------------------------
	// MARK: - Create order
	//
	// Test not owning asset
	createOrderMsg := CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"tomato", 200}},
		Issuer:            addr2,
		Carrier:           addr,
		Receiver:          addr3,
	}
	_, err := keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// Test not enough asset
	createOrderMsg = CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"tomato", 200}},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"eggs", 201}},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"tomato", 50}, {"eggs", 201}},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"tomato", 150}, {"eggs", 50}},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID:                "1",
		TransportedAssets: []TransportedAsset{{"tomato", 200}, {"eggs", 201}},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// Valid order
	createOrderMsg = CreateOrderMsg{
		ID: "1",
		TransportedAssets: []TransportedAsset{
			{"tomato", 50},
			{"eggs", 100},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.Nil(t, err)
	assert.True(t, keeper.coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{{Denom: "eggs", Amount: 100}, {Denom: "tomato", Amount: 50}}))

	// Invalid amount
	createOrderMsg = CreateOrderMsg{
		ID:                "2",
		TransportedAssets: []TransportedAsset{},
		Issuer:            addr,
		Carrier:           addr2,
		Receiver:          addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"tomato", 0},
			{"eggs", 0},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// Create order with existing ID fails
	createOrderMsg = CreateOrderMsg{
		ID: "1",
		TransportedAssets: []TransportedAsset{
			{"tomato", 50},
			{"eggs", 100},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// Insufficient amount after the first order
	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"tomato", 51},
			{"eggs", 0},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"tomato", 0},
			{"eggs", 101},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"tomato", 51},
			{"eggs", 101},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"tomato", 50},
			{"eggs", 101},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"tomato", 51},
			{"eggs", 100},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.NotNil(t, err)

	// Valid amount after the first order
	createOrderMsg = CreateOrderMsg{
		ID: "2",
		TransportedAssets: []TransportedAsset{
			{"eggs", 50},
			{"tomato", 25},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.Nil(t, err)
	assert.True(t, keeper.coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{{Denom: "eggs", Amount: 50}, {Denom: "tomato", Amount: 25}}))

	// Valid amount after the second order
	createOrderMsg = CreateOrderMsg{
		ID: "3",
		TransportedAssets: []TransportedAsset{
			{"eggs", 50},
			{"tomato", 25},
		},
		Issuer:   addr,
		Carrier:  addr2,
		Receiver: addr3,
	}
	_, err = keeper.CreateOrder(ctx, createOrderMsg)
	assert.Nil(t, err)
	assert.True(t, keeper.coinKeeper.GetCoins(ctx, addr).IsZero())

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
	assert.Nil(t, err)

	confirmOrderMsg = ConfirmOrderMsg{
		OrderID: "2",
		Carrier: addr2,
	}
	_, err = keeper.ConfirmOrder(ctx, confirmOrderMsg)
	assert.Nil(t, err)

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
	assert.Nil(t, err)

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

	// Can cancel a confirmed order
	cancelOrderMsg = CancelOrderMsg{
		OrderID: "2",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.Nil(t, err)

	// Can cancel a pending order
	cancelOrderMsg = CancelOrderMsg{
		OrderID: "3",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.Nil(t, err)

	// Cannot cancel a completed order
	cancelOrderMsg = CancelOrderMsg{
		OrderID: "1",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.NotNil(t, err)

	// Cannot cancel a cancelled order
	cancelOrderMsg = CancelOrderMsg{
		OrderID: "2",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.NotNil(t, err)

	cancelOrderMsg = CancelOrderMsg{
		OrderID: "3",
		Issuer:  addr,
	}
	_, err = keeper.CancelOrder(ctx, cancelOrderMsg)
	assert.NotNil(t, err)
}
