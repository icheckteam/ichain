package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgType name to idetify transaction types
const MsgType = "identity"

var _, _ sdk.Msg = &MsgSetTrust{}, &MsgSetCerts{}

// MsgSetTrust struct for set trust
type MsgSetTrust struct {
	Trustor  sdk.AccAddress `json:"trustor"`
	Trusting sdk.AccAddress `json:"trusting"`
	Trust    bool           `json:"trust"`
}

// NewMsgSetTrust ...
func NewMsgSetTrust(trustor, trusting sdk.AccAddress, trust bool) MsgSetTrust {
	return MsgSetTrust{
		Trustor:  trustor,
		Trusting: trusting,
		Trust:    trust,
	}
}

//Type ...
func (msg MsgSetTrust) Type() string { return MsgType }

// GetSigners ...
func (msg MsgSetTrust) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustor}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgSetTrust) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic quick validity check
func (msg MsgSetTrust) ValidateBasic() sdk.Error {
	if msg.Trustor == nil {
		return ErrNilTrustorAddr(DefaultCodespace)
	}
	if msg.Trusting == nil {
		return ErrNilTrustingAddr(DefaultCodespace)
	}

	return nil
}

// MsgSetCerts struct for set certs
type MsgSetCerts struct {
	Certifier sdk.AccAddress `json:"certifier"`
	Recipient sdk.AccAddress `json:"recipient"`
	Values    []CertValue    `json:"values"`
}

// NewMsgSetCerts ...
func NewMsgSetCerts(certifier sdk.AccAddress, recipient sdk.AccAddress, values []CertValue) MsgSetCerts {
	return MsgSetCerts{
		Certifier: certifier,
		Recipient: recipient,
		Values:    values,
	}
}

// Type ...
func (msg MsgSetCerts) Type() string { return MsgType }

// GetSigners ...
func (msg MsgSetCerts) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Certifier}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgSetCerts) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic quick validity check
func (msg MsgSetCerts) ValidateBasic() sdk.Error {
	if msg.Certifier == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil certifier address")
	}
	if msg.Recipient == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil recipient  address")
	}
	for _, value := range msg.Values {
		if err := value.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// Identity
// -----------------------------------------------
// MsgReg
// MsgAddKey
// MsgDelKey

// MsgReg ....
// .......................................................
type MsgReg struct {
	Sender  sdk.AccAddress
	Address sdk.AccAddress
}

// Type ...
func (msg MsgReg) Type() string { return MsgType }

// GetSigners ...
func (msg MsgReg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgReg) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic quick validity check
func (msg MsgReg) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil sender address")
	}
	if msg.Address == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil address  address")
	}
	return nil
}

// MsgAddOwner ...
// .......................................................
type MsgAddOwner struct {
	Sender  sdk.AccAddress
	Address sdk.AccAddress
	Owner   sdk.AccAddress
}

// Type ...
func (msg MsgAddOwner) Type() string { return MsgType }

// GetSigners ...
func (msg MsgAddOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ValidateBasic quick validity check
func (msg MsgAddOwner) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil sender address")
	}
	if msg.Address == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil address  address")
	}
	return nil
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgAddOwner) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// MsgDelOwner ...
// .......................................................
type MsgDelOwner struct {
	Sender  sdk.AccAddress
	Address sdk.AccAddress
	Owner   sdk.AccAddress
}

// Type ...
func (msg MsgDelOwner) Type() string { return MsgType }

// GetSigners ...
func (msg MsgDelOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg MsgDelOwner) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic quick validity check
func (msg MsgDelOwner) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil sender address")
	}
	if msg.Address == nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil address  address")
	}
	return nil
}
