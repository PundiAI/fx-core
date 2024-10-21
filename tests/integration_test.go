package tests

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"

	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v8/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	layer2types "github.com/functionx/fx-core/v8/x/layer2/types"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v8/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
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
	suite.CrosschainTest()

	suite.StakingTest()
	suite.StakingContractTest()
	suite.StakingSharesTest()
	suite.StakingSharesContractTest()
	suite.StakingPrecompileRedelegateTest()
	suite.StakingPrecompileRedelegateByContractTest()
	suite.StakingPrecompileV2()

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

func (suite *IntegrationTest) GetByName(chainName string) CrosschainTestSuite {
	for _, c := range suite.crosschain {
		if c.chainName == chainName {
			return c
		}
	}
	panic(fmt.Sprintf("chain not found %s", chainName))
}
