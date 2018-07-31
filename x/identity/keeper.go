package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *wire.Codec
	codespace sdk.CodespaceType
}

func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: DefaultCodespace,
	}
}

// set the main record holding trust details
func (k Keeper) SetTrust(ctx sdk.Context, trustor, trusting sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(Trust{Trusting: trusting, Trustor: trustor})
	store.Set(KeyTrust(trustor, trusting), bz)
}

// delete cert from the store
func (k Keeper) DeleteTrust(ctx sdk.Context, trustor, trusting sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyTrust(trustor, trusting))
}

// add a trust
func (k Keeper) AddTrust(ctx sdk.Context, msg MsgSetTrust) sdk.Error {
	if msg.Trust == true {
		k.SetTrust(ctx, msg.Trustor, msg.Trusting)
	} else {
		k.DeleteTrust(ctx, msg.Trustor, msg.Trusting)
	}
	return nil
}

func (k Keeper) GetTrust(ctx sdk.Context, trustor, trusting sdk.AccAddress) (trust Trust, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyTrust(trustor, trusting))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(bz, &trust)
	found = true
	return
}

// set the main record holding cert details
func (k Keeper) SetCert(ctx sdk.Context, addr sdk.AccAddress, cert Cert) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(cert)
	store.Set(KeyCert(addr, cert.Property, cert.Certifier), bz)
}

// set the main record holding cert details
func (k Keeper) GetCert(ctx sdk.Context, addr sdk.AccAddress, property string, certifier sdk.AccAddress) (cert Cert, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyCert(addr, property, certifier))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(bz, &cert)
	found = true
	return
}

// delete cert from the store
func (k Keeper) DeleteCert(ctx sdk.Context, addr sdk.AccAddress, property string, certifier sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyCert(addr, property, certifier))
}

// add a trusting
func (k Keeper) AddCerts(ctx sdk.Context, msg MsgSetCerts) sdk.Error {
	for _, value := range msg.Values {
		if value.Confidence == true {
			cert, found := k.GetCert(ctx, msg.Recipient, value.Property, msg.Certifier)
			if !found {
				// new cert
				cert = Cert{
					ID:         value.ID,
					Context:    value.Context,
					Property:   value.Property,
					Certifier:  msg.Certifier,
					Confidence: value.Confidence,
					Data:       value.Data,
					Expires:    value.Expires,
					CreatedAt:  ctx.BlockHeader().Time,
					Revocation: value.Revocation,
				}
			} else if value.Revocation.ID != "" {
				cert.Revocation = value.Revocation
			} else {
				// update cert
				cert.Data = value.Data
				cert.Expires = value.Expires
			}

			// add cert
			k.SetCert(ctx, msg.Recipient, cert)

		} else {
			// delete cert
			k.DeleteCert(ctx, msg.Recipient, value.Property, msg.Certifier)
		}
	}
	return nil
}
