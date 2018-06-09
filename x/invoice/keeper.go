package invoice

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	CodespaceDefault     sdk.CodespaceType = 69
	CodeDuplicateInvoice sdk.CodeType      = 1
)

var (
	PrefixKey             = []byte{0x01}
	ErrorDuplicateInvoice = sdk.NewError(CodespaceDefault, CodeDuplicateInvoice, "Duplicate invoice.")
)

func GetKey(id string) []byte {
	return append(PrefixKey, []byte(id)...)
}

type InvoiceKeeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
	bank     bank.Keeper
}

func (ik InvoiceKeeper) HasInvoice(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(ik.storeKey)
	key := GetKey(id)
	return store.Has(key)
}

func (ik InvoiceKeeper) GetInvoice(ctx sdk.Context, id string) *Invoice {
	store := ctx.KVStore(ik.storeKey)
	key := GetKey(id)
	b := store.Get(key)
	invoice := &Invoice{}

	if err := ik.cdc.UnmarshalBinary(b, invoice); err != nil {
		return nil
	}

	return invoice
}

func (ik InvoiceKeeper) SetInvoice(ctx sdk.Context, invoice Invoice) {
	store := ctx.KVStore(ik.storeKey)
	key := GetKey(invoice.ID)
	b, err := ik.cdc.MarshalBinary(invoice)

	if err != nil {
		panic(err)
	}

	store.Set(key, b)
}

func (ik InvoiceKeeper) CreateInvoice(ctx sdk.Context, msg MsgCreate) sdk.Error {
	if ik.HasInvoice(ctx, msg.ID) {
		return ErrorDuplicateInvoice
	}

	var coins sdk.Coins

	for _, item := range msg.Items {
		coins = append(coins, sdk.Coin{Denom: item.AssetID, Amount: item.Quantity})
	}
	_, _, err := ik.bank.SubtractCoins(ctx, msg.Issuer, coins)

	if err != nil {
		return err
	}

	ik.SetInvoice(ctx, Invoice{
		ID:         msg.ID,
		Issuer:     msg.Issuer,
		Receiver:   msg.Receiver,
		Items:      msg.Items,
		CreateTime: time.Now(),
	})
	return nil
}

func NewInvoiceKeeper(store sdk.StoreKey, cdc *wire.Codec, bank bank.Keeper) InvoiceKeeper {
	return InvoiceKeeper{
		storeKey: store,
		cdc:      cdc,
		bank:     bank,
	}
}
