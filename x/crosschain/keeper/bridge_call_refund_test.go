package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_AddRefundRecord() {
	suite.TestMsgSetOracleSetConfirm()
	receiver := fxtypes.AddressToStr(helpers.GenerateAddress().Bytes(), suite.chainName)
	eventNonce := uint64(1)
	tokens := []types.ERC20Token{
		{
			Contract: fxtypes.AddressToStr(helpers.GenerateAddress().Bytes(), suite.chainName),
			Amount:   sdk.NewInt(100),
		},
		{
			Contract: fxtypes.AddressToStr(helpers.GenerateAddress().Bytes(), suite.chainName),
			Amount:   sdk.NewInt(80),
		},
	}
	err := suite.Keeper().AddRefundRecord(suite.ctx, receiver, eventNonce, tokens)
	suite.NoError(err)

	refundRecord, found := suite.Keeper().GetRefundRecord(suite.ctx, eventNonce)
	suite.True(found)
	suite.Equal(refundRecord.EventNonce, eventNonce)
	suite.Equal(refundRecord.Receiver, receiver)
	suite.Equal(refundRecord.Tokens, tokens)
}
