package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/insurance"
)

type createContractBody struct {
	AccountName string       `json:"name"`
	Password    string       `json:"password"`
	Contract    contractBody `json:"contract"`
	ChainID     string       `json:"chain_id"`
	Sequence    int64        `json:"sequence"`
}

type contractBody struct {
	ID        string      `json:"id"`
	AssetID   string      `json:"asset_id"`
	Expires   time.Time   `json:"expires"`
	Serial    string      `json:"serial"`
	Recipient sdk.Address `json:"recipient"`
}

// CreateContractHandlerFn
func CreateContractHandlerFn(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := createContractBody{}

		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if b.AccountName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("account_name is required"))
			return
		}

		if b.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("password is required"))
			return
		}

		if b.Contract.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("contract.id is required"))
			return
		}

		if b.Contract.Expires.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("contract.expires is required"))
			return
		}

		if b.Contract.Recipient.String() == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("contract.recipient is required"))
			return
		}

		if b.Contract.Serial == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("contract.serial is required"))
			return
		}

		info, err := kb.Get(b.AccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := buildMsgCreateContract(
			info.GetPubKey().Address(),
			b.Contract.Recipient,
			b.Contract.AssetID,
			b.Contract.Serial,
			b.Contract.Expires,
		)

		// sign
		ctx = ctx.WithSequence(b.Sequence)
		txBytes, err := ctx.SignAndBuild(b.AccountName, b.Password, []sdk.Msg{msg}, cdc)
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

func buildMsgCreateContract(issuer sdk.Address, recipient sdk.Address, assetID, serial string, expires time.Time) insurance.MsgCreateContract {
	return insurance.MsgCreateContract{
		ID:        assetID,
		Issuer:    issuer,
		Recipient: recipient,
		Serial:    serial,
		Expires:   expires,
	}
}
