package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	types "github.com/icheckteam/ichain/types"
)

// Keeper ...
type Keeper struct {
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc      *wire.Codec
	am       sdk.AccountMapper
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, am sdk.AccountMapper) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
		am:       am,
	}
}

// ClaimIssue ...
func (k Keeper) Create(ctx sdk.Context, msg CreateMsg) (types.Tags, sdk.Error) {
	allTags := types.EmptyTags()
	claim, err := k.GetClaim(ctx, msg.ID)
	if err != nil {
		return allTags, err
	}

	if claim == nil || !claim.IsOwner(msg.Metadata.Issuer) {
		return allTags, sdk.ErrUnauthorized("")
	}

	key := GetClaimRecordKey(msg.ID)
	store := ctx.KVStore(k.storeKey)
	var b []byte
	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &msg); err != nil {
		return allTags, sdk.ErrInternal(err.Error())
	}

	store.Set(key, b)
	// append tags
	allTags.AppendTag("owner", msg.Metadata.Issuer)
	allTags.AppendTag("owner", msg.Metadata.Recipient)
	return allTags, nil
}

// GetClaim ...
func (k Keeper) GetClaim(ctx sdk.Context, claimID string) (*Claim, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	key := GetClaimRecordKey(claimID)
	claim := &Claim{}
	b := store.Get(key)

	if len(b) == 0 {
		return nil, nil
	}

	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &claim); err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return claim, nil
}

// Revoke ...
func (k Keeper) Revoke(ctx sdk.Context, claimID, revocation string) (types.Tags, sdk.Error) {
	allTags := types.EmptyTags()
	store := ctx.KVStore(k.storeKey)
	key := GetClaimRecordKey(claimID)
	claim, err := k.GetClaim(ctx, claimID)
	if err != nil {
		return nil, err
	}
	var b []byte
	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &claim); err != nil {
		return nil, sdk.ErrUnauthorized(err.Error())
	}

	store.Set(key, b)
	allTags.AppendTag("owner", claim.Metadata.Issuer)
	allTags.AppendTag("owner", claim.Metadata.Recipient)
	return allTags, nil
}
