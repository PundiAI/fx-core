package tests

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app/helpers"
	types2 "github.com/functionx/fx-core/x/ibc/applications/transfer/types"

	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/stretchr/testify/require"

	trontypes "github.com/functionx/fx-core/x/tron/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/x/crosschain/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

type TronCrosschainTestSuite struct {
	CrosschainTestSuite
}

func TestTronCrosschainTestSuite(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	const purseTokenContract = "TUpMhErZL2fhh4sVNULAbNKLokS4GjC1F4"
	const chainName = "tron"
	const purseTokenChannelIBC = "transfer/channel-0"
	purseDenom := types2.DenomTrace{
		Path:      purseTokenChannelIBC,
		BaseDenom: fmt.Sprintf("%s%s", chainName, purseTokenContract),
	}.IBCDenom()
	suite.Run(t, &TronCrosschainTestSuite{
		CrosschainTestSuite: CrosschainTestSuite{
			TestSuite: NewTestSuite(),
			BridgeToken: crosschaintypes.BridgeToken{
				Token:      purseTokenContract,
				Denom:      fmt.Sprintf("%s%s", chainName, purseTokenContract),
				ChannelIbc: "px/transfer/channel-0",
			},
			ibcDenom:   purseDenom,
			ethPrivKey: helpers.GenerateEthKey(),
			chainName:  chainName,
		},
	})
}

func (suite *TronCrosschainTestSuite) TestTronCrosschain() {
	go suite.signPendingValsetRequest()

	suite.setOrchestratorAddress()

	suite.addBridgeTokenClaim()

	suite.externalToFx()

	suite.externalToFxAndIbcTransfer()

	suite.showAllBalance(suite.AdminAddress())

	suite.fxToExternal(5)

	suite.batchRequest()

	suite.confirmBatch()

	suite.sendToExternalAndCancel()
}

func (suite *TronCrosschainTestSuite) tronAddress() tronAddress.Address {
	return tronAddress.PubkeyToAddress(suite.ethPrivKey.PublicKey)
}

func (suite *TronCrosschainTestSuite) sendToExternalAndCancel() {
	suite.T().Helper()
	suite.T().Logf("\n####################      FX to External      ####################\n")
	sendToExternalAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000))
	sendToExternalFee := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(2000))

	suite.BroadcastTx(&types.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Amount:         sendToExternalAmount.Add(sendToExternalFee),
		Sender:         suite.tronAddress().String(),
		Receiver:       suite.AdminAddress().String(),
		TargetIbc:      "",
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	})

	fxAddress := suite.AdminAddress()

	sendToExternalBeforeBalance := suite.getBalanceByAddress(fxAddress, suite.BridgeToken.Denom)
	suite.T().Logf("send-to-External before balance:[%v    %v]", sendToExternalBeforeBalance.Amount.String(), sendToExternalBeforeBalance.Denom)

	sendToExternalHash := suite.BroadcastTx(&types.MsgSendToExternal{
		Sender:    suite.AdminAddress().String(),
		Dest:      suite.tronAddress().String(),
		Amount:    sdk.NewCoin(suite.BridgeToken.Denom, sendToExternalAmount),
		BridgeFee: sdk.NewCoin(suite.BridgeToken.Denom, sendToExternalFee),
		ChainName: suite.chainName,
	})

	sendToExternalAfterBalance := suite.getBalanceByAddress(fxAddress, suite.BridgeToken.Denom)
	suite.T().Logf("send-to-External after balance:[%v    %v]", sendToExternalAfterBalance.Amount.String(), sendToExternalAfterBalance.Denom)
	differentAmount := sendToExternalBeforeBalance.Amount.Sub(sendToExternalAfterBalance.Amount)
	require.True(suite.T(), sendToExternalAmount.Add(sendToExternalFee).Equal(differentAmount), "beforeBalance - afterBalance != sendToExternalFeeAmount+sendToExternalFee",
		sendToExternalBeforeBalance.Amount.String(),
		sendToExternalAfterBalance.Amount.String(),
		sendToExternalAmount.Add(sendToExternalFee).String())

	time.Sleep(3 * time.Second)

	txResponse, err := suite.grpcClient.ServiceClient().GetTx(suite.ctx, &tx.GetTxRequest{Hash: sendToExternalHash})
	require.NoError(suite.T(), err)
	txId, found, err := suite.getSentToExternalTxIdByEvents(txResponse.TxResponse.Logs)
	require.NoError(suite.T(), err)
	require.True(suite.T(), found)
	require.Greater(suite.T(), txId, uint64(0))
	suite.T().Logf("send-to-External txId:[%d]", txId)

	_ = suite.BroadcastTx(&types.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        suite.AdminAddress().String(),
		ChainName:     suite.chainName,
	})

	cancelSendToExternalAfterBalance := suite.getBalanceByAddress(fxAddress, suite.BridgeToken.Denom)
	suite.T().Logf("cancel-send-to-External after balance:[%v    %v]", cancelSendToExternalAfterBalance.Amount.String(), cancelSendToExternalAfterBalance.Denom)
	require.True(suite.T(), sendToExternalBeforeBalance.Equal(cancelSendToExternalAfterBalance), sendToExternalBeforeBalance.String(), cancelSendToExternalAfterBalance.String())
}

