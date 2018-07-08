package asset

import (
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Material defines the total material of new asset
type Material struct {
	AssetID  string  `json:"asset_id"`
	Quantity sdk.Int `json:"quantity"`
}

// Materials - list of materials
type Materials []Material

// SameDenomAs returns true if the two assets are the same asset
func (material Material) SameAssetAs(other Material) bool {
	return (material.AssetID == other.AssetID)
}

// Adds quantities of two assets with same asset
func (material Material) Plus(materialB Material) Material {
	if !material.SameAssetAs(materialB) {
		return material
	}
	return Material{material.AssetID, material.Quantity.Add(materialB.Quantity)}
}

// Plus combines two sets of materials
// CONTRACT: Plus will never return materials where one Material has a 0 quantity.
func (materials Materials) Plus(materialsB Materials) Materials {
	sum := ([]Material)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(materials), len(materialsB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, materialsB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, materials[indexA:]...)
		}
		materialA, materialB := materials[indexA], materialsB[indexB]
		switch strings.Compare(materialA.AssetID, materialB.AssetID) {
		case -1:
			sum = append(sum, materialA)
			indexA++
		case 0:
			if materialA.Quantity.Add(materialB.Quantity).IsZero() {
				// ignore 0 sum coin type
			} else {
				sum = append(sum, materialA.Plus(materialB))
			}
			indexA++
			indexB++
		case 1:
			sum = append(sum, materialB)
			indexB++
		}
	}
}

//----------------------------------------
// Sort interface

//nolint
func (materials Materials) Len() int           { return len(materials) }
func (materials Materials) Less(i, j int) bool { return materials[i].AssetID < materials[j].AssetID }
func (materials Materials) Swap(i, j int)      { materials[i], materials[j] = materials[j], materials[i] }

var _ sort.Interface = Materials{}

// Sort is a helper function to sort the set of materials inplace
func (materials Materials) Sort() Materials {
	sort.Sort(materials)
	return materials
}

// AddMaterials add materials to the asset
func (k Keeper) AddMaterials(ctx sdk.Context, msg MsgAddMaterials) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddMaterials, "addMaterials")
	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(msg.AssetID)
	}

	if asset.Final {
		return nil, ErrInvalidTransaction(fmt.Sprintf("Asset {%s} already final", asset.ID))
	}

	authorized := asset.CheckUpdateAttributeAuthorization(msg.Sender, Property{Name: "materials"})
	if !authorized {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add", msg.Sender))
	}
	// subtract quantity
	materialsToSave := []Asset{}
	for _, material := range msg.Materials {
		m, found := k.GetAsset(ctx, material.AssetID)
		if !found {
			return nil, ErrAssetNotFound(m.ID)
		}
		if m.Final {
			return nil, ErrAssetAlreadyFinal(m.ID)
		}
		if !m.IsOwner(msg.Sender) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to add materials", msg.Sender))
		}
		if m.Quantity.LT(material.Quantity) {
			return nil, ErrInvalidAssetQuantity(m.ID)
		}

		m.Quantity = m.Quantity.Sub(material.Quantity)
		materialsToSave = append(materialsToSave, m)
	}
	msg.Materials = msg.Materials.Sort()
	asset.Materials = asset.Materials.Plus(msg.Materials)
	materialsToSave = append(materialsToSave, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"sender", []byte(msg.Sender.String()),
	)
	for _, meterialToSave := range materialsToSave {
		k.setAsset(ctx, meterialToSave)
		tags = tags.AppendTag("asset_id", []byte(meterialToSave.ID))
	}

	return tags, nil
}
