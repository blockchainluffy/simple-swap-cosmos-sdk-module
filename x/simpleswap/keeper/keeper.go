package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	types "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/simpleswap"
	expectedkeepers "github.com/cosmos/simpleswap/expected_keepers"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema             collections.Schema
	Params             collections.Item[simpleswap.Params]
	Pool               collections.Item[simpleswap.Pool]
	LiquidityProviders collections.Map[string, simpleswap.LiquidityProvider]
	CoinsReserve       collections.Map[string, types.Coin]
	BankKeeper         expectedkeepers.BankKeeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, bankKeeper expectedkeepers.BankKeeper, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:                cdc,
		addressCodec:       addressCodec,
		authority:          authority,
		Params:             collections.NewItem(sb, simpleswap.ParamsKey, "params", codec.CollValue[simpleswap.Params](cdc)),
		Pool:               collections.NewItem(sb, simpleswap.PoolKey, "pool", codec.CollValue[simpleswap.Pool](cdc)),
		LiquidityProviders: collections.NewMap(sb, simpleswap.LiquidityProvidersKey, "liquidity_providers", collections.StringKey, codec.CollValue[simpleswap.LiquidityProvider](cdc)),
		CoinsReserve:       collections.NewMap(sb, simpleswap.CoinsReserveKey, "coins_reserve", collections.StringKey, codec.CollValue[types.Coin](cdc)),
		BankKeeper:         bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
