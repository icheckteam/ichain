package identity

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Register register an identity
func (k Keeper) Register(ctx sdk.Context, msg MsgReg) ([]sdk.Tags, sdk.Error) {
	ownerCount := k.getOwnerCount(ctx, msg.Ident)

	if ownerCount > 0 {
		return nil, ErrIDAlreadyExists(DefaultCodespace, msg.Ident)
	}

	// store data
	k.setOwnerCount(ctx, msg.Ident, 1)
	k.setOwner(ctx, msg.Ident, msg.Sender)
	return nil, nil
}

// AddOwner add an account to identity
func (k Keeper) AddOwner(ctx sdk.Context, msg MsgAddOwner) ([]sdk.Tags, sdk.Error) {
	if !k.hasOwner(ctx, msg.Ident, msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("addr %s unauthorized", msg.Sender))
	}
	ownerCount := k.getOwnerCount(ctx, msg.Ident)
	k.setOwnerCount(ctx, msg.Ident, ownerCount+1)
	k.setOwner(ctx, msg.Ident, msg.Owner)
	return nil, nil
}

// DeleteOwner delete an account of identity
func (k Keeper) DeleteOwner(ctx sdk.Context, msg MsgDelOwner) ([]sdk.Tags, sdk.Error) {
	if !k.hasOwner(ctx, msg.Ident, msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("addr %s unauthorized", msg.Sender))
	}
	ownerCount := k.getOwnerCount(ctx, msg.Ident)
	k.setOwnerCount(ctx, msg.Ident, ownerCount-1)
	k.delOwner(ctx, msg.Ident, msg.Owner)
	return nil, nil
}

// hasOwner check owner of the identity
func (k Keeper) hasOwner(ctx sdk.Context, id sdk.AccAddress, owner sdk.AccAddress) bool {
	if bytes.Equal(id, owner) {
		return true
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyOwner(id, owner))
}

// setOwner
func (k Keeper) setOwner(ctx sdk.Context, id sdk.AccAddress, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyOwner(id, owner), []byte{})
}

// delOwner ...
func (k Keeper) delOwner(ctx sdk.Context, id sdk.AccAddress, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyOwner(id, owner))
}

// getOwnerCount ...
func (k Keeper) getOwnerCount(ctx sdk.Context, id sdk.AccAddress) (count int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyOwnerCount(id))
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinary(bz, &count)
	return
}

// getOwnerCount ...
func (k Keeper) setOwnerCount(ctx sdk.Context, id sdk.AccAddress, num int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(num)
	store.Set(KeyOwnerCount(id), bz)
}

// GetOwners ...
func (k Keeper) GetOwners(ctx sdk.Context, id sdk.AccAddress) []sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	// delete subspace
	iterator := sdk.KVStorePrefixIterator(store, KeyOwners(id))
	owners := []sdk.AccAddress{}
	for ; iterator.Valid(); iterator.Next() {
		addrs := iterator.Key()[1:] // remove prefix bytes
		if len(addrs) == 2*sdk.AddrLen {
			owners = append(owners, sdk.AccAddress(addrs[sdk.AddrLen:]))
		}
	}
	iterator.Close()
	return owners
}
