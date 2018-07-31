package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Reporter struct {
	Addr       sdk.AccAddress `json:"address"`
	Properties []string       `json:"properties"`
	Created    int64          `json:"created"`
}

type Reporters []Reporter

// RevokeReporter delete reporter
func (k Keeper) RevokeReporter(ctx sdk.Context, msg MsgRevokeReporter) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costRevokeReporter, "revokeReporter")
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

	reporter, reporterIndex := asset.GetReporter(msg.Reporter)

	if reporter == nil {
		return nil, ErrInvalidRevokeReporter(msg.Reporter)
	}
	k.removeAssetByReporterIndex(ctx, reporter.Addr, asset.ID)
	asset.Reporters = append(asset.Reporters[:reporterIndex], asset.Reporters[reporterIndex+1:]...)

	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
		TagRecipient, []byte(msg.Reporter.String()),
	)
	return tags, nil
}

func (k Keeper) setAssetByReporterIndex(ctx sdk.Context, reporter sdk.AccAddress, assetId string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(assetId)
	store.Set(GetReporterAssetKey(reporter, assetId), bz)
}

func (k Keeper) removeAssetByReporterIndex(ctx sdk.Context, reporter sdk.AccAddress, assetId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetReporterAssetKey(reporter, assetId))
}
