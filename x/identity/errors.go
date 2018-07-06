package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ABCI Response Codes
// Base SDK reserves 600-700.
const (
	CodeInvalidClaim   sdk.CodeType      = 600
	CodeInvalidExpires sdk.CodeType      = 601
	DefaultCodespace   sdk.CodespaceType = 11
)

func ErrClaimNotFound(claimID string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidClaim, fmt.Sprintf("claim {%s} not found", claimID))
}

func ErrClaimHasPaid(claimID string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidClaim, fmt.Sprintf("claim {%s} has paid", claimID))
}

func ErrInvalidExpires(expires int64) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidExpires, fmt.Sprintf("invalid expires time: %d", expires))
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(codespace, code, msg)
}
