package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// SetTrustHandlerFn
func SetTrustHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m msgSetTrustBody
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

// SetCertsHandlerFn
func SetCertsHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var m msgSetCertsBody
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

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		identityID, err := strconv.Atoi(vars["identityID"])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		msg := identity.MsgSetCerts{
			Certifier:  sdk.AccAddress(info.GetPubKey().Address()),
			IdentityID: int64(identityID),
			Values:     m.Values,
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}

// CreateIdentityHandlerFn
func CreateIdentityHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m msgCreateIdentityBody
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

		info, err := kb.Get(m.BaseReq.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		msg := identity.MsgCreateIdentity{
			Sender: sdk.AccAddress(info.GetPubKey().Address()),
		}
		signAndBuild(ctx, cdc, w, m.BaseReq, msg)
	}
}
