package asset

import (
	"encoding/json"
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

func (msg Property) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(struct {
		Name         string          `json:"name"`
		Type         PropertyType    `json:"type"`
		BytesValue   []byte          `json:"bytes_value,omitempty"`
		StringValue  string          `json:"string_value,omitempty"`
		BooleanValue bool            `json:"boolean_value,omitempty"`
		NumberValue  int64           `json:"number_value,omitempty"`
		EnumValue    []string        `json:"enum_value,omitempty"`
		Location     json.RawMessage `json:"location_value,omitempty"`
	}{
		Name:         msg.Name,
		Type:         msg.Type,
		BytesValue:   msg.BytesValue,
		StringValue:  msg.StringValue,
		BooleanValue: msg.BooleanValue,
		NumberValue:  msg.NumberValue,
		EnumValue:    msg.EnumValue,
		Location:     msg.Location.GetSignBytes(),
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

type Location struct {
	Latitude  int64 `json:"latitude"`
	Longitude int64 `json:"longitude"`
}

func (msg Location) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(struct {
		Latitude  int64 `json:"latitude"`
		Longitude int64 `json:"longitude"`
	}{
		Latitude:  msg.Latitude,
		Longitude: msg.Longitude,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// list all properties
type Properties []Property

func (msg Properties) GetSignBytes() []json.RawMessage {
	props := []json.RawMessage{}
	for _, p := range msg {
		props = append(props, p.GetSignBytes())
	}
	return props
}

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
