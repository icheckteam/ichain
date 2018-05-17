package identity

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(CreateMsg{}, "ichain/ClaimIssueMsg", nil)
	cdc.RegisterConcrete(RevokeMsg{}, "ichain/RevokeMsg", nil)

}
