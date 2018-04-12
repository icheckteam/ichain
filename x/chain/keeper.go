package chain

import (
	"github.com/cosmos-sdk/wire"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper ...
type Keeper struct {
	ck bank.CoinKeeper

	storeKey          sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc               *wire.Codec
	recordIndexNumber int
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper, cdc *wire.Codec) Keeper {
	return Keeper{
		storeKey:          key,
		recordIndexNumber: 0,
		ck:                bankKeeper,
		cdc:               cdc,
	}
}

func (k Keeper) createRecord(ctx sdk.Context, record Record) {
	store := ctx.KVStore(k.storeKey)
	// marshal the record and add to the state
	bz, err := k.cdc.MarshalBinary(record)
	if err != nil {
		panic(err)
	}

	store.Set(GetRecordKey([]byte(record.ID)), bz)
}

// GetRecord get record by IDS
func (k Keeper) GetRecord(ctx sdk.Context, uid string) Record {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetRecordKey([]byte(uid)))
	record := Record{}

	// marshal the record and add to the state
	if err := k.cdc.UnmarshalBinary(b, &record); err != nil {
		panic(err)
	}
	return record
}

// Transfer change owner
func (k Keeper) Transfer(ctx sdk.Context, fromAddress sdk.Address, toAddress sdk.Address, uid string) sdk.Error {
	record := k.GetRecord(ctx, uid)
	if record.ID == "" {
		return ErrUnknownRecord("Record not found")
	}

	// check record owner
	if record.Owner.String() != fromAddress.String() {
		return sdk.ErrUnauthorized(fromAddress.String())
	}
	record.Owner = toAddress
	k.createRecord(ctx, record)
	return nil
}
