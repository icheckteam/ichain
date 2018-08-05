package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// DefaultCodespace ...
	DefaultCodespace sdk.CodespaceType = 12

	// CodeInvalidID ...
	CodeInvalidID sdk.CodeType = 1
	// CodeInvalidGenesis ...
	CodeInvalidGenesis sdk.CodeType = 2
	// CodeInvalidTrustor ...
	CodeInvalidTrustor sdk.CodeType = 3
	// CodeInvalidTrusting ...
	CodeInvalidTrusting sdk.CodeType = 4
	// CodeInvalidInput ...
	CodeInvalidInput sdk.CodeType = 5
)

//----------------------------------------
// Error constructors

// ErrIDAlreadyExists ...
func ErrIDAlreadyExists(codespace sdk.CodespaceType, id sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidID, fmt.Sprintf("id %d already exists", id))
}

// ErrInvalidGenesis ...
func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

// ErrNilTrustorAddr ...
func ErrNilTrustorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrustor, "trustor address is nil")
}

// ErrNilTrustingAddr ...
func ErrNilTrustingAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrusting, "trusting address is nil")
}
