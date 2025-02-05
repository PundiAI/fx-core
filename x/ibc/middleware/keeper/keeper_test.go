package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	ibcmiddleware "github.com/pundiai/fx-core/v8/x/ibc/middleware"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
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
	suite.App.IBCMiddlewareKeeper = suite.App.IBCMiddlewareKeeper.SetCrosschainKeeper(mockCrosschainKeeper{})
	suite.ibcMiddleware = ibcmiddleware.NewIBCMiddleware(suite.App.IBCMiddlewareKeeper, suite.App.IBCKeeper.ChannelKeeper, transfer.NewIBCModule(suite.App.IBCTransferKeeper))
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

var _ types.CrosschainKeeper = (*mockCrosschainKeeper)(nil)

type mockCrosschainKeeper struct{}

func (m mockCrosschainKeeper) IBCCoinToEvm(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) error {
	return nil
}

func (m mockCrosschainKeeper) IBCCoinToBaseCoin(ctx context.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) (string, error) {
	panic("implement me")
}

func (m mockCrosschainKeeper) IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error {
	panic("implement me")
}

func (m mockCrosschainKeeper) AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64) error {
	panic("implement me")
}
