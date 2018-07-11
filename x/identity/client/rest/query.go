package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/icheckteam/ichain/x/identity"
	"github.com/pkg/errors"
)

const storeName = "identity"

func identsByAccountHandlerFn(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, err := sdk.AccAddressFromBech32(vars[RestAccount])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		kvs, err := ctx.QuerySubspace(cdc, identity.KeyIdentitiesByOwnerIndex(address), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query idents. Error: %s", err.Error())))
			return
		}

		idents := make([]identity.Identity, len(kvs))
		for i, kv := range kvs {
			var identID int64
			err = cdc.UnmarshalBinary(kv.Value, &identID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("couldn't decode ident id. Error: %s", err.Error())))
				return
			}

			res, err := ctx.QueryStore(identity.KeyIdentity(identID), "identity")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("couldn't query ident. Error: %s", err.Error())))
				return
			}

			ident := identity.Identity{}
			err = cdc.UnmarshalBinary(res, &ident)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode asset. Error: %s", err.Error())))
				return
			}

			idents[i] = ident
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

		bechCertifier := vars["certifier"]
		property := vars["property"]
		trust := vars["trust"] == "1"
		var err error
		var certifierAddr sdk.AccAddress

		if len(bechCertifier) != 0 {
			certifierAddr, err = sdk.AccAddressFromBech32(bechCertifier)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := errors.Errorf("'%s' needs to be bech32 encoded", "certifier")
				w.Write([]byte(err.Error()))
				return
			}
		}

		identID, err := strconv.Atoi(vars[RestIdentityID])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode ident_id. Error: %s", err.Error())))
			return
		}
		kvs, err := ctx.QuerySubspace(cdc, identity.KeyCerts(int64(identID), vars["property"]), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query idents. Error: %s", err.Error())))
			return
		}

		validators, err := getValidators(ctx, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query validatorss. Error: %s", err.Error())))
			return
		}
		certs := []identity.Cert{}
		for _, kv := range kvs {
			cert := identity.Cert{}
			err = cdc.UnmarshalBinary(kv.Value, &cert)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't encode cert. Error: %s", err.Error())))
				return
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
