package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Keeper manages identity claims
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
func (k Keeper) CreateClaim(ctx sdk.Context, msg MsgCreateClaim) (sdk.Tags, sdk.Error) {
	oldClaim, err := k.GetClaim(ctx, msg.ID)
	if err != nil {
		return nil, err
	}

	if oldClaim != nil && !oldClaim.IsOwner(msg.Metadata.Issuer) {
		return nil, sdk.ErrUnauthorized("")
	}

	claim := Claim{
		ID:       msg.ID,
		Metadata: msg.Metadata,
		Context:  msg.Context,
		Content:  msg.Content,
	}

	k.setClaim(ctx, claim)
	return nil, nil
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
func (k Keeper) RevokeClaim(ctx sdk.Context, msg MsgRevokeClaim) (sdk.Tags, sdk.Error) {
	claim, err := k.GetClaim(ctx, msg.ClaimID)
	if err != nil {
		return nil, err
	}
	if claim == nil || !claim.IsOwner(msg.Owner) {
		return nil, sdk.ErrUnauthorized("")
	}
	claim.Metadata.Revocation = msg.Revocation
	k.setClaim(ctx, *claim)
	return nil, nil
}
