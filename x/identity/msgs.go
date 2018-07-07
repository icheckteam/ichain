package identity

import sdk "github.com/cosmos/cosmos-sdk/types"

type MsgSetTrust struct {
	Trustor  sdk.Address `json:"trustor"`
	Trusting sdk.Address `json:"trusting"`
	Trust    bool        `json:"trust"`
}

type MsgCreateIdentity struct {
	Sender     sdk.Address `json:"sender"`
	IdentityID int64       `json:"identity_id"`
}

type MsgSetCerts struct {
	Certifier  sdk.Address `json:"certifier"`
	IdentityID int64       `json:"identity_id"`
	Values     []CertValue `json:"values"`
}
