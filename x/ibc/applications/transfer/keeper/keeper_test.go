package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	erc20keeper "github.com/functionx/fx-core/v8/x/erc20/keeper"
	fxevmkeeper "github.com/functionx/fx-core/v8/x/evm/keeper"
	fxtransfer "github.com/functionx/fx-core/v8/x/ibc/applications/transfer"
	"github.com/functionx/fx-core/v8/x/ibc/applications/transfer/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context
	cdc codec.Codec

	fxIBCTransferKeeper keeper.Keeper
	bankKeeper          bankkeeper.Keeper
	erc20Keeper         erc20keeper.Keeper
	evmKeeper           *fxevmkeeper.Keeper
	accountKeeper       authkeeper.AccountKeeper
	ibcMiddleware       porttypes.Middleware
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
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(1, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false)

	suite.cdc = suite.app.AppCodec()

	suite.fxIBCTransferKeeper = suite.app.FxTransferKeeper
	suite.bankKeeper = suite.app.BankKeeper
	suite.erc20Keeper = suite.app.Erc20Keeper
	suite.evmKeeper = suite.app.EvmKeeper
	suite.accountKeeper = suite.app.AccountKeeper

	transferIBCModule := transfer.NewIBCModule(suite.fxIBCTransferKeeper.Keeper)
	suite.ibcMiddleware = fxtransfer.NewIBCMiddleware(suite.fxIBCTransferKeeper, transferIBCModule)
}
