package simpleswap

import "cosmossdk.io/errors"

var (
	// ErrDuplicateAddress error if there is a duplicate address
	ErrDuplicateAddress = errors.Register(ModuleName, 2, "duplicate address")
	ErrCoinsNotPresent = errors.Register(ModuleName,3, "white listed coins cannot be nil")
	ErrCoinInvalid = errors.Register(ModuleName, 4, "coin provided is invalid")
	ErrZeroSwapFee = errors.Register(ModuleName, 5, "swap fee cannot be zero")
	ErrZeroAmount = errors.Register(ModuleName, 6, "amount cannot be zero")
	ErrInvalidProviderAddress = errors.Register(ModuleName, 7, "invalid provider address")
	ErrPoolNotInitialized = errors.Register(ModuleName, 8, "pool not initialized")
	ErrInsufficientLiquidity = errors.Register(ModuleName, 9, "required liquidity present in insufficient amount")
	ErrZeroDecimals = errors.Register(ModuleName, 10, "decimals cannot be zero")
	ErrZeroDecimalCoefficient = errors.Register(ModuleName, 11, "decimal coefficient cannot be zero")
	ErrZeroSwapFeeDecimals = errors.Register(ModuleName, 12, "swap fee decimals cannot be zero")
	ErrShareTokenInvalid = errors.Register(ModuleName, 13, "share token is invalid")
	ErrAmountNotEqual = errors.Register(ModuleName, 14, "amounts are not equal")
)
