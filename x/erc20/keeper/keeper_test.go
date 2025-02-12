package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
)

type KeeperTestSuite struct {
	helpers.BaseSuite
}

func TestERC20KeeperTestSuite(t *testing.T) {
	suite.Run(t, &KeeperTestSuite{})
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) GetKeeper() erc20keeper.Keeper {
	return suite.App.Erc20Keeper
}
