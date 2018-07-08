package identity

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgType = "identity"

var _, _, _ sdk.Msg = &MsgSetTrust{}, &MsgCreateIdentity{}, &MsgSetCerts{}

// MsgSetTrust struct for set trust
type MsgSetTrust struct {
	Trustor  sdk.Address `json:"trustor"`
	Trusting sdk.Address `json:"trusting"`
	Trust    bool        `json:"trust"`
}

func NewMsgSetTrust(trustor, trusting sdk.Address, trust bool) MsgSetTrust {
	return MsgSetTrust{
		Trustor:  trustor,
		Trusting: trusting,
		Trust:    trust,
	}
}

//nolint
func (msg MsgSetTrust) Type() string { return MsgType }
func (msg MsgSetTrust) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Trustor}
}

// get the bytes for the message signer to sign on
func (msg MsgSetTrust) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(struct {
		Trustor  string `json:"trustor"`
		Trusting string `json:"trusting"`
		Trust    bool   `json:"trust"`
	}{
		Trustor:  sdk.MustBech32ifyAcc(msg.Trustor),
		Trusting: sdk.MustBech32ifyAcc(msg.Trusting),
		Trust:    msg.Trust,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgSetTrust) ValidateBasic() sdk.Error {
	if msg.Trustor == nil {
		return ErrNilTrustorAddr(DefaultCodespace)
	}
	if msg.Trusting == nil {
		return ErrNilTrustingAddr(DefaultCodespace)
	}

	return nil
}

// MsgCreateIdentity struct for create identity
type MsgCreateIdentity struct {
	Sender sdk.Address `json:"sender"`
}

func NewMsgCreateIdentity(sender sdk.Address) MsgCreateIdentity {
	return MsgCreateIdentity{
		Sender: sender,
	}
}

//nolint
func (msg MsgCreateIdentity) Type() string { return MsgType }
func (msg MsgCreateIdentity) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Sender}
}

// get the bytes for the message signer to sign on
func (msg MsgCreateIdentity) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(struct {
		Sender string `json:"sender"`
	}{
		Sender: sdk.MustBech32ifyAcc(msg.Sender),
	})
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgCreateIdentity) ValidateBasic() sdk.Error {
	if msg.Sender == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil sender address")
	}

	return nil
}

// MsgSetCerts struct for set certs
type MsgSetCerts struct {
	Certifier  sdk.Address `json:"certifier"`
	IdentityID int64       `json:"identity_id"`
	Values     []CertValue `json:"values"`
}

func NewMsgSetCerts(certifier sdk.Address, identityID int64, values []CertValue) MsgSetCerts {
	return MsgSetCerts{
		Certifier:  certifier,
		IdentityID: identityID,
		Values:     values,
	}
}

//nolint
func (msg MsgSetCerts) Type() string { return MsgType }
func (msg MsgSetCerts) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Certifier}
}

// get the bytes for the message signer to sign on
func (msg MsgSetCerts) GetSignBytes() []byte {
	var values []json.RawMessage

	for _, value := range msg.Values {
		values = append(values, value.GetSignBytes())
	}

	b, err := MsgCdc.MarshalJSON(struct {
		Certifier  string            `json:"certifier"`
		IdentityID int64             `json:"identity_id"`
		Values     []json.RawMessage `json:"values"`
	}{
		Certifier:  sdk.MustBech32ifyAcc(msg.Certifier),
		IdentityID: msg.IdentityID,
		Values:     values,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgSetCerts) ValidateBasic() sdk.Error {
	if msg.Certifier == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil certifier address")
	}
	if msg.IdentityID == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil identity id")
	}
	for _, value := range msg.Values {
		if err := value.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
