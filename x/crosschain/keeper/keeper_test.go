package keeper_test

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v7/testutil"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	arbitrumtypes "github.com/functionx/fx-core/v7/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v7/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/crosschain/mock"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	layer2types "github.com/functionx/fx-core/v7/x/layer2/types"
	optimismtypes "github.com/functionx/fx-core/v7/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v7/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

type KeeperTestSuite struct {
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
		suite.Run(t, &KeeperTestSuite{
			moduleName:   moduleName,
			wfxTokenAddr: helpers.GenHexAddress().String(),
		})
	}
}

func (s *KeeperTestSuite) SetupSuite() {
	s.SetupTest()
}

func (s *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(s.moduleName)

	testCtx := testutil.DefaultContextWithDB(s.T(), key, sdk.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	s.ctx = testCtx.Ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Block: &abci.BlockParams{
				MaxGas: types.MaxGasLimit,
			},
		},
	)

	encCfg := testutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	ctrl := gomock.NewController(s.T())
	s.stakingKeeper = mock.NewMockStakingKeeper(ctrl)
	s.stakingMsgServer = mock.NewMockStakingMsgServer(ctrl)
	s.distMsgServer = mock.NewMockDistributionMsgServer(ctrl)
	s.bankKeeper = mock.NewMockBankKeeper(ctrl)
	s.ibcTransferKeeper = mock.NewMockIBCTransferKeeper(ctrl)
	s.erc20Keeper = mock.NewMockErc20Keeper(ctrl)
	s.accountKeeper = mock.NewMockAccountKeeper(ctrl)
	s.evmKeeper = mock.NewMockEVMKeeper(ctrl)

	s.accountKeeper.EXPECT().GetModuleAddress(s.moduleName).Return(authtypes.NewEmptyModuleAccount(s.moduleName).GetAddress()).Times(1)

	s.crosschainKeeper = crosschainkeeper.NewKeeper(
		encCfg.Codec,
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	crosschainRouter := crosschainkeeper.NewRouter()
	crosschainRouter.AddRoute(s.moduleName, crosschainkeeper.NewModuleHandler(s.crosschainKeeper))
	crosschainRouterKeeper := crosschainkeeper.NewRouterKeeper(crosschainRouter)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, crosschainRouterKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = crosschainkeeper.NewMsgServerRouterImpl(crosschainRouterKeeper)

	params := s.CrossChainParams()
	s.NoError(s.crosschainKeeper.SetParams(s.ctx, &params))

	s.crosschainKeeper.AddBridgeToken(s.ctx, s.wfxTokenAddr, fxtypes.DefaultDenom)
}

func (s *KeeperTestSuite) CrossChainParams() types.Params {
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

func (s *KeeperTestSuite) SetOracleSet(nonce, power, height uint64) string {
	s.crosschainKeeper.SetLatestOracleSetNonce(s.ctx, nonce)
	external := helpers.GenExternalAddr(s.moduleName)
	bridgeValidator := types.BridgeValidator{Power: power, ExternalAddress: external}
	s.crosschainKeeper.StoreOracleSet(s.ctx, types.NewOracleSet(nonce, height, types.BridgeValidators{bridgeValidator}))
	return external
}

func (s *KeeperTestSuite) AddBridgeToken(contract string) (bridgeToken *types.BridgeToken) {
	denom := types.NewBridgeDenom(s.moduleName, contract)
	s.crosschainKeeper.AddBridgeToken(s.ctx, contract, denom)
	return &types.BridgeToken{
		Token: contract,
		Denom: denom,
	}
}
