package invoice

import (
	"testing"

	"github.com/icheckteam/ichain/x/asset"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceKeeper(t *testing.T) {
	ctx, invoiceKeeper := createTestInput(t, false, 0)
	invoiceKeeper.assetKeeper.CreateAsset(ctx, asset.MsgCreateAsset{
		AssetID:  "tomato",
		Quantity: 1,
		Sender:   addrs[0],
	})

	msg := NewMsgCreate("1", addrs[1], addrs[0], []Item{Item{"tomato"}})
	_, err := invoiceKeeper.CreateInvoice(ctx, msg)
	assert.NotNil(t, err)

	msg = NewMsgCreate("1", addrs[0], addrs[1], []Item{Item{"tomato"}})
	_, err = invoiceKeeper.CreateInvoice(ctx, msg)
	assert.Nil(t, err)

	invoice := invoiceKeeper.GetInvoice(ctx, "1")
	assert.Equal(t, invoice.ID, "1")
	assert.Equal(t, invoice.Issuer, addrs[0])
	assert.Equal(t, invoice.Receiver, addrs[1])
	assert.Equal(t, invoice.Items, []Item{Item{"tomato"}})

	msg = NewMsgCreate("1", addrs[0], addrs[1], []Item{Item{"tomato"}})
	_, err = invoiceKeeper.CreateInvoice(ctx, msg)
	assert.NotNil(t, err)
}
