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

// NewKeeper ...
func NewKeeper(cdc *wire.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{storeKey, cdc}
}

// CreateRecord create new record
func (k Keeper) CreateRecord(ctx sdk.Context, msg MsgSend) (sdk.Tags, sdk.Error) {
	locationEdges, locationVertices := k.getLocationGraph(msg.Sender, msg.Locations)
	eventEdges, eventVertices := k.getEventGraph(msg.Sender, msg.Events)

	allEdges := append(
		locationEdges,
		eventEdges...,
	)

	allVertices := append(
		locationVertices,
		eventVertices...,
	)

	record := Record{
		ID:       k.getNewRecordID(ctx),
		Sender:   msg.Sender,
		Edges:    allEdges,
		Vertices: allVertices,
	}

	k.setRecord(ctx, record)
	return nil, nil
}

// GetRecord ...
func (k Keeper) GetRecord(ctx sdk.Context, recordID int64) (record *Record) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetRecordKey(recordID))
	if b == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(b, &record)
	return
}

func (k Keeper) getLocationGraph(sender sdk.AccAddress, locations []Location) ([]Edge, []Vertice) {
	totalLocation := 0
	for _, location := range locations {
		totalLocation++
		totalLocation += len(location.Children)
	}
	locationEdges := make([]Edge, totalLocation)
	locationVertices := make([]Vertice, totalLocation)
	var locationEdgesIndex = 0
	for index, location := range locations {
		locationKey := GetLocationKey(sender, location)

		locationVertices[index] = Vertice{
			Key:        locationKey,
			Type:       VertexTypeLocation,
			Attributes: location.Attributes,
		}

		if len(location.Participant) > 0 {
			locationEdgesIndex++
			locationEdges[index] = Edge{
				Key:    GetOwnedByKey(sender, location.Participant, locationKey),
				Source: locationKey,
				Target: GetKey(location.Participant.Bytes()),
				Type:   EdgeTypeOwnerdBy,
			}
		}

		for index2, childLocation := range location.Children {
			childLocation.Attributes = []Attribute{
				{Name: "parent_id", Value: location.ID},
			}
			childLocationKey := GetKey(ChildLocationKey, sender.Bytes(), []byte(childLocation.ID), childLocation.Attributes.Bytes())
			locationVertices[index+index2+1] = Vertice{
				Key:        childLocationKey,
				Type:       VertexTypeLocation,
				Attributes: childLocation.Attributes,
			}
			locationEdges[locationEdgesIndex] = Edge{
				Key:    GetKey(ChildLocationKey, sender.Bytes(), []byte(location.ID), []byte(childLocation.ID)),
				Source: childLocationKey,
				Target: locationKey,
				Type:   EdgeTypeChildLocation,
			}
			locationEdgesIndex++
		}
	}
	return locationEdges, locationVertices
}

func (k Keeper) getActorGraph(sender sdk.AccAddress, actors []Actor) ([]Edge, []Vertice) {
	actorEdges := make([]Edge, len(actors))
	actorVertice := make([]Vertice, len(actors))
	for index, actor := range actors {
		actorKey := GetKey(ActorKey, sender.Bytes(), []byte(actor.ID), actor.Attributes.Bytes())
		actorVertice[index] = Vertice{
			Key:        actorKey,
			Attributes: actor.Attributes,
			Type:       VertexTypeOwner,
		}
		actorEdges[index] = Edge{
			Key:    GetKey(IsKey, sender.Bytes(), actorKey),
			Source: actorKey,
			Target: ActorKey,
			Type:   "IS",
		}
	}
	return actorEdges, actorVertice
}

func (k Keeper) getProductGraph(sender sdk.AccAddress, products []Product) ([]Edge, []Vertice) {
	productEdges := make([]Edge, len(products))
	productVertices := make([]Vertice, len(products))
	for index, product := range products {
		productKey := GetKey(ActorKey, sender.Bytes(), []byte(product.ID), product.Attributes.Bytes())
		productVertices[index] = Vertice{
			Key:        productKey,
			Attributes: product.Attributes,
			Type:       VertexTypeProduct,
		}
		productEdges[index] = Edge{
			Key:    GetKey(IsKey, sender.Bytes(), productKey),
			Source: productKey,
			Target: ProductKey,
			Type:   "IS",
		}
	}
	return productEdges, productVertices
}

func (k Keeper) getBatchGraph(sender sdk.AccAddress, batches []Batch) ([]Edge, []Vertice) {
	batchEdges := make([]Edge, len(batches))
	batchVertices := make([]Vertice, len(batches))
	for index, batch := range batches {
		batchKey := GetKey(ActorKey, sender.Bytes(), []byte(batch.ID), batch.Attributes.Bytes())
		batchVertices[index] = Vertice{
			Key:        batchKey,
			Attributes: batch.Attributes,
			Type:       VertexTypeBatch,
		}
		batchEdges[index] = Edge{
			Key:    GetKey(IsKey, sender.Bytes(), batchKey),
			Source: batchKey,
			Target: BatchKey,
			Type:   "IS",
		}
	}
	return batchEdges, batchVertices
}

