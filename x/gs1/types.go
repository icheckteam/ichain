package epcis

import sdk "github.com/cosmos/cosmos-sdk/types"

// Event ....
type Event struct {
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
	Children   []ChildrenLocation `json:"chilren"`
	Attributes []Attribute        `json:"attributes"`
}

// ChildrenLocation ...
type ChildrenLocation struct {
	ID string `json:"id"`
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
	Key    string `json:"key"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

// Vertice ...
type Vertice struct {
	Key  string `json:"key"`
	Data []byte `json:"data"`
	Type string `json:"type"`
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
	Attributes []Attribute `json:"attributes"`
}
