package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/icheckteam/ichain/x/shipping"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type createOrderBody struct {
	LocalAccountName string                      `json:"account_name"`
	Password         string                      `json:"password"`
	ChainID          string                      `json:"chain_id"`
	Sequence         int64                       `json:"sequence"`
	OrderID          string                      `json:"order_id"`
	Assets           []shipping.TransportedAsset `json:"assets"`
	Carrier          string                      `json:"carrier_address"`
	Receiver         string                      `json:"receiver_address"`
}

// CreateOrderHandlerFn ...
func CreateOrderHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m createOrderBody
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

		if len(m.Carrier) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("carrier_address is required"))
			return
		}

		if len(m.Receiver) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("receiver_address is required"))
			return
		}

		if len(m.Assets) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("assets is required"))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		// build message
		msg := buildCreateOrderMsg(sdk.AccAddress(info.GetPubKey().Address()), m)

		// sign
		ctx = ctx.WithSequence(m.Sequence).WithChainID(m.ChainID)
		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, []sdk.Msg{msg}, cdc)
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

func buildCreateOrderMsg(creator sdk.AccAddress, body createOrderBody) sdk.Msg {
	carrier, _ := sdk.AccAddressFromBech32(body.Carrier)
	receiver, _ := sdk.AccAddressFromBech32(body.Receiver)

	return shipping.CreateOrderMsg{
		ID:                body.OrderID,
		TransportedAssets: body.Assets,
		Issuer:            creator,
		Carrier:           carrier,
		Receiver:          receiver,
	}
}
