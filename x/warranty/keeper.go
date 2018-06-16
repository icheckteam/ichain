package warranty

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/icheckteam/ichain/types"
)

// Keeper manages contracts
type Keeper struct {
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc      *wire.Codec
	bank     bank.Keeper
}

// NewKeeper returns the keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, bank bank.Keeper) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
		bank:     bank,
	}
}

// CreateContract create new a contract
func (k Keeper) CreateContract(ctx sdk.Context, msg MsgCreateContract) sdk.Error {
	c := k.GetContract(ctx, msg.ID)
	if c != nil {
		return types.InvalidTransaction(DefaultCodespace, "Contract already exitsts")
	}

	// subtract coins of the issuer
	coins := sdk.Coins{sdk.Coin{Denom: msg.AssetID, Amount: 1}}
	_, _, err := k.bank.SubtractCoins(ctx, msg.Issuer, coins)
	if err != nil {
		return err
	}

	// save contract to db
	k.setContract(ctx, Contract{
		ID:        msg.ID,
		AssetID:   msg.AssetID,
		Expires:   msg.Expires,
		Issuer:    msg.Issuer,
		Serial:    msg.Serial,
		Recipient: msg.Recipient,
	})
	return nil
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
