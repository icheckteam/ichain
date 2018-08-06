package identity

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Cert ...
type Cert struct {
	Property  string         `json:"property"`
	Certifier sdk.AccAddress `json:"certifier"`
	Owner     sdk.AccAddress `json:"owner"`
	Data      Metadata       `json:"data"`
	CreatedAt int64          `json:"created_at"`
}

// CertValue ...
type CertValue struct {
	Owner      sdk.AccAddress `json:"owner"`
	Property   string         `json:"property"`
	Data       Metadata       `json:"data"`
	Confidence bool           `json:"confidence"`
}

// ValidateBasic quick validity check
func (msg CertValue) ValidateBasic() sdk.Error {
	if len(msg.Property) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil property address")
	}

	if len(msg.Owner) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil owner address")
	}
	return nil
}

// Certs ...
type Certs []Cert

// Metadata struct
type Metadata []byte

// MarshalJSON returns *m as the JSON encoding of m.
func (j Metadata) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (j *Metadata) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// UnmarshalCert ...
func UnmarshalCert(cdc *wire.Codec, value []byte) (cert Cert, err error) {
	err = cdc.UnmarshalBinary(value, &cert)
	return
}
