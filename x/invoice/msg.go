package invoice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/types"
)

type MsgCreate struct {
	ID       string         `json:"id"`
	Issuer   sdk.AccAddress `json:"issuer"`
	Receiver sdk.AccAddress `json:"receiver"`
	Items    []Item         `json:"items"`
}

func (msg MsgCreate) Type() string {
	return "invoice"
}

func (msg MsgCreate) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCreate) ValidateBasic() sdk.Error {
	if msg.ID == "" {
		return types.ErrMissingField(DefaultCodespace, "id")
	}

	if len(msg.Issuer) == 0 {
		return types.ErrMissingField(DefaultCodespace, "issuer")
	}

	if len(msg.Receiver) == 0 {
		return types.ErrMissingField(DefaultCodespace, "receiver")
	}

	if len(msg.Items) == 0 {
		return types.ErrMissingField(DefaultCodespace, "items")
	}

	return nil
}

func (msg MsgCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func NewMsgCreate(id string, issuer, receiver sdk.AccAddress, items []Item) MsgCreate {
	return MsgCreate{
		ID:       id,
		Issuer:   issuer,
		Receiver: receiver,
		Items:    items,
	}
}
