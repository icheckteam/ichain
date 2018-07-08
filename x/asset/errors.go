package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 500 - 599.
const (
	CodeUnknownAsset          sdk.CodeType      = 500
	CodeInvalidTransaction    sdk.CodeType      = 501
	CodeInvalidInput          sdk.CodeType      = 502
	CodeInvalidOutput         sdk.CodeType      = 503
	CodeInvalidAssets         sdk.CodeType      = 504
	CodeMissingField          sdk.CodeType      = 505
	CodeInvalidField          sdk.CodeType      = 506
	CodeInvalidRevokeReporter sdk.CodeType      = 507
	CodeInvalidAssetQuantity  sdk.CodeType      = 508
	CodeAssetAlreadyFinal     sdk.CodeType      = 509
	CodeProposalNotFound      sdk.CodeType      = 510
	DefaultCodespace          sdk.CodespaceType = 10
)

// ErrUnknownAsset ...
func ErrUnknownAsset(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeUnknownAsset, msg)
}

func ErrAssetNotFound(assetID string) sdk.Error {
	return newError(DefaultCodespace, CodeUnknownAsset, fmt.Sprintf("asset {%s} not found", assetID))
}

func ErrAssetAlreadyFinal(assetID string) sdk.Error {
	return newError(DefaultCodespace, CodeAssetAlreadyFinal, fmt.Sprintf("asset {%s} already final", assetID))
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
	return newError(DefaultCodespace, CodeMissingField, fmt.Sprintf("asset {%s} is not enough", assetID))
}

func ErrInvalidAssetRoot(assetID string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidTransaction, fmt.Sprintf("asset {%s} is not root", assetID))
}

// ErrInvalidRevokeReporter is used when the reporter of
// a revoke reporter message is not in the asset's reporter list
func ErrInvalidRevokeReporter(addr sdk.Address) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidRevokeReporter, fmt.Sprintf("address %s is an invalid target for revoking reporter", addr.String()))
}

// InvalidTransaction ...
func ErrInvalidTransaction(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidTransaction, msg)
}

func ErrProposalNotFound(recipient sdk.Address) sdk.Error {
	return newError(DefaultCodespace, CodeProposalNotFound, fmt.Sprintf("proposal %s not found", sdk.MustBech32ifyAcc(recipient)))
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
