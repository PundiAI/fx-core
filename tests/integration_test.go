package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	bsctypes "github.com/functionx/fx-core/v5/x/bsc/types"
	ethtypes "github.com/functionx/fx-core/v5/x/eth/types"
	trontypes "github.com/functionx/fx-core/v5/x/tron/types"
)

type IntegrationTest struct {
	*TestSuite
	crosschain []CrosschainTestSuite
	erc20      Erc20TestSuite
	evm        EvmTestSuite
	staking    StakingSuite
}

func TestIntegrationTest(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skip integration test")
	}

	testSuite := NewTestSuite()
	suite.Run(t, &IntegrationTest{
		TestSuite: testSuite,
		crosschain: []CrosschainTestSuite{
			NewCrosschainWithTestSuite(ethtypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(bsctypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(trontypes.ModuleName, testSuite),
			// NewCrosschainWithTestSuite(polygontypes.ModuleName, testSuite),
			// NewCrosschainWithTestSuite(avalanchetypes.ModuleName, testSuite),
		},
		erc20:   NewErc20TestSuite(testSuite),
		evm:     NewEvmTestSuite(testSuite),
		staking: NewStakingSuite(testSuite),
	})
}

func (suite *IntegrationTest) TestRun() {
	suite.CrossChainTest()
	suite.ERC20Test()
	suite.StakingTest()
	suite.StakingContractTest()
	suite.StakingSharesTest()
	suite.StakingSharesContractTest()
	suite.ERC20IBCChainTokenTest()
	suite.EVMWeb3Test()
	suite.MigrateTestDelegate()
	suite.MigrateTestUnDelegate()
	suite.WFXTest()
	suite.ERC20TokenTest()
	suite.ERC721Test()
	suite.CallContractTest()
	suite.OriginERC20Test()
	suite.OriginERC20IBCChainTokenTest()
	suite.FIP20CodeCheckTest()
	suite.WFXCodeCheckTest()
	suite.ByPassFeeTest()
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
