package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/crosschain/keeper"

	"github.com/functionx/fx-core/app/helpers"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	queryClient types.QueryClient
	oracles     []sdk.AccAddress
	bridgers    []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) Keeper() keeper.Keeper {
	return suite.app.BscKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = helpers.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.CrosschainKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, 3, sdk.NewInt(30000000))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, 3, sdk.NewInt(30000000))

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracles {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.ctx, proposalOracle)
}

func testModuleParams() *types.Params {
	return &types.Params{
		GravityId:                         "test",
		SignedWindow:                      20000,
		ExternalBatchTimeout:              43200000,
		AverageBlockTime:                  5000,
		AverageExternalBlockTime:          3000,
		SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		IbcTransferTimeoutHeight:          10000,
		OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil))),
	}
}
