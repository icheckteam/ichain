package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// setCert set the main record holding cert details
func (k Keeper) setCert(ctx sdk.Context, addr sdk.AccAddress, cert Cert) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(cert)
	store.Set(KeyCert(addr, cert.Property, cert.Certifier), bz)
}

// GetCert  set the main record holding cert details
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

// deleteCert delete a cert from the store
func (k Keeper) deleteCert(ctx sdk.Context, addr sdk.AccAddress, property string, certifier sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyCert(addr, property, certifier))
}

// AddCerts add all certs
func (k Keeper) AddCerts(ctx sdk.Context, msg MsgSetCerts) (sdk.Tags, sdk.Error) {
	// check owner to add certificate
	if !k.hasOwner(ctx, msg.Issuer, msg.Sender) {
		return nil, sdk.ErrUnauthorized(fmt.Sprintf("addr %s unauthorized to add", msg.Sender))
	}

	for _, value := range msg.Values {
		if value.Confidence == true {
			cert, found := k.GetCert(ctx, value.Owner, value.Property, msg.Issuer)
			if !found {
				// new cert
				cert = Cert{
					Property:  value.Property,
					Owner:     value.Owner,
					Certifier: msg.Issuer,
					Data:      value.Data,
					CreatedAt: ctx.BlockHeader().Time,
				}
			} else {
				// update cert
				cert.Data = value.Data
			}

			// add cert
			k.setCert(ctx, value.Owner, cert)

		} else {
			// delete cert
			k.deleteCert(ctx, value.Owner, value.Property, msg.Issuer)
		}
	}
	return nil, nil
}

// GetCerts ...
func (k Keeper) GetCerts(ctx sdk.Context, id sdk.AccAddress) Certs {
	store := ctx.KVStore(k.storeKey)

	// delete subspace
	iterator := sdk.KVStorePrefixIterator(store, KeyCerts(id))
	certs := Certs{}
	for ; iterator.Valid(); iterator.Next() {
		cert := Cert{}
		k.cdc.MustUnmarshalBinary(iterator.Value(), &cert)
		certs = append(certs, cert)
	}
	iterator.Close()
	return certs
}
