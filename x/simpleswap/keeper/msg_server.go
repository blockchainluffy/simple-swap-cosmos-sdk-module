package keeper

import (
	"context"
	"fmt"
	"strings"

	maths "math"

	math "cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/simpleswap"
)

type msgServer struct {
	k Keeper
}

var _ simpleswap.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) simpleswap.MsgServer {
	return &msgServer{k: keeper}
}

// UpdateParams params is defining the handler for the MsgUpdateParams message.
func (ms msgServer) UpdateParams(ctx context.Context, msg *simpleswap.MsgUpdateParams) (*simpleswap.MsgUpdateParamsResponse, error) {
	if _, err := ms.k.addressCodec.StringToBytes(msg.Authority); err != nil {
		return nil, fmt.Errorf("invalid authority address: %w", err)
	}

	if authority := ms.k.GetAuthority(); !strings.EqualFold(msg.Authority, authority) {
		return nil, fmt.Errorf("unauthorized, authority does not match the module's authority: got %s, want %s", msg.Authority, authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := ms.k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &simpleswap.MsgUpdateParamsResponse{}, nil
}

// AddLiquidity is defining the handler for the MsgAddLiquidity message.
func (ms msgServer) AddLiquidity(ctx context.Context, msg *simpleswap.MsgAddLiquidity) (*simpleswap.MsgAddLiquidityResponse, error) {

	// Check if the amount is zero
	if msg.Token.Amount.Int64() == 0 {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 400,
		}, simpleswap.ErrZeroAmount
	}

	// Check if the liquidity provider address is valid
	_, err := ms.k.addressCodec.StringToBytes(msg.LiquidityProvider)
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 401,
		}, simpleswap.ErrInvalidProviderAddress
	}

	// Get the current pool state
	currentPoolState, err := ms.k.Pool.Get(ctx)
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Check if the coin being added is present in the whitelist
	params, err := ms.k.Params.Get(ctx)
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	isInputCoinPresent := false
	whitelistedCoins := params.GetWhitelistedCoins()
	for _, coin := range whitelistedCoins {
		if coin.Denom == msg.Token.Denom {
			isInputCoinPresent = true
		}
	}

	if !isInputCoinPresent {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 400,
		}, fmt.Errorf("error: %w, for the denom: %s", simpleswap.ErrCoinInvalid, msg.Token.Denom)
	}

	globallyAccruedFeesCurrent := currentPoolState.GetTotalAccruedFees()

	// Get the liquidity provider address in AccAddress format
	addr, err := types.AccAddressFromBech32(msg.LiquidityProvider)
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, fmt.Errorf("error in converting address to AccAddress: %w", err)
	}

	// Transfer the stablecoin from the liquidity provider to the module account
	err = ms.k.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, simpleswap.ModuleName, types.NewCoins(msg.Token))
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Get the liquidity provider and update the stable coins
	liquidityProvider, err := ms.k.LiquidityProviders.Get(ctx, msg.LiquidityProvider)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			liquidityProvider.StableCoin = &types.Coin{
				Denom:  msg.Token.Denom,
				Amount: math.ZeroInt(),
			}
			liquidityProvider.PoolShare = &types.Coin{
				Denom:  currentPoolState.ShareToken.Denom,
				Amount: math.ZeroInt(),
			}
		} else {
			return &simpleswap.MsgAddLiquidityResponse{
				StatusCode: 500,
			}, err
		}
	}

	coin := liquidityProvider.StableCoin
	coin.Amount = coin.Amount.Add(msg.Token.Amount)

	// Fetch the pool share for the LP
	poolShare := liquidityProvider.PoolShare

	// Coins to be minted
	coinsToMint := types.Coin{
		Denom:  poolShare.Denom,
		Amount: msg.Token.Amount,
	}

	accruedFeesGloballyRecordedByLP := liquidityProvider.GloballyAccruedFees
	accruedFeesByLP := liquidityProvider.AccruedFees

	// See if the accrued fees by liquidity provider is equal to the globally accrued fees
	if accruedFeesGloballyRecordedByLP != globallyAccruedFeesCurrent {
		// Calculate the fees and update the Pool
		diff := globallyAccruedFeesCurrent - accruedFeesGloballyRecordedByLP
		accruedFeesByLP += (diff * poolShare.Amount.Int64()) / currentPoolState.TotalLiquidity
	}

	// Update the pool share
	poolShare.Amount = poolShare.Amount.Add(coinsToMint.Amount)

	// Update the liquidity provider
	if err := ms.k.LiquidityProviders.Set(ctx, msg.LiquidityProvider, simpleswap.LiquidityProvider{
		StableCoin:          coin,
		PoolShare:           poolShare,
		AccruedFees:         accruedFeesByLP,
		GloballyAccruedFees: globallyAccruedFeesCurrent,
	}); err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Check if the coin is already present in the Coins Reserve
	coinReserves, err := ms.k.CoinsReserve.Get(ctx, msg.Token.Denom)
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Update the coin amount
	coinReserves.Amount = coinReserves.Amount.Add(msg.Token.Amount)

	// Update the Coins Reserve
	if err := ms.k.CoinsReserve.Set(ctx, coin.Denom, coinReserves); err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Update the pool
	if err := ms.k.Pool.Set(ctx, simpleswap.Pool{
		TotalAccruedFees:  currentPoolState.TotalAccruedFees,
		TotalLiquidity:    currentPoolState.TotalLiquidity + msg.Token.Amount.Int64(),
		Decimals:          currentPoolState.Decimals,
		ShareToken:        currentPoolState.ShareToken,
		SwapFeePercentage: currentPoolState.SwapFeePercentage,
	}); err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// MINT AND SEND the share token to the liquidity provider
	err = ms.k.BankKeeper.MintCoins(ctx, simpleswap.ModuleName, types.NewCoins(coinsToMint))
	if err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	if err := ms.k.BankKeeper.SendCoinsFromModuleToAccount(ctx, simpleswap.ModuleName, addr, types.NewCoins(coinsToMint)); err != nil {
		return &simpleswap.MsgAddLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	return &simpleswap.MsgAddLiquidityResponse{
		StatusCode: 200,
	}, nil
}

