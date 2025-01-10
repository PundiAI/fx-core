package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/client/grpc"
	"github.com/pundiai/fx-core/v8/testutil"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/testutil/network"
)

type IntegrationTest struct {
	network       *network.Network
	numValidator  int
	enableLogging bool

	*FxCoreSuite
}

func TestIntegrationTest(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skip integration test")
	}

	suite.Run(t, &IntegrationTest{FxCoreSuite: &FxCoreSuite{EthSuite: &EthSuite{}}})
}

func (suite *IntegrationTest) SetupSuite() {
	suite.T().Log("setting up integration test suite")

	suite.numValidator = 1
	suite.enableLogging = false

	timeoutCommit := 50 * time.Millisecond
	if suite.numValidator > 1 {
		timeoutCommit = 500 * time.Millisecond
	}

	ibcGenesisOpt := func(config *network.Config) {
		config.GenesisState = testutil.IbcGenesisState(config.Codec, config.GenesisState)
	}
	bankGenesisOpt := func(config *network.Config) {
		config.GenesisState = testutil.BankGenesisState(config.Codec, config.GenesisState)
	}
	govGenesisOpt := func(config *network.Config) {
		votingPeriod := time.Millisecond
		if suite.numValidator > 1 {
			votingPeriod = time.Duration(suite.numValidator*5) * timeoutCommit
		}
		config.GenesisState = testutil.GovGenesisState(config.Codec, config.GenesisState, votingPeriod)
	}
	slashingGenesisOpt := func(config *network.Config) {
		signedBlocksWindow := int64(10)
		minSignedPerWindow := sdkmath.LegacyNewDecWithPrec(2, 1)
		downtimeJailDuration := 5 * time.Second
		config.GenesisState = testutil.SlashingGenesisState(config.Codec, config.GenesisState, signedBlocksWindow, minSignedPerWindow, downtimeJailDuration)
	}

	cfg := testutil.DefaultNetworkConfig(ibcGenesisOpt, bankGenesisOpt, govGenesisOpt, slashingGenesisOpt)
	cfg.TimeoutCommit = timeoutCommit
	cfg.NumValidators = suite.numValidator
	cfg.EnableJSONRPC = true
	if suite.enableLogging {
		cfg.EnableTMLogging = true
	}

	suite.network = network.New(suite.T(), cfg)

	_, err := suite.network.WaitForHeight(3)
	suite.Require().NoError(err)

	suite.FxCoreSuite.EthSuite.ctx = suite.network.GetContext()
	suite.FxCoreSuite.EthSuite.ethCli = suite.network.Validators[0].JSONRPCClient

	suite.FxCoreSuite.codec = suite.network.Config.Codec
	suite.FxCoreSuite.validators = suite.GetAllValSigners()
	suite.FxCoreSuite.grpcCli = suite.GRPCClient(suite.ctx)
	suite.FxCoreSuite.gasPrices = suite.GetGasPrices()
	suite.FxCoreSuite.defDenom = suite.network.Config.BondDenom
	suite.FxCoreSuite.timeoutCommit = timeoutCommit
	suite.FxCoreSuite.waitForHeightFunc = suite.network.WaitForHeight
}

func (suite *IntegrationTest) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *IntegrationTest) TestRun() {
	suite.CrosschainTest()

	suite.StakingTest()
	suite.StakingSharesTest()
	suite.StakingPrecompileRedelegateTest()
	suite.StakingPrecompileV2()

	suite.StakingContractTest()
	suite.StakingSharesContractTest()
	suite.StakingPrecompileRedelegateByContractTest()

	suite.MigrateTestDelegate()
	suite.MigrateTestUnDelegate()

	suite.EVMWeb3Test()
	suite.WFXTest()
	suite.ERC20TokenTest()
	suite.ERC721Test()
	suite.CallContractTest()
	suite.ERC20CodeTest()
	suite.WFXCodeTest()

	suite.ByPassFeeTest()
}

func (suite *IntegrationTest) GetAllValSigners() []*helpers.Signer {
	signers := make([]*helpers.Signer, 0, len(suite.network.Config.Mnemonics))
	for _, mnemonics := range suite.network.Config.Mnemonics {
		privKey, err := helpers.PrivKeyFromMnemonic(mnemonics, hd.Secp256k1Type, 0, 0)
		suite.Require().NoError(err)
		signers = append(signers, helpers.NewSigner(privKey))
	}
	return signers
}

func (suite *IntegrationTest) GetGasPrices() sdk.Coins {
	gasPrices, err := sdk.ParseCoinsNormalized(suite.network.Config.MinGasPrices)
	suite.Require().NoError(err)
	if gasPrices.Len() <= 0 {
		// Let me know if you use sdk.newCoins sanitizeCoins will remove all zero coins
		gasPrices = sdk.Coins{suite.NewCoin(sdkmath.ZeroInt())}
	}
	return gasPrices
}

func (suite *IntegrationTest) GRPCClient(ctx context.Context) *grpc.Client {
	validator := suite.network.Validators[0]
	if validator.ClientCtx.GRPCClient != nil {
		return grpc.NewClient(validator.ClientCtx)
	}
	grpcUrl := fmt.Sprintf("http://%s", validator.AppConfig.GRPC.Address)
	client, err := grpc.DailClient(grpcUrl, ctx)
	suite.Require().NoError(err)
	return client
}
