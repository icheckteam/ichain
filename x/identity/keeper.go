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
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

// ClaimIssue ...
func (k Keeper) Create(ctx sdk.Context, claim Claim) (types.Tags, sdk.Error) {
	allTags := types.EmptyTags()
	oldClaim, err := k.GetClaim(ctx, claim.ID)
	if err != nil {
		return allTags, err
	}

	if oldClaim != nil && !oldClaim.IsOwner(claim.Metadata.Issuer) {
		return allTags, sdk.ErrUnauthorized("")
	}

	k.setClaim(ctx, claim)

	err = k.addClaimsAccount(ctx, claim.Metadata.Recipient, claim.ID)
	if err != nil {
		return allTags, err
	}
	err = k.addClaimsAccount(ctx, claim.Metadata.Issuer, claim.ID)
	if err != nil {
		return allTags, err
	}

	// append tags
	allTags.AppendTag("owner", claim.Metadata.Issuer)
	allTags.AppendTag("owner", claim.Metadata.Recipient)
	return allTags, nil
}

func (k Keeper) setClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	key := GetClaimRecordKey(claim.ID)

	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(claim)
	if err != nil {
		panic(err)
	}

	store.Set(key, bz)
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
func (k Keeper) Revoke(ctx sdk.Context, addr sdk.Address, claimID, revocation string) (types.Tags, sdk.Error) {
	allTags := types.EmptyTags()
	claim, err := k.GetClaim(ctx, claimID)
	if err != nil {
		return nil, err
	}
	if claim == nil || !claim.IsOwner(addr) {
		return allTags, sdk.ErrUnauthorized("")
	}
	claim.Metadata.Revocation = revocation
	k.setClaim(ctx, *claim)
	allTags.AppendTag("owner", claim.Metadata.Issuer)
	allTags.AppendTag("owner", claim.Metadata.Recipient)
	return allTags, nil
}

func (k Keeper) getClaimsAccount(ctx sdk.Context, addr sdk.Address) ([]string, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	key := GetClaimsAccountKey(addr)
	claimIDS := []string{}
	b := store.Get(key)
	if len(b) == 0 {
		return nil, nil
	}
	// marshal the claim and add to the state
	if err := k.cdc.UnmarshalBinary(b, &claimIDS); err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return claimIDS, nil
}

func (k Keeper) setClaimsAccount(ctx sdk.Context, addr sdk.Address, ids []string) {
	store := ctx.KVStore(k.storeKey)
	key := GetClaimsAccountKey(addr)

	bz, err := k.cdc.MarshalBinary(ids)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) addClaimsAccount(ctx sdk.Context, addr sdk.Address, id string) sdk.Error {
	claimIDS, err := k.getClaimsAccount(ctx, addr)
	if err != nil {
		return err
	}
	claimIDS = append(claimIDS, id)
	k.setClaimsAccount(ctx, addr, claimIDS)
	return nil
}
