package testutil_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/testutil/network"
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

	encCfg := app.MakeEncodingConfig()
	cfg := testutil.DefaultNetworkConfig(encCfg, func(config *network.Config) {
		// config.EnableTMLogging = true
	})

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

	gotHeight, err := suite.network.WaitForHeightWithTimeout(height, time.Second*3)
	suite.NoError(err, "expected to reach %d blocks; got %d", height, gotHeight)

	latestHeight, err = suite.network.LatestHeight()
	suite.NoError(err, "latest height failed")
	suite.GreaterOrEqual(latestHeight, gotHeight)
}

func (suite *IntegrationTestSuite) TestValidatorInfo() {
	suite.Equal(sdk.GetConfig().GetCoinType(), uint32(sdk.CoinType))

	suite.Equal(len(suite.network.Config.Mnemonics), suite.network.Config.NumValidators)
	for i := 0; i < suite.network.Config.NumValidators; i++ {

		validator := suite.network.Validators[i]

		mnemonic := suite.network.Config.Mnemonics[i]
		keySeedFileName := filepath.Join(validator.ClientCtx.KeyringDir, "key_seed.json")
		if _, err := os.Stat(keySeedFileName); err == nil {
			file, err := os.ReadFile(keySeedFileName)
			suite.NoError(err)

			var data map[string]string
			suite.NoError(json.Unmarshal(file, &data))
			suite.Equal(mnemonic, data["secret"])
		}

		key, err := validator.ClientCtx.Keyring.Key(validator.Ctx.Config.Moniker)
		suite.NoError(err)
		addr, err := key.GetAddress()
		suite.NoError(err)
		suite.Equal(addr, validator.Address)
		suite.Equal(key.GetType().String(), "local")

		keyringAlgos1, _ := validator.ClientCtx.Keyring.SupportedAlgorithms()
		suite.Equal(keyringAlgos1, hd2.SupportedAlgorithms)

		privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, hd.Secp256k1Type, 0, 0)
		suite.NoError(err)
		suite.Equal(validator.Address.Bytes(), privKey.PubKey().Address().Bytes())
	}
}
