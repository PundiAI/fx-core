package keeper_test

import (
	"testing"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	ibcmiddleware "github.com/functionx/fx-core/v8/x/ibc/middleware"
	"github.com/functionx/fx-core/v8/x/ibc/middleware/keeper"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	ibcMiddlewareKeeper keeper.Keeper
	ibcTransferKeeper   ibctransferkeeper.Keeper
	bankKeeper          bankkeeper.Keeper
	ibcMiddleware       porttypes.Middleware
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()

	suite.ibcMiddlewareKeeper = suite.App.IBCMiddlewareKeeper
	suite.ibcTransferKeeper = suite.App.IBCTransferKeeper
	suite.bankKeeper = suite.App.BankKeeper

	transferIBCModule := transfer.NewIBCModule(suite.App.IBCTransferKeeper)
	suite.ibcMiddleware = ibcmiddleware.NewIBCMiddleware(suite.ibcMiddlewareKeeper, suite.App.IBCKeeper.ChannelKeeper, transferIBCModule)
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}
