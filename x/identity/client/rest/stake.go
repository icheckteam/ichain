package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
)

func getValidators(ctx context.CoreContext, cdc *wire.Codec) ([]stake.Validator, error) {
	kvs, err := ctx.QuerySubspace(cdc, stake.ValidatorsKey, "stake")
	if err != nil {
		return nil, err
	}

	// parse out the validators
	validators := make([]stake.Validator, len(kvs))
	for i, kv := range kvs {

		addr := kv.Key[1:]
		validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
		if err != nil {
			return nil, err
		}

		validators[i] = validator
	}

	return validators, nil
}