func (ms msgServer) SwapLiquidity(ctx context.Context, msg *simpleswap.MsgSwapLiquidity) (*simpleswap.MsgSwapLiquidityResponse, error) {

	// Check if the amount is zero
	if msg.Input.Amount.Int64() == 0 {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 400,
		}, simpleswap.ErrZeroAmount
	}

	// Check if the input and output amount are equal
	if msg.Input.Amount.Int64() != msg.Output.Amount.Int64() {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 400,
		}, simpleswap.ErrAmountNotEqual
	}

	// Check if the liquidity provider address is valid
	_, err := ms.k.addressCodec.StringToBytes(msg.Trader)
	if err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 401,
		}, simpleswap.ErrInvalidProviderAddress
	}

	// Get the current pool state
	currentPoolState, err := ms.k.Pool.Get(ctx)
	if err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Check if the coins being swapped are present in the whitelist
	params, err := ms.k.Params.Get(ctx)
	if err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	isInputCoinPresent := false
	isOutputCoinPresent := false
	whitelistedCoins := params.GetWhitelistedCoins()
	for _, coin := range whitelistedCoins {
		if coin.Denom == msg.Input.Denom {
			isInputCoinPresent = true
		} else if coin.Denom == msg.Output.Denom {
			isOutputCoinPresent = true
		}
	}

	if !isInputCoinPresent || !isOutputCoinPresent {
		if !isInputCoinPresent {
			return &simpleswap.MsgSwapLiquidityResponse{
				StatusCode: 400,
			}, fmt.Errorf("error: %w, for the denom: %s", simpleswap.ErrCoinInvalid, msg.Input.Denom)
		}
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 400,
		}, fmt.Errorf("error: %w, for the denom: %s", simpleswap.ErrCoinInvalid, msg.Output.Denom)
			
	}

	// Get the liquidity provider address in AccAddress format
	addr, err := types.AccAddressFromBech32(msg.Trader)
	if err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, fmt.Errorf("error in converting address to AccAddress: %w", err)
	}

	// Check if the Liquidity Provider has the required input token
	coins := ms.k.BankKeeper.SpendableCoin(ctx, addr, msg.Input.Denom)
	if coins.Amount.LT(msg.Input.Amount) {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 409,
		}, fmt.Errorf("error: %w for the denom: %s", simpleswap.ErrInsufficientLiquidity, msg.Input.Denom)
	}

	// Check if the required output token is present in required quantity
	coinsReserveOutputToken, err := ms.k.CoinsReserve.Get(ctx, msg.Output.Denom)
	if err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	if coinsReserveOutputToken.Amount.LT(msg.Output.Amount) {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 409,
		}, fmt.Errorf("error: %w for the denom: %s", simpleswap.ErrInsufficientLiquidity, msg.Output.Denom)
	}

	// Transfer the input token from the trader to the module account
	err = ms.k.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, simpleswap.ModuleName, types.NewCoins(msg.Input))
	if err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Get the coins reserve and Update the coins reserve for the input coins provided by the trader
	coinsReserveInputToken, err := ms.k.CoinsReserve.Get(ctx, msg.Input.Denom)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			coinsReserveInputToken = types.Coin{
				Denom:  msg.Input.Denom,
				Amount: math.ZeroInt(),
			}
		} else {
			return &simpleswap.MsgSwapLiquidityResponse{
				StatusCode: 500,
			}, err
		}
	}

	coinsReserveOutputToken.Amount = coinsReserveOutputToken.Amount.Sub(msg.Output.Amount)
	if err := ms.k.CoinsReserve.Set(ctx, msg.Output.Denom, coinsReserveOutputToken); err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Calculate Fees and charge it from the output token
	// TODO: Explain the below conversions and calculations
	swapFeePercentage := currentPoolState.SwapFeePercentage
	swapFee := msg.Output.Amount.Mul(math.NewInt(int64(swapFeePercentage))).Quo(math.NewInt(int64(maths.Pow10(int(currentPoolState.Decimals))) * int64(100)))

	// Deduct the swap fee from the output token
	msg.Output.Amount = msg.Output.Amount.Sub(swapFee)

	// Update the Coins Reserve For Input Token
	coinsReserveInputToken.Amount = coinsReserveInputToken.Amount.Add(msg.Input.Amount)
	if err := ms.k.CoinsReserve.Set(ctx, msg.Input.Denom, coinsReserveInputToken); err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Update the fees in the Pool
	currentPoolState.TotalAccruedFees += swapFee.Int64()
	if err := ms.k.Pool.Set(ctx, simpleswap.Pool{
		TotalAccruedFees:  currentPoolState.TotalAccruedFees,
		TotalLiquidity:    currentPoolState.TotalLiquidity,
		Decimals:          currentPoolState.Decimals,
		ShareToken:        currentPoolState.ShareToken,
		SwapFeePercentage: currentPoolState.SwapFeePercentage,
	}); err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}
	// Send the Required Token to the liquidity provider
	if err := ms.k.BankKeeper.SendCoinsFromModuleToAccount(ctx, simpleswap.ModuleName, addr, types.NewCoins(msg.Output)); err != nil {
		return &simpleswap.MsgSwapLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	return &simpleswap.MsgSwapLiquidityResponse{
		StatusCode: 200,
	}, err
}