func (suite *TronCrosschainTestSuite) getBalanceByAddress(accAddr sdk.AccAddress, denom string) *sdk.Coin {
	balanceResp, err := suite.grpcClient.BankQuery().Balance(suite.ctx, banktypes.NewQueryBalanceRequest(accAddr, denom))
	require.NoError(suite.T(), err)
	return balanceResp.Balance
}

func (suite *TronCrosschainTestSuite) getSentToExternalTxIdByEvents(logs sdk.ABCIMessageLogs) (uint64, bool, error) {
	for _, eventLog := range logs {
		for _, event := range eventLog.Events {
			if event.Type != types.EventTypeSendToExternal {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != types.AttributeKeyOutgoingTxID {
					continue
				}
				result, err := strconv.ParseUint(attribute.Value, 10, 64)
				if err != nil {
					return 0, false, err
				}
				return result, true, nil
			}
		}
	}
	return 0, false, nil
}

func (suite *TronCrosschainTestSuite) addBridgeTokenClaim() {
	suite.T().Helper()
	suite.T().Logf("\n####################      Add bridge token claim      ####################\n")
	bridgeToken, err := crosschaintypes.NewQueryClient(suite.grpcClient).TokenToDenom(suite.ctx, &types.QueryTokenToDenomRequest{ChainName: suite.chainName, Token: suite.BridgeToken.Token})

	if err != nil && !strings.Contains(err.Error(), "bridge token is not exist") {
		suite.T().Fatal(err)
	}
	if err == nil && bridgeToken.Denom == suite.BridgeToken.Denom {
		suite.T().Logf("bridge token already exists!tokenContract:[%v], denom:[%v], channelIbc:[%v]", suite.BridgeToken.Token, bridgeToken.Denom, bridgeToken.ChannelIbc)
		return
	}
	fxOriginatedTokenClaimMsg := &types.MsgBridgeTokenClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Name:           "USDT",
		Symbol:         "USDT",
		Decimals:       18,
		BridgerAddress: suite.AdminAddress().String(),
		ChannelIbc:     "",
		ChainName:      suite.chainName,
	}
	suite.BroadcastTx(fxOriginatedTokenClaimMsg)
	suite.T().Logf("\n")
}

