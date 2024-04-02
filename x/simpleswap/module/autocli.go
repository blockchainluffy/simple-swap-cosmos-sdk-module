package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	simpleswapv1 "github.com/cosmos/simpleswap/api/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: simpleswapv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use: 	 "params",
					Short:   "Get the simpleswap module parameters",
				},
				{
					RpcMethod: "Pool",
					Use:       "pool",
					Short:     "Get the simpleswap module pool",
				},
				{
					RpcMethod: "LiquidityProvider",
					Use:       "liquidity-provider lpAddress",
					Short:     "Get the simpleswap module liquidity provider",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "lpAddress"},
					},
				},
				{
					RpcMethod: "CoinReserve",
					Use:       "coin-reserve coin",
					Short:     "Get the reserve of a coin",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "coin"},
					},
				},
				{
					RpcMethod: "CoinReserves",
					Use:       "coin-reserves",
					Short:     "Get all the reserves",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: simpleswapv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "AddLiquidity",
					Use:       "add-liquidity liquidityProvider amount token",
					Short:     "Add liquidity to the pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "liquidityProvider"},
						{ProtoField: "amount"},
						{ProtoField: "token"},
					},
				},
				{
					RpcMethod: "SwapLiquidity",
					Use:       "swap-liquidity trader input output",
					Short:     "Swap liquidity from one token to another",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "trader"},
						{ProtoField: "input"},
						{ProtoField: "output"},
					},
				},
				{
					RpcMethod: "RemoveLiquidity",
					Use:       "remove-liquidity liquidityProvider amount token",
					Short:     "Remove liquidity from the pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "liquidityProvider"},
						{ProtoField: "token"},
					},
				},
			},
		},
	}
}
