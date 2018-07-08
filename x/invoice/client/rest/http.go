package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"

	cosmosContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/invoice"
)

var (
	StoreName            = "invoice"
	ErrorInvoiceNotFound = errors.New("Invoice not found.")
	ErrorUnauthorized    = errors.New("Unauthorized.")
)

type decodeRequest func(context.Context, *http.Request) (interface{}, error)

type endpoint func(context.Context, interface{}) (interface{}, error)

type createInvoiceRequest struct {
	Account  string         `json:"account"`
	Password string         `json:"password"`
	ChainID  string         `json:"chain_id"`
	Sequence int64          `json:"sequence"`
	ID       string         `json:"id"`
	Receiver string         `json:"receiver"`
	Items    []invoice.Item `json:"items"`
}

type getInvoiceRequest struct {
	ID string `json:"id"`
}

func decodeCreateInvoiceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createInvoiceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetInvoiceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		return nil, ErrorInvoiceNotFound
	}

	return getInvoiceRequest{ID: id}, nil
}

func code(err error) int {
	switch err {
	case ErrorInvoiceNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func encodeSuccess(ctx context.Context, w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func encodeError(ctx context.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(code(err))
	encodeSuccess(ctx, w, map[string]interface{}{
		"error": err.Error(),
	})
}

func makeHandle(e endpoint, d decodeRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := d(ctx, r)

		if err != nil {
			encodeError(ctx, w, err)
			return
		}

		res, err := e(ctx, req)

		if err != nil {
			encodeError(ctx, w, err)
			return
		}

		encodeSuccess(ctx, w, res)
	}
}

func makeCreateInvoiceEndpoint(ctx cosmosContext.CoreContext, cdc *wire.Codec, kb keys.Keybase) endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createInvoiceRequest)

		info, err := kb.Get(req.Account)

		if err != nil {
			return nil, ErrorUnauthorized
		}

		issuer := info.GetPubKey().Address()
		receiver, _ := sdk.GetAccAddressHex(req.Receiver)
		msg := invoice.NewMsgCreate(
			req.ID,
			issuer,
			receiver,
			req.Items,
		)

		ctx = ctx.WithSequence(req.Sequence).WithChainID(req.ChainID)
		txBytes, err := ctx.SignAndBuild(req.Account, req.Password, []sdk.Msg{msg}, cdc)

		if err != nil {
			return nil, ErrorUnauthorized
		}

		tx, err := ctx.BroadcastTx(txBytes)

		if err != nil {
			return nil, err
		}

		return json.MarshalIndent(tx, "", "  ")
	}
}

func makeGetInvoiceEndpoint(ctx cosmosContext.CoreContext, cdc *wire.Codec) endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getInvoiceRequest)
		key := invoice.GetKey(req.ID)
		res, err := ctx.QueryStore(key, StoreName)

		var invoice invoice.Invoice
		err = cdc.UnmarshalBinary(res, &invoice)

		if err != nil {
			return nil, err
		}

		output, err := cdc.MarshalJSON(invoice)

		if err != nil {
			return nil, err
		}

		return output, nil
	}
}

func RegisterHTTPHandle(r *mux.Router, ctx cosmosContext.CoreContext, cdc *wire.Codec, kb keys.Keybase) {
	createInvoiceEndpoint := makeCreateInvoiceEndpoint(ctx, cdc, kb)
	r.HandleFunc("/invoices", makeHandle(createInvoiceEndpoint, decodeCreateInvoiceRequest))

	getInvoiceEndpoint := makeGetInvoiceEndpoint(ctx, cdc)
	r.HandleFunc("/invoices/{id}", makeHandle(getInvoiceEndpoint, decodeGetInvoiceRequest))
}
