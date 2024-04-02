package simpleswap

import "cosmossdk.io/collections"

const ModuleName = "simpleswap"

var (
	ParamsKey  = collections.NewPrefix(0)
	PoolKey    = collections.NewPrefix(1)
	LiquidityProvidersKey = collections.NewPrefix(2)
	CoinsReserveKey = collections.NewPrefix(3)
)
