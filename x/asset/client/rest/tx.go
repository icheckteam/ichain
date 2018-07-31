package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/client/errors"
	"github.com/icheckteam/ichain/x/asset"
)

// CreateProposalHandlerFn
func createProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m msgCreateCreateProposalBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			return err
		}

		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}

		if m.Recipient == "" {
			return errors.New("recipient is required")
		}

		switch m.Role {
		case asset.RoleOwner, asset.RoleReporter:
			break
		default:
			return errors.New("invalid role")
		}

		if m.Role == asset.RoleReporter && len(m.Properties) == 0 {
			return errors.New("properties is required")
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		address, err := sdk.AccAddressFromBech32(m.Recipient)
		if err != nil {
			return err
		}
		msg := asset.MsgCreateProposal{
			Sender:     sdk.AccAddress(info.GetPubKey().Address()),
			Properties: m.Properties,
			Role:       m.Role,
			AssetID:    vars["id"],
			Recipient:  address,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

// AnswerProposalHandlerFn
func answerProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m msgAnswerProposalBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

		if err != nil {
			return err
		}

		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}

		switch m.Response {
		case asset.StatusAccepted, asset.StatusCancel, asset.StatusRejected:
			break
		default:
			return errors.New("invalid response")
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		recipient, err := sdk.AccAddressFromBech32(vars["recipient"])
		if err != nil {
			return err
		}
		msg := asset.MsgAnswerProposal{
			Sender:    sdk.AccAddress(info.GetPubKey().Address()),
			Recipient: recipient,
			Response:  m.Response,
			AssetID:   vars["id"],
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil

	})
}
