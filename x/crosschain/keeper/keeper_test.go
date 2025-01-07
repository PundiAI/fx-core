package keeper_test

import (
	"os"
	"testing"

	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	"github.com/pundiai/fx-core/v8/x/crosschain/mock"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

type KeeperMockSuite struct {
	suite.Suite

	ctx       sdk.Context
	chainName string

	queryClient types.QueryClient
	// msgServer   types.MsgServer

	crosschainKeeper  crosschainkeeper.Keeper
	stakingKeeper     *mock.MockStakingKeeper
	stakingMsgServer  *mock.MockStakingMsgServer
	distMsgServer     *mock.MockDistributionMsgServer
	bankKeeper        *mock.MockBankKeeper
	ibcTransferKeeper *mock.MockIBCTransferKeeper
	erc20Keeper       *mock.MockErc20Keeper
	accountKeeper     *mock.MockAccountKeeper
	evmKeeper         *mock.MockEVMKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	modules := []string{
		trontypes.ModuleName,
		ethtypes.ModuleName,
	}
	if os.Getenv("TEST_CROSSCHAIN") == "true" {
		modules = fxtypes.GetSupportChains()
	}
	for _, moduleName := range modules {
		suite.Run(t, &KeeperMockSuite{chainName: moduleName})
	}
}

func (s *KeeperMockSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(s.chainName)

	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	s.ctx = testCtx.Ctx.WithConsensusParams(
		tmproto.ConsensusParams{
			Block: &tmproto.BlockParams{
				MaxGas: types.MaxGasLimit,
			},
		},
	)

	myApp := helpers.NewApp()

	ctrl := gomock.NewController(s.T())
	s.stakingKeeper = mock.NewMockStakingKeeper(ctrl)
	s.stakingMsgServer = mock.NewMockStakingMsgServer(ctrl)
	s.distMsgServer = mock.NewMockDistributionMsgServer(ctrl)
	s.bankKeeper = mock.NewMockBankKeeper(ctrl)
	s.ibcTransferKeeper = mock.NewMockIBCTransferKeeper(ctrl)
	s.erc20Keeper = mock.NewMockErc20Keeper(ctrl)
	s.accountKeeper = mock.NewMockAccountKeeper(ctrl)
	s.evmKeeper = mock.NewMockEVMKeeper(ctrl)

	s.accountKeeper.EXPECT().GetModuleAddress(s.chainName).Return(authtypes.NewEmptyModuleAccount(s.chainName).GetAddress()).Times(1)

	s.crosschainKeeper = crosschainkeeper.NewKeeper(
		myApp.AppCodec(),
		s.chainName,
		key,
		s.stakingKeeper,
		s.stakingMsgServer,
		s.distMsgServer,
		s.bankKeeper,
		s.ibcTransferKeeper,
		s.erc20Keeper,
		s.accountKeeper,
		s.evmKeeper,
		contract.NewBridgeFeeQuoteKeeper(nil, contract.BridgeFeeAddress),
		contract.NewERC20TokenKeeper(nil),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	crosschainRouter := crosschainkeeper.NewRouter()
	crosschainRouter.AddRoute(s.chainName, crosschainkeeper.NewModuleHandler(s.crosschainKeeper))
	crosschainRouterKeeper := crosschainkeeper.NewRouterKeeper(crosschainRouter)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, myApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, crosschainRouterKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	// s.msgServer = crosschainkeeper.NewMsgServerRouterImpl(crosschainRouterKeeper)
}

func (s *KeeperMockSuite) SetupSubTest() {
	s.SetupTest()
}
