package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeUnknownAsset       sdk.CodeType = 500
	CodeInvalidTransaction sdk.CodeType = 501
	CodeInvalidInput       sdk.CodeType = 502
	CodeInvalidOutput      sdk.CodeType = 503
	CodeInvalidAssets      sdk.CodeType = 504
	CodeMissingField       sdk.CodeType = 505
)

// ErrUnknownAsset ...
func ErrUnknownAsset(msg string) sdk.Error {
	return newError(CodeUnknownAsset, msg)
}

// ErrMissingField ...
func ErrMissingField(field string) sdk.Error {
	return newError(CodeMissingField, fmt.Sprintf("missing %s", field))
}

// InvalidTransaction ...
func InvalidTransaction(msg string) sdk.Error {
	return newError(CodeInvalidTransaction, msg)
}

//----------------------------------------
// Error constructors

func ErrInvalidInput(msg string) sdk.Error {
	return newError(CodeInvalidInput, msg)
}

func ErrInvalidAssets(msg string) sdk.Error {
	return newError(CodeInvalidAssets, msg)
}

func ErrNoInputs() sdk.Error {
	return newError(CodeInvalidInput, "")
}

func ErrInvalidOutput(msg string) sdk.Error {
	return newError(CodeInvalidOutput, msg)
}

func ErrNoOutputs() sdk.Error {
	return newError(CodeInvalidOutput, "")
}

// CodeToDefaultMsg NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeUnknownAsset:
		return "Unknown asset"
	case CodeInvalidInput:
		return "Invalid input assets"
	case CodeInvalidOutput:
		return "Invalid output assets"
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
