package keeper_test

import (
	"testing"

	"github.com/cosmos/simpleswap"
)

func (s *KeeperTestSuite) TestQueryParams() {
	require := s.Require()
	req := simpleswap.QueryParamsRequest{}
	resp, err := s.queryClient.Params(s.ctx, &req)
	require.NoError(err)
	require.Equal(simpleswap.DefaultParams(), resp.Params)
}

func (s *KeeperTestSuite) TestQueryPool() {
	require := s.Require()

	resp, err := s.queryClient.Pool(s.ctx, &simpleswap.QueryPoolRequest{})
	require.NoError(err)
	require.Equal(simpleswap.Pool{}, resp.Pool)
}

func (s *KeeperTestSuite) TestQueryLiquidityProvider() {
	t := s.T()

	t.Run("no liquidity provider", func(t *testing.T) {
		require := s.Require()

		_, err := s.queryClient.LiquidityProvider(s.ctx, &simpleswap.QueryLiquidityProviderRequest{
			LpAddress: s.addrs[0].String(),
		})
		require.Error(err)
	})
}

func (s *KeeperTestSuite) TestQueryCoinReserve() {
	t := s.T()

	t.Run("empty coins reserve", func(t *testing.T) {
		require := s.Require()

		_, err := s.queryClient.CoinReserve(s.ctx, &simpleswap.QueryCoinReserveRequest{
			CoinDenom: "invalid",
		})
		require.Error(err)
	})
}

func (s *KeeperTestSuite) TestQueryCoinsReserve() {
	t := s.T()

	t.Run("empty coins reserve", func(t *testing.T) {
		require := s.Require()

		_, err := s.queryClient.CoinReserves(s.ctx, &simpleswap.QueryCoinReservesRequest{})
		require.Error(err)
	})
}
