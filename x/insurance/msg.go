package insurance

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/types"
)

const msgType = "insurance"

// MsgCreateContract
// --------------------------------------------------
type MsgCreateContract struct {
	ID        string         `json:"id"`
	Issuer    sdk.AccAddress `json:"issuer"`
	Recipient sdk.AccAddress `json:"recipient"`
	Expires   time.Time      `json:"expires"`
	Serial    string         `json:"serial"`
	AssetID   string         `json:"asset_id"`
}

// nolint ...
func (msg MsgCreateContract) Type() string                            { return msgType }
func (msg MsgCreateContract) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateContract) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Issuer} }

// Implements Msg.
func (msg MsgCreateContract) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgCreateContract) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "id")
	}

	if len(msg.Issuer) == 0 {
		return types.ErrMissingField(DefaultCodespace, "issuer")
	}

	if len(msg.Recipient) == 0 {
		return types.ErrMissingField(DefaultCodespace, "recipient")
	}

	if msg.Expires.IsZero() {
		return types.ErrMissingField(DefaultCodespace, "expires")
	}

	if len(msg.AssetID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "asset_id")
	}

	return nil
}

// MsgCreateClaim
// --------------------------------------------------
type MsgCreateClaim struct {
	ContractID string         `json:"contract_id"`
	Issuer     sdk.AccAddress `json:"issuer"`
	Recipient  sdk.AccAddress `json:"recipient"`
}

func NewMsgCreateClaim(issuer, recipient sdk.AccAddress, contractID string) MsgCreateClaim {
	return MsgCreateClaim{
		ContractID: contractID,
		Issuer:     issuer,
		Recipient:  recipient,
	}
}

// nolint ...
func (msg MsgCreateClaim) Type() string                            { return msgType }
func (msg MsgCreateClaim) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgCreateClaim) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Issuer} }
func (msg MsgCreateClaim) String() string {
	return fmt.Sprintf(`MsgCreateClaim{%v->%v}`, msg.ContractID, msg.Issuer)
}

// Implements Msg.
func (msg MsgCreateClaim) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgCreateClaim) ValidateBasic() sdk.Error {
	if len(msg.ContractID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "contract_id")
	}

	if len(msg.Issuer) == 0 {
		return types.ErrMissingField(DefaultCodespace, "issuer")
	}
	if len(msg.Recipient) == 0 {
		return types.ErrMissingField(DefaultCodespace, "recipient")
	}
	return nil
}

// MsgCompleteClaim
// --------------------------------------------------
type MsgProcessClaim struct {
	ContractID string         `json:"contract_id"`
	Issuer     sdk.AccAddress `json:"issuer"`
	Status     ClaimStatus    `json:"status"`
}

// nolint ...
func (msg MsgProcessClaim) Type() string                            { return msgType }
func (msg MsgProcessClaim) Get(key interface{}) (value interface{}) { return nil }
func (msg MsgProcessClaim) GetSigners() []sdk.AccAddress            { return []sdk.AccAddress{msg.Issuer} }
func (msg MsgProcessClaim) String() string {
	return fmt.Sprintf(`MsgProcessClaim{%v->%v}`, msg.ContractID, msg.Issuer)
}

// Implements Msg.
func (msg MsgProcessClaim) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgProcessClaim) ValidateBasic() sdk.Error {
	if len(msg.ContractID) == 0 {
		return types.ErrMissingField(DefaultCodespace, "contract_id")
	}

	if len(msg.Issuer) == 0 {
		return types.ErrMissingField(DefaultCodespace, "issuer")
	}
	switch msg.Status {
	case ClaimStatusPending,
		ClaimStatusRejected,
		ClaimStatusReimbursement,
		ClaimStatusTheftConfirmed,
		ClaimStatusClaimRepair:
		break
	default:
		return types.ErrInvalidField(DefaultCodespace, "status")
	}
	return nil
}
