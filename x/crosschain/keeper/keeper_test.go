package keeper_test

import (
	"crypto/ecdsa"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/functionx/fx-core/app/helpers"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/crosschain/keeper"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/crosschain/types"
)

type KeeperTestSuite struct {
	suite.Suite
	sync.Mutex

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
	methodFinder := reflect.TypeOf(new(KeeperTestSuite))
	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)
		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}
		t.Run(method.Name, func(subT *testing.T) {
			mySuite := new(KeeperTestSuite)
			mySuite.SetT(t)
			mySuite.SetupTest()
			method.Func.Call([]reflect.Value{reflect.ValueOf(mySuite)})
		})
	}
}

func (suite *KeeperTestSuite) Msg() types.MsgServer {
	return keeper.NewMsgServerImpl(suite.Keeper())
}

func (suite *KeeperTestSuite) Keeper() keeper.Keeper {
	return suite.app.BscKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(types.MaxOracleSize, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.CrosschainKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient

	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.externals = genEthKey(types.MaxOracleSize)
	suite.delegateAmount = sdk.NewInt(10 * 1e3).MulRaw(1e18)
	for i := 0; i < types.MaxOracleSize; i++ {
		suite.validator = append(suite.validator, valAccounts[i].GetAddress().Bytes())
	}
	suite.chainName = "bsc"

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
