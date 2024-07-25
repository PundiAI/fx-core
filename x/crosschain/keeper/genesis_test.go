package keeper_test

import (
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
)

// Tests that batches and transactions are preserved during chain restart
func (suite *KeeperTestSuite) TestBatchAndTxImportExport() {
	bridgeTokens := make([]types.BridgeToken, 10)
	for i := 0; i < len(bridgeTokens); i++ {
		contractAddress := helpers.GenHexAddress().Hex()
		bridgeToken := types.BridgeToken{
			Token: contractAddress,
			Denom: types.NewBridgeDenom(suite.chainName, contractAddress),
		}
		bridgeTokens[i] = bridgeToken
		denom, err := suite.Keeper().SetIbcDenomTrace(suite.ctx, bridgeToken.Token, "")
		suite.Require().NoError(err)
		suite.Require().Equal(denom, bridgeToken.Denom)
		suite.Keeper().AddBridgeToken(suite.ctx, bridgeToken.Token, denom) // nolint:staticcheck

		for _, bridger := range suite.bridgerAddrs {
			voucher := sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(9990))
			err := suite.app.BankKeeper.MintCoins(suite.ctx, suite.chainName, sdk.NewCoins(voucher))
			suite.Require().NoError(err)

			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, suite.chainName, bridger, sdk.NewCoins(voucher))
			suite.Require().NoError(err)
		}
	}

	// CREATE TRANSACTIONS
	// ==================
	numTxs := 1000 // should end up with 1000 txs per contract
	txs := make(types.OutgoingTransferTxs, numTxs)
	fees := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	amounts := []int{51, 52, 53, 54, 55, 56, 57, 58, 59, 60}
	for i := 0; i < numTxs; i++ {
		// Pick fee, amoun sender, receiver, and contract for the ith transaction
		// Sender and contract will always match up (they must since sender i controls the whole balance of the ith token)
		// Receivers should get a balance of many token types since i % len(receivers) is usually different than i % len(contracts)
		fee := fees[i%len(fees)] // fee for this transaction
		amount := amounts[i%len(amounts)]
		sender := suite.bridgerAddrs[i%len(suite.bridgerAddrs)]
		receiver := crypto.PubkeyToAddress(suite.externalPris[i%len(suite.externalPris)].PublicKey).String()
		bridgeToken := bridgeTokens[i%len(bridgeTokens)]

		amountToken := sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(int64(amount)))
		feeToken := sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(int64(fee)))

		// add transaction to the pool
		nextTxID, err := suite.Keeper().AddToOutgoingPool(suite.ctx, sender, receiver, amountToken, feeToken)
		suite.Require().NoError(err)

		txs[i] = &types.OutgoingTransferTx{
			Id:          nextTxID,
			Sender:      sender.String(),
			DestAddress: receiver,
			Token:       types.NewERC20Token(amountToken.Amount, bridgeToken.Token),
			Fee:         types.NewERC20Token(feeToken.Amount, bridgeToken.Token),
		}
	}

	suite.Keeper().SetLastObservedBlockHeight(suite.ctx, 10, 10)

	// CREATE BATCHES
	// ==================
	// Want to create batches for half of the transactions for each contract
	// with 100 tx in each batch, 1000 txs per contrac we want 5 batches per contract to batch 500 txs per contract
	for i, bridgeToken := range bridgeTokens {
		suite.ctx = suite.ctx.WithBlockHeight(int64(50 + i))
		batch, err := suite.Keeper().BuildOutgoingTxBatch(suite.ctx, bridgeToken.Token, bridgeToken.Token, 100, sdkmath.NewInt(1), sdkmath.NewInt(1))
		suite.Require().NoError(err)
		suite.Require().EqualValues(100, len(batch.Transactions))
		suite.Require().EqualValues(50+i, batch.Block)
		if suite.chainName == ethtypes.ModuleName {
			suite.Require().True(batch.BatchTimeout > 2800)
		} else {
			suite.Require().True(batch.BatchTimeout > 14000)
		}
		suite.Require().EqualValues(1+i, batch.BatchNonce)
		suite.Require().Equal(bridgeToken.Token, batch.TokenContract)
		suite.Require().Equal(bridgeToken.Token, batch.FeeReceive)
	}

	// export
	suite.checkAllTransactionsExist(txs)
	genesisState := keeper.ExportGenesis(suite.ctx, suite.Keeper())

	// clear data
	storeKey := suite.app.GetKey(suite.chainName)
	store := suite.ctx.KVStore(storeKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	unbatched := suite.Keeper().GetUnbatchedTransactions(suite.ctx)
	suite.Require().Empty(unbatched)
	batches := suite.Keeper().GetOutgoingTxBatches(suite.ctx)
	suite.Require().Empty(batches)

	// import
	keeper.InitGenesis(suite.ctx, suite.Keeper(), genesisState)
	suite.checkAllTransactionsExist(txs)
}

// Requires that all transactions in txs exist in keeper
func (suite *KeeperTestSuite) checkAllTransactionsExist(txs types.OutgoingTransferTxs) {
	unbatched := suite.Keeper().GetUnbatchedTransactions(suite.ctx)
	batches := suite.Keeper().GetOutgoingTxBatches(suite.ctx)
	// Collect all txs into an array
	var gotTxs types.OutgoingTransferTxs
	gotTxs = append(gotTxs, unbatched...)
	for _, batch := range batches {
		gotTxs = append(gotTxs, batch.Transactions...)
	}
	suite.Require().Equal(len(txs), len(gotTxs))
	// Sort both arrays for simple searching
	sort.Slice(gotTxs, func(i, j int) bool {
		return gotTxs[i].Id < gotTxs[j].Id
	})
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Id < txs[j].Id
	})
	// Actually check that the txs all exist, iterate on txs in case some got lost in the import/export step
	for i, exp := range txs {
		suite.Require().Equal(exp.Id, gotTxs[i].Id)
		suite.Require().Equal(exp.Fee.String(), gotTxs[i].Fee.String())
		suite.Require().Equal(exp.Token.String(), gotTxs[i].Token.String())
		suite.Require().Equal(exp.DestAddress, gotTxs[i].DestAddress)
		suite.Require().Equal(exp.Sender, gotTxs[i].Sender)
	}
}
