package shipping

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeDuplicateOrderID sdk.CodeType = iota + 600
	CodeUnknownOrder
	CodeDuplicateAddress
	CodeInvalidAssetAmount
	DefaultCodespace sdk.CodespaceType = 10
)

// ErrDuplicateOrder ...
func ErrDuplicateOrder(orderID string) sdk.Error {
	return newError(DefaultCodespace, CodeDuplicateOrderID, fmt.Sprintf("duplicate order id %s", orderID))
}

// ErrUnknownOrder ...
func ErrUnknownOrder(orderID string) sdk.Error {
	return newError(DefaultCodespace, CodeUnknownOrder, fmt.Sprintf("order id %s not found", orderID))
}

// ErrDuplicateAddress ...
func ErrDuplicateAddress() sdk.Error {
	return newError(DefaultCodespace, CodeDuplicateAddress, "issuer, carrier or receiver address is duplicated")
}

// ErrInavlidAssetAmount ...
func ErrInavlidAssetAmount() sdk.Error {
	return newError(DefaultCodespace, CodeInvalidAssetAmount, "asset amount cannot be zero")
}

// CodeToDefaultMsg NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code sdk.CodeType) string {
	switch code {

	default:
		return fmt.Sprintf("Unknown code %d", code)
	}
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	// TODO capture stacktrace if ENV is set.
	if msg == "" {
		msg = CodeToDefaultMsg(code)
	}
	return sdk.NewError(codespace, code, msg)
}
