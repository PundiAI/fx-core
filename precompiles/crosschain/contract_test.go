package crosschain_test

import (
	"testing"

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
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
)

type CrosschainPrecompileTestSuite struct {
	helpers.BaseSuite

	signer         *helpers.Signer
	crosschainAddr common.Address
	chainName      string

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

	chainNames := crosschaintypes.GetSupportChains()
	suite.chainName = chainNames[tmrand.Intn(len(chainNames))]
}

func (suite *CrosschainPrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *CrosschainPrecompileTestSuite) SetOracle(online bool) crosschaintypes.Oracle {
	oracle := crosschaintypes.Oracle{
		OracleAddress:     helpers.GenAccAddress().String(),
		BridgerAddress:    helpers.GenAccAddress().String(),
		ExternalAddress:   helpers.GenExternalAddr(suite.chainName),
		DelegateAmount:    sdkmath.NewInt(1e18).MulRaw(1000),
		StartHeight:       1,
		Online:            online,
		DelegateValidator: sdk.ValAddress(helpers.GenAccAddress()).String(),
		SlashTimes:        0,
	}
	keeper := suite.App.CrosschainKeepers.GetKeeper(suite.chainName)
	keeper.SetOracle(suite.Ctx, oracle)
	keeper.SetOracleAddrByExternalAddr(suite.Ctx, oracle.ExternalAddress, oracle.GetOracle())
	keeper.SetOracleAddrByBridgerAddr(suite.Ctx, oracle.GetBridger(), oracle.GetOracle())
	return oracle
}
