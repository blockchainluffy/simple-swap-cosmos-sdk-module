# `SimpleSwap`

A module to support swap functionality for coins with the same prices.

## Concepts

The `SimpleSwap` module enables users to swap coins with identical prices, such as ETH and WETH. This module is designed to be used by other modules to provide swap functionality. Liquidity providers can offer liquidity in any of the stablecoins whitelisted in the module parameters. By providing liquidity, they receive share tokens. The tokens provided by liquidity providers are held in the module account. Users can swap any coins as long as there is available liquidity. A fee percentage (e.g., 0.3%) is configured in the parameters, and this swap fee is given to liquidity providers. Liquidity providers can withdraw their liquidity and collected fees at any time.

## State

The `SimpleSwap` module's state consists of the following data, which is stored in the KVStore using the collection keys defined in the `keys.go` file:

1. `Pool`: A struct that contains information about the pool.
2. `Params`: A struct that contains module parameter information.
3. `LiquidityProviders`: A map that contains information about liquidity providers.
4. `CoinsReserve`: A map that contains information about the coins reserve.

## State Transitions

The state transition operations are defined in the `tx.proto` file located in the `/proto/cosmos/simpleswap/v1` directory.

## Messages

The `SimpleSwap` module defines the following messages:

1. `MsgAddLiquidity`: A message to add liquidity to the pool.
2. `MsgRemoveLiquidity`: A message to remove liquidity from the pool.
3. `MsgSwapLiquidity`: A message to swap coins.

## Client

You can find CLI commands for the `SimpleSwap` module in the `/module/autocli.go` file.

## Params

The module parameters are as follows:

1. `WhitelistCoins`: A list of coins that are allowed to be used in the module.
2. `SwapFeePercentage`: The fee percentage charged on swaps.
3. `Decimals`: The number of decimal places for the coins.
4. `ShareToken`: The share token given to liquidity providers.

## Assumptions

1. The module assumes that the coins provided by liquidity providers have the same price.
2. The module assumes that liquidity providers can only provide coins of a single denomination.
3. The module assumes a swap fee of 0.3% for swap exchanges.
4. The module has set the default whitelist coins to be `ETH` and `WETH`.
5. The collected fees by the liquidity providers are returned in the same denomination provided by the liquidity providers.

## Future Improvements

1. Add more coins to the whitelist to expand the available options for swapping.
2. Add events to the module to provide descriptive information about state changes and aid in debugging.
3. Increase test coverage to handle all edge cases and ensure robustness.
4. Enhance error handling to make the module more resilient.
5. Implement functionality for users to provide liquidity in multiple coins simultaneously or for liquidity providers to offer multiple coins.
6. Introduce functionality to swap multiple coins at once for improved convenience.
7. Add a feature to distribute fees to the liquidity providers based on amount of their liquidity used.
8. Add a feature to distribute fees to a governance fund or other entity.
9. Swap Fees for each type of exchange can be configured specifically, this will provide fees to be transferred in the same denomination it is collected.

## Appendix

This Cosmos SDK module was generated using [https://github.com/cosmosregistry/example](https://github.com/cosmosregistry/example).