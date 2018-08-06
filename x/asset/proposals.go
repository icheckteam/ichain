package asset

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//--------------------------------------------------

// Proposal is an invitation to manage an asset
type Proposal struct {
	Role       ProposalRole   `json:"role"`       // The role assigned to the recipient
	Status     ProposalStatus `json:"status"`     // The response of the recipient
	Properties []string       `json:"properties"` // The asset's attributes name that the recipient is authorized to update
	Issuer     sdk.AccAddress `json:"issuer"`     // The proposal issuer
	Recipient  sdk.AccAddress `json:"recipient"`  // The recipient of the proposal
}

// ValidateAnswer ...
func (p Proposal) ValidateAnswer(msg MsgAnswerProposal) sdk.Error {
	if p.Role != msg.Role {
		return ErrInvalidRole("")
	}

	switch msg.Response {
	case StatusCancel:
		if !bytes.Equal(msg.Sender, p.Issuer) {
			return sdk.ErrUnauthorized("Only the issuing can cancel a proposal")
		}
		return nil
	case StatusRejected, StatusAccepted:
		if !bytes.Equal(msg.Sender, p.Recipient) {
			return sdk.ErrUnauthorized("Only the recipient can rejected/accepted a proposal")
		}
		return nil
	default:
		return ErrInvalidField("response")
	}
}

// Proposals is a sclice of Proposal
type Proposals []Proposal

// ProposalRole defines the authority of the proposal's recipient
type ProposalRole int

const (
	// RoleReporter is authorized to update the asset's attributes
	// whose name is included in the proposal's properties field
	RoleReporter ProposalRole = iota + 1

	// RoleOwner has the same authorization as RoleReporter
	// but also authorized to make proposal to other recipient
	RoleOwner
)

// ProposalStatus define the status of the proposal
type ProposalStatus int

// All available status of the proposal
const (
	StatusPending  ProposalStatus = iota // The recipient has not answered
	StatusAccepted                       // The recipient accepted the proposal
	StatusCancel                         // The issuer cancel the proposal
	StatusRejected                       // the recipient reject the proposal
)

// SetProposal ...
func (k Keeper) SetProposal(ctx sdk.Context, assetID string, proposal Proposal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(proposal)
	store.Set(GetProposalKey(assetID, proposal.Recipient), bz)
}

// DeleteProposal ...
func (k Keeper) DeleteProposal(ctx sdk.Context, assetID string, recipient sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetProposalKey(assetID, recipient))
}

// GetProposal ...
func (k Keeper) GetProposal(ctx sdk.Context, assetID string, recipient sdk.AccAddress) (proposal Proposal, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetProposalKey(assetID, recipient))
	if b == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(b, &proposal)
	found = true
	return
}

// AddProposal ...
func (k Keeper) AddProposal(ctx sdk.Context, msg MsgCreateProposal) (sdk.Tags, sdk.Error) {
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(asset.ID)
	}

	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add", msg.Sender))
	}

	proposal := Proposal{
		Role:       msg.Role,
		Status:     StatusPending,
		Properties: msg.Properties,
		Issuer:     msg.Sender,
		Recipient:  msg.Recipient,
	}
	k.SetProposal(ctx, asset.ID, proposal)
	k.setProposalAccountIndex(ctx, msg.Recipient, asset.ID)
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagRecipient, []byte(msg.Recipient.String()),
		TagSender, []byte(msg.Sender.String()),
	)

	return tags, nil
}

// AnswerProposal ...
func (k Keeper) AnswerProposal(ctx sdk.Context, msg MsgAnswerProposal) (sdk.Tags, sdk.Error) {
	proposal, found := k.GetProposal(ctx, msg.AssetID, msg.Recipient)
	if !found {
		return nil, ErrProposalNotFound(msg.Recipient)
	}
	// validate answer msg
	if err := proposal.ValidateAnswer(msg); err != nil {
		return nil, err
	}
	// delete proposal
	k.DeleteProposal(ctx, msg.AssetID, proposal.Recipient)
	k.removeProposalAccountIndex(ctx, msg.Recipient, msg.AssetID)
	asset, _ := k.GetAsset(ctx, msg.AssetID)
	if !asset.IsOwner(proposal.Issuer) {
		// Only delete the proposal
		return nil, nil
	}

	if msg.Response == StatusAccepted {
		switch proposal.Role {
		case RoleOwner:
			// update owner
			asset.Owner = proposal.Recipient
			k.DeleteReporters(ctx, asset.ID)
			k.setAssetByAccountIndex(ctx, asset.ID, proposal.Recipient)
			break
		case RoleReporter:
			k.SetReporter(ctx, asset.ID, Reporter{
				Properties: proposal.Properties,
				Created:    ctx.BlockHeader().Time,
				Addr:       proposal.Recipient,
			})
			k.setAssetByReporterIndex(ctx, proposal.Recipient, asset.ID)
			break
		default:
			break
		}
		k.setAsset(ctx, asset)
	}

	tags := sdk.NewTags(
		TagAsset, []byte(msg.AssetID),
		TagSender, []byte(msg.Recipient.String()),
	)
	return tags, nil
}

func (k Keeper) setProposalAccountIndex(ctx sdk.Context, addr sdk.AccAddress, recordID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetProposalAccountKey(addr, recordID), []byte{})
}
func (k Keeper) removeProposalAccountIndex(ctx sdk.Context, addr sdk.AccAddress, recordID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetProposalAccountKey(addr, recordID))
}
