package asset

import (
	"sort"
	"strings"
)

// Material defines the total material of new asset
type Material struct {
	AssetID  string `json:"asset_id"`
	Quantity int64  `json:"quantity"`
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
	return Material{material.AssetID, material.Quantity + materialB.Quantity}
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
			if materialA.Quantity+materialB.Quantity == 0 {
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
