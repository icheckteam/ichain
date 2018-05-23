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

type answerProposalBody struct {
	LocalAccountName string `json:"account_name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	Sequence         int64  `json:"sequence"`

	AssetID  string               `json:"asset_id"`
	Response asset.ProposalStatus `json:"response"`
}

// AnswerProposalHandlerFn  REST handler
func AnswerProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m answerProposalBody
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

		if len(m.AssetID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("asset_id is required"))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := buildAnswerProposalMsg(info.PubKey.Address(), m.AssetID, m.Response)
		if err != nil { // XXX rechecking same error ?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// sign
		ctx = ctx.WithSequence(m.Sequence).
			WithChainID(m.ChainID)
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

func buildAnswerProposalMsg(issuer sdk.Address, assetID string, response asset.ProposalStatus) sdk.Msg {
	return asset.AnswerProposalMsg{
		Recipient: issuer,
		Response:  response,
		AssetID:   assetID,
	}
}
