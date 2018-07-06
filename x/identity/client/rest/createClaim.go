package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/identity"
	"github.com/tendermint/go-crypto/keys"
)

type createClaimBody struct {
	baseBody

	// claim ...
	ClaimID   string           `json:"claim_id"`
	Context   string           `json:"context"`
	Content   identity.Content `json:"content"`
	Recipient string           `json:"recipient"`
	Expires   int64            `json:"expires"`
	Fee       sdk.Coins        `json:"fee"`
}

func CreateClaimHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m createClaimBody
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.LocalAccountName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("name is required"))
			return
		}

		if m.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("password is required"))
			return
		}

		if m.ClaimID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("claim_id is required"))
			return
		}

		if len(m.Context) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("context is required"))
			return
		}

		if len(m.Content) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("content is required"))
			return
		}

		if len(m.Recipient) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("metadata.recipient is required"))
			return
		}

		if m.Expires == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("metadata.expires is required"))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		recipient, err := sdk.GetAccAddressBech32(m.Recipient)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := identity.MsgCreateClaim{
			ClaimID:   m.ClaimID,
			Issuer:    info.PubKey.Address(),
			Recipient: recipient,
			Context:   m.Context,
			Content:   m.Content,
			Fee:       m.Fee,
			Expires:   m.Expires,
		}
		signAndBuild(ctx, cdc, w, m.baseBody, msg)
	}
}
