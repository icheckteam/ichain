package identity

import (
	"bytes"

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

func (k Keeper) NewIdentity(ctx sdk.Context, owner sdk.AccAddress) Identity {
	return Identity{
		ID:    k.getNewIdentityID(ctx),
		Owner: owner,
	}
}

// set the main record holding identity details
func (k Keeper) SetIdentity(ctx sdk.Context, identity Identity) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(identity)
	store.Set(KeyIdentity(identity.ID), bz)
}

func (k Keeper) SetClaimedIdentity(ctx sdk.Context, account sdk.AccAddress, identity Identity) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(identity)
	store.Set(KeyClaimedIdentity(account), bz)
}

func (k Keeper) DeleteClaimedIdentity(ctx sdk.Context, account sdk.AccAddress, identityID int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyClaimedIdentity(account))
}

func (k Keeper) HasClaimedIdentity(ctx sdk.Context, account sdk.AccAddress, identityID int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyClaimedIdentity(account))
}

// Get Identity from store by identityID
func (k Keeper) GetIdentity(ctx sdk.Context, identityID int64) (Identity, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyIdentity(identityID))
	if bz == nil {
		return Identity{}, false
	}

	var identity Identity
	k.cdc.MustUnmarshalBinary(bz, &identity)

	return identity, true
}

// set the main record holding
func (k Keeper) SetIdentityByOwnerIndex(ctx sdk.Context, identity Identity) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(identity.ID)
	store.Set(KeyIdentityByOwnerIndex(identity.Owner, identity.ID), bz)
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
func (k Keeper) SetCert(ctx sdk.Context, identity int64, cert Cert) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(cert)
	store.Set(KeyCert(identity, cert.Property, cert.Certifier), bz)
}

// set the main record holding cert details
func (k Keeper) GetCert(ctx sdk.Context, identity int64, property string, certifier sdk.AccAddress) (cert Cert, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyCert(identity, property, certifier))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(bz, &cert)
	found = true
	return
}

// delete cert from the store
func (k Keeper) DeleteCert(ctx sdk.Context, identity int64, property string, certifier sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyCert(identity, property, certifier))
}

// add a trusting
func (k Keeper) AddCerts(ctx sdk.Context, msg MsgSetCerts) sdk.Error {
	ident, found := k.GetIdentity(ctx, msg.IdentityID)
	if !found {
		return ErrUnknownIdentity(k.codespace, msg.IdentityID)
	}
	for _, value := range msg.Values {
		if value.Confidence == true {
			cert, found := k.GetCert(ctx, msg.IdentityID, value.Property, msg.Certifier)
			if !found {
				// new cert
				cert = Cert{
					Property:   value.Property,
					Certifier:  msg.Certifier,
					Confidence: value.Confidence,
					Data:       value.Data,
					Type:       value.Type,
				}
			} else {
				// update cert
				cert.Data = value.Data
				cert.Type = value.Type
			}

			// add cert
			k.SetCert(ctx, msg.IdentityID, cert)

			// special handling for owner
			if value.Property == "owner" {
				if bytes.Equal(msg.Certifier, ident.Owner) {
					k.SetClaimedIdentity(ctx, msg.Certifier, ident)
				}
			}

		} else {
			// delete cert
			k.DeleteCert(ctx, msg.IdentityID, value.Property, msg.Certifier)

			if value.Property == "owner" {
				if bytes.Equal(msg.Certifier, ident.Owner) {
					k.DeleteClaimedIdentity(ctx, msg.Certifier, ident.ID)
				}
			}
		}
	}
	return nil
}

func (k Keeper) getNewIdentityID(ctx sdk.Context) (identityID int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextIdentityID)
	if bz == nil {
		identityID = 1
	} else {
		k.cdc.MustUnmarshalBinary(bz, &identityID)
	}
	bz = k.cdc.MustMarshalBinary(identityID + 1)
	store.Set(KeyNextIdentityID, bz)
	return identityID
}

// AddIdentity add new an identity
func (k Keeper) AddIdentity(ctx sdk.Context, msg MsgCreateIdentity) sdk.Error {
	ident := k.NewIdentity(ctx, msg.Sender)
	k.SetIdentity(ctx, ident)
	k.SetIdentityByOwnerIndex(ctx, ident)
	return nil
}
