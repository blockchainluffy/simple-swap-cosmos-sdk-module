package keeper

import (
	"context"

	types "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/simpleswap"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *simpleswap.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	// Get The params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	shareToken := &types.Coin{
		Denom:  params.ShareToken.Denom,
		Amount: params.ShareToken.Amount,
	}
	// Set the pool
	if err := k.Pool.Set(ctx, simpleswap.Pool{
		Decimals:          params.Decimals,
		ShareToken:        shareToken,
		SwapFeePercentage: params.SwapFeePercentage,
	}); err != nil {
		return err
	}

	// Set the whitelisted coins
	for _, coin := range params.WhitelistedCoins {
		if err := k.CoinsReserve.Set(ctx, coin.Denom, *coin); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*simpleswap.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &simpleswap.GenesisState{
		Params: params,
	}, nil
}
