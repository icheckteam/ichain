package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddMaterials add materials to the asset
func (k Keeper) AddMaterials(ctx sdk.Context, msg MsgAddMaterials) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddMaterials, "addMaterials")
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}

	if err := asset.ValidateAddMaterial(msg.Sender); err != nil {
		return nil, err
	}
	// subtract quantity
	materialsToSave := []Asset{}
	for _, amount := range msg.Amount {
		m, found := k.GetAsset(ctx, amount.Denom)
		if !found {
			return nil, ErrAssetNotFound(amount.Denom)
		}

		if err := m.ValidateSubtractQuantity(msg.Sender, amount.Amount); err != nil {
			return nil, err
		}

		m.Quantity = m.Quantity.Sub(amount.Amount)
		materialsToSave = append(materialsToSave, m)
	}
	asset.Materials = asset.Materials.Plus(msg.Amount.Sort())
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
	)
	for index, meterialToSave := range materialsToSave {
		// subtract inventory
		k.subtractInventory(ctx, asset.Owner, sdk.Coin{
			Denom:  asset.GetRoot(),
			Amount: msg.Amount[index].Amount,
		})

		k.setAsset(ctx, meterialToSave)
		tags = tags.AppendTag(TagAsset, []byte(meterialToSave.ID))
	}
	k.setAsset(ctx, asset)
	tags = tags.AppendTag(TagAsset, []byte(asset.ID))
	return tags, nil
}
