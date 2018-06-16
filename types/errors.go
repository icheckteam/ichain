package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeInvalidTransaction sdk.CodeType      = 501
	CodeMissingField       sdk.CodeType      = 502
	CodeInvalidField       sdk.CodeType      = 503
	DefaultCodespace       sdk.CodespaceType = 10
)

// ErrMissingField ...
func ErrMissingField(codespace sdk.CodespaceType, field string) sdk.Error {
	return newError(codespace, CodeMissingField, fmt.Sprintf("missing %s", field))
}

// ErrInvalidField ...
func ErrInvalidField(codespace sdk.CodespaceType, field string) sdk.Error {
	return newError(codespace, CodeInvalidField, fmt.Sprintf("field %s has invalid value", field))
}

// InvalidTransaction ...
func InvalidTransaction(codespace sdk.CodespaceType, msg string) sdk.Error {
	return newError(codespace, CodeInvalidTransaction, msg)
}

// CodeToDefaultMsg NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code sdk.CodeType) string {
	return fmt.Sprintf("Unknown code %d", code)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	// TODO capture stacktrace if ENV is set.
	if msg == "" {
		msg = CodeToDefaultMsg(code)
	}
	return sdk.NewError(codespace, code, msg)
}
