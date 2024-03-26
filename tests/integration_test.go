package tests

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"

	arbitrumtypes "github.com/functionx/fx-core/v7/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v7/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	layer2types "github.com/functionx/fx-core/v7/x/layer2/types"
	optimismtypes "github.com/functionx/fx-core/v7/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v7/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

type IntegrationTest struct {
	*TestSuite
	crosschain []CrosschainTestSuite
	erc20      Erc20TestSuite
	evm        EvmTestSuite
	staking    StakingSuite
	precompile PrecompileTestSuite
}

func TestIntegrationTest(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skip integration test")
	}

	testSuite := NewTestSuite()
	testingSuite := &IntegrationTest{
		TestSuite: testSuite,
		crosschain: []CrosschainTestSuite{
			NewCrosschainWithTestSuite(ethtypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(bsctypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(trontypes.ModuleName, testSuite),
		},
		erc20:      NewErc20TestSuite(testSuite),
		evm:        NewEvmTestSuite(testSuite),
		staking:    NewStakingSuite(testSuite),
		precompile: NewPrecompileTestSuite(testSuite),
	}
	if runtime.GOOS == "linux" {
		evmChainModules := []string{
			polygontypes.ModuleName,
			avalanchetypes.ModuleName,
			arbitrumtypes.ModuleName,
			optimismtypes.ModuleName,
			layer2types.ModuleName,
		}
		for _, module := range evmChainModules {
			testingSuite.crosschain = append(
				testingSuite.crosschain,
				NewCrosschainWithTestSuite(module, testSuite),
			)
		}
	}

	suite.Run(t, testingSuite)
}

func (suite *IntegrationTest) TestRun() {
	suite.CrossChainTest()
	suite.OriginalCrossChainTest()

	suite.PrecompileTransferCrossChainTest()
	suite.PrecompileCrossChainTest()
	suite.PrecompileCancelSendToExternalTest()
	suite.PrecompileIncreaseBridgeFeeTest()
	suite.PrecompileCrossChainConvertedDenomTest()

	suite.ERC20TokenOriginTest()
	suite.ERC20IBCChainTokenOriginTest()
	suite.ERC20TokenERC20Test()
	suite.ERC20IBCChainTokenERC20Test()

	suite.StakingTest()
	suite.StakingContractTest()
	suite.StakingSharesTest()
	suite.StakingSharesContractTest()
	suite.StakingPrecompileRedelegateTest()
	suite.StakingPrecompileRedelegateByContractTest()

	suite.MigrateTestDelegate()
	suite.MigrateTestUnDelegate()

	suite.EVMWeb3Test()
	suite.WFXTest()
	suite.ERC20TokenTest()
	suite.ERC721Test()
	suite.CallContractTest()
	suite.FIP20CodeCheckTest()
	suite.WFXCodeCheckTest()

	suite.ByPassFeeTest()
}

func (suite *IntegrationTest) GetCrossChainByName(chainName string) CrosschainTestSuite {
	for _, c := range suite.crosschain {
		if c.chainName == chainName {
			return c
		}
	}
	panic("chain not found")
}

type IntegrationMultiNodeTest struct {
	*TestSuiteMultiNode
	staking StakingSuite
	authz   AuthzSuite
	slasing SlashingSuite
}

func TestIntegrationMultiNodeTest(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skip integration test")
	}

	testSuiteMultiNode := NewTestSuiteMultiNode()
	suite.Run(t, &IntegrationMultiNodeTest{
		TestSuiteMultiNode: testSuiteMultiNode,
		staking:            NewStakingSuite(testSuiteMultiNode.TestSuite),
		authz:              NewAuthzSuite(testSuiteMultiNode.TestSuite),
		slasing:            NewSlashingSuite(testSuiteMultiNode.TestSuite),
	})
}

func (suite *IntegrationMultiNodeTest) TestRun() {
	suite.StakingEditPubKey()
	suite.StakingGrantPrivilege()
	suite.StakingEditPubKeyJailBlock()
}
