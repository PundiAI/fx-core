package tests

import (
	"math/big"
	"slices"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v7/client"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	precompilescrosschain "github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
)

type PrecompileTestSuite struct {
	EvmTestSuite
	privKey cryptotypes.PrivKey
}

func NewPrecompileTestSuite(ts *TestSuite) PrecompileTestSuite {
	return PrecompileTestSuite{
		EvmTestSuite: NewEvmTestSuite(ts),
		privKey:      helpers.NewEthPrivKey(),
	}
}

func (suite *PrecompileTestSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *PrecompileTestSuite) HexAddress() common.Address {
	return common.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *PrecompileTestSuite) TransferCrossChain(token common.Address, recipient string, amount, fee *big.Int, target string) *ethtypes.Transaction {
	privateKey := suite.privKey
	beforeBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	pack, err := fxtypes.GetFIP20().ABI.Pack("transferCrossChain", recipient, amount, fee, fxtypes.MustStrToByte32(target))
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err, target)
	suite.SendTransaction(ethTx)
	afterBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	suite.Require().True(new(big.Int).Sub(beforeBalanceOf, afterBalanceOf).Cmp(new(big.Int).Add(amount, fee)) == 0)
	return ethTx
}

func (suite *PrecompileTestSuite) TransferCrossChainAndCheckPendingTx(token common.Address, recipient string, amount, fee *big.Int, chainName string) *ethtypes.Transaction {
	pendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chainName,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	ethTx := suite.TransferCrossChain(token, recipient, amount, fee, chainName)
	newPendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chainName,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	suite.Equal(1, len(newPendingTxs.UnbatchedTransfers)-len(pendingTxs.UnbatchedTransfers))
	tx := newPendingTxs.UnbatchedTransfers[0]
	suite.Equal(tx.Token.Amount, sdkmath.NewIntFromBigInt(amount))
	suite.Equal(tx.Fee.Amount, sdkmath.NewIntFromBigInt(fee))
	suite.Equal(tx.Sender, suite.AccAddress().String())
	return ethTx
}

func (suite *PrecompileTestSuite) CrossChainAndResponse(token common.Address, recipient string, amount, fee *big.Int, target string) *ethtypes.Transaction {
	privateKey := suite.privKey
	crossChainContract := precompilescrosschain.GetAddress()
	suite.ApproveERC20(privateKey, token, crossChainContract, big.NewInt(0).Add(amount, fee))

	beforeBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	pack, err := precompilescrosschain.GetABI().Pack("crossChain", token, recipient, amount, fee, fxtypes.MustStrToByte32(target), "")
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &crossChainContract, nil, pack)
	suite.Require().NoError(err, target)
	suite.SendTransaction(ethTx)
	afterBalanceOf := suite.BalanceOf(token, common.BytesToAddress(privateKey.PubKey().Address().Bytes()))
	suite.Require().True(new(big.Int).Sub(beforeBalanceOf, afterBalanceOf).Cmp(new(big.Int).Add(amount, fee)) == 0)
	return ethTx
}

func (suite *PrecompileTestSuite) CrossChainAndCheckPendingTx(token common.Address, recipient string, amount, fee *big.Int, chainName string) (*ethtypes.Transaction, uint64) {
	pendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chainName,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	ethTx := suite.CrossChainAndResponse(token, recipient, amount, fee, chainName)
	newPendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chainName,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	suite.Equal(1, len(newPendingTxs.UnbatchedTransfers)-len(pendingTxs.UnbatchedTransfers))

	var currentTx *crosschaintypes.OutgoingTransferTx
	for _, newTx := range newPendingTxs.UnbatchedTransfers {
		found := false
		for _, tx := range pendingTxs.UnbatchedTransfers {
			if tx.Id == newTx.Id {
				found = true
				break
			}
		}
		if !found {
			currentTx = newTx
			break
		}
	}

	suite.NotEmpty(currentTx)
	suite.Equal(currentTx.Token.Amount, sdkmath.NewIntFromBigInt(amount))
	suite.Equal(currentTx.Fee.Amount, sdkmath.NewIntFromBigInt(fee))
	suite.Equal(currentTx.Sender, suite.AccAddress().String())
	return ethTx, currentTx.Id
}

func (suite *PrecompileTestSuite) CrossChain(token common.Address, recipient string, amount, fee *big.Int, chainName string) uint64 {
	_, txId := suite.CrossChainAndCheckPendingTx(token, recipient, amount, fee, chainName)
	return txId
}

func (suite *PrecompileTestSuite) CancelSendToExternal(chain string, txId uint64) *ethtypes.Transaction {
	privateKey := suite.privKey
	crossChainContract := precompilescrosschain.GetAddress()
	pack, err := precompilescrosschain.GetABI().Pack(precompilescrosschain.CancelSendToExternalMethodName, chain, big.NewInt(int64(txId)))
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &crossChainContract, nil, pack)
	suite.Require().NoError(err, chain)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *PrecompileTestSuite) CancelSendToExternalAndCheckPendingTx(chain string, txId uint64) *ethtypes.Transaction {
	txIdEqual := func(tx *crosschaintypes.OutgoingTransferTx) bool { return tx.Id == txId }
	pendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chain,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	suite.True(len(pendingTxs.UnbatchedTransfers) > 0)
	suite.True(slices.ContainsFunc(pendingTxs.UnbatchedTransfers, txIdEqual))
	ethTx := suite.CancelSendToExternal(chain, txId)
	newPendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chain,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	suite.Equal(len(pendingTxs.UnbatchedTransfers)-1, len(newPendingTxs.UnbatchedTransfers))
	suite.False(slices.ContainsFunc(newPendingTxs.UnbatchedTransfers, txIdEqual))
	return ethTx
}

func (suite *PrecompileTestSuite) IncreaseBridgeFee(chain string, txId uint64, token common.Address, fee *big.Int) *ethtypes.Transaction {
	privateKey := suite.privKey
	crossChainContract := precompilescrosschain.GetAddress()
	suite.ApproveERC20(privateKey, token, crossChainContract, fee)
	pack, err := precompilescrosschain.GetABI().Pack(precompilescrosschain.IncreaseBridgeFeeMethodName, chain, big.NewInt(int64(txId)), token, fee)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &crossChainContract, nil, pack)
	suite.Require().NoError(err, chain)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *PrecompileTestSuite) IncreaseBridgeFeeCheckPendingTx(chain string, txId uint64, token common.Address, fee *big.Int) *ethtypes.Transaction {
	txIdEqual := func(tx *crosschaintypes.OutgoingTransferTx) bool { return tx.Id == txId }
	pendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chain,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	tx := pendingTxs.UnbatchedTransfers[slices.IndexFunc(pendingTxs.UnbatchedTransfers, txIdEqual)]
	suite.Equal(txId, tx.Id)

	ethTx := suite.IncreaseBridgeFee(chain, txId, token, fee)
	newPendingTxs, err := suite.GRPCClient().CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     chain,
			SenderAddress: suite.AccAddress().String(),
		})
	suite.NoError(err)
	newTx := newPendingTxs.UnbatchedTransfers[slices.IndexFunc(newPendingTxs.UnbatchedTransfers, txIdEqual)]
	suite.Equal(tx.Id, newTx.Id)
	suite.Equal(tx.Fee.Amount.Add(sdkmath.NewIntFromBigInt(fee)), newTx.Fee.Amount)
	return ethTx
}
