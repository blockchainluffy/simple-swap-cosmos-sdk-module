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
 export PATH=$PATH:/Users/amandeepsingh/go/bin
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