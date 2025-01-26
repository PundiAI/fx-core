package keeper_test

import (
	"testing"

	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	ibcmiddleware "github.com/pundiai/fx-core/v8/x/ibc/middleware"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	ibcMiddleware porttypes.Middleware
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()

	suite.ibcMiddleware = ibcmiddleware.NewIBCMiddleware(suite.App.IBCMiddlewareKeeper, suite.App.IBCKeeper.ChannelKeeper, transfer.NewIBCModule(suite.App.IBCTransferKeeper))
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}
