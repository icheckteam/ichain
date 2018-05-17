package asset

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	addr := sdk.Address([]byte("addr1"))

	asset := Asset{
		ID:       "asset1",
		Issuer:   addr,
		Name:     "asset 1",
		Quantity: 100,
	}

	// Test register asset
	keeper.RegisterAsset(ctx, asset)
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))
	keeper.RegisterAsset(ctx, asset)
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))

	// Test add quantity
	keeper.AddQuantity(ctx, AddQuantityMsg{ID: asset.ID, Issuer: addr, Quantity: 50})
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 150}}))

	// Test subtract quantity
	keeper.SubtractQuantity(ctx, SubtractQuantityMsg{ID: asset.ID, Issuer: addr, Quantity: 50})
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))
	keeper.SubtractQuantity(ctx, SubtractQuantityMsg{ID: asset.ID, Issuer: addr, Quantity: 102})
	assert.True(t, keeper.bank.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.Coin{Denom: asset.ID, Amount: 100}}))

	// Test update attributes
	attrs := []Attribute{Attribute{Name: "weight", NumberValue: 100}}
	keeper.UpdateAttribute(ctx, UpdateAttrMsg{ID: asset.ID, Issuer: addr, Attributes: attrs})
	newAsset := keeper.GetAsset(ctx, asset.ID)
	assert.True(t, newAsset.Attributes[0].Name == "weight")
	assert.True(t, newAsset.Attributes[0].NumberValue == 100)

	attrs = []Attribute{Attribute{Name: "weight", NumberValue: 101}}
	keeper.UpdateAttribute(ctx, UpdateAttrMsg{ID: asset.ID, Issuer: addr, Attributes: attrs})
	newAsset = keeper.GetAsset(ctx, asset.ID)
	assert.True(t, newAsset.Attributes[0].Name == "weight")
	assert.True(t, newAsset.Attributes[0].NumberValue == 101)
}
