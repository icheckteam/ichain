package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Keeper manages the identity and certificate
type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *wire.Codec
	codespace sdk.CodespaceType
}

// NewKeeper ...
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: DefaultCodespace,
	}
}

// SetTrust set the main record holding trust details
func (k Keeper) SetTrust(ctx sdk.Context, trustor, trusting sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyTrust(trustor, trusting), []byte{})
}

// DeleteTrust delete cert from the store
func (k Keeper) DeleteTrust(ctx sdk.Context, trustor, trusting sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyTrust(trustor, trusting))
}

// AddTrust add a trust
func (k Keeper) AddTrust(ctx sdk.Context, msg MsgSetTrust) (sdk.Tags, sdk.Error) {
	if msg.Trust == true {
		k.SetTrust(ctx, msg.Trustor, msg.Trusting)
	} else {
		k.DeleteTrust(ctx, msg.Trustor, msg.Trusting)
	}
	return nil, nil
}

func (k Keeper) hasTrust(ctx sdk.Context, trustor, trusting sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyTrust(trustor, trusting))
}
