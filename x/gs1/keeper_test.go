package gs1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	ctx, keeper := createTestInput(t, false, 0)
	_, err := keeper.CreateRecord(ctx, MsgSend{
		Sender:   addrs[0],
		Receiver: addrs[0],
		Actors: []Actor{
			{Addr: addrs[0], ID: "urn:ichain:object:actor:id:Company_A"},
			{Addr: addrs[1], ID: "urn:ichain:object:actor:id:Company_B"},
		},

		Locations: []Location{
			{
				ID:          "urn:epc:id:sgln:Building_A",
				Participant: addrs[0],
				Attributes: Attributes{
					Attribute{Name: "category", Value: "Company"},
				},
			},
			{
				ID:          "urn:epc:id:sgln:Building_B",
				Participant: addrs[1],
				Attributes: Attributes{
					Attribute{Name: "category", Value: "Company"},
				},
				Children: ChildrenLocations{
					ChildrenLocation{ID: "urn:epc:id:sgln:Building_B_1"},
					ChildrenLocation{ID: "urn:epc:id:sgln:Building_B_2"},
				},
			},
		},

		Products: []Product{
			{
				ID: "urn:ichain:object:product:id:Product_1",
				Attributes: []Attribute{
					{Name: "productId", Value: "urn:ichain:object:product:id:Product_1"},
				},
			},
			{
				ID: "urn:ichain:object:product:id:Product_2",
			},
		},

		Batches: []Batch{
			Batch{
				ID:        "urn:epc:id:sgtin:Batch_1",
				ProductID: "urn:ichain:object:product:id:Product_1",
			},
		},

		Events: []Event{
			Event{
				Action: "OBSERVE",
				EpcList: Epcs{
					"urn:epc:id:sgtin:Batch_1",
				},
				ReadPoint: ReadPoint{
					ID: "urn:epc:id:sgln:Building_B_2",
				},
			},
		},
	})

	record := keeper.GetRecord(ctx, 1)

	b, _ := json.MarshalIndent(record, "", " ")

	fmt.Printf("%s", b)

	assert.True(t, err != nil)
}
