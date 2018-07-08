package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

func getValidators(ctx context.CoreContext, cdc *wire.Codec) []stake.Validator {
	return nil
}
