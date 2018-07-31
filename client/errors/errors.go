package errors

import (
	"encoding/json"
	"errors"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/icheckteam/ichain/x/identity"
)

var ErrUnauthorized = New("Unauthorized")

func New(err string) error {
	return errors.New(err)
}

var errorCodes map[sdk.CodeType]string = map[sdk.CodeType]string{
	asset.CodeInvalidTransaction:    "invalid_transaction",
	asset.CodeAssetAlreadyFinal:     "asset_already_final",
	asset.CodeInvalidField:          "invalid_field",
	asset.CodeInvalidRevokeReporter: "invalid_reporter",
	asset.CodeInvalidAssets:         "invalid_assets",
	asset.CodeMissingField:          "missing_fields",
	asset.CodeProposalNotFound:      "proposal_not_found",
	asset.CodeUnknownAsset:          "unknow_asset",
	bank.CodeInvalidInput:           "invalid_input",
	bank.CodeInvalidOutput:          "invalid_output",
	sdk.CodeInvalidCoins:            "invalid_coin",
	sdk.CodeMemoTooLarge:            "memo_too_large",
	sdk.CodeInsufficientCoins:       "insufficient_coins",
}

var codespaces map[sdk.CodespaceType]string = map[sdk.CodespaceType]string{
	asset.DefaultCodespace:    "asset",
	bank.DefaultCodespace:     "auth",
	identity.DefaultCodespace: "stake",
}

type Error struct {
	CodeSpace string `json:"code_space"`
	Code      string `json:"code"`
	Message   string `json:"msg"`
}

type Result struct {
	Error Error `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	result := Result{
		Error: Error{},
	}
	switch e := err.(type) {
	case sdk.Error:
		result.Error.Code = errorCodes[e.Code()]
		result.Error.CodeSpace = codespaces[e.Codespace()]
	default:
		result.Error.Message = e.Error()
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}