func (k Keeper) getEventGraph(sender sdk.AccAddress, events []Event) ([]Edge, []Vertice) {
	eventEdges := []Edge{}
	eventVertices := []Vertice{}
	for _, event := range events {
		eventKey := GetEventKey(sender, event)

		eventVertices = append(eventVertices, Vertice{
			Key:  eventKey,
			Type: VertexTypeEvent,
		})

		if len(event.BizLocation.ID) > 0 {
			eventEdges = append(eventEdges, Edge{
				Key:    GetKey(AtKey, sender.Bytes(), eventKey, []byte(event.BizLocation.ID)),
				Source: eventKey,
				Target: []byte(event.BizLocation.ID),
				Type:   EdgeTypeAt,
			})
		}
		if len(event.ReadPoint.ID) > 0 {
			eventEdges = append(eventEdges, Edge{
				Key:    GetKey(ReadPointKey, sender.Bytes(), eventKey, []byte(event.ReadPoint.ID)),
				Source: eventKey,
				Target: []byte(event.ReadPoint.ID),
				Type:   EdgeTypeReadPoint,
			})
		}
		if len(event.EpcList) > 0 {
			for _, epc := range event.EpcList {
				eventEdges = append(eventEdges, Edge{
					Key:    GetKey(EventBatch, sender.Bytes(), eventKey, []byte(epc)),
					Source: eventKey,
					Target: []byte(epc),
					Type:   EdgeTypeEventBatch,
				})
				eventEdges = append(eventEdges, Edge{
					Key:    GetKey(EventBatch, sender.Bytes(), []byte(epc), eventKey),
					Source: []byte(epc),
					Target: eventKey,
					Type:   EdgeTypeEventBatch,
				})
			}
		}

		if len(event.ChildEPCs) > 0 {
			var edgeType string
			if event.Action == ActionTypeAdd {
				edgeType = EdgeTypeAddedBatch
			} else if event.Action == ActionTypeDelete {
				edgeType = EdgeTypeRemovedBatch
			}
			for _, epc := range event.ChildEPCs {
				eventEdges = append(eventEdges, Edge{
					Key:    GetKey(EventBatch, sender.Bytes(), eventKey, []byte(epc)),
					Source: eventKey,
					Type:   edgeType,
					Target: []byte(epc),
				})
			}
		}

		if len(event.ParentID) > 0 {
			eventEdges = append(eventEdges, Edge{
				Key:    GetKey(EventBatch, sender.Bytes(), eventKey, []byte(event.ParentID)),
				Source: eventKey,
				Type:   EdgeTypePallet,
				Target: []byte(event.ParentID),
			})
			eventEdges = append(eventEdges, Edge{
				Key:    GetKey(EventBatch, sender.Bytes(), []byte(event.ParentID), eventKey),
				Source: []byte(event.ParentID),
				Target: eventKey,
				Type:   EdgeTypePallet,
			})
		}

		if len(event.InputEPCList) > 0 {
			for _, epc := range event.InputEPCList {
				eventEdges = append(eventEdges, Edge{
					Key:    GetKey(EventBatch, sender.Bytes(), eventKey, []byte(epc)),
					Source: eventKey,
					Type:   EdgeTypeInputBatch,
					Target: []byte(epc),
				})
			}
		}

		if len(event.OutputEPCList) > 0 {
			for _, epc := range event.OutputEPCList {
				eventEdges = append(eventEdges, Edge{
					Key:    GetKey(EventBatch, sender.Bytes(), eventKey, []byte(epc)),
					Source: eventKey,
					Type:   EdgeTypeOutputBatch,
					Target: []byte(epc),
				})
				eventEdges = append(eventEdges, Edge{
					Key:    GetKey(EventBatch, sender.Bytes(), []byte(epc), eventKey),
					Source: []byte(epc),
					Target: eventKey,
					Type:   EdgeTypeOutputBatch,
				})
			}
		}
	}
	return eventEdges, eventVertices
}

func (k Keeper) setRecord(ctx sdk.Context, record Record) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(record)
	store.Set(GetRecordKey(record.ID), bz)
}

// getNewRecordID ...
func (k Keeper) getNewRecordID(ctx sdk.Context) (number int64) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(RecordIDKey)
	if b == nil {
		number = 0
	} else {
		k.cdc.MustUnmarshalBinary(b, &number)
	}
	number++
	bz := k.cdc.MustMarshalBinary(number)
	store.Set(RecordIDKey, bz)
	return
}
