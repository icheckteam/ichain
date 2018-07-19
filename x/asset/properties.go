package asset

import (
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Property property of the asset
type Property struct {
	Name         string       `json:"name"`
	Type         PropertyType `json:"type"`
	BytesValue   []byte       `json:"bytes_value,omitempty"`
	StringValue  string       `json:"string_value,omitempty"`
	BooleanValue bool         `json:"boolean_value,omitempty"`
	NumberValue  int64        `json:"number_value,omitempty"`
	EnumValue    []string     `json:"enum_value,omitempty"`
	Location     Location     `json:"location_value,omitempty"`
}

func PropertyTypeToString(t PropertyType) string {
	switch t {
	case PropertyTypeBoolean:
		return "boolean"
	case PropertyTypeBytes:
		return "bytes"
	case PropertyTypeEnum:
		return "enum"
	case PropertyTypeLocation:
		return "location"
	case PropertyTypeNumber:
		return "number"
	case PropertyTypeString:
		return "string"
	default:
		return "Unknown"
	}
}

func (p Property) GetValue() interface{} {
	switch p.Type {
	case PropertyTypeBoolean:
		return p.BooleanValue
	case PropertyTypeBytes:
		return p.BytesValue
	case PropertyTypeEnum:
		return p.EnumValue
	case PropertyTypeLocation:
		return p.Location
	case PropertyTypeNumber:
		return p.NumberValue
	case PropertyTypeString:
		return p.StringValue
	default:
		return "Unknown"
	}
}

type Location struct {
	Latitude  int64 `json:"latitude"`
	Longitude int64 `json:"longitude"`
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

	if err := asset.ValidateUpdateProperties(msg.Sender, msg.Properties); err != nil {
		return nil, err
	}

	// update all Properties
	msg.Properties = msg.Properties.Sort()
	asset.Properties = asset.Properties.Adds(msg.Properties...)
	// save asset to store
	k.setAsset(ctx, asset)
	tags := sdk.NewTags(
		TagAsset, []byte(asset.ID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}
