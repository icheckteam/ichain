package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *wire.Codec
	codespace sdk.CodespaceType
	// The ValidatorSet to get information about validators
	vs sdk.ValidatorSet
}

func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

func (k Keeper) NewIdentity(ctx sdk.Context, owner sdk.Address) Identity {
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

func (k Keeper) DeleteIdentity(ctx sdk.Context, identityID int64) {
	// delete identity
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyIdentity(identityID))
	// delete all certs
	certs := k.GetCerts(ctx, identityID)
	for _, cert := range certs {
		k.DeleteCert(ctx, identityID, cert.Certifier)
	}
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
	store.Set(KeyIdentitiesByOwnerIndex(identity.Owner, identity.ID), bz)
}

// set the main record holding cert details
func (k Keeper) SetCert(ctx sdk.Context, identity int64, cert Cert) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(cert)
	store.Set(KeyCert(identity, cert.Certifier), bz)
}

func (k Keeper) GetCerts(ctx sdk.Context, identity int64) (certs []Cert) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyCerts(identity))

	i := 0
	for ; ; i++ {
		if !iterator.Valid() {
			iterator.Close()
			break
		}
		bz := iterator.Value()
		var cert Cert
		k.cdc.MustUnmarshalBinary(bz, &cert)
		certs = append(certs, cert)
		iterator.Next()
	}
	return certs
}

// delete cert from the store
func (k Keeper) DeleteCert(ctx sdk.Context, identity int64, certifier sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyCert(identity, certifier))
}

// set the main record holding trust details
func (k Keeper) SetTrust(ctx sdk.Context, trustor, trusting sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(Trust{Trusting: trusting, Trustor: trustor})
	store.Set(KeyTrust(trustor, trusting), bz)
	store.Set(KeyTrusting(trustor, trusting), bz)
}

// delete cert from the store
func (k Keeper) DeleteTrust(ctx sdk.Context, trustor, trusting sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyTrust(trustor, trusting))
	store.Delete(KeyTrusting(trustor, trusting))
}

func (k Keeper) HasTrusting(ctx sdk.Context, trustor, trusting sdk.Address) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyTrusting(trustor, trusting))
}

// add a trusting
func (k Keeper) AddTrust(ctx sdk.Context, msg MsgAddTrust) sdk.Error {
	k.SetTrust(ctx, msg.Trustor, msg.Trusting)
	return nil
}

func (k Keeper) IsTrust(ctx sdk.Context, certifier sdk.Address) bool {
	validator := k.vs.Validator(ctx, certifier)
	if validator == nil {
		return k.HasTrusting(ctx, validator.GetOwner(), certifier)
	}
	return true
}

// add a trusting
func (k Keeper) AddCerts(ctx sdk.Context, msg MsgAddCerts) sdk.Error {
	_, found := k.GetIdentity(ctx, msg.IdentityID)
	if !found {
		return ErrUnknownIdentity(k.codespace, msg.IdentityID)
	}
	trust := k.IsTrust(ctx, msg.Certifier)
	for _, value := range msg.Values {
		if value.Confidence == true {
			// add cert
			k.SetCert(ctx, msg.IdentityID, Cert{
				Certifier:  msg.Certifier,
				Confidence: value.Confidence,
				Data:       value.Data,
				Type:       value.Type,
				Trust:      trust,
			})
		} else {
			// delete cert
			k.DeleteCert(ctx, msg.IdentityID, msg.Certifier)
		}
	}
	return nil
}

func (k Keeper) getNewIdentityID(ctx sdk.Context) (identityID int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextIdentityID)
	if bz == nil {
		return 1
	}
	k.cdc.MustUnmarshalBinary(bz, &identityID)
	bz = k.cdc.MustMarshalBinary(identityID + 1)
	store.Set(KeyNextIdentityID, bz)
	return identityID
}

func (k Keeper) AddIdentity(ctx sdk.Context, msg MsgCreateIdent) sdk.Error {
	k.SetIdentity(ctx, k.NewIdentity(ctx, msg.Sender))
	return nil
}
