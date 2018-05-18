package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeMissingField sdk.CodeType      = 505
	CodeInvalidField sdk.CodeType      = 506
	DefaultCodespace sdk.CodespaceType = 10
)

// ErrMissingField ...
func ErrMissingField(field string) sdk.Error {
	return newError(DefaultCodespace, CodeMissingField, fmt.Sprintf("missing %s", field))
}

// ErrInvalidField ...
func ErrInvalidField(field string) sdk.Error {
	return newError(DefaultCodespace, CodeMissingField, fmt.Sprintf("field %s has invalid value", field))
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
