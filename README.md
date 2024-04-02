# SimpleSwap - Liquidity Pool and Token Swap Module for Mini

SimpleSwap is a liquidity pool and token swap module designed for our personalize blockchain `minid`. It allows liquidity providers to provide liquidity in stablecoins that are whitelisted on the module's parameters. Liquidity providers are rewarded with share tokens when they provide liquidity, and the tokens they provide are stored in the module's account.

Users can swap any supported coins as long as there is liquidity available in the pool. A configurable fee percentage (e.g., 0.3%) is set in the module's parameters, and the swap fees collected are distributed to the liquidity providers.

Liquidity providers have the flexibility to withdraw their liquidity and collected fees at any time, giving them control over their assets.

This README provides instructions on setting up and running the `minid` chain with the SimpleSwap module, as well as example CLI commands to interact with the local chain.

## Installation

To install and run Mini with the MiniSwap module, follow these steps:

1. *Clone the repository:*
 ```bash
 git clone https://github.com/punit-j/simple-swap-cosmos-sdk-module.git
 ```
2. To ensure the dependencies are up to date, run the following command:
 ```bash
 go mod tidy
 ```
3. Build and install the `minid` binary by running:
 ```bash
 make install
 ```
 This will install the `minid` binary in your `$GOBIN` directory.
4. Add the path to the `minid` binary to your system's `$PATH` environment variable. This allows you to run the `minid` command from any directory.
Add the below command at the end of your .bashrc or .zshrc file.
 ```bash
 export PATH=$PATH:/Users/punit-j/go/bin
```
> P.S. Run the command `which minid`, in a new terminal to verify that `minid` is now correctly in you path.
5. Run the following command to install all the dependencies, intialize your genesis and start the chain:
 ```bash
 make .PHONY
 ```

## Usage

Once you have installed `minid` and  with the SimpleSwap module, you can use the following commands to interact with the local chain:

1. Start the Mini chain: Run the command `make .PHONY` in the terminal to start the chain.
2. Get account addresses: Run the command `minid keys list` in the terminal to retrieve the addresses of the user accounts.
3. Perform operations: Use the commands provided in the `sample_commands.sh` file to perform various operations, such as providing liquidity, swapping tokens, and withdrawing liquidity and fees.
4. Explore available commands: To see a list of available commands for the SimpleSwap module, run `minid tx simpleswap --help` for transaction-based commands and `minid query simpleswap --help` for query-based commands.
5. You can also change the account holders' names in the `init.sh` script file if you'd like to use different account names.


## Additional Resources
For more information on the Cosmos SDK and building blockchain applications, refer to the following resources:

* [Cosmos SDK Documentation](https://docs.cosmos.network/): The official documentation for the Cosmos SDK.
* [GitHub Repository](https://github.com/cosmos/cosmos-sdk): The GitHub repository for the Cosmos SDK, which includes the latest version and additional resources.
* [Cosmos SDK Tutorials](https://tutorials.cosmos.network/): A collection of tutorials to help you get started with the Cosmos SDK.

## Appendix
This chain was generated using the [chain-minimal](https://github.com/cosmosregistry/chain-minimal) repository. It serves as a minimal template for building your own blockchain using the Cosmos SDK.

---

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
4. The module has set the default whitelist coins to be `ETH`, `WETH` and `stkETH`.
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