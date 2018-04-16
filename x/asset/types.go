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
	Sender     sdk.Address
	RecordID   string
	RecordName string
}

// NewAssetCreateMsg new record create msg
func NewAssetCreateMsg(sender sdk.Address, recordID, recordName string) AssetCreateMsg {
	return AssetCreateMsg{
		Sender:     sender,
		RecordID:   recordID,
		RecordName: recordName,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = AssetCreateMsg{}

// nolint
func (msg AssetCreateMsg) Type() string                            { return msgType }
func (msg AssetCreateMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg AssetCreateMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg AssetCreateMsg) String() string {
	return fmt.Sprintf("AssetCreateMsg{Sender: %v, RecordID: %s, RecordName: %s}", msg.Sender, msg.RecordID, msg.RecordName)
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
	RecordID string
}

// nolint
func (msg TransferMsg) Type() string                            { return msgType }
func (msg TransferMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg TransferMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg TransferMsg) String() string {
	return fmt.Sprintf("TransferMsg{Sender: %v, To: %s, RecordID: %s}", msg.Sender, msg.To, msg.RecordID)
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

// Asset record
// NewAsset("0001", "Cu cai", "0x199...", 100, {"weight": 100})
type Asset struct {
	ID         string
	Name       string
	Issuer     sdk.Address
	Quantity   uint64
	Properties map[string]interface{}
}
