package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/icheckteam/ichain/types"
)

const (
	costGetAsset         sdk.Gas = 10
	costRegisterAsset    sdk.Gas = 100
	costHasAsset         sdk.Gas = 10
	costSubtractCoin     sdk.Gas = 10
	costAddQuantity      sdk.Gas = 10
	costUpdateAttributes sdk.Gas = 10
	costCreateProposal   sdk.Gas = 10
	costRevokeProposal   sdk.Gas = 10
	costAnswerProposal   sdk.Gas = 10
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
func (k Keeper) CreateAsset(ctx sdk.Context, msg MsgCreateAsset) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costRegisterAsset, "registerAsset")
	if k.Has(ctx, msg.AssetID) {
		return nil, InvalidTransaction(fmt.Sprintf("Asset already exists: {%s}", msg.AssetID))
	}

	asset := Asset{
		ID:       msg.AssetID,
		Name:     msg.Name,
		Issuer:   msg.Issuer,
		Owner:    msg.Issuer,
		Quantity: msg.Quantity,
		Parent:   msg.Parent,
	}

	// update asset info
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"issuer", types.AddrToBytes(asset.Issuer),
		"owner", types.AddrToBytes(asset.Owner),
	)
	return tags, nil
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
	ctx.GasMeter().ConsumeGas(costHasAsset, "hasAsset")
	store := ctx.KVStore(k.storeKey)
	assetKey := GetAssetKey(assetID)
	return store.Has(assetKey)
}

// GetAsset get asset by IDS
func (k Keeper) GetAsset(ctx sdk.Context, assetID string) *Asset {
	ctx.GasMeter().ConsumeGas(costGetAsset, "getAsset")
	store := ctx.KVStore(k.storeKey)
	assetBytes := store.Get(GetAssetKey(assetID))
	asset := &Asset{}
	if err := k.cdc.UnmarshalBinary(assetBytes, asset); err != nil {
		return nil
	}
	return asset
}

// UpdateAttribute ...
func (k Keeper) UpdatePropertipes(ctx sdk.Context, msg MsgUpdatePropertipes) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costUpdateAttributes, "updateAttributes")
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return nil, ErrAssetNotFound(msg.ID)
	}
	// check role permissions
	for _, attr := range msg.Propertipes {
		authorized := asset.CheckUpdateAttributeAuthorization(msg.Issuer, attr)
		if !authorized {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
		}
	}

	// update all propertipes
	asset.Propertipes = asset.Propertipes.Adds(msg.Propertipes...)
	// save asset to store
	k.setAsset(ctx, *asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
	)
	return tags, nil
}

// AddQuantity ...
func (k Keeper) AddQuantity(ctx sdk.Context, msg AddQuantityMsg) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddQuantity, "addQuantity")
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return nil, ErrAssetNotFound(msg.ID)
	}
	if !asset.IsIssuer(msg.Issuer) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}
	asset.Quantity += msg.Quantity
	k.setAsset(ctx, *asset)
	tags := sdk.NewTags("asset_id", []byte(asset.ID))
	return tags, nil
}

// SubtractQuantity ...
func (k Keeper) SubtractQuantity(ctx sdk.Context, msg SubtractQuantityMsg) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costSubtractCoin, "subtractQuantity")
	asset := k.GetAsset(ctx, msg.ID)
	if asset == nil {
		return nil, ErrAssetNotFound(msg.ID)
	}
	if !asset.IsIssuer(msg.Issuer) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
	}

	if asset.Quantity < msg.Quantity {
		return nil, ErrInvalidAssetQuantity(asset.ID)
	}
	asset.Quantity -= msg.Quantity
	k.setAsset(ctx, *asset)
	tags := sdk.NewTags("asset_id", []byte(asset.ID))
	return tags, nil
}

// CreateProposal validates and adds a new proposal to the asset,
// or update a propsal if there already exists one for the recipient
func (k Keeper) CreateProposal(ctx sdk.Context, msg CreateProposalMsg) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costCreateProposal, "createProposal")
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
	tags := sdk.NewTags(
		"sender", types.AddrToBytes(msg.Issuer),
		"asset_id", []byte(asset.ID),
		"action", []byte("createProposal"),
	)
	return tags, nil
}

// RevokeProposal delete some properties from an existing proposal
// and will delete the proposal if there is no property left
func (k Keeper) RevokeProposal(ctx sdk.Context, msg RevokeProposalMsg) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costRevokeProposal, "revokeProposal")
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
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
	)
	return tags, nil
}

// AnswerProposal update the status of the proposal of the recipient if the answer is valid
func (k Keeper) AnswerProposal(ctx sdk.Context, msg AnswerProposalMsg) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAnswerProposal, "answerProposal")
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
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
	)
	return tags, nil
}
