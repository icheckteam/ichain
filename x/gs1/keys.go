package gs1

import (
	"crypto/md5"
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
	RecordIDKey      = []byte{0x13}
)

// GetRecordKey ...
func GetRecordKey(recordID int64) []byte {
	id := make([]byte, 8)
	binary.LittleEndian.PutUint64(id, uint64(recordID))
	return append(RecordKey, id...)
}

// GetLocationKey ...
func GetLocationKey(sender sdk.AccAddress, location Location) []byte {
	return GetKey(LocationKey, sender.Bytes(), []byte(location.ID))
}

// GetChildLocationKey ...
func GetChildLocationKey(sender sdk.AccAddress, location Location, child ChildrenLocation) []byte {
	return GetKey(ChildLocationKey, sender.Bytes(), []byte(location.ID), []byte(child.ID))
}

// GetActorKey ...
func GetActorKey(sender sdk.AccAddress, actor Actor) []byte {
	return GetKey(ActorKey, sender.Bytes(), actor.Addr.Bytes())
}

// GetProductKey ...
func GetProductKey(sender sdk.AccAddress, product Product) []byte {
	return GetKey(ProductKey, sender.Bytes(), []byte(product.ID))
}

// GetBatchKey ...
func GetBatchKey(sender sdk.AccAddress, batch Batch) []byte {
	return GetKey(BatchKey, sender.Bytes(), []byte(batch.ID), []byte(batch.ProductID))
}

// GetEventKey ...
func GetEventKey(sender sdk.AccAddress, event Event) []byte {
	timeB := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeB, uint64(event.Time.Unix()))
	return GetKey(EventKey, sender.Bytes(), timeB)
}

// GetEventBatchKey ...
func GetEventBatchKey(sender sdk.AccAddress, event Event, batchID string) []byte {
	return GetKey(EventBatch, GetEventKey(sender, event), []byte(batchID))
}

// GetBatchProductKey ...
func GetBatchProductKey(sender sdk.AccAddress, batch Vertice, productID string) []byte {
	return GetKey(BatchProduct, batch.Key, []byte(productID))
}

// GetIsKey ...
func GetIsKey(sender sdk.AccAddress, vertice Vertice, prefix []byte) []byte {
	return GetKey(IsKey, vertice.Key, prefix)
}

// GetOwnedByKey ...
func GetOwnedByKey(sender sdk.AccAddress, actor sdk.AccAddress, locationKey []byte) []byte {
	return GetKey(OwnedByKey, sender.Bytes(), actor.Bytes(), locationKey)
}

// GetAtKey ...
func GetAtKey(sender sdk.AccAddress, actor sdk.AccAddress, bizLocation BizLocation) []byte {
	return GetKey(AtKey, sender.Bytes(), actor.Bytes(), []byte(bizLocation.ID))
}

// GetReadPointKey ...
func GetReadPointKey(sender sdk.AccAddress, actor sdk.AccAddress, readPoint ReadPoint) []byte {
	return GetKey(ReadPointKey, sender.Bytes(), actor.Bytes(), []byte(readPoint.ID))
}

// GetKey ...
func GetKey(args ...[]byte) []byte {
	var key []byte
	for _, arg := range args {
		key = append(key, arg...)
	}
	sum := md5.Sum(key)
	return Key(sum[:])
}
