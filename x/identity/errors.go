package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultCodespace sdk.CodespaceType = 12

	CodeUnknownIdentity sdk.CodeType = 1
	CodeInvalidGenesis  sdk.CodeType = 2
)

//----------------------------------------
// Error constructors

func ErrUnknownIdentity(codespace sdk.CodespaceType, identityID int64) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownIdentity, fmt.Sprintf("Unknown proposal - %d", identityID))
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}
