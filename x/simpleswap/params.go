package simpleswap

import (
	"cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
)

// DefaultParams returns default module parameters.
func DefaultParams() Params {
	shareToken := types.NewCoin("USDT", math.NewInt(100))
	return Params{
		WhitelistedCoins: []*types.Coin{
			{
				Denom:  "ETH",
				Amount: math.ZeroInt(),
			},
			{
				Denom:  "WETH",
				Amount: math.ZeroInt(),
			},
			{
				Denom:  "stkETH",
				Amount: math.ZeroInt(),
			},
		},
		Decimals:          6,
		ShareToken:        &shareToken,
		SwapFeePercentage: 30000,
	}
}

// Validate does the sanity check on the params.
func (p Params) Validate() error {
	// Sanity check goes here.
	// Check if the Whitelisted coins are present and valid
	if len(p.WhitelistedCoins) == 0 {
		return ErrCoinsNotPresent
	}

	for _, coin := range p.WhitelistedCoins {
		if coin.Denom == "" {
			return ErrCoinInvalid
		}
	}

	if p.SwapFeePercentage == 0 {
		return ErrZeroSwapFee
	}

	if p.Decimals == 0 {
		return ErrZeroDecimals
	}

	if p.ShareToken.Denom == "" {
		return ErrShareTokenInvalid
	}

	return nil
}
