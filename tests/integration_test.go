package tests

import (
	"testing"

	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"

	"github.com/stretchr/testify/suite"
)

type IntegrationTest struct {
	*TestSuite
	crosschain []CrosschainTestSuite
	erc20      Erc20TestSuite
	evm        EvmTestSuite
}

func TestIntegrationTest(t *testing.T) {
	testSuite := NewTestSuite()
	suite.Run(t, &IntegrationTest{
		TestSuite: testSuite,
		crosschain: []CrosschainTestSuite{
			NewCrosschainWithTestSuite(bsctypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(polygontypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(trontypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(avalanchetypes.ModuleName, testSuite),
			NewCrosschainWithTestSuite(ethtypes.ModuleName, testSuite),
		},
		erc20: NewErc20TestSuite(testSuite),
		evm:   NewEvmTestSuite(testSuite),
	})
}

func (suite *IntegrationTest) TestRun() {
	suite.CrossChainTest()
	suite.ERC20Test()
	suite.EVMWeb3Test()
	suite.MigrateTestDelegate()
	suite.MigrateTestUnDelegate()
	suite.WFXTest()
	suite.ERC20TokenTest()
	suite.ERC721Test()
}
