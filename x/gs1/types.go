package gs1

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
const (
	// edge
	EdgeTypeOwnerdBy      = "OWNED_BY"
	EdgeTypeChildLocation = "CHILD_LOCATION"
	EdgeTypeLocation      = "LOCATION"
	EdgeTypeAddedBatch    = "ADDED_BATCH"
	EdgeTypeRemovedBatch  = "REMOVED_BATCH"
	EdgeTypePallet        = "PALLET"
	EdgeTypeOutputBatch   = "OUTPUT_BATCH"
	EdgeTypeAt            = "AT"
	EdgeTypeReadPoint     = "READ_POINT"
	EdgeTypeEventBatch    = "EVENT_BATCH"
	EdgeTypeInputBatch    = "INPUT_BATCH"

	// VertexType
	VertexTypeLocation      = "LOCATION"
	VertexTypeChildLocation = "CHILD_LOCATION"
	VertexTypeOwner         = "OWNER"
	VertexTypeProduct       = "PRODUCT"
	VertexTypeBatch         = "BATCH"
	VertexTypeEvent         = "EVENT"

	ActionTypeAdd     = "ADD"
	ActionTypeDelete  = "DELETE"
	ActionTypeObserve = "OBSERVE"
)

// Event ....
type Event struct {
	Time          time.Time   `json:"time"`
	Action        ActionType  `json:"action"` // ADD/DELETE/OBSERVE
	EpcList       Epcs        `json:"epc_list"`
	ParentID      Epc         `json:"parent_id"`
	ChildEPCs     Epcs        `json:"child_epcs"`
	ReadPoint     ReadPoint   `json:"read_point"`
	BizStep       Epc         `json:"biz_step"`
	BizLocation   BizLocation `json:"biz_location"`
	InputEPCList  Epcs        `json:"input_epc_list"`
	OutputEPCList Epcs        `json:"output_epc_list"`
}

// Epc ...
type Epc string

// Epcs slice epc
type Epcs []Epc

// ActionType ...
type ActionType string

// ReadPoint ...
type ReadPoint struct {
	ID string `json:"id"`
}

// BizLocation ...
type BizLocation struct {
	ID string `json:"id"`
}

// Location ...
type Location struct {
	ID          string            `json:"id"`
	Children    ChildrenLocations `json:"chilren"`
	Participant sdk.AccAddress    `json:"participant"`
	Attributes  Attributes        `json:"attributes"`
}

// ChildrenLocation ...
type ChildrenLocation struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

// ChildrenLocations ...
type ChildrenLocations []ChildrenLocation

// Record ...
type Record struct {
	ID       int64          `json:"id"`
	Sender   sdk.AccAddress `json:"sender"`
	Edges    []Edge         `json:"edges"`
	Vertices []Vertice      `json:"vertices"`
}

// Edge ...
type Edge struct {
	Key    Key    `json:"key"`
	Source Key    `json:"source"`
	Target Key    `json:"target"`
	Type   string `json:"type"`
}

// Vertice ...
type Vertice struct {
	Key        Key        `json:"key"`
	Attributes Attributes `json:"attributes"`
	Type       string     `json:"type"`
}

// Attribute ...
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Actor ...
type Actor struct {
	ID         string         `json:"id"`
	Addr       sdk.AccAddress `json:"addr"`
	Attributes Attributes     `json:"attributes"`
}

// Product ...
type Product struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

// Batch ...
type Batch struct {
	ID         string     `json:"id"`
	ProductID  string     `json:"product_id"`
	Attributes Attributes `json:"attributes"`
}

// Attributes slice attributge
type Attributes []Attribute

// Bytes ...
func (attrs Attributes) Bytes() []byte {
	b, err := msgCdc.MarshalJSON(attrs)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Key ...
type Key []byte

func (k Key) String() string {
	return string(k)
}

// MarshalJSON to JSON using Bech32
func (k Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%x", k))
}
