package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/simpleswap"
)

var _ simpleswap.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) simpleswap.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// Params defines the handler for the Query/Params RPC method.
func (qs queryServer) Params(ctx context.Context, req *simpleswap.QueryParamsRequest) (*simpleswap.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &simpleswap.QueryParamsResponse{Params: simpleswap.Params{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &simpleswap.QueryParamsResponse{Params: params}, nil
}

// Pool defines the handler for the Query/Pool RPC method.
func (qs queryServer) Pool(ctx context.Context, req *simpleswap.QueryPoolRequest) (*simpleswap.QueryPoolResponse, error) {
	pool, err := qs.k.Pool.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &simpleswap.QueryPoolResponse{Pool: simpleswap.Pool{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &simpleswap.QueryPoolResponse{Pool: pool}, nil
}

// LiquidityProvider defines the handler for the Query/LiquidityProvider RPC method.
func (qs queryServer) LiquidityProvider(ctx context.Context, req *simpleswap.QueryLiquidityProviderRequest) (*simpleswap.QueryLiquidityProviderResponse, error) {
	lp, err := qs.k.LiquidityProviders.Get(ctx, req.LpAddress)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &simpleswap.QueryLiquidityProviderResponse{
				LiquidityProvider: simpleswap.LiquidityProvider{},
			}, fmt.Errorf("liquidity provider not found for %s", req.LpAddress)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &simpleswap.QueryLiquidityProviderResponse{LiquidityProvider: lp}, nil
}

// CoinReserves defines the handler for the Query/CoinReserves RPC method.
func (qs queryServer) CoinReserve(ctx context.Context, req *simpleswap.QueryCoinReserveRequest) (*simpleswap.QueryCoinReserveResponse, error) {
	reserves, err := qs.k.CoinsReserve.Get(ctx, req.CoinDenom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &simpleswap.QueryCoinReserveResponse{
				CoinReserve: types.Coin{},
			}, fmt.Errorf("coin reserve not found for %s", req.CoinDenom)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &simpleswap.QueryCoinReserveResponse{CoinReserve: reserves}, nil
}


func (qs queryServer) CoinReserves(ctx context.Context, req *simpleswap.QueryCoinReservesRequest) (*simpleswap.QueryCoinReservesResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &simpleswap.QueryCoinReservesResponse{
				CoinReserves: []types.Coin{},
			}, err
		}
	}

	whitelsitedCoins := params.WhitelistedCoins
	var reserves []types.Coin 
	for _, coin := range whitelsitedCoins {
		reserve, err := qs.k.CoinsReserve.Get(ctx, coin.Denom)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return &simpleswap.QueryCoinReservesResponse{
					CoinReserves: []types.Coin{},
				}, fmt.Errorf("coin reserve not found for %s", coin.Denom)
			}

			return nil, status.Error(codes.Internal, err.Error())
		}
		reserves = append(reserves, reserve)
	}

	return &simpleswap.QueryCoinReservesResponse{CoinReserves: reserves}, nil
}