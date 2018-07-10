package invoice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/icheckteam/ichain/x/asset"
)

const (
	CodespaceDefault     sdk.CodespaceType = 69
	CodeDuplicateInvoice sdk.CodeType      = 1
)

var (
	PrefixKey             = []byte{0x01}
	AccountInvoiceKey     = []byte{0x02}
	ErrorDuplicateInvoice = sdk.NewError(CodespaceDefault, CodeDuplicateInvoice, "Duplicate invoice.")
)

func GetKey(id string) []byte {
	return append(PrefixKey, []byte(id)...)
}

func GetAccountInvoiceKey(addr sdk.AccAddress, contractID string) []byte {
	return append(GetAccountInvoicesKey(addr), []byte(contractID)...)
}

func GetAccountInvoicesKey(addr sdk.AccAddress) []byte {
	return append(AccountInvoiceKey, []byte(addr.String())...)
}

type InvoiceKeeper struct {
	storeKey    sdk.StoreKey
	cdc         *wire.Codec
	assetKeeper asset.Keeper
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

func (ik InvoiceKeeper) setInvoiceByAccountIndex(ctx sdk.Context, invoice Invoice) {
	store := ctx.KVStore(ik.storeKey)
	b, _ := ik.cdc.MarshalBinary(invoice.ID)
	store.Set(GetAccountInvoiceKey(invoice.Issuer, invoice.ID), b)
	if len(invoice.Receiver) > 0 {
		store.Set(GetAccountInvoiceKey(invoice.Receiver, invoice.ID), b)
	}

}

func (ik InvoiceKeeper) removeInvoiceByAccountIndex(ctx sdk.Context, invoice Invoice) {
	store := ctx.KVStore(ik.storeKey)
	store.Delete(GetAccountInvoiceKey(invoice.Issuer, invoice.ID))
	if len(invoice.Receiver) > 0 {
		store.Delete(GetAccountInvoiceKey(invoice.Receiver, invoice.ID))
	}
}

func (ik InvoiceKeeper) CreateInvoice(ctx sdk.Context, msg MsgCreate) (sdk.Tags, sdk.Error) {
	if ik.HasInvoice(ctx, msg.ID) {
		return nil, ErrorDuplicateInvoice
	}

	assetItems := []asset.Asset{}
	for _, item := range msg.Items {
		aitem, found := ik.assetKeeper.GetAsset(ctx, item.AssetID)
		if !found {
			return nil, asset.ErrAssetNotFound(item.AssetID)
		}
		if aitem.Final {
			return nil, asset.ErrAssetAlreadyFinal(aitem.ID)
		}

		if !aitem.IsOwner(msg.Issuer) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("%v not unauthorized to create", msg.Issuer))
		}
	}

	invoice := Invoice{
		ID:         msg.ID,
		Issuer:     msg.Issuer,
		Receiver:   msg.Receiver,
		Items:      msg.Items,
		CreateTime: ctx.BlockHeader().Time,
	}

	ik.SetInvoice(ctx, invoice)
	ik.setInvoiceByAccountIndex(ctx, invoice)
	tags := sdk.NewTags(
		"sender", []byte(msg.Receiver.String()),
	)
	for _, item := range assetItems {
		item.Final = true
		ik.assetKeeper.SetAsset(ctx, item)
		tags = tags.AppendTag("asset_id", []byte(item.ID))
	}

	if len(msg.Receiver) > 0 {
		tags = tags.AppendTag("recipient", []byte(msg.Receiver.String()))
	}

	return tags, nil
}

func NewInvoiceKeeper(store sdk.StoreKey, cdc *wire.Codec, assetKeeper asset.Keeper) InvoiceKeeper {
	return InvoiceKeeper{
		storeKey:    store,
		cdc:         cdc,
		assetKeeper: assetKeeper,
	}
}
