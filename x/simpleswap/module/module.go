package module

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/simpleswap"
	"github.com/cosmos/simpleswap/keeper"
)

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasGenesis     = AppModule{}
	_ appmodule.AppModule   = AppModule{}
)

// ConsensusVersion defines the current module consensus version.
const ConsensusVersion = 1

type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

func NewAppModuleBasic(m AppModule) module.AppModuleBasic {
	return module.CoreAppModuleBasicAdaptor(m.Name(), m)
}

// Name returns the simpleswap module's name.
func (AppModule) Name() string { return simpleswap.ModuleName }

// RegisterLegacyAminoCodec registers the simpleswap module's types on the LegacyAmino codec.
// New modules do not need to support Amino.
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	simpleswap.RegisterLegacyAminoCodec(cdc)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the simpleswap module.
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := simpleswap.RegisterQueryHandlerClient(context.Background(), mux, simpleswap.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces registers interfaces and implementations of the simpleswap module.
func (AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	simpleswap.RegisterInterfaces(registry)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	simpleswap.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	simpleswap.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))

	// Register in place module state migration migrations
	// m := keeper.NewMigrator(am.keeper)
	// if err := cfg.RegisterMigration(simpleswap.ModuleName, 1, m.Migrate1to2); err != nil {
	// 	panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", simpleswap.ModuleName, err))
	// }
}

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(simpleswap.NewGenesisState())
}

// ValidateGenesis performs genesis state validation for the circuit module.
func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data simpleswap.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", simpleswap.ModuleName, err)
	}

	return data.Validate()
}

// InitGenesis performs genesis initialization for the simpleswap module.
// It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState simpleswap.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	if err := am.keeper.InitGenesis(ctx, &genesisState); err != nil {
		panic(fmt.Sprintf("failed to initialize %s genesis state: %v", simpleswap.ModuleName, err))
	}

	fmt.Println("CheckPoint1")
}

// ExportGenesis returns the exported genesis state as raw bytes for the circuit
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to export %s genesis state: %v", simpleswap.ModuleName, err))
	}

	return cdc.MustMarshalJSON(gs)
}
func (AppModule) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   simpleswap.ModuleName,
		Short: "Transaction commands for simpleswap module",
		Args:  cobra.ExactArgs(1),
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(
		addLiquidityCmd(),
		swapLiquidityCmd(),
		removeLiquidityCmd(),
	)
	return cmd
}

func addLiquidityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity [provider] [token]",
		Short: "Add liquidity to the pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			providerAddress := clientCtx.GetFromAddress()

			fmt.Println("context address: ", providerAddress.String())
			fmt.Println(args[0])
			if providerAddress.String() != args[0] {
				return simpleswap.ErrInvalidProviderAddress
			}

			token, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := &simpleswap.MsgAddLiquidity{
				LiquidityProvider: providerAddress.String(),
				Token:             token,
			}

			fmt.Println("msg: ", msg.LiquidityProvider, " ", msg.Token.Denom, " ", msg.Token.Amount.String())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func swapLiquidityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-liquidity [trader] [input] [output]",
		Short: "Swap liquidity from one token to another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			providerAddress := clientCtx.GetFromAddress()

			if providerAddress.String() != args[0] {
				return simpleswap.ErrInvalidProviderAddress
			}

			input, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			output, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := &simpleswap.MsgSwapLiquidity{
				Trader: providerAddress.String(),
				Input:             input,
				Output:            output,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func removeLiquidityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity [liquidityProvider] [token]",
		Short: "Remove liquidity from the pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			providerAddress := clientCtx.GetFromAddress()

			if providerAddress.String() != args[0] {
				return simpleswap.ErrInvalidProviderAddress
			}

			token, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := &simpleswap.MsgRemoveLiquidity{
				LiquidityProvider: providerAddress.String(),
				Token:             token,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func (AppModule) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   simpleswap.ModuleName,
		Short: "Querying commands for the simpleswap module",
	}

	cmd.AddCommand(
		getParamsCmd(),
		getPoolCmd(),
		getLiquidityProviderCmd(),
		getCoinReserveCmd(),
		getCoinReservesCmd(),
	)
	return cmd
}

func getParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Get the simpleswap module parameters",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := simpleswap.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &simpleswap.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func getPoolCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Short: "Get the simpleswap module pool",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := simpleswap.NewQueryClient(clientCtx)

			res, err := queryClient.Pool(context.Background(), &simpleswap.QueryPoolRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func getLiquidityProviderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "liquidity-provider [address]",
		Short: "Get the simpleswap module liquidity provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := simpleswap.NewQueryClient(clientCtx)

			res, err := queryClient.LiquidityProvider(context.Background(), &simpleswap.QueryLiquidityProviderRequest{LpAddress: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func getCoinReserveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "coin-reserve [denom]",
		Short: "Get the simpleswap module coin reserve",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := simpleswap.NewQueryClient(clientCtx)

			res, err := queryClient.CoinReserve(context.Background(), &simpleswap.QueryCoinReserveRequest{CoinDenom: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func getCoinReservesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "coin-reserves",
		Short: "Get all the reserves",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := simpleswap.NewQueryClient(clientCtx)

			res, err := queryClient.CoinReserves(context.Background(), &simpleswap.QueryCoinReservesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
