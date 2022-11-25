package network_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v3/testutil"
	"github.com/functionx/fx-core/v3/testutil/network"
)

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := testutil.DefaultNetworkConfig()
	cfg.NumValidators = 1

	baseDir, err := os.MkdirTemp(s.T().TempDir(), cfg.ChainID)
	s.Require().NoError(err)
	s.T().Logf("created temporary directory: %s", baseDir)

	s.network, err = network.New(s.T(), baseDir, cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestNetwork_Liveness() {
	h, err := s.network.WaitForHeightWithTimeout(10, time.Minute)
	s.Require().NoError(err, "expected to reach 10 blocks; got %d", h)

	latestHeight, err := s.network.LatestHeight()
	s.Require().NoError(err, "latest height failed")
	s.Require().GreaterOrEqual(latestHeight, h)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
