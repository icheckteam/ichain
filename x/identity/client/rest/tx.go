package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// SetTrustHandlerFn ...
func SetTrustHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m msgSetTrustBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

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

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		trusting, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		msg := identity.MsgSetTrust{
			Trust:    m.Trust,
			Trusting: trusting,
			Trustor:  sdk.AccAddress(info.GetPubKey().Address()),
		}

		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}

// SetCertsHandlerFn ...
func SetCertsHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m msgSetCertsBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)

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

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// assign owner
		for i := range m.Values {
			m.Values[i].Owner = address
		}
		certifier := sdk.AccAddress(info.GetPubKey().Address())
		msg := identity.MsgSetCerts{
			Sender: certifier,
			Issuer: certifier,
			Values: m.Values,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}

func registerHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) http.HandlerFunc {
	return withErr(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m msgRegBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			return err
		}
		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}
		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		ident, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			return err
		}
		msg := identity.MsgReg{
			Sender: sdk.AccAddress(info.GetPubKey().Address()),
			Ident:  ident,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

func addOwnerHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) http.HandlerFunc {
	return withErr(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m msgAddOwnerBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			return err
		}
		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}
		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		ident, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			return err
		}
		msg := identity.MsgAddOwner{
			Sender: sdk.AccAddress(info.GetPubKey().Address()),
			Ident:  ident,
			Owner:  m.Owner,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}

func delOwnerHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) http.HandlerFunc {
	return withErr(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		var m msgDelOwnerBody
		body, err := ioutil.ReadAll(r.Body)
		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			return err
		}
		err = m.BaseReq.Validate()
		if err != nil {
			return err
		}
		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			return err
		}
		ident, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			return err
		}
		owner, err := sdk.AccAddressFromBech32(vars["owner"])
		if err != nil {
			return err
		}
		msg := identity.MsgDelOwner{
			Sender: sdk.AccAddress(info.GetPubKey().Address()),
			Ident:  ident,
			Owner:  owner,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
		return nil
	})
}
