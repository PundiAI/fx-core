package testutil_test

import (
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/testutil"
	"github.com/functionx/fx-core/v3/testutil/network"
)

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.T().Log("setting up integration test suite")

	cfg := testutil.DefaultNetworkConfig()
	// cfg.EnableTMLogging = true

	baseDir, err := os.MkdirTemp(suite.T().TempDir(), cfg.ChainID)
	suite.Require().NoError(err)
	suite.network, err = network.New(suite.T(), baseDir, cfg)
	suite.Require().NoError(err)

	_, err = suite.network.WaitForHeight(1)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *IntegrationTestSuite) TestNetworkLiveness() {
	height := int64(3)
	latestHeight, err := suite.network.LatestHeight()
	suite.NoError(err, "latest height failed")
	suite.GreaterOrEqual(height, latestHeight)

	height, err = suite.network.WaitForHeightWithTimeout(height, time.Second*3)
	suite.NoError(err, "expected to reach 200 blocks; got %d", height)

	latestHeight, err = suite.network.LatestHeight()
	suite.NoError(err, "latest height failed")
	suite.GreaterOrEqual(latestHeight, height)
}

func (suite *IntegrationTestSuite) TestValidatorInfo() {
	suite.Equal(sdk.GetConfig().GetCoinType(), uint32(sdk.CoinType))

	suite.Equal(len(suite.network.Config.Mnemonics), suite.network.Config.NumValidators)
	for i := 0; i < suite.network.Config.NumValidators; i++ {

		validator := suite.network.Validators[i]
		// keyringDir := validator.ClientCtx.KeyringDir
		// file, err := os.ReadFile(filepath.Join(keyringDir, "key_seed.json"))
		// suite.NoError(err)

		// var data map[string]string
		// err = json.Unmarshal(file, &data)
		// suite.NoError(err)

		mnemonic := suite.network.Config.Mnemonics[i]
		// suite.Equal(mnemonic, data["secret"])

		info, err := validator.ClientCtx.Keyring.Key(validator.Moniker)
		suite.NoError(err)
		suite.Equal(info.GetAddress(), validator.Address)
		suite.Equal(info.GetAlgo(), hd.PubKeyType(suite.network.Config.SigningAlgo))
		suite.Equal(info.GetType().String(), "local")

		keyringAlgos1, _ := validator.ClientCtx.Keyring.SupportedAlgorithms()
		suite.Equal(keyringAlgos1, hd2.SupportedAlgorithms)

		privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, hd.Secp256k1Type, 0, 0)
		suite.NoError(err)
		suite.Equal(validator.Address.Bytes(), privKey.PubKey().Address().Bytes())
	}
}
