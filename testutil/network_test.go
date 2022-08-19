package testutil_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/functionx/fx-core/v2/testutil"

	"github.com/functionx/fx-core/v2/testutil/network"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v2/app/helpers"
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
	cfg.NumValidators = 1
	cfg.Mnemonics = append(cfg.Mnemonics, helpers.NewMnemonic())

	baseDir, err := os.MkdirTemp(suite.T().TempDir(), cfg.ChainID)
	suite.NoError(err)
	suite.network, err = network.New(suite.T(), baseDir, cfg)
	suite.NoError(err)

	_, err = suite.network.WaitForHeight(1)
	suite.NoError(err)

	//_, err := suite.network.WaitForHeight(1)
	//suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *IntegrationTestSuite) TestSuite() {
	suite.Require().Equal(sdk.GetConfig().GetCoinType(), uint32(sdk.CoinType))

	validator := suite.network.Validators[0]
	keyringDir := validator.ClientCtx.KeyringDir
	file, err := os.ReadFile(filepath.Join(keyringDir, "key_seed.json"))
	suite.Require().NoError(err)

	var data map[string]string
	err = json.Unmarshal(file, &data)
	suite.Require().NoError(err)

	mnemonic := suite.network.Config.Mnemonics[0]
	suite.Equal(mnemonic, data["secret"])

	info, err := validator.ClientCtx.Keyring.Key("node0")
	suite.NoError(err)
	suite.Equal(info.GetAddress(), validator.Address)
	suite.Equal(info.GetAlgo(), hd.Secp256k1Type)
	suite.Equal(info.GetType().String(), "local")

	keyringAlgos1, _ := validator.ClientCtx.Keyring.SupportedAlgorithms()
	suite.Equal(keyringAlgos1, hd2.SupportedAlgorithms)

	privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, hd.Secp256k1Type, 0, 0)
	suite.NoError(err)
	suite.Require().Equal(validator.Address.String(), sdk.AccAddress(privKey.PubKey().Address().Bytes()).String())
}
