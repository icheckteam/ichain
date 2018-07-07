package identity

import sdk "github.com/cosmos/cosmos-sdk/types"

type MsgAddTrust struct {
	Trustor  sdk.Address `json:"trustor"`
	Trusting sdk.Address `json:"trusting"`
	Trust    bool        `json:"trust"`
}

type MsgCreateIdent struct {
	Sender     sdk.Address `json:"sender"`
	IdentityID int64       `json:"identity_id"`
}

type MsgAddCerts struct {
	Certifier  sdk.Address `json:"certifier"`
	IdentityID int64       `json:"identity_id"`
	Values     []CertValue `json:"values"`
}
