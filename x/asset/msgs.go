package asset

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

// RegisterMsg A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type RegisterMsg struct {
	Issuer   sdk.Address `json:"issuer"`
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Quantity int64       `json:"quantity"`

	// Company info
	Company string `json:"company"`
	Email   string `json:"email"`
}

// NewRegisterMsg new record create msg
func NewRegisterMsg(issuer sdk.Address, id, name string, quantity int64, company, email string) RegisterMsg {
	return RegisterMsg{
		Issuer:   issuer,
		ID:       id,
		Quantity: quantity,
		Name:     name,
		Company:  company,
		Email:    email,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = RegisterMsg{}

// nolint ...
func (msg RegisterMsg) Type() string                            { return msgType }
func (msg RegisterMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg RegisterMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg RegisterMsg) String() string {
	return fmt.Sprintf(`RegisterMsg{
		Issuer:   %s,
		ID:       %s,
		Quantity: %d,
		Name:     %s,
		Company:  %s,
		Email:    %s,
	}`,
		msg.Issuer, msg.ID, msg.Quantity,
		msg.Name, msg.Company, msg.Email,
	)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}

	if len(msg.ID) == 0 {
		return ErrMissingField("asset_id")
	}

	if len(msg.Name) == 0 {
		return ErrMissingField("asset_name")
	}

	if msg.Quantity == 0 {
		return ErrMissingField("asset_quantity")
	}

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg RegisterMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// UpdateAttrMsg ...
// ---------------------------------------------------------------
type UpdateAttrMsg struct {
	Issuer     sdk.Address `json:"issuer"`
	ID         string      `json:"id"`
	Attributes []Attribute `json:"attribute"`
}

func (msg UpdateAttrMsg) Type() string                            { return msgType }
func (msg UpdateAttrMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg UpdateAttrMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg UpdateAttrMsg) String() string {
	return fmt.Sprintf("UpdateAttrMsg{%s->%v}", msg.ID, msg.Attributes)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg UpdateAttrMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}
	if len(msg.ID) == 0 {
		return ErrMissingField("id")
	}
	if len(msg.Attributes) == 0 {
		return ErrMissingField("name")
	}
	for _, attr := range msg.Attributes {
		if attr.Type > 5 {
			return ErrMissingField("attributes")
		}
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg UpdateAttrMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// AddQuantityMsg ...
// ---------------------------------------------------------------
type AddQuantityMsg struct {
	Issuer   sdk.Address `json:"issuer"`
	ID       string      `json:"id"`
	Quantity int64       `json:"quantity"`
}

func (msg AddQuantityMsg) Type() string                            { return msgType }
func (msg AddQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AddQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg AddQuantityMsg) String() string {
	return fmt.Sprintf("AddQuantityMsg{Sender: %v, quantity: %v}", msg.Issuer, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AddQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}
	if len(msg.ID) == 0 {
		return ErrMissingField("id")
	}
	if msg.Quantity <= 0 {
		return ErrMissingField("quantity")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg AddQuantityMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// AddQuantityMsg ...
// ---------------------------------------------------------------
type SubtractQuantityMsg struct {
	Issuer   sdk.Address `json:"issuer"`
	ID       string      `json:"id"`
	Quantity int64       `json:"quantity"`
}

func (msg SubtractQuantityMsg) Type() string                            { return msgType }
func (msg SubtractQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg SubtractQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg SubtractQuantityMsg) String() string {
	return fmt.Sprintf("SubtractQuantityMsg{Issuer: %v, quantity: %v}", msg.Issuer, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg SubtractQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}
	if len(msg.ID) == 0 {
		return ErrMissingField("id")
	}
	if msg.Quantity <= 0 {
		return ErrMissingField("quantity")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg SubtractQuantityMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// CreateProposalMsg ...
type CreateProposalMsg struct {
	AssetID     string      `json:"asset_id"`
	Issuer      sdk.Address `json:"issuer"`
	Recipient   sdk.Address `json:"recipient"`
	Propertipes []string    `json:"propertipes"`
}

func (msg CreateProposalMsg) Type() string                            { return msgType }
func (msg CreateProposalMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CreateProposalMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg CreateProposalMsg) String() string {
	return fmt.Sprintf(`
		CreateProposalMsg{
			AssetID: %s, 
			Issuer: %v,
			Recipient:%v,
			Propertipes:%v
		}	
	`, msg.AssetID, msg.Issuer, msg.Recipient, msg.Propertipes)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg CreateProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Issuer) == 0 {
		return sdk.ErrUnknownAddress(msg.Issuer.String()).Trace("")
	}
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	if len(msg.Propertipes) == 0 {
		return ErrMissingField("propertipes")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg CreateProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// AnswerProposalMsg ...
type AnswerProposalMsg struct {
	AssetID   string      `json:"asset_id"`
	Recipient sdk.Address `json:"recipient"`
	Response  int         `json:"response"`
}

func (msg AnswerProposalMsg) Type() string                            { return msgType }
func (msg AnswerProposalMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AnswerProposalMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Recipient} }
func (msg AnswerProposalMsg) String() string {
	return fmt.Sprintf(`
		AnswerProposalMsg{
			AssetID: %s, 
			Recipient: %v,
			Response:%v,
		}	
	`, msg.AssetID, msg.Recipient, msg.Response)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AnswerProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.AssetID) == 0 {
		return ErrMissingField("asset_id")
	}
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	if msg.Response > 2 {
		return ErrMissingField("response")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg AnswerProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// RevokeProposalMsg ...
type RevokeProposalMsg struct {
	AssetID     string      `json:"asset_id"`
	Issuer      sdk.Address `json:"issuer"`
	Recipient   sdk.Address `json:"recipient"`
	Propertipes []string    `json:"propertipes"`
}

func (msg RevokeProposalMsg) Type() string                            { return msgType }
func (msg RevokeProposalMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg RevokeProposalMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Issuer} }
func (msg RevokeProposalMsg) String() string {
	return fmt.Sprintf(`
		RevokeProposalMsg{
			AssetID: %s, 
			Issuer: %s,
			Recipient: %v,
			Propertipes:%v,
		}	
	`, msg.AssetID, msg.Issuer, msg.Recipient, msg.Propertipes)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg RevokeProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Recipient) == 0 {
		return ErrMissingField("recipient")
	}
	if len(msg.Issuer) == 0 {
		return ErrMissingField("Issuer")
	}
	if len(msg.Propertipes) == 0 {
		return ErrMissingField("propertipes")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg RevokeProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
