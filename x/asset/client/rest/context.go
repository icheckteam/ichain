package rest

import "github.com/cosmos/cosmos-sdk/client/context"

func withContext(ctx context.CoreContext, gas int64) (context.CoreContext, error) {
	var err error
	if gas == 0 {
		gas = 20000
	}

	// sign
	ctx, err = context.EnsureSequence(ctx)

	if err != nil {
		return ctx, err
	}

	return ctx.WithGas(gas), nil
}
