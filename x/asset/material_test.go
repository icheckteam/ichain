package asset

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMaterialType(t *testing.T) {
	materials := Materials{Material{AssetID: "1", Quantity: sdk.NewInt(1)}}
	materials = materials.Plus(Materials{Material{AssetID: "1", Quantity: sdk.NewInt(1)}})
	assert.Equal(t, materials[0].Quantity, sdk.NewInt(2))
}
