package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Reporter struct {
	Addr       sdk.Address `json:"address"`
	Properties []string    `json:"properties"`
	Created    int64       `json:"created"`
}

type Reporters []Reporter

// CreateReporter validates and adds a new reporter to the asset,
// or update a reporter if there already exists one for the reporter
func (k Keeper) CreateReporter(ctx sdk.Context, msg MsgCreateReporter) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costCreateReporter, "createReporter")
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}
	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Sender))
	}

	reporter, reporterIndex := asset.GetReporter(msg.Reporter)
	if reporter != nil {
		// Update reporter
		reporter.Properties = msg.Properties
		reporter.Created = ctx.BlockHeader().Time
		asset.Reporters[reporterIndex] = *reporter
	} else {
		// Add new reporter
		reporter = &Reporter{
			Addr:       msg.Reporter,
			Properties: msg.Properties,
			Created:    ctx.BlockHeader().Time,
		}
		asset.Reporters = append(asset.Reporters, *reporter)
	}

	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"sender", []byte(msg.Sender.String()),
		"recipient", []byte(msg.Reporter.String()),
	)
	return tags, nil
}

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
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Sender))
	}

	reporter, reporterIndex := asset.GetReporter(msg.Reporter)

	if reporter == nil {
		return nil, ErrInvalidRevokeReporter(msg.Reporter)
	}

	asset.Reporters = append(asset.Reporters[:reporterIndex], asset.Reporters[reporterIndex+1:]...)

	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"sender", []byte(msg.Sender.String()),
		"recipient", []byte(msg.Reporter.String()),
	)
	return tags, nil
}
