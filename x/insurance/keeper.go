package insurance

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/types"
	"github.com/icheckteam/ichain/x/asset"
)

// Keeper manages contracts
type Keeper struct {
	storeKey    sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc         *wire.Codec
	assetKeeper asset.Keeper
}

// NewKeeper returns the keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, assetKeeper asset.Keeper) Keeper {
	return Keeper{
		storeKey:    key,
		cdc:         cdc,
		assetKeeper: assetKeeper,
	}
}

// CreateContract create new a contract
func (k Keeper) CreateContract(ctx sdk.Context, msg MsgCreateContract) (sdk.Tags, sdk.Error) {
	if k.hasContract(ctx, msg.ID) {
		return nil, types.InvalidTransaction(DefaultCodespace, "Contract already exitsts")
	}

	a, found := k.assetKeeper.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, asset.ErrAssetNotFound(msg.AssetID)
	}
	if a.Final {
		return nil, asset.ErrAssetAlreadyFinal(a.ID)
	}

	if !a.IsOwner(msg.Issuer) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Issuer))
	}

	c := Contract{
		ID:        msg.ID,
		AssetID:   msg.AssetID,
		Expires:   msg.Expires,
		Issuer:    msg.Issuer,
		Serial:    msg.Serial,
		Recipient: msg.Recipient,
	}

	// save contract to db
	k.setContract(ctx, c)
	k.setContractByAccountIndex(ctx, c.Issuer, c.ID)
	k.setContractByAccountIndex(ctx, c.Recipient, c.ID)
	a.Final = true
	k.assetKeeper.SetAsset(ctx, a)

	tags := sdk.NewTags(
		"asset_id", []byte(msg.AssetID),
		"sender", []byte(msg.Issuer.String()),
		"recipient", []byte(msg.Recipient.String()),
	)

	return tags, nil
}

// CreateClaim create new a claim
func (k Keeper) CreateClaim(ctx sdk.Context, msg MsgCreateClaim) sdk.Error {
	c := k.GetContract(ctx, msg.ContractID)
	if c == nil {
		return types.InvalidTransaction(DefaultCodespace, "Contract not found")
	}

	if c.ValidateCreateClaim(msg.Issuer) == false {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create claim", msg.Issuer))
	}

	c.Claim = &Claim{
		Status:    ClaimStatusPending,
		Recipient: msg.Recipient,
	}
	k.setContract(ctx, *c)
	k.setContractByAccountIndex(ctx, msg.Recipient, c.ID)
	return nil
}

// ProcessClaim process claim
func (k Keeper) ProcessClaim(ctx sdk.Context, msg MsgProcessClaim) sdk.Error {
	c := k.GetContract(ctx, msg.ContractID)
	if c == nil {
		return types.InvalidTransaction(DefaultCodespace, "Contract not found")
	}

	if c.Claim == nil {
		return types.InvalidTransaction(DefaultCodespace, "Claim not found")
	}

	if !c.ValidateClaimProcess(msg.Issuer, msg.Status) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create claim", msg.Issuer))
	}
	c.Claim.Status = msg.Status
	k.setContract(ctx, *c)

	switch c.Claim.Status {
	case ClaimStatusClaimRepair, ClaimStatusReimbursement, ClaimStatusRejected, ClaimStatusTheftConfirmed:
		k.removeContractByAccountIndex(ctx, msg.Issuer, c.ID)
		break
	default:
		break
	}

	return nil
}

func (k Keeper) setContract(ctx sdk.Context, c Contract) {
	store := ctx.KVStore(k.storeKey)
	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(c)
	if err != nil {
		panic(err)
	}
	store.Set(GetContractKey(c.ID), bz)
}

func (k Keeper) setContractByAccountIndex(ctx sdk.Context, addr sdk.Address, contractID string) {
	store := ctx.KVStore(k.storeKey)
	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(contractID)
	if err != nil {
		panic(err)
	}
	store.Set(GetAccountContractKey(addr, contractID), bz)
}

func (k Keeper) removeContractByAccountIndex(ctx sdk.Context, addr sdk.Address, contractID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetAccountContractKey(addr, contractID))
}

// GetContract get contract by ID
func (k Keeper) GetContract(ctx sdk.Context, contractID string) *Contract {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetContractKey(contractID))
	c := &Contract{}

	if err := k.cdc.UnmarshalBinary(b, c); err != nil {
		return nil
	}
	return c
}

// hasContract
func (k Keeper) hasContract(ctx sdk.Context, contractID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetContractKey(contractID))
}
