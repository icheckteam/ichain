package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
)

const storeName = "identity"

type IdentOutput struct {
	ID    int64  `json:"id"`    // id of the identity
	Owner string `json:"owner"` // owner of the identity
}

func bech32IdentOutput(ident identity.Identity) IdentOutput {
	return IdentOutput{
		ID:    ident.ID,
		Owner: sdk.MustBech32ifyAcc(ident.Owner),
	}
}

type CertOutput struct {
	ID         string            `json:"id"`
	Property   string            `json:"property"`
	Certifier  string            `json:"certifier"`
	Type       string            `json:"type"`
	Trust      bool              `json:"trust"`
	Data       identity.Metadata `json:"data"`
	Confidence bool              `json:"confidence"`
}

func bech32CertOutput(cert identity.Cert) CertOutput {
	return CertOutput{
		ID:         cert.ID,
		Property:   sdk.MustBech32ifyAcc(cert.Property),
		Certifier:  sdk.MustBech32ifyAcc(cert.Certifier),
		Type:       cert.Type,
		Data:       cert.Data,
		Confidence: cert.Confidence,
		Trust:      cert.Trust,
	}
}

func identsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kvs, err := ctx.QuerySubspace(cdc, identity.IdentitiesKey, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query idents. Error: %s", err.Error())))
			return
		}

		idents := make([]IdentOutput, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			ident := identity.Identity{}
			err = cdc.UnmarshalBinary(addr, &ident)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
				return
			}

			bech32Ident := bech32IdentOutput(ident)
			idents[i] = bech32Ident
		}

		output, err := cdc.MarshalJSON(idents)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func certsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		identID, err := strconv.Atoi(vars["ident_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode ident_id. Error: %s", err.Error())))
			return
		}
		kvs, err := ctx.QuerySubspace(cdc, identity.KeyCerts(int64(identID)), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query idents. Error: %s", err.Error())))
			return
		}

		certs := make([]CertOutput, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			cert := identity.Cert{}
			err = cdc.UnmarshalBinary(addr, &cert)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
				return
			}

			bech32Cert := bech32CertOutput(cert)
			certs[i] = bech32Cert
		}

		output, err := cdc.MarshalJSON(certs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
