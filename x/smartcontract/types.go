package smartcontract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Contract
type Contract struct {
	ID      string
	Issuer  sdk.Address
	Email   string
	Website string
}

// Deploy contract
// ----------------------------------
