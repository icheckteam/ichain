package gs1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Keeper ...
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
}

// CreateRecord create new record
func (k Keeper) CreateRecord(ctx sdk.Context, msg MsgSend) (sdk.Tags, sdk.Error) {
	locationEdges, locationVertices := k.getLocationGraph(msg.Sender, msg.Locations)
	k.setRecord(ctx, Record{
		Sender:   msg.Sender,
		Edges:    locationEdges,
		Vertices: locationVertices,
	})
	return nil, nil
}

func (k Keeper) getLocationGraph(sender sdk.AccAddress, locations []Location) ([]Edge, []Vertice) {
	totalLocation := 0
	for _, location := range locations {
		totalLocation++
		totalLocation += len(location.Children)
	}
	locationEdges := make([]Edge, totalLocation)
	locationVertices := make([]Vertice, totalLocation)

	for _, location := range locations {
		locationKey := GetLocationKey(sender, location)
		if len(location.Participant) > 0 {
			locationEdges = append(locationEdges, Edge{
				Key:    GetOwnedByKey(sender, location.Participant, locationKey),
				Source: locationKey,
				Target: location.Participant,
				Type:   EdgeTypeOwnerdBy,
			})
		}
		locationVertices = append(locationVertices, Vertice{
			Key:        locationKey,
			Type:       VertexTypeLocation,
			Attributes: location.Attributes,
		})

		for _, childLocation := range location.Children {
			childLocation.Attributes = []Attribute{
				{Name: "parent_id", Value: location.ID},
			}
			childLocationKey := GetKey(ChildLocationKey, sender.Bytes(), []byte(childLocation.ID), childLocation.Attributes.Bytes())
			locationVertices = append(locationVertices, Vertice{
				Key:        childLocationKey,
				Type:       VertexTypeLocation,
				Attributes: childLocation.Attributes,
			})

			locationEdges = append(locationEdges, Edge{
				Key:    GetKey(ChildLocationKey, sender.Bytes(), []byte(location.ID), []byte(childLocation.ID)),
				Source: childLocationKey,
				Target: locationKey,
				Type:   EdgeTypeChildLocation,
			})
		}
	}
	return locationEdges, locationVertices
}

func (k Keeper) setRecord(ctx sdk.Context, record Record) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(record)
	store.Set(GetRecordKey(record.ID), bz)
}
