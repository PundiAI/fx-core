package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"
	fxibctesting "github.com/functionx/fx-core/v3/x/ibc/testing"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	// chainA/chainB is fxApp, chainC is simApp
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain

	queryClient types.QueryClient
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *KeeperTestSuite) DoSetupTest(t *testing.T) {
	suite.coordinator = fxibctesting.NewCoordinator(t, 2, 1)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(3))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.GetApp(suite.chainA.App).InterfaceRegistry())
	transfertypes.RegisterQueryServer(queryHelper, suite.GetApp(suite.chainA.App).IBCTransferKeeper)
	suite.queryClient = transfertypes.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) GetApp(testingApp ibctesting.TestingApp) *app.App {
	result, ok := testingApp.(*app.App)
	suite.Require().True(ok)
	return result
}
