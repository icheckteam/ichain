package epcis

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
var (
	// Keys for store prefixes
	LocationKey      = []byte{0x00}
	ChildLocationKey = []byte{0x01}
	ActorKey         = []byte{0x02}
	ProductKey       = []byte{0x03}
	BatchKey         = []byte{0x04}
	EventKey         = []byte{0x05}
	EventBatch       = []byte{0x06}
	BatchProduct     = []byte{0x07}
	IsKey            = []byte{0x08}
	OwnedByKey       = []byte{0x09}
	AtKey            = []byte{0x10}
	ReadPointKey     = []byte{0x11}
	RecordKey        = []byte{0x12}
)

// GetRecordKey ...
func GetRecordKey(recordID string) []byte {
	return append(
		RecordKey,
		[]byte(recordID)...,
	)
}

// GetLocationKey ...
func GetLocationKey(sender sdk.AccAddress, location Location) []byte {
	return append(
		append(LocationKey, sender.Bytes()...),
		[]byte(location.ID)...,
	)
}

// GetChildLocationKey ...
func GetChildLocationKey(sender sdk.AccAddress, location Location, child ChildrenLocation) []byte {
	return append(
		append(ChildLocationKey, sender.Bytes()...),
		append(
			[]byte(location.ID),
			[]byte(child.ID)...,
		)...,
	)
}

// GetActorKey ...
func GetActorKey(sender sdk.AccAddress, actor Actor) []byte {
	return append(
		append(ActorKey, sender.Bytes()...),
		actor.Addr.Bytes()...,
	)
}

// GetProductKey ...
func GetProductKey(sender sdk.AccAddress, product Product) []byte {
	return append(
		append(ProductKey, sender.Bytes()...),
		[]byte(product.ID)...,
	)
}

// GetBatchKey ...
func GetBatchKey(sender sdk.AccAddress, batch Batch) []byte {
	return append(
		append(BatchKey, sender.Bytes()...),
		append(
			[]byte(batch.ID),
			[]byte(batch.ProductID)...,
		)...,
	)
}

// GetEventKey ...
func GetEventKey(sender sdk.AccAddress, event Event) []byte {
	timeB := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeB, uint64(event.Time.Unix()))
	return append(
		append(BatchKey, sender.Bytes()...),
		timeB...,
	)
}

// GetEventBatchKey ...
func GetEventBatchKey(sender sdk.AccAddress, event Event, batchID string) []byte {
	return append(
		GetEventKey(sender, event),
		append(
			EventBatch,
			append(
				sender.Bytes(),
				[]byte(batchID)...,
			)...,
		)...,
	)
}

// GetBatchProductKey ...
func GetBatchProductKey(sender sdk.AccAddress, batch Vertice, productID string) []byte {
	return append(
		BatchProduct,
		append(
			batch.Key,
			[]byte(productID)...,
		)...,
	)
}

// GetIsKey ...
func GetIsKey(sender sdk.AccAddress, vertice Vertice, prefix []byte) []byte {
	return append(
		IsKey,
		append(
			vertice.Key,
			prefix...,
		)...,
	)
}

// GetOwnedByKey ...
func GetOwnedByKey(sender sdk.AccAddress, actor sdk.AccAddress, locationKey []byte) []byte {
	return append(
		OwnedByKey,
		append(
			actor.Bytes(),
			locationKey...,
		)...,
	)
}

// GetAtKey ...
func GetAtKey(sender sdk.AccAddress, actor sdk.AccAddress, bizLocation BizLocation) []byte {
	return append(
		AtKey,
		append(
			actor.Bytes(),
			[]byte(bizLocation.ID)...,
		)...,
	)
}

// GetReadPointKey ...
func GetReadPointKey(sender sdk.AccAddress, actor sdk.AccAddress, readPoint ReadPoint) []byte {
	return append(
		ReadPointKey,
		append(
			actor.Bytes(),
			[]byte(readPoint.ID)...,
		)...,
	)
}
