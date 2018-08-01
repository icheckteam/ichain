package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddMaterials add materials to the asset
func (k Keeper) AddMaterials(ctx sdk.Context, msg MsgAddMaterials) (sdk.Tags, sdk.Error) {
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}

	if !asset.IsOwner(msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add", msg.Sender))
	}

	// validate material amount
	cached := map[string]Asset{}
	for _, amount := range msg.Amount {
		m, found := k.GetAsset(ctx, amount.Denom)
		if !found {
			return nil, ErrAssetNotFound(amount.Denom)
		}
		if !m.IsOwner(msg.Sender) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add", msg.Sender))
		}

		cached[m.ID] = m
	}

	// new tags ...
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
	)

	// update record and material
	for _, amount := range msg.Amount {
		record := cached[amount.Denom]
		record.Quantity = record.Quantity.Sub(amount.Amount)
		k.SetAsset(ctx, record)
		k.AddMaterial(ctx, msg.AssetID, Material{Amount: amount.Amount, RecordID: amount.Denom})
		tags = tags.AppendTag(TagAsset, []byte(amount.Denom))
	}

	return tags, nil
}

// Material ...
type Material struct {
	RecordID string
	Amount   sdk.Int
}

// SetMaterial ...
func (k Keeper) setMaterial(ctx sdk.Context, recordID string, material Material) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(material)
	store.Set(GetMaterialKey(recordID, material.RecordID), bz)
}

// AddMaterial ...
func (k Keeper) AddMaterial(ctx sdk.Context, recordID string, input Material) {
	material, found := k.GetMaterial(ctx, recordID, input.RecordID)
	if !found {
		material = Material{
			RecordID: input.RecordID,
			Amount:   input.Amount,
		}
	} else {
		material.Amount = material.Amount.Add(input.Amount)
	}
	k.setMaterial(ctx, recordID, material)

}

// GetMaterial ...
func (k Keeper) GetMaterial(ctx sdk.Context, recordID string, materialID string) (material Material, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetMaterialKey(recordID, materialID))
	if b == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinary(b, &material)
	return material, true
}

// GetMaterials ...
func (k Keeper) GetMaterials(ctx sdk.Context, recordID string) (materials []Material) {
	store := ctx.KVStore(k.storeKey)
	materialsPrefixKey := GetMaterialsKey(recordID)
	iterator := sdk.KVStorePrefixIterator(store, materialsPrefixKey)
	i := 0
	for ; ; i++ {
		if !iterator.Valid() {
			break
		}
		material := Material{}
		k.cdc.MustUnmarshalBinary(iterator.Value(), &material)
		materials = append(materials, material)
		iterator.Next()
	}

	iterator.Close()
	return materials
}
