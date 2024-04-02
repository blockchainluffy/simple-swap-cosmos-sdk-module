#!/usr/bin/env bash

rm -r ~/.minid || true
MINID_BIN=$(which minid)
# configure minid
$MINID_BIN config set client chain-id demo
$MINID_BIN config set client keyring-backend test
$MINID_BIN keys add alice
$MINID_BIN keys add bob
$MINID_BIN keys add validator
$MINID_BIN keys add traderA
$MINID_BIN keys add traderB
$MINID_BIN keys add traderC
$MINID_BIN keys add simpleswap
$MINID_BIN init test --chain-id demo --default-denom mini
# update genesis
$MINID_BIN genesis add-genesis-account validator 10000000mini --keyring-backend test
$MINID_BIN genesis add-genesis-account alice 10000mini,10000000000ETH --keyring-backend test
$MINID_BIN genesis add-genesis-account bob 10000mini,10000000000WETH --keyring-backend test
$MINID_BIN genesis add-genesis-account traderA 2000mini,1000000000ETH --keyring-backend test
$MINID_BIN genesis add-genesis-account traderB 2000mini,1000000000WETH --keyring-backend test
$MINID_BIN genesis add-genesis-account traderC 2000mini,1000000000stkETH --keyring-backend test
$MINID_BIN genesis add-genesis-account simpleswap 20000mini --keyring-backend test --module-name simpleswap
# create default validator
$MINID_BIN genesis gentx validator 1000000mini --chain-id demo
$MINID_BIN genesis collect-gentxs
