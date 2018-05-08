package asset

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const msgType = "asset"

// AssetCreateMsg A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type AssetCreateMsg struct {
	Sender    sdk.Address
	AssetID   string
	AssetName string
	Quantity  int64
}

// NewAssetCreateMsg new record create msg
func NewAssetCreateMsg(sender sdk.Address, assetID, assetName string, quantity int64) AssetCreateMsg {
	return AssetCreateMsg{
		Sender:    sender,
		AssetID:   assetID,
		Quantity:  quantity,
		AssetName: assetID,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = AssetCreateMsg{}

// nolint ...
func (msg AssetCreateMsg) Type() string                            { return msgType }
func (msg AssetCreateMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AssetCreateMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg AssetCreateMsg) String() string {
	return fmt.Sprintf("AssetCreateMsg{Sender: %v, RecordID: %s, RecordName: %s}", msg.Sender, msg.AssetID, msg.AssetName)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AssetCreateMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg AssetCreateMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// TransferMsg A really msg record create type, these fields are can be entirely arbitrary and
// custom to your message
type TransferMsg struct {
	Sender   sdk.Address
	To       sdk.Address
	AssetID  string
	Quantity int64
}

// nolint
func (msg TransferMsg) Type() string                            { return msgType }
func (msg TransferMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg TransferMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg TransferMsg) String() string {
	return fmt.Sprintf("TransferMsg{Sender: %v, To: %s, AssetID: %s, Quantity: %d}", msg.Sender, msg.To, msg.AssetID, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg TransferMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.To) == 0 {
		return sdk.ErrUnknownAddress(msg.To.String()).Trace("")
	}

	return nil
}

// GetSignBytes Get the bytes for the message signer to sign on
func (msg TransferMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// UpdateAttrMsg ...
// ---------------------------------------------------------------
type UpdateAttrMsg struct {
	Sender         sdk.Address
	AssetID        string
	AttributeName  string
	AttributeValue interface{}
}

func (msg UpdateAttrMsg) Type() string                            { return msgType }
func (msg UpdateAttrMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg UpdateAttrMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg UpdateAttrMsg) String() string {
	return fmt.Sprintf("UpdateAttrMsg{Sender: %v}", msg.Sender)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg UpdateAttrMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return sdk.ErrUnknownAddress(msg.AssetID).Trace("")
	}
	if len(msg.AttributeName) == 0 {
		return sdk.ErrInternal("Attribute name is required").Trace("")
	}
	if msg.AttributeValue == 0 {
		return sdk.ErrInternal("Attribute value is required").Trace("")
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
	Sender   sdk.Address
	AssetID  string
	Quantity int64
}

func (msg AddQuantityMsg) Type() string                            { return msgType }
func (msg AddQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AddQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg AddQuantityMsg) String() string {
	return fmt.Sprintf("AddQuantityMsg{Sender: %v, quantity: %v}", msg.Sender, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg AddQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return sdk.ErrUnknownAddress(msg.AssetID).Trace("")
	}
	if msg.Quantity <= 0 {
		return sdk.ErrInternal("Quantity is required").Trace("")
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
	Sender   sdk.Address
	AssetID  string
	Quantity int64
}

func (msg SubtractQuantityMsg) Type() string                            { return msgType }
func (msg SubtractQuantityMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg SubtractQuantityMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg SubtractQuantityMsg) String() string {
	return fmt.Sprintf("SubtractQuantityMsg{Sender: %v, quantity: %v}", msg.Sender, msg.Quantity)
}

// ValidateBasic Validate Basic is used to quickly disqualify obviously invalid messages quickly
func (msg SubtractQuantityMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrUnknownAddress(msg.Sender.String()).Trace("")
	}
	if len(msg.AssetID) == 0 {
		return sdk.ErrUnknownAddress(msg.AssetID).Trace("")
	}
	if msg.Quantity <= 0 {
		return sdk.ErrInternal("Quantity is required").Trace("")
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

//----------------------------------------
// Input

// Transaction Output
type Input struct {
	Address sdk.Address `json:"address"`
	Assets  sdk.Coins   `json:"assets"`
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() sdk.Error {
	if len(in.Address) == 0 {
		return sdk.ErrInvalidAddress(in.Address.String())
	}
	if !in.Assets.IsValid() {
		return sdk.ErrInvalidCoins(in.Assets.String())
	}
	if !in.Assets.IsPositive() {
		return sdk.ErrInvalidCoins(in.Assets.String())
	}
	return nil
}

func (in Input) String() string {
	return fmt.Sprintf("Input{%v,%v}", in.Address, in.Assets)
}

// NewInput - create a transaction input, used with SendMsg
func NewInput(addr sdk.Address, assets sdk.Coins) Input {
	input := Input{
		Address: addr,
		Assets:  assets,
	}
	return input
}

//----------------------------------------
// Output

// Transaction Output
type Output struct {
	Address sdk.Address `json:"address"`
	Assets  sdk.Coins   `json:"assets"`
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() sdk.Error {
	if len(out.Address) == 0 {
		return sdk.ErrInvalidAddress(out.Address.String())
	}
	if !out.Assets.IsValid() {
		return sdk.ErrInvalidCoins(out.Assets.String())
	}
	if !out.Assets.IsPositive() {
		return sdk.ErrInvalidCoins(out.Assets.String())
	}
	return nil
}

func (out Output) String() string {
	return fmt.Sprintf("Output{%v,%v}", out.Address, out.Assets)
}

// NewOutput - create a transaction output, used with SendMsg
func NewOutput(addr sdk.Address, assets sdk.Coins) Output {
	output := Output{
		Address: addr,
		Assets:  assets,
	}
	return output
}
