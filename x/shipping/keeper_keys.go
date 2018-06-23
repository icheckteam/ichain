package shipping

import sdk "github.com/cosmos/cosmos-sdk/types"

// TODO remove some of these prefixes once have working multistore

//nolint
var (
	// Keys for store prefixes
	OrderKey         = []byte{0x06} // prefix for each key to a order
	AccountOrdersKey = []byte{0x01}
)

// GetOrderKey get the key for the record with address
func GetOrderKey(uid []byte) []byte {
	return append(OrderKey, uid...)
}

// GetAccountOrderKey get the key for an account for a order
func GetAccountOrderKey(addr sdk.Address, claimID string) []byte {
	return append(GetAccountOrdersKey(addr), []byte(claimID)...)
}

// GetAccountOrdersKey get the key for an account for all orders
func GetAccountOrdersKey(addr sdk.Address) []byte {
	return append(AccountOrdersKey, []byte(addr.String())...)
}
