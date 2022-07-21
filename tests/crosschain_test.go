package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/functionx/fx-core/v2/app/helpers"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	crosschaintypes "github.com/functionx/fx-core/v2/x/crosschain/types"
)

type CrosschainTestSuite struct {
	TestSuite
	params            crosschaintypes.Params
	chainName         string
	oraclePrivKey     cryptotypes.PrivKey
	bridgerFxPrivKey  cryptotypes.PrivKey
	bridgerExtPrivKey *ecdsa.PrivateKey
	privKey           cryptotypes.PrivKey
}

func NewCrosschainTestSuite(chainName string) CrosschainTestSuite {
	return CrosschainTestSuite{
		TestSuite:         NewTestSuite(),
		chainName:         chainName,
		oraclePrivKey:     helpers.NewPriKey(),
		bridgerFxPrivKey:  helpers.NewPriKey(),
		bridgerExtPrivKey: helpers.GenerateEthKey(),
		privKey:           helpers.NewEthPrivKey(),
	}
}

func (suite *CrosschainTestSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()

	suite.Send(suite.OracleAddr(), helpers.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.BridgerFxAddr(), helpers.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.params = suite.QueryParams()
}

func (suite *CrosschainTestSuite) OracleAddr() sdk.AccAddress {
	return suite.oraclePrivKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) BridgerExtAddr() string {
	return ethCrypto.PubkeyToAddress(suite.bridgerExtPrivKey.PublicKey).String()
}

func (suite *CrosschainTestSuite) BridgerFxAddr() sdk.AccAddress {
	return suite.bridgerFxPrivKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) AccAddr() sdk.AccAddress {
	return suite.privKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) HexAddr() gethcommon.Address {
	return gethcommon.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *CrosschainTestSuite) CrosschainQuery() crosschaintypes.QueryClient {
	return suite.GRPCClient().CrosschainQuery()
}

func (suite *CrosschainTestSuite) QueryParams() crosschaintypes.Params {
	response, err := suite.CrosschainQuery().Params(suite.ctx,
		&crosschaintypes.QueryParamsRequest{ChainName: suite.chainName})
	suite.NoError(err)
	return response.Params
}

func (suite *CrosschainTestSuite) queryFxLastEventNonce() uint64 {
	lastEventNonce, err := suite.CrosschainQuery().LastEventNonceByAddr(suite.ctx,
		&crosschaintypes.QueryLastEventNonceByAddrRequest{
			ChainName:      suite.chainName,
			BridgerAddress: suite.BridgerFxAddr().String(),
		},
	)
	suite.NoError(err)
	return lastEventNonce.EventNonce + 1
}

func (suite *CrosschainTestSuite) queryObserverExternalBlockHeight() uint64 {
	response, err := suite.CrosschainQuery().LastObservedBlockHeight(suite.ctx,
		&crosschaintypes.QueryLastObservedBlockHeightRequest{
			ChainName: suite.chainName,
		},
	)
	suite.NoError(err)
	return response.ExternalBlockHeight
}

func (suite *CrosschainTestSuite) BondedOracle() {
	response, err := suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerFxAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Error(err, crosschaintypes.ErrNoFoundOracle)
	suite.Nil(response)

	suite.BroadcastTx(suite.oraclePrivKey, &crosschaintypes.MsgBondedOracle{
		OracleAddress:    suite.OracleAddr().String(),
		BridgerAddress:   suite.BridgerFxAddr().String(),
		ExternalAddress:  suite.BridgerExtAddr(),
		ValidatorAddress: suite.ValAddress().String(),
		DelegateAmount:   suite.params.DelegateThreshold,
		ChainName:        suite.chainName,
	})

	response, err = suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerFxAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)
	suite.T().Log("oracle", response.Oracle)
}

