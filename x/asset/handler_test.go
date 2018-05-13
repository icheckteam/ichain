package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRegister(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 1000)
	addr := addrs[0]

	var msg = RegisterMsg{
		Issuer:   addr,
		ID:       "a23423423423426546546",
		Name:     "adad",
		Quantity: 2,
	}

	got := handleRegisterAsset(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected no error on handleRegisterAsset")

	got = handleRegisterAsset(ctx, keeper, msg)
	require.False(t, got.IsOK(), "expected no error on handleRegisterAsset")

	asset := keeper.GetAsset(ctx, msg.ID)
	require.True(t, asset != nil)
	assert.Equal(t, msg.Quantity, asset.Quantity)
	assert.Equal(t, msg.Quantity, keeper.bank.GetCoins(ctx, addr)[0].Amount)
	assert.Equal(t, msg.ID, keeper.bank.GetCoins(ctx, addr)[0].Denom)

	// Test handle add quantity
	got = handleAddQuantity(ctx, keeper, AddQuantityMsg{
		Issuer:   addr,
		ID:       msg.ID,
		Quantity: 10,
	})
	require.True(t, got.IsOK(), "expected no error on handleAddQuantity")

	asset = keeper.GetAsset(ctx, msg.ID)
	require.True(t, asset != nil)
	assert.Equal(t, int64(12), asset.Quantity)
	assert.Equal(t, int64(12), keeper.bank.GetCoins(ctx, addr)[0].Amount)
	assert.Equal(t, msg.ID, keeper.bank.GetCoins(ctx, addr)[0].Denom)

	// Test handle add quantity
	got = handleSubtractQuantity(ctx, keeper, SubtractQuantityMsg{
		Issuer:   addr,
		ID:       msg.ID,
		Quantity: 10,
	})
	require.True(t, got.IsOK(), "expected no error on handleSubtractQuantity")

	asset = keeper.GetAsset(ctx, msg.ID)
	require.True(t, asset != nil)
	assert.Equal(t, int64(2), asset.Quantity)
	assert.Equal(t, int64(2), keeper.bank.GetCoins(ctx, addr)[0].Amount)
	assert.Equal(t, msg.ID, keeper.bank.GetCoins(ctx, addr)[0].Denom)

	got = handleSubtractQuantity(ctx, keeper, SubtractQuantityMsg{
		Issuer:   addr,
		ID:       msg.ID,
		Quantity: 2,
	})
	require.True(t, got.IsOK(), "expected no error on handleSubtractQuantity")

	asset = keeper.GetAsset(ctx, msg.ID)
	require.True(t, asset != nil)
	assert.Equal(t, int64(0), asset.Quantity)
	assert.Equal(t, true, keeper.bank.GetCoins(ctx, addr).IsZero())

	got = handleSubtractQuantity(ctx, keeper, SubtractQuantityMsg{
		Issuer:   addr,
		ID:       msg.ID,
		Quantity: 2,
	})
	require.False(t, got.IsOK(), "expected no error on handleSubtractQuantity")

	got = handleUpdateAttr(ctx, keeper, UpdateAttrMsg{
		ID:     msg.ID,
		Issuer: addr,
		Attributes: []Attribute{
			attr,
		},
	})
	require.True(t, got.IsOK(), "expected no error on handleUpdateAttr")
	asset = keeper.GetAsset(ctx, msg.ID)
	require.True(t, asset.Attributes[0].Name == "weight")
	require.True(t, asset.Attributes[0].NumberValue == 100)

	got = handleUpdateAttr(ctx, keeper, UpdateAttrMsg{
		ID:     msg.ID,
		Issuer: addr,
		Attributes: []Attribute{
			Attribute{
				Name: "location",
				Type: 3,
				Location: Location{
					Latitude:  1,
					Longitude: 1,
				},
			},
		},
	})
	require.True(t, got.IsOK(), "expected no error on handleUpdateAttr")
	asset = keeper.GetAsset(ctx, msg.ID)
	require.True(t, asset.Attributes[1].Name == "location")
	require.True(t, asset.Attributes[1].Location.Latitude == 1)

	got = handleUpdateAttr(ctx, keeper, UpdateAttrMsg{
		ID:     msg.ID,
		Issuer: addrs[1],
		Attributes: []Attribute{
			attr,
		},
	})
	require.False(t, got.IsOK(), "expected no error on handleUpdateAttr")
}
