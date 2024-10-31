package precompile_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
)

type CrosschainPrecompileTestSuite struct {
	helpers.BaseSuite

	signer         *helpers.Signer
	crosschainAddr common.Address

	helpers.CrosschainPrecompileSuite
}

func TestCrosschainPrecompileTestSuite(t *testing.T) {
	testingSuite := new(CrosschainPrecompileTestSuite)
	testingSuite.crosschainAddr = common.HexToAddress(contract.CrosschainAddress)
	suite.Run(t, testingSuite)
}

func TestCrosschainPrecompileTestSuite_Contract(t *testing.T) {
	suite.Run(t, new(CrosschainPrecompileTestSuite))
}

func (suite *CrosschainPrecompileTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.Ctx = suite.Ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.Ctx = suite.Ctx.WithBlockGasMeter(storetypes.NewGasMeter(1e18))

	suite.signer = suite.AddTestSigner(10_000)

	if suite.crosschainAddr.String() != contract.StakingAddress {
		crosschainContract, err := suite.App.EvmKeeper.DeployContract(suite.Ctx, suite.signer.Address(), contract.MustABIJson(testscontract.CrosschainTestMetaData.ABI), contract.MustDecodeHex(testscontract.CrosschainTestMetaData.Bin))
		suite.Require().NoError(err)
		suite.crosschainAddr = crosschainContract
	}

	suite.CrosschainPrecompileSuite = helpers.NewCrosschainPrecompileSuite(suite.Require(), suite.signer, suite.App.EvmKeeper, suite.crosschainAddr)
}

func (suite *CrosschainPrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *CrosschainPrecompileTestSuite) Commit() {
	header := suite.Ctx.BlockHeader()
	_, err := suite.App.EndBlocker(suite.Ctx)
	suite.Require().NoError(err)
	_, err = suite.App.Commit()
	suite.Require().NoError(err)
	// after commit Ctx header
	header.Height++

	// begin block
	header.Time = time.Now().UTC()
	header.Height++
	suite.Ctx = suite.Ctx.WithBlockHeader(header)
	_, err = suite.App.BeginBlocker(suite.Ctx)
	suite.Require().NoError(err)
	suite.Ctx = suite.Ctx.WithBlockHeight(header.Height)
}

func (suite *CrosschainPrecompileTestSuite) CrosschainKeepers() map[string]crosschainkeeper.Keeper {
	value := reflect.ValueOf(suite.App.CrosschainKeepers)
	keepers := make(map[string]crosschainkeeper.Keeper)
	for i := 0; i < value.NumField(); i++ {
		res := value.Field(i).MethodByName("GetGravityID").Call([]reflect.Value{reflect.ValueOf(suite.Ctx)})
		gravityID := res[0].String()
		chainName := strings.TrimSuffix(strings.TrimPrefix(gravityID, "fx-"), "-bridge")
		if chainName == "bridge-eth" {
			keepers["eth"] = value.Field(i).Interface().(crosschainkeeper.Keeper)
		} else {
			keepers[chainName] = value.Field(i).Interface().(crosschainkeeper.Keeper)
		}
	}
	return keepers
}

func (suite *CrosschainPrecompileTestSuite) GenerateModuleName() string {
	keepers := suite.CrosschainKeepers()
	modules := make([]string, 0, len(keepers))
	for m := range keepers {
		modules = append(modules, m)
	}
	if len(modules) == 0 {
		return ""
	}
	return modules[tmrand.Intn(len(modules))]
}

func (suite *CrosschainPrecompileTestSuite) GenerateOracles(moduleName string, online bool, num int) []Oracle {
	keeper := suite.CrosschainKeepers()[moduleName]
	oracles := make([]Oracle, 0, num)
	for i := 0; i < num; i++ {
		oracle := crosschaintypes.Oracle{
			OracleAddress:     helpers.GenAccAddress().String(),
			BridgerAddress:    helpers.GenAccAddress().String(),
			ExternalAddress:   helpers.GenExternalAddr(moduleName),
			DelegateAmount:    sdkmath.NewInt(1e18).MulRaw(1000),
			StartHeight:       1,
			Online:            online,
			DelegateValidator: sdk.ValAddress(helpers.GenAccAddress()).String(),
			SlashTimes:        0,
		}
		keeper.SetOracle(suite.Ctx, oracle)
		keeper.SetOracleAddrByExternalAddr(suite.Ctx, oracle.ExternalAddress, oracle.GetOracle())
		keeper.SetOracleAddrByBridgerAddr(suite.Ctx, oracle.GetBridger(), oracle.GetOracle())
		oracles = append(oracles, Oracle{
			moduleName: moduleName,
			oracle:     oracle,
		})
	}
	return oracles
}

func (suite *CrosschainPrecompileTestSuite) GenerateRandOracle(moduleName string, online bool) Oracle {
	oracles := suite.GenerateOracles(moduleName, online, 1)
	return oracles[0]
}

type Oracle struct {
	moduleName string
	oracle     crosschaintypes.Oracle
}

func (o Oracle) GetExternalHexAddr() common.Address {
	return crosschaintypes.ExternalAddrToHexAddr(o.moduleName, o.oracle.ExternalAddress)
}