func (suite *CrosschainTestSuite) SendUpdateChainOraclesProposal() (proposalId uint64) {
	proposal, err := govtypes.NewMsgSubmitProposal(
		&crosschaintypes.UpdateChainOraclesProposal{
			Title:       fmt.Sprintf("Update %s cross chain oracle", suite.chainName),
			Description: "foo",
			Oracles:     []string{suite.OracleAddr().String()},
			ChainName:   suite.chainName,
		},
		sdk.NewCoins(helpers.NewCoin(sdk.NewInt(10_000).MulRaw(1e18))),
		suite.OracleAddr(),
	)
	suite.NoError(err)
	return suite.BroadcastProposalTx(suite.oraclePrivKey, proposal)
}

func (suite *CrosschainTestSuite) SendOracleSetConfirm() {
	timeoutCtx, cancel := context.WithTimeout(suite.ctx, suite.network.Config.TimeoutCommit)
	defer cancel()
	for {
		time.Sleep(10 * time.Millisecond)
		queryResponse, err := suite.CrosschainQuery().LastPendingOracleSetRequestByAddr(
			timeoutCtx,
			&crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
				BridgerAddress: suite.BridgerFxAddr().String(),
				ChainName:      suite.chainName,
			},
		)
		if err != nil {
			suite.Require().ErrorContains(err, "oracle")
			continue
		}
		for _, valset := range queryResponse.OracleSets {
			checkpoint, err := valset.GetCheckpoint(suite.params.GravityId)
			suite.NoError(err)

			signature, err := crosschaintypes.NewEthereumSignature(checkpoint, suite.bridgerExtPrivKey)
			suite.NoError(err)

			suite.BroadcastTx(suite.bridgerFxPrivKey, &crosschaintypes.MsgOracleSetConfirm{
				Nonce:           valset.Nonce,
				BridgerAddress:  suite.BridgerFxAddr().String(),
				ExternalAddress: suite.BridgerExtAddr(),
				Signature:       hex.EncodeToString(signature),
				ChainName:       suite.chainName,
			})
		}
		if len(queryResponse.OracleSets) > 0 {
			break
		}
	}
}

func (suite *CrosschainTestSuite) AddBridgeTokenClaim(name, symbol string, decimals uint64, token, channelIBC string) string {
	bridgeToken, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.ErrorContains(err, "bridge token")
	suite.Nil(bridgeToken)

	suite.BroadcastTx(suite.bridgerFxPrivKey, &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: suite.BridgerFxAddr().String(),
		ChannelIbc:     hex.EncodeToString([]byte(channelIBC)),
		ChainName:      suite.chainName,
	})

	bridgeToken, err = suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	suite.T().Log("bridge token", bridgeToken)
	return bridgeToken.Denom
}

func (suite *CrosschainTestSuite) SendToFxClaim(token string, amount sdk.Int, targetIbc string) {
	suite.BroadcastTx(suite.bridgerFxPrivKey, &crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Amount:         amount,
		Sender:         suite.HexAddr().Hex(),
		Receiver:       suite.AccAddr().String(),
		TargetIbc:      hex.EncodeToString([]byte(targetIbc)),
		BridgerAddress: suite.BridgerFxAddr().String(),
		ChainName:      suite.chainName,
	})
	bridgeToken, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	balances := suite.QueryBalances(suite.AccAddr())
	suite.True(balances.IsAllGTE(sdk.NewCoins(sdk.NewCoin(bridgeToken.Denom, amount))))
}

func (suite *CrosschainTestSuite) SendToExternal(count int, amount sdk.Coin) uint64 {
	msgList := make([]sdk.Msg, 0, count)
	for i := 0; i < count; i++ {
		msgList = append(msgList, &crosschaintypes.MsgSendToExternal{
			Sender:    suite.AccAddr().String(),
			Dest:      suite.HexAddr().Hex(),
			Amount:    amount.SubAmount(sdk.NewInt(1).ModRaw(1e18)),
			BridgeFee: sdk.NewCoin(amount.Denom, sdk.NewInt(1).ModRaw(1e18)),
			ChainName: suite.chainName,
		})
	}
	txResponse := suite.BroadcastTx(suite.privKey, msgList...)
	for _, eventLog := range txResponse.Logs {
		for _, event := range eventLog.Events {
			if event.Type != crosschaintypes.EventTypeSendToExternal {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != crosschaintypes.AttributeKeyOutgoingTxID {
					continue
				}
				txId, err := strconv.ParseUint(attribute.Value, 10, 64)
				suite.NoError(err)
				return txId
			}
		}
	}
	return 0
}