func (suite *TronCrosschainTestSuite) signPendingValsetRequest() {
	suite.T().Helper()
	defer func() {
		suite.T().Logf("sign pending valset request defer ....\n")
		if err := recover(); err != nil {
			suite.T().Fatal(err)
		}
	}()
	gravityId := suite.queryGravityId()
	requestParams := &types.QueryLastPendingOracleSetRequestByAddrRequest{
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		queryResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).LastPendingOracleSetRequestByAddr(suite.ctx, requestParams)
		if err != nil {
			suite.T().Logf("query last pending valset request is err!params:%+v, errors:%v\n", requestParams, err)
			continue
		}
		valsets := queryResponse.OracleSets
		if len(valsets) <= 0 {
			continue
		}
		for _, valset := range valsets {
			checkpoint, err := trontypes.GetCheckpointOracleSet(valset, gravityId)
			require.NoError(suite.T(), err)
			suite.T().Logf("need confirm valset: nonce:%v ExternalAddress:%v\n", valset.Nonce, suite.tronAddress().Hex())
			signature, err := trontypes.NewTronSignature(checkpoint, suite.ethPrivKey)
			if err != nil {
				suite.T().Log(err)
				continue
			}
			suite.BroadcastTx(
				&types.MsgOracleSetConfirm{
					Nonce:           valset.Nonce,
					BridgerAddress:  suite.AdminAddress().String(),
					ExternalAddress: suite.tronAddress().String(),
					Signature:       hex.EncodeToString(signature),
					ChainName:       suite.chainName,
				},
			)
		}
	}
}

func (suite *TronCrosschainTestSuite) queryGravityId() string {
	suite.T().Helper()
	chainParamsResp, err := crosschaintypes.NewQueryClient(suite.grpcClient).Params(suite.ctx, &types.QueryParamsRequest{ChainName: suite.chainName})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.T().Logf("chain params result:%+v\n", chainParamsResp.Params)
	return chainParamsResp.Params.GravityId
}

func (suite *TronCrosschainTestSuite) confirmBatch() {
	suite.T().Helper()

	gravityId := suite.queryGravityId()
	orchestrator := suite.AdminAddress()
	for {
		lastPendingBatchRequestResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).LastPendingBatchRequestByAddr(suite.ctx,
			&types.QueryLastPendingBatchRequestByAddrRequest{BridgerAddress: orchestrator.String(), ChainName: suite.chainName})
		if err != nil {
			suite.T().Fatal(err)
		}
		outgoingTxBatch := lastPendingBatchRequestResponse.Batch
		if outgoingTxBatch == nil {
			break
		}
		checkpoint, err := trontypes.GetCheckpointConfirmBatch(outgoingTxBatch, gravityId)
		if err != nil {
			suite.T().Fatal(err)
		}
		signatureBytes, err := trontypes.NewTronSignature(checkpoint, suite.ethPrivKey)
		if err != nil {
			suite.T().Fatal(err)
		}

		err = trontypes.ValidateTronSignature(checkpoint, signatureBytes, suite.tronAddress().String())
		if err != nil {
			suite.T().Fatal(err)
		}
		suite.BroadcastTx(
			&types.MsgConfirmBatch{
				Nonce:           outgoingTxBatch.BatchNonce,
				TokenContract:   outgoingTxBatch.TokenContract,
				BridgerAddress:  suite.AdminAddress().String(),
				ExternalAddress: suite.tronAddress().String(),
				Signature:       hex.EncodeToString(signatureBytes),
				ChainName:       suite.chainName,
			},
		)
		suite.T().Logf("\n")
		time.Sleep(2 * time.Second)

		suite.BroadcastTx(
			&types.MsgSendToExternalClaim{
				EventNonce:     suite.queryFxLastEventNonce(),
				BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
				BatchNonce:     outgoingTxBatch.BatchNonce,
				TokenContract:  outgoingTxBatch.TokenContract,
				BridgerAddress: suite.AdminAddress().String(),
				ChainName:      suite.chainName,
			},
		)
	}
}

func (suite *TronCrosschainTestSuite) batchRequest() {
	suite.T().Helper()

	batchFeeResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).BatchFees(suite.ctx, &types.QueryBatchFeeRequest{ChainName: suite.chainName})
	if err != nil {
		suite.T().Fatal(err)
	}
	orchestrator := suite.AdminAddress()
	feeReceive := suite.tronAddress().String()
	msgList := make([]sdk.Msg, 0, len(batchFeeResponse.BatchFees))
	for _, batchToken := range batchFeeResponse.BatchFees {
		if batchToken.TotalTxs >= 5 {
			denomResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).TokenToDenom(suite.ctx, &types.QueryTokenToDenomRequest{
				Token:     batchToken.TokenContract,
				ChainName: suite.chainName,
			})
			if err != nil {
				suite.T().Fatal(err)
			}
			if strings.HasPrefix(denomResponse.Denom, batchToken.TokenContract) {
				suite.T().Logf("warn!!! not found token contract, tokenContract:[%v], erc20ToDenom response:[%v]\n", batchToken.TokenContract, denomResponse.Denom)
				continue
			}

			msgList = append(msgList, &types.MsgRequestBatch{
				Sender:     orchestrator.String(),
				Denom:      denomResponse.Denom,
				MinimumFee: batchToken.TotalFees,
				FeeReceive: feeReceive,
				ChainName:  suite.chainName,
			})
		}
	}
	if len(msgList) <= 0 {
		return
	}
	suite.BroadcastTx(msgList...)
	suite.T().Logf("\n")
}

