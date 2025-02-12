package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/erc20/keeper"
)

type KeeperTestSuite struct {
	helpers.BaseSuite
	erc20TokenSuite helpers.ERC20TokenSuite

	signer *helpers.Signer
}

func TestERC20KeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.erc20TokenSuite = helpers.NewERC20Suite(suite.Require(), suite.App.EvmKeeper)

	suite.signer = suite.AddTestSigner()
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) GetKeeper() keeper.Keeper {
	return suite.App.Erc20Keeper
}
