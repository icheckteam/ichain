package identity

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identity
type Identity struct {
	ID    int64       `json:"id"`    // id of the identity
	Owner sdk.Address `json:"owner"` // owner of the identity
}

type Cert struct {
	ID         string      `json:"id"`
	Certifier  sdk.Address `json:"certifier"`
	Type       string      `json:"type"`
	Trust      bool        `json:"trust"`
	Data       Metadata    `json:"data"`
	Confidence bool        `json:"confidence"`
}

type CertValue struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Data       Metadata `json:"data"`
	Confidence bool     `json:"confidence"`
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
	Trusting sdk.Address `json:"trusting"`
}
