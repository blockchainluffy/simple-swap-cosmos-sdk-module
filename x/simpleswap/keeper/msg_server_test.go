package keeper_test

import (
	"fmt"
	"testing"

	math "cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/simpleswap"
	// expectedkeepers "github.com/cosmos/simpleswap/expected_keepers"
	// "github.com/golang/mock/gomock"
)

func (s *KeeperTestSuite) TestUpdateParams() {
	require := s.Require()

	testCases := []struct {
		name         string
		request      *simpleswap.MsgUpdateParams
		expectErrMsg string
	}{
		{
			name: "set invalid authority (not an address)",
			request: &simpleswap.MsgUpdateParams{
				Authority: "foo",
			},
			expectErrMsg: "invalid authority address",
		},
		{
			name: "set invalid authority (not defined authority)",
			request: &simpleswap.MsgUpdateParams{
				Authority: s.addrs[1].String(),
			},
			expectErrMsg: fmt.Sprintf("unauthorized, authority does not match the module's authority: got %s, want %s", s.addrs[1].String(), s.simpleSwapKeeper.GetAuthority()),
		},
		{
			name: "set valid params",
			request: &simpleswap.MsgUpdateParams{
				Authority: s.simpleSwapKeeper.GetAuthority(),
				Params: simpleswap.Params{
					WhitelistedCoins: []*types.Coin{
						{
							Denom:  "ETH",
							Amount: math.ZeroInt(),
						},
						{
							Denom:  "WETH",
							Amount: math.ZeroInt(),
						},
					},
					Decimals:          6,
					ShareToken:        &types.Coin{Denom: "LP", Amount: math.NewInt(100)},
					SwapFeePercentage: 3,
				},
			},
			expectErrMsg: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := s.msgServer.UpdateParams(s.ctx, tc.request)
			if tc.expectErrMsg != "" {
				require.Error(err)
				require.ErrorContains(err, tc.expectErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAddLiquidity() {
	t := s.T()
	// t.Run("add zero amount", func(t *testing.T) {
	// 	require := s.Require()

	// 	_, err := s.msgServer.AddLiquidity(s.ctx, &simpleswap.MsgAddLiquidity{
	// 		LiquidityProvider: s.addrs[0].String(),
	// 		Token:             types.Coin{Denom: "ETH", Amount: math.ZeroInt()},
	// 	})
	// 	require.Error(err)
	// 	require.ErrorContains(err, "amount cannot be zero")
	// })

	// t.Run("add liquidity with invalid provider address", func(t *testing.T) {
	// 	require := s.Require()

	// 	_, err := s.msgServer.AddLiquidity(s.ctx, &simpleswap.MsgAddLiquidity{
	// 		LiquidityProvider: "foo",
	// 		Token:             types.Coin{Denom: "ETH", Amount: math.NewInt(100)},
	// 	})
	// 	require.Error(err)
	// 	require.ErrorContains(err, "invalid provider address")
	// })

	t.Run("add liquidity with valid provider address first time", func(t *testing.T) {
		require := s.Require()

		
		response, err := s.msgServer.AddLiquidity(s.ctx, &simpleswap.MsgAddLiquidity{
			LiquidityProvider: s.addrs[1].String(),
			Token:             types.Coin{Denom: "ETH", Amount: math.NewInt(100)},
		})
		require.NoError(err)
		require.Equal(&simpleswap.MsgAddLiquidityResponse{StatusCode: 200}, response)
		s.bankKeeper.EXPECT().MintCoins(s.ctx, simpleswap.ModuleName, types.NewCoins(types.Coin{Denom: "LP", Amount: math.NewInt(100)})).Return(nil).Times(1)
		s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(s.ctx, s.addrs[1], simpleswap.ModuleName, types.NewCoins(types.Coin{Denom: "ETH", Amount: math.NewInt(100)})).Return(nil).Times(1)

		lp, err := s.simpleSwapKeeper.LiquidityProviders.Get(s.ctx, s.addrs[1].String())
		require.NoError(err)
		require.Equal("ETH", lp.StableCoin.Denom)
		require.Equal(math.NewInt(100), lp.StableCoin.Amount)
		require.Equal("LP", lp.PoolShare.Denom)
		require.Equal(math.NewInt(100), lp.PoolShare.Amount)

		pool, err := s.simpleSwapKeeper.Pool.Get(s.ctx)
		require.NoError(err)
		require.Equal(100, pool.TotalLiquidity)

		coinsReserve, err := s.simpleSwapKeeper.CoinsReserve.Get(s.ctx, "ETH")
		require.NoError(err)
		require.Equal(math.NewInt(100), coinsReserve.Amount)
	})

	// t.Run("add liquidity with valid provider address second time", func(t *testing.T) {
	// 	require := s.Require()

	// 	_, err := s.msgServer.AddLiquidity(s.ctx, &simpleswap.MsgAddLiquidity{
	// 		LiquidityProvider: s.addrs[2].String(),
	// 		Token:             types.Coin{Denom: "WETH", Amount: math.NewInt(50)},
	// 	})
	// 	require.NoError(err)

	// 	s.bankKeeper.EXPECT().MintCoins(s.ctx, simpleswap.ModuleName, types.NewCoins(types.Coin{Denom: "LP", Amount: math.NewInt(50)})).Return(nil).Times(1)
	// 	s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(s.ctx, s.addrs[1], simpleswap.ModuleName, types.NewCoins(types.Coin{Denom: "WETH", Amount: math.NewInt(50)})).Return(nil).Times(1)

	// 	lp, err := s.simpleSwapKeeper.LiquidityProviders.Get(s.ctx, s.addrs[2].String())
	// 	require.NoError(err)
	// 	require.Equal("WETH", lp.StableCoin.Denom)
	// 	require.Equal(math.NewInt(50), lp.StableCoin.Amount)
	// 	require.Equal("LP", lp.PoolShare.Denom)
	// 	require.Equal(math.NewInt(50), lp.PoolShare.Amount)

	// 	pool, err := s.simpleSwapKeeper.Pool.Get(s.ctx)
	// 	require.NoError(err)
	// 	require.Equal(150, pool.TotalLiquidity)

	// 	coinsReserve, err := s.simpleSwapKeeper.CoinsReserve.Get(s.ctx, "WETH")
	// 	require.NoError(err)
	// 	require.Equal(math.NewInt(50), coinsReserve.Amount)
	// })
}


