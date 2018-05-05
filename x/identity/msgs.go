package idetify

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgType name to idetify transaction types
const MsgType = "idetify"

// MsgClaim ..
type MsgClaim struct {
	Claim
}

// NOLINT
func (msg MsgClaim) Type() string                            { return MsgType } //TODO update "stake/declarecandidacy"
func (msg MsgClaim) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgClaim) GetSigners() []sdk.Address               { return []sdk.Address{msg.Metadata.Issuer} }
func (msg MsgClaim) String() string {
	return fmt.Sprintf("MsgClaim{Issuer: %v}", msg.Metadata.Issuer) // XXX fix
}
