package keeper_test

import (
	"crypto/ecdsa"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app/helpers"

	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/crosschain/keeper"

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
	for _, moduleName := range []string{bsctypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName} {
		methodFinder := reflect.TypeOf(new(KeeperTestSuite))
		for i := 0; i < methodFinder.NumMethod(); i++ {
			method := methodFinder.Method(i)
			if !strings.HasPrefix(method.Name, "Test") {
				continue
			}
			t.Run(fmt.Sprintf("%s_%s", method.Name, moduleName), func(subT *testing.T) {
				mySuite := new(KeeperTestSuite)
				mySuite.SetT(subT)
				mySuite.chainName = moduleName
				mySuite.SetupTest()
				method.Func.Call([]reflect.Value{reflect.ValueOf(mySuite)})
			})
		}
	}

}

func (suite *KeeperTestSuite) MsgServer() types.MsgServer {
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
	if suite.chainName == trontypes.ModuleName {
		types.InitMsgValidatorBasicRouter()
		types.RegisterValidatorBasic(trontypes.ModuleName, trontypes.TronMsgValidate{})
	}

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(types.MaxOracleSize, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.CrosschainKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.externals = genEthKey(types.MaxOracleSize)
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

func genEthKey(count int) []*ecdsa.PrivateKey {
	var ethKeys []*ecdsa.PrivateKey
	for i := 0; i < count; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		ethKeys = append(ethKeys, key)
	}
	return ethKeys
}
