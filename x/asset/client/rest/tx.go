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

func addAssetQuantityHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m addAssetQuantityBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message
		msg := asset.MsgAddQuantity{
			Sender:   sdk.AccAddress(info.GetPubKey().Address()),
			AssetID:  vars["id"],
			Quantity: m.Quantity,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

// AddMaterialsHandlerFn  REST handler
func addMaterialsHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m addMaterialsBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}
		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}

		msg := asset.MsgAddMaterials{
			AssetID: vars["id"],
			Sender:  sdk.AccAddress(info.GetPubKey().Address()),
			Amount:  m.Amount,
		}

		// sign
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

// FinalizeHandlerFn ...
func finalizeHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		var m finalizeBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message
		msg := asset.MsgFinalize{
			Sender:  sdk.AccAddress(info.GetPubKey().Address()),
			AssetID: vars["id"],
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

// Create asset REST handler
func createAssetHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		var m createAssetBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}

		// build message
		msg := asset.MsgCreateAsset{
			AssetID:    m.AssetID,
			Name:       m.Name,
			Parent:     m.Parent,
			Properties: m.Properties,
			Sender:     sdk.AccAddress(info.GetPubKey().Address()),
			Quantity:   m.Quantity,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

func revokeReporterHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m revokeReporterBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			return err
		}

		// build message

		msg := asset.MsgRevokeReporter{
			Sender:   sdk.AccAddress(info.GetPubKey().Address()),
			Reporter: address,
			AssetID:  vars["id"],
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

func subtractQuantityBodyHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m subtractAssetQuantityBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message

		msg := asset.MsgSubtractQuantity{
			Sender:   sdk.AccAddress(info.GetPubKey().Address()),
			AssetID:  vars["id"],
			Quantity: m.Quantity,
		}

		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

func updateAttributeHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return withErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m updateAttributeBody
		if err := validateAndGetDecodeBody(r, cdc, &m); err != nil {
			return err
		}

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		// build message

		msg := asset.MsgUpdateProperties{
			AssetID:    vars["id"],
			Properties: m.Properties,
			Sender:     sdk.AccAddress(info.GetPubKey().Address()),
		}

		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
