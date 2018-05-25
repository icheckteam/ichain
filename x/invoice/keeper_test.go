package invoice

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceKeeper(t *testing.T) {
	ctx, invoiceKeeper := createTestInput(t, false, 0)

	coins := sdk.Coins{
		{Denom: "tomato", Amount: 100},
		{Denom: "jav", Amount: 200},
	}
	coins = coins.Sort()
	invoiceKeeper.bank.AddCoins(ctx, addrs[0], coins)

	msg := NewMsgCreate("1", addrs[1], addrs[0], []Item{Item{"jav", 100}})
	err := invoiceKeeper.CreateInvoice(ctx, msg)
	assert.NotNil(t, err)

	msg = NewMsgCreate("1", addrs[0], addrs[1], []Item{Item{"jav", 300}})
	err = invoiceKeeper.CreateInvoice(ctx, msg)
	assert.NotNil(t, err)

	msg = NewMsgCreate("1", addrs[0], addrs[1], []Item{Item{"jav", 10}})
	err = invoiceKeeper.CreateInvoice(ctx, msg)
	assert.Nil(t, err)

	invoice := invoiceKeeper.GetInvoice(ctx, "1")
	assert.Equal(t, invoice.ID, "1")
	assert.Equal(t, invoice.Issuer, addrs[0])
	assert.Equal(t, invoice.Receiver, addrs[1])
	assert.Equal(t, invoice.Items, []Item{Item{"jav", 10}})
}
