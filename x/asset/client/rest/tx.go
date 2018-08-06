package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/asset"
)

// CreateProposalHandlerFn
func createProposalHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m msgCreateCreateProposalBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
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
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
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
			Role:      m.Role,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil

	})
}