func (suite *TronCrosschainTestSuite) fxToExternal(count int) {
	suite.T().Helper()
	suite.T().Logf("\n####################      FX to External      ####################\n")
	sendToFxBeforeBalance, err := suite.grpcClient.BankQuery().Balance(suite.ctx, &banktypes.QueryBalanceRequest{
		Address: suite.AdminAddress().String(),
		Denom:   suite.BridgeToken.Denom,
	})
	require.NoError(suite.T(), err)
	sendToExternalAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(1900))
	sendToExternalFee := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(100))
	totalSendToExternalAmount := sdk.ZeroInt()
	msgList := make([]sdk.Msg, 0, count)
	for i := 0; i < count; i++ {
		msgList = append(msgList, &types.MsgSendToExternal{
			Sender:    suite.AdminAddress().String(),
			Dest:      suite.tronAddress().String(),
			Amount:    sdk.NewCoin(suite.BridgeToken.Denom, sendToExternalAmount),
			BridgeFee: sdk.NewCoin(suite.BridgeToken.Denom, sendToExternalFee),
			ChainName: suite.chainName,
		})
		totalSendToExternalAmount = totalSendToExternalAmount.Add(sendToExternalAmount).Add(sendToExternalFee)
	}
	suite.BroadcastTx(msgList...)
	sendToFxBeforeAfter, err := suite.grpcClient.BankQuery().Balance(suite.ctx, &banktypes.QueryBalanceRequest{
		Address: suite.AdminAddress().String(),
		Denom:   suite.BridgeToken.Denom,
	})
	require.NoError(suite.T(), err)
	differentAmount := sendToFxBeforeBalance.Balance.Amount.Sub(sendToFxBeforeAfter.Balance.Amount)
	require.True(suite.T(), totalSendToExternalAmount.Equal(differentAmount), "beforeBalance - afterBalance !=  totalSendToExternalAmount",
		sendToFxBeforeBalance.Balance.Amount.String(),
		sendToFxBeforeAfter.Balance.Amount.String(),
		totalSendToExternalAmount.String(),
	)
	suite.T().Logf("\n")
}

func (suite *TronCrosschainTestSuite) showAllBalance(address sdk.AccAddress) {
	suite.T().Helper()
	suite.T().Logf("\n####################      Query AdminAddress Balance      ####################\n")
	queryAllBalancesResponse, err := suite.grpcClient.BankQuery().AllBalances(suite.ctx, banktypes.NewQueryAllBalancesRequest(address, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: true,
	}))
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.T().Logf("address: [%v] all balance\n", address.String())
	for _, balance := range queryAllBalancesResponse.Balances {
		suite.T().Logf("denom:%v, amount:%v\n", balance.Denom, balance.Amount.String())
	}
	suite.T().Logf("\n")
}

