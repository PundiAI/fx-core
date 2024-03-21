package keeper_test

import (
	"testing"

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

func (suite *KeeperTestSuite) TestKeeper_IterateRefundRecordByNonce() {
	testCases := []struct {
		name             string
		eventNonces      []uint64
		startNonce       uint64
		expectEventNonce []uint64
	}{
		{
			name:             "only 1 - start 0",
			eventNonces:      []uint64{1},
			startNonce:       uint64(0),
			expectEventNonce: []uint64{1},
		},
		{
			name:             "only 1 - start 1",
			eventNonces:      []uint64{1},
			startNonce:       uint64(1),
			expectEventNonce: []uint64{1},
		},
		{
			name:             "out-of-order",
			eventNonces:      []uint64{6, 2, 5},
			startNonce:       uint64(1),
			expectEventNonce: []uint64{2, 5, 6},
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			suite.SetupTest()
			for _, nonce := range testCase.eventNonces {
				suite.Keeper().SetRefundRecord(suite.ctx, &types.RefundRecord{EventNonce: nonce})
			}
			actualEventNonces := make([]uint64, 0, len(testCase.expectEventNonce))
			suite.Keeper().IterateRefundRecordByNonce(suite.ctx, testCase.startNonce, func(record *types.RefundRecord) bool {
				actualEventNonces = append(actualEventNonces, record.EventNonce)
				return false
			})

			suite.EqualValues(testCase.expectEventNonce, actualEventNonces)
		})
	}
}
