package epcis

import sdk "github.com/cosmos/cosmos-sdk/types"

// Event ....
type Event struct {
	Action        string   `json:"action"` // ADD/DELETE/OBSERVE
	EpcList       []string `json:"epc_list"`
	ParentID      string   `json:"parent_id"`
	ChildEPCs     []string `json:"child_epcs"`
	ReadPoint     string   `json:"read_point"`
	BizLocation   string   `json:"biz_location"`
	InputEPCList  string   `json:"input_epc_list"`
	OutputEPCList string   `json:"output_epc_list"`
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
