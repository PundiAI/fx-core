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

	"github.com/functionx/fx-core/v8/testutil"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v8/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	"github.com/functionx/fx-core/v8/x/crosschain/mock"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	layer2types "github.com/functionx/fx-core/v8/x/layer2/types"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v8/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

type KeeperMockSuite struct {
	suite.Suite

	ctx          sdk.Context
	moduleName   string
	wfxTokenAddr string

	queryClient types.QueryClient
	msgServer   types.MsgServer

	crosschainKeeper  crosschainkeeper.Keeper
	stakingKeeper     *mock.MockStakingKeeper
	stakingMsgServer  *mock.MockStakingMsgServer
	distMsgServer     *mock.MockDistributionMsgServer
	bankKeeper        *mock.MockBankKeeper
	ibcTransferKeeper *mock.MockIBCTransferKeeper
	erc20Keeper       *mock.MockErc20Keeper
	accountKeeper     *mock.MockAccountKeeper
	evmKeeper         *mock.MockEVMKeeper
	evmErc20Keeper    *mock.MockEvmERC20Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	mustTestModule := []string{
		trontypes.ModuleName,
		ethtypes.ModuleName,
	}
	subModules := mustTestModule
	if os.Getenv("TEST_CROSSCHAIN") == "true" {
		subModules = append(subModules, []string{
			bsctypes.ModuleName,
			polygontypes.ModuleName,
			avalanchetypes.ModuleName,
			arbitrumtypes.ModuleName,
			optimismtypes.ModuleName,
			layer2types.ModuleName,
		}...)
	}
	for _, moduleName := range subModules {
		suite.Run(t, &KeeperMockSuite{
			moduleName:   moduleName,
			wfxTokenAddr: helpers.GenHexAddress().String(),
		})
	}
}

func (s *KeeperMockSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(s.moduleName)

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
	s.evmErc20Keeper = mock.NewMockEvmERC20Keeper(ctrl)

	s.accountKeeper.EXPECT().GetModuleAddress(s.moduleName).Return(authtypes.NewEmptyModuleAccount(s.moduleName).GetAddress()).Times(1)

	s.crosschainKeeper = crosschainkeeper.NewKeeper(
		myApp.AppCodec(),
		s.moduleName,
		key,
		s.stakingKeeper,
		s.stakingMsgServer,
		s.distMsgServer,
		s.bankKeeper,
		s.ibcTransferKeeper,
		s.erc20Keeper,
		s.accountKeeper,
		s.evmKeeper,
		s.evmErc20Keeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	crosschainRouter := crosschainkeeper.NewRouter()
	crosschainRouter.AddRoute(s.moduleName, crosschainkeeper.NewModuleHandler(s.crosschainKeeper))
	crosschainRouterKeeper := crosschainkeeper.NewRouterKeeper(crosschainRouter)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, myApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, crosschainRouterKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = crosschainkeeper.NewMsgServerRouterImpl(crosschainRouterKeeper)

	params := s.CrossChainParams()
	params.EnableSendToExternalPending = true
	s.NoError(s.crosschainKeeper.SetParams(s.ctx, &params))

	bridgeDenom := types.NewBridgeDenom(s.moduleName, s.wfxTokenAddr)
	s.crosschainKeeper.AddBridgeToken(s.ctx, bridgeDenom, fxtypes.DefaultDenom)
	s.crosschainKeeper.AddBridgeToken(s.ctx, fxtypes.DefaultDenom, bridgeDenom)
}

func (s *KeeperMockSuite) SetupSubTest() {
	s.SetupTest()
}

func (s *KeeperMockSuite) CrossChainParams() types.Params {
	switch s.moduleName {
	case ethtypes.ModuleName:
		return ethtypes.DefaultGenesisState().Params
	case bsctypes.ModuleName:
		return bsctypes.DefaultGenesisState().Params
	case polygontypes.ModuleName:
		return polygontypes.DefaultGenesisState().Params
	case trontypes.ModuleName:
		return trontypes.DefaultGenesisState().Params
	case avalanchetypes.ModuleName:
		return avalanchetypes.DefaultGenesisState().Params
	case optimismtypes.ModuleName:
		return optimismtypes.DefaultGenesisState().Params
	case arbitrumtypes.ModuleName:
		return arbitrumtypes.DefaultGenesisState().Params
	case layer2types.ModuleName:
		return layer2types.DefaultGenesisState().Params
	default:
		panic("module not support")
	}
}
