package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeUnknownAsset           sdk.CodeType      = 500
	CodeInvalidTransaction     sdk.CodeType      = 501
	CodeInvalidInput           sdk.CodeType      = 502
	CodeInvalidOutput          sdk.CodeType      = 503
	CodeInvalidAssets          sdk.CodeType      = 504
	CodeMissingField           sdk.CodeType      = 505
	CodeInvalidField           sdk.CodeType      = 506
	CodeInvalidRevokeRecipient sdk.CodeType      = 507
	CodeInvalidAssetQuantity   sdk.CodeType      = 508
	DefaultCodespace           sdk.CodespaceType = 10
)

// ErrUnknownAsset ...
func ErrUnknownAsset(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeUnknownAsset, msg)
}

func ErrAssetNotFound(assetID string) sdk.Error {
	return newError(DefaultCodespace, CodeUnknownAsset, fmt.Sprintf("asset id %s not found", assetID))
}

// ErrMissingField ...
func ErrMissingField(field string) sdk.Error {
	return newError(DefaultCodespace, CodeMissingField, fmt.Sprintf("missing %s", field))
}

// ErrInvalidField ...
func ErrInvalidField(field string) sdk.Error {
	return newError(DefaultCodespace, CodeMissingField, fmt.Sprintf("field %s has invalid value", field))
}

// ErrInvalidAssetQuantity ...
func ErrInvalidAssetQuantity(assetID string) sdk.Error {
	return newError(DefaultCodespace, CodeMissingField, fmt.Sprintf("asset quantity is not enough: {%s}", assetID))
}

// ErrInvalidRevokeRecipient is used when the recipient of
// a revoke proposal message is not in the asset's proposal list
func ErrInvalidRevokeRecipient(addr sdk.Address) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidRevokeRecipient, fmt.Sprintf("address %s is an invalid target for revoking proposal", addr.String()))
}

// InvalidTransaction ...
func ErrInvalidTransaction(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidTransaction, msg)
}

// CodeToDefaultMsg NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeUnknownAsset:
		return "Unknown asset"
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
