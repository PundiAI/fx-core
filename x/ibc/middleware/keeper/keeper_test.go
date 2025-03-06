package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	ibcmiddleware "github.com/pundiai/fx-core/v8/x/ibc/middleware"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	erc20TokenSuite helpers.ERC20TokenSuite
	ibcMiddleware   ibcmiddleware.IBCMiddleware
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.erc20TokenSuite = helpers.NewERC20Suite(suite.Require(), suite.App.EvmKeeper)
	suite.ibcMiddleware = ibcmiddleware.NewIBCMiddleware(suite.App.IBCMiddlewareKeeper, suite.App.IBCKeeper.ChannelKeeper, transfer.NewIBCModule(suite.App.IBCTransferKeeper))
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

var _ types.CrosschainKeeper = (*mockCrosschainKeeper)(nil)

type mockCrosschainKeeper struct{}

func (m mockCrosschainKeeper) IBCCoinToEvm(ctx sdk.Context, holder string, ibcCoin sdk.Coin) error {
	return nil
}

func (m mockCrosschainKeeper) IBCCoinToBaseCoin(ctx context.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) (bool, string, error) {
	panic("implement me")
}

func (m mockCrosschainKeeper) IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error {
	return nil
}

func (m mockCrosschainKeeper) AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64) error {
	return nil
}
