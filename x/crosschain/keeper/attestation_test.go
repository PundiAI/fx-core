package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func (suite *KeeperTestSuite) TestCleanupTimeoutBatches() {
	batch1 := types.OutgoingTxBatch{BatchNonce: 1, BatchTimeout: 100, TokenContract: helpers.GenExternalAddr(suite.chainName), Block: 100, Transactions: []*types.OutgoingTransferTx{{Id: 1}}}
	batch2 := types.OutgoingTxBatch{BatchNonce: 2, BatchTimeout: 2000, TokenContract: helpers.GenExternalAddr(suite.chainName), Block: 2000, Transactions: []*types.OutgoingTransferTx{{Id: 2}}}
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, &batch1))
	suite.Require().NoError(suite.Keeper().StoreBatch(suite.Ctx, &batch2))
	suite.NotEmpty(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, batch1.TokenContract, batch1.BatchNonce))
	suite.NotEmpty(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, batch2.TokenContract, batch2.BatchNonce))

	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1000)
	suite.Keeper().SetLastObservedBlockHeight(suite.Ctx, 1000, 1000)
	suite.SetAutoIncrementID(types.KeyLastOutgoingBatchID, 3)
	suite.Require().NoError(suite.Keeper().CleanupTimedOutBatches(suite.Ctx))

	suite.Empty(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, batch1.TokenContract, batch1.BatchNonce))
	suite.NotEmpty(suite.Keeper().GetOutgoingTxBatch(suite.Ctx, batch2.TokenContract, batch2.BatchNonce))

	batch3 := suite.Keeper().GetOutgoingTxBatch(suite.Ctx, batch1.TokenContract, 3)
	suite.NotEmpty(batch3)
	suite.Equal(uint64(3), batch3.BatchNonce)
	suite.Equal(uint64(suite.Ctx.BlockHeight()), batch3.Block)
	suite.Equal(batch1.Transactions[0].Id, batch3.Transactions[0].Id)
}

func (suite *KeeperTestSuite) TestCleanupTimeoutBridgeCall() {
	originBridgeToken := suite.AddBridgeToken(suite.chainName, fxtypes.DefaultSymbol, false)
	testBridgeToken := suite.AddBridgeToken(suite.chainName, helpers.NewRandSymbol(), true)

	outCall1 := types.OutgoingBridgeCall{
		Refund: helpers.GenExternalAddr(suite.chainName),
		Tokens: types.ERC20Tokens{
			{Contract: originBridgeToken.Contract, Amount: helpers.NewRandAmount()},
			{Contract: testBridgeToken.Contract, Amount: helpers.NewRandAmount()},
		},
		Nonce:   1,
		Timeout: 100,
	}
	outCall2 := types.OutgoingBridgeCall{
		Refund: helpers.GenExternalAddr(suite.chainName),
		Tokens: types.ERC20Tokens{
			{Contract: originBridgeToken.Contract, Amount: helpers.NewRandAmount()},
			{Contract: testBridgeToken.Contract, Amount: helpers.NewRandAmount()},
		},
		Nonce:   2,
		Timeout: 100,
	}
	outCall3 := types.OutgoingBridgeCall{
		Refund: helpers.GenExternalAddr(suite.chainName),
		Tokens: types.ERC20Tokens{
			{Contract: originBridgeToken.Contract, Amount: helpers.NewRandAmount()},
			{Contract: testBridgeToken.Contract, Amount: helpers.NewRandAmount()},
		},
		Nonce:   3,
		Timeout: 2000,
	}
	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &outCall1)
	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &outCall2)
	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &outCall3)
	suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, types.NewOriginTokenKey(suite.chainName, outCall1.Nonce), outCall1.Tokens[0].Amount))
	suite.Keeper().SetOutgoingBridgeCallQuoteInfo(suite.Ctx, outCall2.Nonce, types.QuoteInfo{Id: 1})
	suite.MintTokenToModule(ethtypes.ModuleName, sdk.NewCoin(fxtypes.DefaultDenom, outCall1.Tokens[0].Amount))

	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1000)
	suite.Keeper().SetLastObservedBlockHeight(suite.Ctx, 1000, 1000)
	suite.SetAutoIncrementID(types.KeyLastBridgeCallID, 4)
	suite.Require().NoError(suite.Keeper().CleanupTimeOutBridgeCall(suite.Ctx))

	_, found := suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, outCall1.Nonce)
	suite.False(found)
	addr1 := fxtypes.ExternalAddrToHexAddr(suite.chainName, outCall1.Refund)
	suite.AssertBalance(addr1.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, outCall1.Tokens[0].Amount))
	testErc20 := suite.GetERC20Token(testBridgeToken.Denom)
	suite.Equal(outCall1.Tokens[1].Amount.String(), suite.erc20TokenSuite.WithContract(testErc20.GetERC20Contract()).BalanceOf(suite.Ctx, addr1).String())

	_, found = suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, outCall2.Nonce)
	suite.False(found)
	outCall4, found := suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, 4)
	suite.True(found)
	suite.Equal(outCall2.Refund, outCall4.Refund)
	suite.Equal(outCall2.Tokens, outCall4.Tokens)
	suite.Equal(uint64(suite.Ctx.BlockHeight()), outCall4.BlockHeight)

	_, found = suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, outCall3.Nonce)
	suite.True(found)
}
