package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	baseapp "github.com/cosmos/cosmos-sdk/baseapp"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/simpleswap"
	expectedkeepers "github.com/cosmos/simpleswap/expected_keepers"
	simpleswapKeeper "github.com/cosmos/simpleswap/keeper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	simpleSwapKeeper simpleswapKeeper.Keeper
	bankKeeper       *expectedkeepers.MockBankKeeper
	msgServer        simpleswap.MsgServer
	queryClient      simpleswap.QueryClient

	addrs []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	key := storetypes.NewKVStoreKey(simpleswap.ModuleName)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	storeService := runtime.NewKVStoreService(key)
	addrs := simtestutil.CreateIncrementalAccounts(3)
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	bankKeeper := expectedkeepers.NewMockBankKeeper(ctrl)
	k := simpleswapKeeper.NewKeeper(encCfg.Codec, addresscodec.NewBech32Codec("cosmos"), storeService, bankKeeper, addrs[0].String())
	k.Params.Set(ctx, simpleswap.DefaultParams())

	
	s.ctx = ctx
	s.bankKeeper = bankKeeper
	s.simpleSwapKeeper = k
	simpleswap.RegisterInterfaces(encCfg.InterfaceRegistry)
	s.addrs = addrs
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	simpleswap.RegisterQueryServer(queryHelper, simpleswapKeeper.NewQueryServerImpl(k))
	s.queryClient = simpleswap.NewQueryClient(queryHelper)
	s.msgServer = simpleswapKeeper.NewMsgServerImpl(k)
}
