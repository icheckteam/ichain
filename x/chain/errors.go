package chain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeUnknownRecord      sdk.CodeType = 500
	CodeInvalidTransaction sdk.CodeType = 501
)

// ErrUnknownRecord ...
func ErrUnknownRecord(msg string) sdk.Error {
	return newError(CodeUnknownRecord, msg)
}

// InvalidTransaction ...
func InvalidTransaction(msg string) sdk.Error {
	return newError(CodeInvalidTransaction, msg)
}

// CodeToDefaultMsg NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeUnknownRecord:
		return "Unknown Record"
	default:
		return fmt.Sprintf("Unknown code %d", code)
	}
}

func newError(code sdk.CodeType, msg string) sdk.Error {
	// TODO capture stacktrace if ENV is set.
	if msg == "" {
		msg = CodeToDefaultMsg(code)
	}
	return sdk.NewError(code, msg)
}
