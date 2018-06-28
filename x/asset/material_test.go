package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaterialType(t *testing.T) {
	materials := Materials{Material{AssetID: "1", Quantity: 1}}
	materials = materials.Plus(Materials{Material{AssetID: "1", Quantity: 1}})
	assert.Equal(t, materials[0].Quantity, int64(2))
}
