package identity

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identity
type Identity struct {
	ID       int64          `json:"id"`    // id of the identity
	Owner    sdk.AccAddress `json:"owner"` // owner of the identity
	Defaults bool           `json:"defaults"`
}

type Cert struct {
	ID         string         `json:"id"`
	Property   string         `json:"property"`
	Certifier  sdk.AccAddress `json:"certifier"`
	Type       string         `json:"type"`
	Trust      bool           `json:"trust"`
	Data       Metadata       `json:"data"`
	Confidence bool           `json:"confidence"`
}

type CertValue struct {
	Property   string   `json:"property"`
	Type       string   `json:"type"`
	Data       Metadata `json:"data"`
	Confidence bool     `json:"confidence"`
}

// quick validity check
func (msg CertValue) ValidateBasic() sdk.Error {
	if len(msg.Property) == 0 {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil property address")
	}
	return nil
}

func (msg CertValue) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(struct {
		Property   string   `json:"property"`
		Type       string   `json:"type"`
		Data       Metadata `json:"data"`
		Confidence bool     `json:"confidence"`
	}{
		Property:   msg.Property,
		Type:       msg.Type,
		Data:       msg.Data,
		Confidence: msg.Confidence,
	})
	if err != nil {
		panic(err)
	}
	return b
}

type Certs []Cert

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

type Trust struct {
	Trustor  sdk.AccAddress `json:"trustor"`
	Trusting sdk.AccAddress `json:"trusting"`
}
