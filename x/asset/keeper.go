package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/bank"
)

// Keeper ...
type Keeper struct {
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc      *wire.Codec
	bank     bank.Keeper
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, bank bank.Keeper) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
		bank:     bank,
	}
}

// Register register new asset
func (k Keeper) RegisterAsset(ctx sdk.Context, asset Asset) (sdk.Coins, sdk.Tags, sdk.Error) {
	if asset.ID == "icc" {
		return nil, nil, InvalidTransaction("Asset already exists")
	}

	if k.Has(ctx, asset.ID) {
		return nil, nil, InvalidTransaction("Asset already exists")
	}
	// update asset info
	k.setAsset(ctx, asset)

	// add coin ...
	return k.bank.AddCoins(ctx, asset.Issuer, sdk.Coins{
		sdk.Coin{Denom: asset.ID, Amount: asset.Quantity},
	})
}

func (k Keeper) setAsset(ctx sdk.Context, asset Asset) {
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey(asset.ID)

	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(asset)
	if err != nil {
		panic(err)
	}

	store.Set(assetKey, bz)
}

// Has asset
func (k Keeper) Has(ctx sdk.Context, assetID string) bool {
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey(assetID)
	return store.Has(assetKey)
}

// GetAsset get asset by IDS
func (k Keeper) GetAsset(ctx sdk.Context, assetID string) *Asset {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetAssetKey(assetID))
	asset := &Asset{}

	// marshal the record and add to the state
	if err := k.cdc.UnmarshalBinary(b, asset); err != nil {
		return nil
	}
	return asset
}

// UpdateAttribute ...
func (k Keeper) UpdateAttribute(ctx sdk.Context, msg UpdateAttrMsg) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return nil, ErrAssetNotFound(msg.ID)
	}

	for _, attr := range msg.Attributes {
		authorized := asset.CheckUpdateAttributeAuthorization(msg.Issuer, attr)
		if !authorized {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
		}
		setAttribute(asset, attr)
	}

	k.setAsset(ctx, *asset)
	allTags.AppendTag("owner", msg.Issuer.Bytes())
	allTags.AppendTag("asset_id", []byte(msg.ID))
	return allTags, nil
}

// AddQuantity ...
func (k Keeper) AddQuantity(ctx sdk.Context, msg AddQuantityMsg) (sdk.Coins, sdk.Tags, sdk.Error) {
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return nil, nil, ErrUnknownAsset("Asset not found")
	}
	if !asset.IsOwner(msg.Issuer) {
		return nil, nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}
	if len(msg.Materials) > 0 {
		coins := sdk.Coins{}
		for _, material := range msg.Materials {
			coins = append(coins, sdk.Coin{Denom: material.AssetID, Amount: material.Quantity})
		}
		_, _, err := k.bank.SubtractCoins(ctx, asset.Issuer, coins)
		if err != nil {
			return nil, nil, err
		}
		asset.Materials = asset.Materials.Plus(msg.Materials.Sort()).Sort()
	}
	asset.Quantity += msg.Quantity
	k.setAsset(ctx, *asset)
	// add coin ...
	return k.bank.AddCoins(ctx, asset.Issuer, sdk.Coins{
		sdk.Coin{Denom: asset.ID, Amount: msg.Quantity},
	})
}

// SubtractQuantity ...
func (k Keeper) SubtractQuantity(ctx sdk.Context, msg SubtractQuantityMsg) (sdk.Coins, sdk.Tags, sdk.Error) {
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return nil, nil, ErrUnknownAsset("Asset not found")
	}
	if !asset.IsOwner(msg.Issuer) {
		return nil, nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}

	// add coin ...
	coins, tags, err := k.bank.SubtractCoins(ctx, asset.Issuer, sdk.Coins{
		sdk.Coin{Denom: asset.ID, Amount: msg.Quantity},
	})

	if err != nil {
		return nil, nil, err
	}
	asset.Quantity -= msg.Quantity
	k.setAsset(ctx, *asset)
	return coins, tags, err
}

func setAttribute(a *Asset, attr Attribute) {
	for index, oldAttr := range a.Attributes {
		if oldAttr.Name == attr.Name {
			a.Attributes[index] = attr
			return
		}
	}
	a.Attributes = append(a.Attributes, attr)
}

// CreateProposal validates and adds a new proposal to the asset,
// or update a propsal if there already exists one for the recipient
func (k Keeper) CreateProposal(ctx sdk.Context, msg CreateProposalMsg) (sdk.Tags, sdk.Error) {
	switch msg.Role {
	case RoleOwner, RoleReporter:
		break
	default:
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Issuer))
	}

	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return nil, ErrUnknownAsset("Asset not found")
	}

	proposal, proposalIndex, authorized := asset.ValidatePropossal(msg.Issuer, msg.Recipient)
	if !authorized {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Issuer))
	}

	if proposal != nil {
		// Update proposal
		proposal.Role = msg.Role
		proposal.AddProperties(msg.Propertipes)
		asset.Proposals[proposalIndex] = *proposal
	} else {
		// Add new proposal
		proposal = &Proposal{
			Role:       msg.Role,
			Status:     StatusPending,
			Properties: msg.Propertipes,
			Issuer:     msg.Issuer,
			Recipient:  msg.Recipient,
		}
		asset.Proposals = append(asset.Proposals, *proposal)
	}

	k.setAsset(ctx, *asset)
	return nil, nil
}

// RevokeProposal delete some properties from an existing proposal
// and will delete the proposal if there is no property left
func (k Keeper) RevokeProposal(ctx sdk.Context, msg RevokeProposalMsg) (sdk.Tags, sdk.Error) {
	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return nil, ErrUnknownAsset("Asset not found")
	}

	proposal, proposalIndex, authorized := asset.ValidatePropossal(msg.Issuer, msg.Recipient)
	if !authorized {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to revoke", msg.Issuer))
	}

	if proposal == nil {
		return nil, ErrInvalidRevokeRecipient(msg.Recipient)
	}

	proposal.RemoveProperties(msg.Propertipes)

	if len(proposal.Properties) > 0 {
		// Update proposal
		asset.Proposals[proposalIndex] = *proposal
	} else {
		// Remove proposal
		i := proposalIndex
		asset.Proposals = append(asset.Proposals[:i], asset.Proposals[i+1:]...)
	}

	k.setAsset(ctx, *asset)
	return nil, nil
}

// AnswerProposal update the status of the proposal of the recipient if the answer is valid
func (k Keeper) AnswerProposal(ctx sdk.Context, msg AnswerProposalMsg) (sdk.Tags, sdk.Error) {
	asset := k.GetAsset(ctx, msg.AssetID)
	if asset == nil {
		return nil, ErrUnknownAsset("Asset not found")
	}

	proposal, proposalIndex, authorized := asset.ValidateProposalAnswer(msg.Recipient, ProposalStatus(msg.Response))

	if !authorized {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to answer", msg.Recipient))
	}

	proposal.Status = ProposalStatus(msg.Response)
	asset.Proposals[proposalIndex] = *proposal

	k.setAsset(ctx, *asset)
	return nil, nil
}
