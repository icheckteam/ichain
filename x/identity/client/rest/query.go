package rest

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
)

const storeName = "identity"

func queryCertsByOwner(ctx context.CoreContext, cdc *wire.Codec, owner sdk.AccAddress, vars map[string]string) (identity.Certs, error) {
	bechCertifier := vars["certifier"]
	property := vars["property"]
	trust := vars["trust"] == "1"
	var err error
	var certifierAddr sdk.AccAddress

	if len(bechCertifier) != 0 {
		certifierAddr, err = sdk.AccAddressFromBech32(bechCertifier)
		if err != nil {
			return nil, err
		}
	}

	kvs, err := ctx.QuerySubspace(cdc, identity.KeyCerts(owner), storeName)
	if err != nil {
		return nil, err
	}

	validators, err := getValidators(ctx, cdc)
	if err != nil {
		return nil, err
	}
	certs := []identity.Cert{}
	for _, kv := range kvs {
		cert := identity.Cert{}
		err = cdc.UnmarshalBinary(kv.Value, &cert)
		if err != nil {
			return nil, err
		}

		if len(bechCertifier) != 0 {
			if !bytes.Equal(certifierAddr, cert.Certifier) {
				continue
			}
		}

		if len(property) != 0 {
			if property != cert.Property {
				continue
			}
		}

		// check trust
		for _, validator := range validators {
			if bytes.Equal(validator.Owner, cert.Certifier) {
				cert.Trust = true
				break
			}
		}

		if cert.Trust == false {
			for _, validator := range validators {
				if hasTrust(ctx, cdc, validator.Owner, cert.Certifier) {
					cert.Trust = true
					break
				}
			}
		}

		if len(vars["trust"]) != 0 {
			if cert.Trust != trust {
				continue
			}
		}

		certs = append(certs, cert)
	}
	return certs, nil
}

func queryCertsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		certs, err := queryCertsByOwner(ctx, cdc, address, vars)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query certs. Error: %s", err.Error())))
			return
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

func trustsHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode address. Error: %s", err.Error())))
			return
		}
		kvs, err := ctx.QuerySubspace(cdc, identity.KeyTrusts(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query trusts. Error: %s", err.Error())))
			return
		}

		trusts := make([]identity.Trust, len(kvs))
		for i, kv := range kvs {
			trust := identity.Trust{}
			err = cdc.UnmarshalBinary(kv.Value, &trust)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode trust. Error: %s", err.Error())))
				return
			}

			trusts[i] = trust
		}

		output, err := cdc.MarshalJSON(trusts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func hasTrust(ctx context.CoreContext, cdc *wire.Codec, trustor, trusting sdk.AccAddress) bool {
	res, err := ctx.QueryStore(identity.KeyTrust(trustor, trusting), "identity")
	if err != nil {
		panic(err)
	}

	if len(res) > 0 {
		return true
	}

	return false
}
