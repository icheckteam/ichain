package identity

import (
	"bytes"
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim struct
type Claim struct {
	ID         string      `json:"id"`
	Issuer     sdk.Address `json:"issuer"`
	Recipient  sdk.Address `json:"recipient"`
	Context    string      `json:"context"`
	Content    Content     `json:"content"`
	Fee        sdk.Coins   `json:"fee"`
	Paid       bool        `json:"paid"`
	CreateTime int64       `json:"create_time"`
	Expires    int64       `json:"expires"`
	Revocation string      `json:"revocation"`
}

func (c Claim) IsOwner(addr sdk.Address) bool {
	return bytes.Equal(c.Issuer, addr)
}

func (c Claim) GetContent() (map[string]interface{}, error) {
	content := map[string]interface{}{}
	err := json.Unmarshal(c.Content, &content)
	return content, err
}

type Content []byte

// MarshalJSON returns *m as the JSON encoding of m.
func (j Content) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (j *Content) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

type ClaimRest struct {
	ID         string    `json:"id"`
	Issuer     string    `json:"issuer"`
	Recipient  string    `json:"recipient"`
	Context    string    `json:"context"`
	Content    Content   `json:"content"`
	Fee        sdk.Coins `json:"fee"`
	Paid       bool      `json:"paid"`
	CreateTime int64     `json:"create_time"`
	Expires    int64     `json:"expires"`
	Revocation string    `json:"revocation"`
}

func ClaimToRest(claim Claim) ClaimRest {
	return ClaimRest{
		ID:         claim.ID,
		Issuer:     sdk.MustBech32ifyAcc(claim.Issuer),
		Recipient:  sdk.MustBech32ifyAcc(claim.Recipient),
		Context:    claim.Context,
		Content:    claim.Content,
		Expires:    claim.Expires,
		CreateTime: claim.CreateTime,
		Revocation: claim.Revocation,
	}
}