func (suite *TronCrosschainTestSuite) externalToFx() {
	suite.T().Helper()
	suite.T().Logf("\n####################      External to FX      ####################\n")
	sendToFxBeforeBalance, err := suite.grpcClient.BankQuery().Balance(suite.ctx, &banktypes.QueryBalanceRequest{
		Address: suite.AdminAddress().String(),
		Denom:   suite.BridgeToken.Denom,
	})
	require.NoError(suite.T(), err)
	sendToFxAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000))
	suite.BroadcastTx(&types.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Amount:         sendToFxAmount,
		Sender:         suite.tronAddress().String(),
		Receiver:       suite.AdminAddress().String(),
		TargetIbc:      "",
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	})
	sendToFxBeforeAfter, err := suite.grpcClient.BankQuery().Balance(suite.ctx, &banktypes.QueryBalanceRequest{
		Address: suite.AdminAddress().String(),
		Denom:   suite.BridgeToken.Denom,
	})
	require.NoError(suite.T(), err)
	differentAmount := sendToFxBeforeAfter.Balance.Amount.Sub(sendToFxBeforeBalance.Balance.Amount)
	require.True(suite.T(), sendToFxAmount.Equal(differentAmount), "beforeBalance + sendToFxAmount != afterBalance",
		sendToFxBeforeBalance.Balance.Amount.String(),
		sendToFxAmount.String(),
		sendToFxBeforeAfter.Balance.Amount.String())
	suite.T().Logf("\n")
}

func (suite *TronCrosschainTestSuite) externalToFxAndIbcTransfer() {
	suite.T().Helper()
	suite.T().Logf("\n####################      External to FX to Pundix      ####################\n")

	sendToFxBeforeBalance, err := suite.grpcClient.BankQuery().Balance(suite.ctx, &banktypes.QueryBalanceRequest{
		Address: suite.AdminAddress().String(),
		Denom:   suite.BridgeToken.Denom,
	})
	require.NoError(suite.T(), err)
	sendToFxAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000))
	suite.BroadcastTx(&types.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Amount:         sendToFxAmount,
		Sender:         suite.tronAddress().String(),
		Receiver:       suite.AdminAddress().String(),
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	})
	sendToFxBeforeAfter, err := suite.grpcClient.BankQuery().Balance(suite.ctx, &banktypes.QueryBalanceRequest{
		Address: suite.AdminAddress().String(),
		Denom:   suite.BridgeToken.Denom,
	})
	require.NoError(suite.T(), err)
	//differentAmount := sendToFxBeforeAfter.Balance.Amount.Sub(sendToFxBeforeBalance.Balance.Amount)
	require.True(suite.T(), sendToFxBeforeAfter.Balance.Amount.Equal(sendToFxBeforeBalance.Balance.Amount), "externalToFxAndIbcTransfer beforeBalance  != afterBalance",
		sendToFxBeforeBalance.Balance.Amount.String(),
		//sendToFxAmount.String(),
		sendToFxBeforeAfter.Balance.Amount.String())
	suite.T().Logf("\n")
}

func (suite *TronCrosschainTestSuite) setOrchestratorAddress() {
	suite.T().Helper()

	fxAddress := suite.AdminAddress()
	if err := trontypes.ValidateTronAddress(suite.tronAddress().String()); err != nil {
		suite.T().Fatal(err, "external address is invalid", suite.tronAddress().String())
	}
	queryOracleResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).GetOracleByAddr(suite.ctx, &types.QueryOracleByAddrRequest{
		OracleAddress: fxAddress.String(),
		ChainName:     suite.chainName,
	})
	if queryOracleResponse != nil && queryOracleResponse.GetOracle() != nil {
		oracle := queryOracleResponse.GetOracle()
		suite.T().Logf("already set orchestrator address! oracle:[%v], orchestrator:[%v], externalAddress:[%v]\n", oracle.OracleAddress, oracle.BridgerAddress, oracle.ExternalAddress)
		return
	}

	if err != nil {
		if !strings.Contains(err.Error(), "No oracleAddr") {
			suite.T().Fatal(err)
		}
		suite.T().Logf("not found validator!!error msg:%v\n", err.Error())
	}
	chainParams, err := crosschaintypes.NewQueryClient(suite.grpcClient).Params(suite.ctx, &types.QueryParamsRequest{ChainName: suite.chainName})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.BroadcastTx(&types.MsgBondedOracle{
		OracleAddress:   fxAddress.String(),
		BridgerAddress:  fxAddress.String(),
		ExternalAddress: suite.tronAddress().String(),
		DelegateAmount:  chainParams.Params.DelegateThreshold,
		ChainName:       suite.chainName,
	})
	suite.T().Logf("\n")
}
