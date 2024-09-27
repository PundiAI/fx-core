package keeper_test

import (
	"testing"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	erc20keeper "github.com/functionx/fx-core/v8/x/erc20/keeper"
	fxevmkeeper "github.com/functionx/fx-core/v8/x/evm/keeper"
	ibcmiddleware "github.com/functionx/fx-core/v8/x/ibc/middleware"
	"github.com/functionx/fx-core/v8/x/ibc/middleware/keeper"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	ibcMiddlewareKeeper keeper.Keeper
	ibcTransferKeeper   ibctransferkeeper.Keeper
	bankKeeper          bankkeeper.Keeper
	erc20Keeper         erc20keeper.Keeper
	evmKeeper           *fxevmkeeper.Keeper
	accountKeeper       authkeeper.AccountKeeper
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
	suite.erc20Keeper = suite.App.Erc20Keeper
	suite.evmKeeper = suite.App.EvmKeeper
	suite.accountKeeper = suite.App.AccountKeeper

	transferIBCModule := transfer.NewIBCModule(suite.App.IBCTransferKeeper)
	suite.ibcMiddleware = ibcmiddleware.NewIBCMiddleware(suite.ibcMiddlewareKeeper, suite.App.IBCKeeper.ChannelKeeper, transferIBCModule)
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}
