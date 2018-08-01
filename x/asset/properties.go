package asset

import (
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

// ValidateBasic ...
func (p Property) ValidateBasic() sdk.Error {
	if p.Name == "" {
		return ErrMissingField("properties[$].name")
	}
	switch p.Type {
	case PropertyTypeBoolean,
		PropertyTypeBytes,
		PropertyTypeEnum,
		PropertyTypeLocation,
		PropertyTypeNumber,
		PropertyTypeString:
		break
	default:
		return ErrInvalidField("properties")
	}
	return nil
}

type Location struct {
	Latitude  int64 `json:"latitude"`
	Longitude int64 `json:"longitude"`
}

// Properties list all properties
type Properties []Property

// ValidateBasic ...
func (props Properties) ValidateBasic() sdk.Error {
	for _, p := range props {
		if err := p.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
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

// UpdateProperties ...
func (k Keeper) UpdateProperties(ctx sdk.Context, msg MsgUpdateProperties) (sdk.Tags, sdk.Error) {
	record, found := k.GetAsset(ctx, msg.AssetID)
	if !found {
		return nil, ErrAssetNotFound(record.ID)
	}

	if err := k.ValidateUpdateProperties(ctx, record, msg.Sender, msg.Properties); err != nil {
		return nil, err
	}

	k.SetProperties(ctx, msg.AssetID, msg.Properties)
	tags := sdk.NewTags(
		TagAsset, []byte(record.ID),
		TagSender, []byte(msg.Sender.String()),
	)
	return tags, nil
}

// SetProperty ....
func (k Keeper) SetProperty(ctx sdk.Context, recordID string, property Property) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(property)
	store.Set(GetPropertyKey(recordID, property.Name), bz)
}

// SetProperties ...
func (k Keeper) SetProperties(ctx sdk.Context, recordID string, props Properties) {
	for _, prop := range props {
		k.SetProperty(ctx, recordID, prop)
	}
}

// GetProperties ...
func (k Keeper) GetProperties(ctx sdk.Context, recordID string) (reporters Properties) {
	store := ctx.KVStore(k.storeKey)

	// delete subspace
	iterator := sdk.KVStorePrefixIterator(store, GetPropertiesKey(recordID))
	for ; iterator.Valid(); iterator.Next() {
		p := Property{}
		k.cdc.MustUnmarshalBinary(iterator.Value(), &p)
		reporters = append(reporters, p)
	}
	iterator.Close()
	return
}
