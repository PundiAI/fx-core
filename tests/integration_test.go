package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type IntegrationTest struct {
	*TestSuite
	crosschain []CrosschainTestSuite
	erc20      Erc20TestSuite
	evm        EvmTestSuite
	staking    StakingSuite
}

func TestIntegrationTest(t *testing.T) {
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
}