func (ms msgServer) RemoveLiquidity(ctx context.Context, msg *simpleswap.MsgRemoveLiquidity) (*simpleswap.MsgRemoveLiquidityResponse, error) {
	// Check if the amount is zero
	if msg.Token.Amount.Int64() == 0 {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 400,
		}, simpleswap.ErrZeroAmount
	}

	// Check if the liquidity provider address is valid
	_, err := ms.k.addressCodec.StringToBytes(msg.LiquidityProvider)
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 401,
		}, simpleswap.ErrInvalidProviderAddress
	}

	// Get the current pool state
	currentPoolState, err := ms.k.Pool.Get(ctx)
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Get the liquidity provider
	liquidityProvider, err := ms.k.LiquidityProviders.Get(ctx, msg.LiquidityProvider)
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// TODO: Add check for same denomination

	// Convert the address to AccAddress
	addr, err := types.AccAddressFromBech32(msg.LiquidityProvider)
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, fmt.Errorf("error in converting address to AccAddress: %w", err)
	}

	// Check if the coins Reserve has the required amount of coins
	coinsReserve, err := ms.k.CoinsReserve.Get(ctx, msg.Token.Denom)
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Check if the liquidity provider has provided the requested amount of stable coins to the pool
	if liquidityProvider.StableCoin.Amount.LT(msg.Token.Amount) {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 409,
		}, fmt.Errorf("error: %w for the denom: %s, with the User: %s", simpleswap.ErrInsufficientLiquidity, msg.Token.Denom, msg.LiquidityProvider)
	}

	// Calculate if the liquidity provider has received his share of the fees totally
	accruedFeesGlobally := currentPoolState.GetTotalAccruedFees()
	globallyAccruedFeesRecordedByUser := liquidityProvider.GloballyAccruedFees

	// Calculate the fees and update the liquidity provider
	// TODO: Explain the below calculations
	if accruedFeesGlobally != globallyAccruedFeesRecordedByUser {
		diff := accruedFeesGlobally - globallyAccruedFeesRecordedByUser
		liquidityProvider.AccruedFees += (diff * liquidityProvider.PoolShare.Amount.Int64()) / currentPoolState.TotalLiquidity
	}

	// Calculate the accrued fees to the liquidity provider
	accruedFees := liquidityProvider.AccruedFees

	// Check if the coins Reserve has the required amount of coins
	if coinsReserve.Amount.LT(msg.Token.Amount) {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 409,
		}, fmt.Errorf("error: %w for the denom: %s", simpleswap.ErrInsufficientLiquidity, msg.Token.Denom)
	}

	// Update the coins reserve for the input coins provided by the liquidity provider
	coinsReserve.Amount = coinsReserve.Amount.Sub(msg.Token.Amount)

	if err := ms.k.CoinsReserve.Set(ctx, coinsReserve.Denom, coinsReserve); err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Update the pool
	if err := ms.k.Pool.Set(ctx, simpleswap.Pool{
		TotalAccruedFees:  currentPoolState.TotalAccruedFees,
		TotalLiquidity:    currentPoolState.TotalLiquidity - msg.Token.Amount.Int64(),
		Decimals:          currentPoolState.Decimals,
		ShareToken:        currentPoolState.ShareToken,
		SwapFeePercentage: currentPoolState.SwapFeePercentage,
	}); err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Get the pool share for the LP
	poolShare := liquidityProvider.PoolShare

	// Transfer LP coins from LP to Module Accounts inorder to burn them
	err = ms.k.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, simpleswap.ModuleName, types.NewCoins(types.NewCoin(poolShare.Denom, poolShare.Amount.Sub(msg.Token.Amount))))
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// Burn the share token	from the liquidity provider
	err = ms.k.BankKeeper.BurnCoins(ctx, simpleswap.ModuleName, types.NewCoins(types.NewCoin(poolShare.Denom, poolShare.Amount.Sub(msg.Token.Amount))))
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	// We check if the liquidity provider is removing all of its liquidity
	// If not we update the liquidity provider
	// If yes we remove the liquidity provider from the store
	if msg.Token.Amount.LT(liquidityProvider.StableCoin.Amount) {
		if err := ms.k.LiquidityProviders.Set(ctx, msg.LiquidityProvider, simpleswap.LiquidityProvider{
			StableCoin: &types.Coin{
				Denom:  msg.Token.Denom,
				Amount: liquidityProvider.StableCoin.Amount.Sub(msg.Token.Amount),
			},
			PoolShare:           poolShare,
			AccruedFees:         int64(0),
			GloballyAccruedFees: accruedFeesGlobally,
		}); err != nil {
			return &simpleswap.MsgRemoveLiquidityResponse{
				StatusCode: 500,
			}, err
		}
	} else {
		// Remove the liquidity provider
		if err := ms.k.LiquidityProviders.Remove(ctx, msg.LiquidityProvider); err != nil {
			return &simpleswap.MsgRemoveLiquidityResponse{
				StatusCode: 500,
			}, err
		}
	}
	
	// Add the accrued fees to the output coin
	msg.Token.Amount = msg.Token.Amount.Add(math.NewInt(accruedFees))

	// Transfer the stable coins from the module account to the liquidity provider
	err = ms.k.BankKeeper.SendCoinsFromModuleToAccount(ctx, simpleswap.ModuleName, addr, types.NewCoins(msg.Token))
	if err != nil {
		return &simpleswap.MsgRemoveLiquidityResponse{
			StatusCode: 500,
		}, err
	}

	return &simpleswap.MsgRemoveLiquidityResponse{
		StatusCode: 200,
	}, nil
}
