package testutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/testutil/network"
)

type NetworkTestSuite struct {
	suite.Suite

	network *network.Network
}

func TestNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkTestSuite))
}

func (suite *NetworkTestSuite) SetupSuite() {
	suite.T().Log("setting up integration test suite")

	cfg := testutil.DefaultNetworkConfig(func(config *network.Config) {
		config.EnableTMLogging = true
	})

	suite.network = network.New(suite.T(), cfg)

	_, err := suite.network.WaitForHeight(1)
	suite.Require().NoError(err)
}

func (suite *NetworkTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create a network!
	suite.network.Cleanup()
}

func (suite *NetworkTestSuite) TestNetworkLiveness() {
	height := int64(3)
	latestHeight, err := suite.network.LatestHeight()
	suite.Require().NoError(err, "latest height failed")
	suite.GreaterOrEqual(height, latestHeight)

	gotHeight, err := suite.network.WaitForHeightWithTimeout(height, time.Second*3)
	suite.Require().NoError(err, "expected to reach %d blocks; got %d", height, gotHeight)

	latestHeight, err = suite.network.LatestHeight()
	suite.Require().NoError(err, "latest height failed")
	suite.GreaterOrEqual(latestHeight, gotHeight)
}

func (suite *NetworkTestSuite) TestValidatorInfo() {
	suite.Equal(sdk.GetConfig().GetCoinType(), uint32(sdk.CoinType))

	suite.Equal(len(suite.network.Config.Mnemonics), suite.network.Config.NumValidators)
	for i := 0; i < suite.network.Config.NumValidators; i++ {
		validator := suite.network.Validators[i]
		mnemonic := suite.network.Config.Mnemonics[i]

		key, err := validator.ClientCtx.Keyring.Key(validator.Ctx.Config.Moniker)
		suite.Require().NoError(err)
		addr, err := key.GetAddress()
		suite.Require().NoError(err)
		suite.Equal(addr, validator.Address)
		suite.Equal("local", key.GetType().String())

		keyringAlgos1, _ := validator.ClientCtx.Keyring.SupportedAlgorithms()
		suite.Equal(keyringAlgos1, hd2.SupportedAlgorithms)

		privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, hd.Secp256k1Type, 0, 0)
		suite.Require().NoError(err)
		suite.Equal(validator.Address.Bytes(), privKey.PubKey().Address().Bytes())
	}
}

func (suite *NetworkTestSuite) TestValidatorsPower() {
	for _, val := range suite.network.Validators {
		result, err := val.RPCClient.Validators(context.Background(), nil, nil, nil)
		suite.Require().NoError(err)
		suite.Equal(4, result.Total)
		suite.Equal(len(result.Validators), result.Total)
		var totalProposerPriority int64
		for _, validator := range result.Validators {
			totalProposerPriority += validator.ProposerPriority
			suite.Equal(int64(100), validator.VotingPower)
		}
		suite.Equal(int64(0), totalProposerPriority)
	}
}
