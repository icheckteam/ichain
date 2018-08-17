package gs1

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
const (
	// edge
	EdgeTypeOwnerdBy = "OWNED_BY"

	// VertexType
	VertexTypeLocation    = "LOCATION"
	EdgeTypeChildLocation = "CHILD_LOCATION"
)

// Event ....
type Event struct {
	Time          time.Time   `json:"time"`
	Action        ActionType  `json:"action"` // ADD/DELETE/OBSERVE
	EpcList       []Epc       `json:"epc_list"`
	ParentID      Epc         `json:"parent_id"`
	ChildEPCs     []Epc       `json:"child_epcs"`
	ReadPoint     ReadPoint   `json:"read_point"`
	BizLocation   BizLocation `json:"biz_location"`
	InputEPCList  []Epc       `json:"input_epc_list"`
	OutputEPCList []Epc       `json:"output_epc_list"`
}

// Epc ...
type Epc string

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
	ID          string             `json:"id"`
	Children    []ChildrenLocation `json:"chilren"`
	Participant sdk.AccAddress     `json:"participant"`
	Attributes  Attributes         `json:"attributes"`
}

// ChildrenLocation ...
type ChildrenLocation struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

// Record ...
type Record struct {
	ID       string         `json:"id"`
	Sender   sdk.AccAddress `json:"sender"`
	Edges    []Edge         `json:"edges"`
	Vertices []Vertice      `json:"vertices"`
}

// Edge ...
type Edge struct {
	Key    []byte `json:"key"`
	Source []byte `json:"source"`
	Target []byte `json:"target"`
	Type   string `json:"type"`
}

// Vertice ...
type Vertice struct {
	Key        []byte      `json:"key"`
	Attributes []Attribute `json:"attributes"`
	Type       string      `json:"type"`
}

// Attribute ...
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Actor ...
type Actor struct {
	Addr       sdk.AccAddress `json:"addr"`
	Attributes []Attribute    `json:"attributes"`
}

// Product ...
type Product struct {
	ID         string      `json:"id"`
	Attributes []Attribute `json:"attributes"`
}

// Batch ...
type Batch struct {
	ID         string      `json:"id"`
	ProductID  string      `json:"product_id"`
	Attributes []Attribute `json:"attributes"`
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
