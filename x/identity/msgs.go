package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgType = "identity"

var _, _, _ sdk.Msg = &MsgSetTrust{}, &MsgCreateIdentity{}, &MsgSetCerts{}

// MsgSetTrust struct for set trust
type MsgSetTrust struct {
	Trustor  sdk.AccAddress `json:"trustor"`
	Trusting sdk.AccAddress `json:"trusting"`
	Trust    bool           `json:"trust"`
}

func NewMsgSetTrust(trustor, trusting sdk.AccAddress, trust bool) MsgSetTrust {
	return MsgSetTrust{
		Trustor:  trustor,
		Trusting: trusting,
		Trust:    trust,
	}
}

//nolint
func (msg MsgSetTrust) Type() string { return MsgType }
func (msg MsgSetTrust) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustor}
}

// get the bytes for the message signer to sign on
func (msg MsgSetTrust) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
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
	Sender sdk.AccAddress `json:"sender"`
}

func NewMsgCreateIdentity(sender sdk.AccAddress) MsgCreateIdentity {
	return MsgCreateIdentity{
		Sender: sender,
	}
}

//nolint
func (msg MsgCreateIdentity) Type() string { return MsgType }
func (msg MsgCreateIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// get the bytes for the message signer to sign on
func (msg MsgCreateIdentity) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
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
	Certifier  sdk.AccAddress `json:"certifier"`
	IdentityID int64          `json:"identity_id"`
	Values     []CertValue    `json:"values"`
}

func NewMsgSetCerts(certifier sdk.AccAddress, identityID int64, values []CertValue) MsgSetCerts {
	return MsgSetCerts{
		Certifier:  certifier,
		IdentityID: identityID,
		Values:     values,
	}
}

//nolint
func (msg MsgSetCerts) Type() string { return MsgType }
func (msg MsgSetCerts) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Certifier}
}

// get the bytes for the message signer to sign on
func (msg MsgSetCerts) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
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
