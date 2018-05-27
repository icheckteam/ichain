package smartcontract

import sdk "github.com/cosmos/cosmos-sdk/types"

// MsgDeployContract deploy contract to blockchain
type MsgDeployContract struct {
	ContractID string
	Issuer     sdk.Address
	Email      string
	Website    string
}
