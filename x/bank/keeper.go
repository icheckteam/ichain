package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/icheckteam/ichain/types"
)

const moduleName = "bank"

// Keeper manages transfers between accounts
type Keeper struct {
	am sdk.AccountMapper
}

// NewKeeper returns a new Keeper
func NewKeeper(am sdk.AccountMapper) Keeper {
	return Keeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper Keeper) GetCoins(ctx sdk.Context, addr sdk.Address) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// SetCoins sets the coins at the addr.
func (keeper Keeper) SetCoins(ctx sdk.Context, addr sdk.Address, amt sdk.Coins) sdk.Error {
	return setCoins(ctx, keeper.am, addr, amt)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr sdk.Address, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (ck Keeper) SubtractCoins(ctx sdk.Context, addr sdk.Address, amt sdk.Coins) (sdk.Coins, types.Tags, sdk.Error) {
	return subtractCoins(ctx, ck.am, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (ck Keeper) AddCoins(ctx sdk.Context, addr sdk.Address, amt sdk.Coins) (sdk.Coins, types.Tags, sdk.Error) {
	return addCoins(ctx, ck.am, addr, amt)
}

// SendCoins moves coins from one account to another
func (ck Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) (types.Tags, sdk.Error) {
	return sendCoins(ctx, ck.am, fromAddr, toAddr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (ck Keeper) InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (types.Tags, sdk.Error) {
	return inputOutputCoins(ctx, ck.am, inputs, outputs)
}

//______________________________________________________________________________________________

func getCoins(ctx sdk.Context, am sdk.AccountMapper, addr sdk.Address) sdk.Coins {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, am sdk.AccountMapper, addr sdk.Address, amt sdk.Coins) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	acc.SetCoins(amt)
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am sdk.AccountMapper, addr sdk.Address, amt sdk.Coins) bool {
	return getCoins(ctx, am, addr).IsGTE(amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func subtractCoins(ctx sdk.Context, am sdk.AccountMapper, addr sdk.Address, amt sdk.Coins) (sdk.Coins, types.Tags, sdk.Error) {
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Minus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := types.NewTags("owner", addr.Bytes())
	for _, coin := range amt {
		tags.AppendTag("asset_id", []byte(coin.Denom))
	}
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am sdk.AccountMapper, addr sdk.Address, amt sdk.Coins) (sdk.Coins, types.Tags, sdk.Error) {
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Plus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := types.NewTags("owner", addr.Bytes())

	for _, coin := range amt {
		tags.AppendTag("asset_id", []byte(coin.Denom))
	}

	return newCoins, tags, err
}

// SendCoins moves coins from one account to another
// NOTE: Make sure to revert state changes from tx on error
func sendCoins(ctx sdk.Context, am sdk.AccountMapper, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) (types.Tags, sdk.Error) {
	_, subTags, err := subtractCoins(ctx, am, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, am, toAddr, amt)
	if err != nil {
		return nil, err
	}

	return subTags.AppendTags(addTags), nil
}

func inputOutputCoins(ctx sdk.Context, am sdk.AccountMapper, inputs []Input, outputs []Output) (types.Tags, sdk.Error) {
	allTags := types.EmptyTags()
	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}
