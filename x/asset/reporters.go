package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reporter ...
type Reporter struct {
	Addr       sdk.AccAddress `json:"address"`
	Properties []string       `json:"properties"`
	Created    int64          `json:"created"`
}

// Reporters list all reporters
type Reporters []Reporter

// RevokeReporter delete reporter
func (k Keeper) RevokeReporter(ctx sdk.Context, msg MsgRevokeReporter) (sdk.Tags, sdk.Error) {
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}
	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to revoke", msg.Sender))
	}

	_, found = k.GetReporter(ctx, msg.AssetID, msg.Reporter)

	if !found {
		return nil, ErrInvalidRevokeReporter(msg.Reporter)
	}

	k.DeleteReporter(ctx, msg.AssetID, msg.Reporter)

	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
		TagRecipient, []byte(msg.Reporter.String()),
	)
	return tags, nil
}

func (k Keeper) setAssetByReporterIndex(ctx sdk.Context, reporter sdk.AccAddress, recordID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetReporterAssetKey(reporter, recordID), []byte{})
}

func (k Keeper) removeAssetByReporterIndex(ctx sdk.Context, reporter sdk.AccAddress, recordID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetReporterAssetKey(reporter, recordID))
}

// SetReporter ...
func (k Keeper) SetReporter(ctx sdk.Context, recordID string, reporter Reporter) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(reporter)
	store.Set(GetReporterKey(recordID, reporter.Addr), bz)
}

// DeleteReporter ...
func (k Keeper) DeleteReporter(ctx sdk.Context, recordID string, reporter sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetReporterKey(recordID, reporter))
}

// GetReporter ...
func (k Keeper) GetReporter(ctx sdk.Context, recordID string, addr sdk.AccAddress) (reporter Reporter, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetReporterKey(recordID, addr))
	if b == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinary(b, &reporter)
	return reporter, true
}

// DeleteReporters ...
func (k Keeper) DeleteReporters(ctx sdk.Context, recordID string) {
	store := ctx.KVStore(k.storeKey)

	// delete subspace
	iterator := sdk.KVStorePrefixIterator(store, GetReportersKey(recordID))
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	iterator.Close()
}

// GetReporters ...
func (k Keeper) GetReporters(ctx sdk.Context, recordID string) (reporters []Reporter) {
	store := ctx.KVStore(k.storeKey)

	// delete subspace
	iterator := sdk.KVStorePrefixIterator(store, GetReportersKey(recordID))
	for ; iterator.Valid(); iterator.Next() {
		reporter := Reporter{}
		k.cdc.MustUnmarshalBinary(iterator.Value(), &reporter)
		reporters = append(reporters, reporter)
	}
	iterator.Close()
	return
}
