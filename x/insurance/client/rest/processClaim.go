package rest

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/insurance"
)

type processClaimBody struct {
	AccountName string                `json:"name"`
	Password    string                `json:"password"`
	ChainID     string                `json:"chain_id"`
	Sequence    int64                 `json:"sequence"`
	ContractID  string                `json:"contract_id"`
	Status      insurance.ClaimStatus `json:"status"`
}

// ProcessClaimHandlerFn
func ProcessClaimHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := processClaimBody{}

		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if b.AccountName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("account_name is required"))
			return
		}

		if b.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("password is required"))
			return
		}

		if b.Status == insurance.ClaimStatusPending {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("status is required"))
			return
		}

		if b.ContractID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("contract_id is required"))
			return
		}

		info, err := kb.Get(b.AccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := buildMsgProcessClaim(
			info.GetPubKey().Address(),
			b.ContractID,
			b.Status,
		)

		// sign
		ctx = ctx.WithSequence(b.Sequence)
		txBytes, err := ctx.SignAndBuild(b.AccountName, b.Password, []sdk.Msg{msg}, cdc)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		// send
		res, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func buildMsgProcessClaim(issuer sdk.Address, contractID string, status insurance.ClaimStatus) insurance.MsgProcessClaim {
	return insurance.MsgProcessClaim{ContractID: contractID, Issuer: issuer, Status: status}
}
