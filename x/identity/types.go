package identity

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Cert ...
type Cert struct {
	Property  string         `json:"property"`
	Certifier sdk.AccAddress `json:"certifier"`
	Data      Metadata       `json:"data"`
	CreatedAt int64          `json:"created_at"`
}

// CertValue ...
type CertValue struct {
	Property   string   `json:"property"`
	Data       Metadata `json:"data"`
	Confidence bool     `json:"confidence"`
	Expires    int64    `json:"expires"`
}

// ValidateBasic quick validity check
func (msg CertValue) ValidateBasic() sdk.Error {
	if len(msg.Property) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil property address")
	}
	return nil
}

// GetSignBytes ...
func (msg CertValue) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
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
