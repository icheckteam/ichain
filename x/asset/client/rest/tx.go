package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
)

// CreateProposalHandlerFn
func CreateProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m msgCreateCreateProposalBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = m.BaseReq.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.Recipient == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("recipient is required"))
			return
		}

		switch m.Role {
		case asset.RoleOwner, asset.RoleReporter:
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid role"))
			return
		}

		if m.Role == asset.RoleReporter && len(m.Properties) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("properties is required"))
			return
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		address, err := sdk.GetAccAddressBech32(m.Recipient)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		msg := asset.MsgCreateProposal{
			Sender:     info.GetPubKey().Address(),
			Properties: m.Properties,
			Role:       m.Role,
			AssetID:    vars["id"],
			Recipient:  address,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}

// AnswerProposalHandlerFn
func AnswerProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m msgAnswerProposalBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = m.BaseReq.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		switch m.Response {
		case asset.StatusAccepted, asset.StatusCancel, asset.StatusRejected:
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid reponse"))
			return
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		recipient, err := sdk.GetAccAddressBech32(vars["recipient"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		msg := asset.MsgAnswerProposal{
			Sender:    info.GetPubKey().Address(),
			Recipient: recipient,
			Response:  m.Response,
			AssetID:   vars["id"],
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}
