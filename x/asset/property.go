package asset

import (
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Property property of the asset
type Property struct {
	Name         string       `json:"name"`
	Type         PropertyType `json:"type"`
	BytesValue   []byte       `json:"bytes_value"`
	StringValue  string       `json:"string_value"`
	BooleanValue bool         `json:"boolean_value"`
	NumberValue  int64        `json:"number_value"`
	EnumValue    []string     `json:"enum_value"`
	Location     Location     `json:"location_value"`
	Precision    int          `json:"precision"`
}

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// list all properties
type Properties []Property

func (propertiesA Properties) Adds(othersB ...Property) Properties {
	sum := ([]Property)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(propertiesA), len(othersB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, othersB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, propertiesA[indexA:]...)
		}
		propertyA, propertyB := propertiesA[indexA], othersB[indexB]
		switch strings.Compare(propertyA.Name, propertyB.Name) {
		case -1:
			sum = append(sum, propertyA)
			indexA++
		case 0:
			sum = append(sum, propertyB)
			indexA++
			indexB++
		case 1:
			indexB++
			sum = append(sum, propertyB)
		}
	}
}

//----------------------------------------
// Sort interface

//nolint
func (properties Properties) Len() int           { return len(properties) }
func (properties Properties) Less(i, j int) bool { return properties[i].Name < properties[j].Name }
func (properties Properties) Swap(i, j int) {
	properties[i], properties[j] = properties[j], properties[i]
}

var _ sort.Interface = Properties{}

// Sort is a helper function to sort the set of materials inplace
func (properties Properties) Sort() Properties {
	sort.Sort(properties)
	return properties
}

// PropertyType define the type of the property
type PropertyType int

// All avaliable type ò the attribute
const (
	PropertyTypeBytes PropertyType = iota + 1
	PropertyTypeString
	PropertyTypeBoolean
	PropertyTypeNumber
	PropertyTypeEnum
	PropertyTypeLocation
)

// UpdateAttribute ...
func (k Keeper) UpdateProperties(ctx sdk.Context, msg MsgUpdateProperties) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costUpdateProperties, "updateProperties")

	asset, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(asset.ID)
	}
	if asset.Final {
		return nil, ErrAssetAlreadyFinal(asset.ID)
	}

	// check role permissions
	for _, attr := range msg.Properties {
		authorized := asset.CheckUpdateAttributeAuthorization(msg.Issuer, attr)
		if !authorized {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to transfer", msg.Issuer))
		}
	}

	// update all Properties
	msg.Properties = msg.Properties.Sort()
	asset.Properties = asset.Properties.Adds(msg.Properties...)
	// save asset to store
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		"asset_id", []byte(asset.ID),
		"sender", []byte(msg.Issuer.String()),
	)
	return tags, nil
}