func (suite *CrosschainTestSuite) SendToExternalAndCancel(token, denom string, amount sdk.Int) {
	coin := sdk.NewCoin(denom, amount)
	suite.SendToFxClaim(token, coin.Amount, "")

	txId := suite.SendToExternal(1, coin)
	suite.Greater(txId, uint64(0))

	suite.SendCancelSendToExternal(txId)
}

func (suite *CrosschainTestSuite) SendCancelSendToExternal(txId uint64) {
	suite.BroadcastTx(suite.privKey, &crosschaintypes.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        suite.AccAddr().String(),
		ChainName:     suite.chainName,
	})
}

func (suite *CrosschainTestSuite) SendBatchRequest(minTxs uint64) {
	msgList := make([]sdk.Msg, 0)
	for {
		batchFeeResponse, err := suite.CrosschainQuery().BatchFees(suite.ctx, &crosschaintypes.QueryBatchFeeRequest{ChainName: suite.chainName})
		suite.NoError(err)
		for _, batchToken := range batchFeeResponse.BatchFees {
			if batchToken.TotalTxs >= minTxs {
				denomResponse, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
					Token:     batchToken.TokenContract,
					ChainName: suite.chainName,
				})
				suite.NoError(err)

				msgList = append(msgList, &crosschaintypes.MsgRequestBatch{
					Sender:     suite.BridgerFxAddr().String(),
					Denom:      denomResponse.Denom,
					MinimumFee: batchToken.TotalFees,
					FeeReceive: suite.HexAddr().String(),
					ChainName:  suite.chainName,
				})
			}
		}
		if len(msgList) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	suite.BroadcastTx(suite.bridgerFxPrivKey, msgList...)
}

func (suite *CrosschainTestSuite) SendConfirmBatch() {
	timeoutCtx, cancel := context.WithTimeout(suite.ctx, suite.network.Config.TimeoutCommit)
	defer cancel()
	for {
		response, err := suite.CrosschainQuery().LastPendingBatchRequestByAddr(
			timeoutCtx,
			&crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{
				BridgerAddress: suite.BridgerFxAddr().String(),
				ChainName:      suite.chainName,
			},
		)
		suite.NoError(err)

		if response.Batch == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		outgoingTxBatch := response.Batch
		checkpoint, err := outgoingTxBatch.GetCheckpoint(suite.params.GravityId)
		suite.NoError(err)

		signatureBytes, err := crosschaintypes.NewEthereumSignature(checkpoint, suite.bridgerExtPrivKey)
		suite.NoError(err)

		err = crosschaintypes.ValidateEthereumSignature(checkpoint, signatureBytes, suite.BridgerExtAddr())
		suite.NoError(err)

		suite.BroadcastTx(suite.bridgerFxPrivKey,
			&crosschaintypes.MsgConfirmBatch{
				Nonce:           outgoingTxBatch.BatchNonce,
				TokenContract:   outgoingTxBatch.TokenContract,
				BridgerAddress:  suite.BridgerFxAddr().String(),
				ExternalAddress: suite.BridgerExtAddr(),
				Signature:       hex.EncodeToString(signatureBytes),
				ChainName:       suite.chainName,
			},
			&crosschaintypes.MsgSendToExternalClaim{
				EventNonce:     suite.queryFxLastEventNonce(),
				BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
				BatchNonce:     outgoingTxBatch.BatchNonce,
				TokenContract:  outgoingTxBatch.TokenContract,
				BridgerAddress: suite.BridgerFxAddr().String(),
				ChainName:      suite.chainName,
			},
		)
		break
	}
}
