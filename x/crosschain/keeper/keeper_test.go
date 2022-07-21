package keeper_test

import (
	"crypto/ecdsa"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	tronkeeper "github.com/functionx/fx-core/v2/x/tron/keeper"

	bsctypes "github.com/functionx/fx-core/v2/x/bsc/types"
	polygontypes "github.com/functionx/fx-core/v2/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v2/x/tron/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v2/app/helpers"

	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v2/app"
	"github.com/functionx/fx-core/v2/x/crosschain/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v2/x/crosschain/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app            *app.App
	ctx            sdk.Context
	oracles        []sdk.AccAddress
	bridgers       []sdk.AccAddress
	externals      []*ecdsa.PrivateKey
	validator      []sdk.ValAddress
	chainName      string
	delegateAmount sdk.Int
	queryClient    types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	compile, err := regexp.Compile("^Test")
	require.NoError(t, err)
	for _, moduleName := range []string{bsctypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName} {
		methodFinder := reflect.TypeOf(new(KeeperTestSuite))
		for i := 0; i < methodFinder.NumMethod(); i++ {
			method := methodFinder.Method(i)
			if !compile.MatchString(method.Name) {
				continue
			}
			t.Run(method.Name, func(subT *testing.T) {
				mySuite := &KeeperTestSuite{chainName: moduleName}
				mySuite.SetT(subT)
				mySuite.SetupTest()
				method.Func.Call([]reflect.Value{reflect.ValueOf(mySuite)})
			})
		}
	}
}

func (suite *KeeperTestSuite) MsgServer() types.MsgServer {
	if suite.chainName == trontypes.ModuleName {
		return tronkeeper.NewMsgServerImpl(suite.app.TronKeeper)
	}
	return keeper.NewMsgServerImpl(suite.Keeper())
}

func (suite *KeeperTestSuite) Keeper() keeper.Keeper {
	switch suite.chainName {
	case bsctypes.ModuleName:
		return suite.app.BscKeeper
	case polygontypes.ModuleName:
		return suite.app.PolygonKeeper
	case trontypes.ModuleName:
		return suite.app.TronKeeper.Keeper
	default:
		panic("invalid chain name")
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(types.MaxOracleSize, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.CrosschainKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.externals = helpers.CreateMultiEthKey(types.MaxOracleSize)
	suite.delegateAmount = sdk.NewInt(10 * 1e3).MulRaw(1e18)
	for i := 0; i < types.MaxOracleSize; i++ {
		suite.validator = append(suite.validator, valAccounts[i].GetAddress().Bytes())
	}

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracles {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.ctx, proposalOracle)
}
