package invoice

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/types"
)

type MsgCreate struct {
	ID       string      `json:"id"`
	Issuer   sdk.Address `json:"issuer"`
	Receiver sdk.Address `json:"receiver"`
	Items    []Item      `json:"items"`
}

func (msg MsgCreate) Type() string {
	return "invoice"
}

func (msg MsgCreate) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgCreate) ValidateBasic() sdk.Error {
	if msg.ID == "" {
		return types.ErrMissingField("id")
	}

	if len(msg.Issuer) == 0 {
		return types.ErrMissingField("issuer")
	}

	if len(msg.Receiver) == 0 {
		return types.ErrMissingField("receiver")
	}

	if len(msg.Items) == 0 {
		return types.ErrMissingField("items")
	}

	return nil
}

func (msg MsgCreate) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Issuer}
}

func NewMsgCreate(id string, issuer, receiver sdk.Address, items []Item) MsgCreate {
	return MsgCreate{
		ID:       id,
		Issuer:   issuer,
		Receiver: receiver,
		Items:    items,
	}
}
