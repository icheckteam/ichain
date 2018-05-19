package rest

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/warranty"
	"github.com/tendermint/go-crypto/keys"
)

type createClaimBody struct {
	AccountName string    `json:"name"`
	Password    string    `json:"password"`
	Claim       claimBody `json:"contract"`
	ChainID     string    `json:"chain_id"`
	Sequence    int64     `json:"sequence"`
}

type claimBody struct {
	ContractID string `json:"contract_id"`
	Recipient  sdk.Address
}

// CreateClaimHandlerFn
func CreateClaimHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := createClaimBody{}

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

		if b.Claim.ContractID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("claim.contract_id is required"))
			return
		}

		if b.Claim.Recipient.String() == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("contract.recipient is required"))
			return
		}

		info, err := kb.Get(b.AccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := buildMsgCreateClaim(
			info.PubKey.Address(),
			b.Claim.Recipient,
			b.Claim.ContractID,
		)

		// sign
		ctx = ctx.WithSequence(b.Sequence)
		txBytes, err := ctx.SignAndBuild(b.AccountName, b.Password, msg, cdc)
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

func buildMsgCreateClaim(issuer sdk.Address, recipient sdk.Address, contractID string) warranty.MsgCreateClaim {
	return warranty.NewMsgCreateClaim(issuer, recipient, contractID)
}
