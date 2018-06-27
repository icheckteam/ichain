package identity

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	costCreateClaim sdk.Gas = 100
	costRevokeClaim sdk.Gas = 10
)

// Keeper manages identity claims
type Keeper struct {
	storeKey   sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc        *wire.Codec
	coinKeeper bank.Keeper
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, coinKeeper bank.Keeper) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		coinKeeper: coinKeeper,
	}
}

// ClaimIssue ...
func (k Keeper) CreateClaim(ctx sdk.Context, msg MsgCreateClaim) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costCreateClaim, "createClaim")
	oldClaim := k.GetClaim(ctx, msg.ClaimID)
	if oldClaim != nil && !oldClaim.IsOwner(msg.Issuer) {
		return nil, sdk.ErrUnauthorized("")
	}
	claim := Claim{
		ID:         msg.ClaimID,
		Issuer:     msg.Issuer,
		Recipient:  msg.Recipient,
		Context:    msg.Context,
		Content:    msg.Content,
		CreateTime: ctx.BlockHeader().Time,
		Expires:    msg.Expires,
	}

	k.setClaim(ctx, claim)
	k.setClaimByRecipientIndex(ctx, claim)
	k.setClaimByIssuerIndex(ctx, claim)

	allTags := sdk.NewTags(
		"sender", []byte(msg.Issuer.String()),
		"recipient", []byte(msg.Recipient.String()),
	)
	return allTags, nil
}

// set claim
func (k Keeper) setClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(claim)
	store.Set(GetClaimKey(claim.ID), bz)
}

func (k Keeper) removeClaim(ctx sdk.Context, claimID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetClaimKey(claimID))
}

// set claim
func (k Keeper) setClaimByRecipientIndex(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(claim.ID)
	store.Set(GetAccountClaimKey(claim.Recipient, claim.ID), bz)
}

func (k Keeper) removeClaimByRecipientIndex(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetIssuerClaimKey(claim.Recipient, claim.ID))
}

func (k Keeper) setClaimByIssuerIndex(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(claim.ID)
	store.Set(GetIssuerClaimKey(claim.Issuer, claim.ID), bz)
}

func (k Keeper) removeClaimByIssuerIndex(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetAccountClaimKey(claim.Issuer, claim.ID))
}

// GetClaim ...
func (k Keeper) GetClaim(ctx sdk.Context, claimID string) *Claim {
	claim := &Claim{}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetClaimKey(claimID))
	if b == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinary(b, claim)
	return claim
}

// Revoke ...
func (k Keeper) RevokeClaim(ctx sdk.Context, msg MsgRevokeClaim) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costRevokeClaim, "revokeClaim")
	claim := k.GetClaim(ctx, msg.ClaimID)

	if claim == nil {
		return nil, ErrClaimNotFound(msg.ClaimID)
	}

	if !bytes.Equal(claim.Issuer, msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("address %s not unauthorized to revoke", msg.Sender))
	}

	claim.Revocation = msg.Revocation
	k.setClaim(ctx, *claim)
	allTags := sdk.NewTags(
		"sender", []byte(msg.Sender.String()),
	)
	return allTags, nil
}

func (k Keeper) AnswerClaim(ctx sdk.Context, msg MsgAnswerClaim) (sdk.Tags, sdk.Error) {
	claim := k.GetClaim(ctx, msg.ClaimID)

	if claim == nil {
		return nil, ErrClaimNotFound(msg.ClaimID)
	}

	if claim.Paid == true {
		return nil, ErrClaimHasPaid(claim.ID)
	}

	if !bytes.Equal(claim.Recipient, msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("address %s not unauthorized to answer", msg.Sender))
	}
	allTags := sdk.NewTags(
		"sender", []byte(msg.Sender.String()),
	)
	if msg.Response == 0 {
		// reject the claim
		k.removeClaim(ctx, claim.ID)
		// remove index by account
		k.removeClaimByRecipientIndex(ctx, *claim)
		k.removeClaimByIssuerIndex(ctx, *claim)
	} else if len(claim.Fee) > 0 {
		// approve the claim
		_, tags, err := k.coinKeeper.SubtractCoins(ctx, msg.Sender, claim.Fee)
		if err != nil {
			return nil, err
		}
		_, tags2, err := k.coinKeeper.AddCoins(ctx, claim.Issuer, claim.Fee)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags).AppendTags(tags2)
	}
	claim.Paid = true
	k.setClaim(ctx, *claim)
	return allTags, nil
}
