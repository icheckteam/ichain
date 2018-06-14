package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/tendermint/go-crypto/keys"
)

type createProposalBody struct {
	baseBody
	Msg asset.CreateProposalMsg `json:"msg"`
}

// CreateProposalHandlerFn CreateProposal REST handler
func CreateProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m createProposalBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.LocalAccountName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("account_name is required"))
			return
		}

		if m.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("password is required"))
			return
		}

		if len(m.Msg.Recipient) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("recipient is required"))
			return
		}

		if len(m.Msg.AssetID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("asset_id is required"))
			return
		}

		if len(m.Msg.Propertipes) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("propertipes is required"))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		m.Msg.Issuer = info.PubKey.Address()
		msg := m.Msg

		ctx = ctx.WithGas(m.Gas)
		ctx = ctx.WithAccountNumber(m.AccountNumber)
		ctx = ctx.WithSequence(m.Sequence)

		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, msg, cdc)
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

func buildCreateProposalMsg(issuer, recipient sdk.Address, assetID string, propertipes []string, role asset.ProposalRole) sdk.Msg {
	return asset.CreateProposalMsg{
		Issuer:      issuer,
		Recipient:   recipient,
		AssetID:     assetID,
		Propertipes: propertipes,
		Role:        role,
	}
}
