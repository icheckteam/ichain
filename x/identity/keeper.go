package identity

import (
	"github.com/cosmos-sdk/wire"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper ....
type Keeper struct {
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc      *wire.Codec
}

// ClaimIssue ...
func (k Keeper) ClaimIssue(ctx sdk.Context, msg ClaimIssueMsg) sdk.Error {
	key := GetClaimRecordKey(msg.ID)
	store := ctx.KVStore(k.storeKey)
	if store.Has(key) == false {
		return sdk.ErrInternal("Claim already exists")
	}

	var b []byte
	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &msg); err != nil {
		return sdk.ErrInternal(err.Error())
	}

	store.Set(key, b)
	return nil
}

// GetClaim ...
func (k Keeper) GetClaim(ctx sdk.Context, claimID string) (*Claim, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	key := GetClaimRecordKey(claimID)
	claim := &Claim{}
	b := store.Get(key)
	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &claim); err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return claim, nil
}

// Revoke ...
func (k Keeper) Revoke(ctx sdk.Context, claimID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	key := GetClaimRecordKey(claimID)
	claim, err := k.GetClaim(ctx, claimID)
	if err != nil {
		return err
	}
	var b []byte
	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &claim); err != nil {
		return sdk.ErrInternal(err.Error())
	}

	store.Set(key, b)
	return nil
}
