package tests

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app/helpers"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/x/ibc/applications/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

type CrosschainTestSuite struct {
	TestSuite
	crosschaintypes.BridgeToken
	ibcDenom string

	ethPrivKey *ecdsa.PrivateKey
	chainName  string
}

func TestCrosschainTestSuite(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	const purseTokenContract = "0xFBBbB4f7B1e5bCb0345c5A5a61584B2547d5D582"
	const chainName = "bsc"
	const purseTokenChannelIBC = "transfer/channel-0"
	purseDenom := types.DenomTrace{
		Path:      purseTokenChannelIBC,
		BaseDenom: fmt.Sprintf("%s%s", chainName, purseTokenContract),
	}.IBCDenom()
	suite.Run(t, &CrosschainTestSuite{
		TestSuite: NewTestSuite(),
		BridgeToken: crosschaintypes.BridgeToken{
			Token:      purseTokenContract,
			Denom:      fmt.Sprintf("%s%s", chainName, purseTokenContract),
			ChannelIbc: "px/transfer/channel-0",
		},
		ibcDenom:   purseDenom,
		ethPrivKey: helpers.GenerateEthKey(),
		chainName:  chainName,
	})
}

func (suite *CrosschainTestSuite) TestBSCCrosschain() {

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

func (suite *CrosschainTestSuite) ethAddress() gethcommon.Address {
	return ethCrypto.PubkeyToAddress(suite.ethPrivKey.PublicKey)
}

func (suite *CrosschainTestSuite) queryFxLastEventNonce() uint64 {
	suite.T().Helper()
	lastEventNonce, err := crosschaintypes.NewQueryClient(suite.grpcClient).LastEventNonceByAddr(suite.ctx,
		&crosschaintypes.QueryLastEventNonceByAddrRequest{
			ChainName:      suite.chainName,
			BridgerAddress: suite.AdminAddress().String(),
		})
	suite.Require().NoError(err)
	return lastEventNonce.EventNonce + 1
}

func (suite *CrosschainTestSuite) queryObserver() *crosschaintypes.QueryLastObservedBlockHeightResponse {
	suite.T().Helper()
	height, err := crosschaintypes.NewQueryClient(suite.grpcClient).LastObservedBlockHeight(suite.ctx,
		&crosschaintypes.QueryLastObservedBlockHeightRequest{
			ChainName: suite.chainName,
		})
	suite.Require().NoError(err)
	return height
}

func (suite *CrosschainTestSuite) sendToExternalAndCancel() {
	suite.T().Helper()
	suite.T().Logf("\n####################      FX to ETH      ####################\n")
	sendToEthAmount, _ := sdk.NewIntFromString("20000000000000000000")
	sendToEthFee, _ := sdk.NewIntFromString("30000000000000000000")
	suite.BroadcastTx(&crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Amount:         sendToEthAmount.Add(sendToEthFee),
		Sender:         suite.ethAddress().Hex(),
		Receiver:       suite.AdminAddress().String(),
		TargetIbc:      "",
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	})

	fxAddress := suite.AdminAddress()
	sendToEthBeforeBalance := suite.getBalanceByAddress(fxAddress, suite.ibcDenom)
	suite.T().Logf("send-to-eth before balance:[%v    %v]", sendToEthBeforeBalance.Amount.String(), sendToEthBeforeBalance.Denom)

	sendToEthHash := suite.BroadcastTx(&crosschaintypes.MsgSendToExternal{
		Sender:    suite.AdminAddress().String(),
		Dest:      suite.ethAddress().Hex(),
		Amount:    sdk.NewCoin(suite.ibcDenom, sendToEthAmount),
		BridgeFee: sdk.NewCoin(suite.ibcDenom, sendToEthFee),
		ChainName: suite.chainName,
	})

	sendToEthAfterBalance := suite.getBalanceByAddress(fxAddress, suite.ibcDenom)
	suite.T().Logf("send-to-eth after balance:[%v    %v]", sendToEthAfterBalance.Amount.String(), sendToEthAfterBalance.Denom)

	time.Sleep(3 * time.Second)

	txResponse, err := suite.grpcClient.ServiceClient().GetTx(suite.ctx, &tx.GetTxRequest{Hash: sendToEthHash})
	suite.NoError(err)
	txId, found, err := suite.getSentToExternalTxIdByEvents(txResponse.TxResponse.Logs)
	suite.NoError(err)
	require.True(suite.T(), found)
	require.Greater(suite.T(), txId, uint64(0))
	suite.T().Logf("send-to-eth txId:[%d]", txId)
	_ = suite.BroadcastTx(&crosschaintypes.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        suite.AdminAddress().String(),
		ChainName:     suite.chainName,
	})

	cancelSendToEthAfterBalance := suite.getBalanceByAddress(fxAddress, suite.ibcDenom)
	suite.T().Logf("cancel-send-to-eth after balance:[%v    %v]", cancelSendToEthAfterBalance.Amount.String(), cancelSendToEthAfterBalance.Denom)
	require.True(suite.T(), sendToEthBeforeBalance.Equal(cancelSendToEthAfterBalance))
}

func (suite *CrosschainTestSuite) getBalanceByAddress(accAddr sdk.AccAddress, denom string) *sdk.Coin {
	balanceResp, err := suite.grpcClient.BankQuery().Balance(suite.ctx, banktypes.NewQueryBalanceRequest(accAddr, denom))
	suite.NoError(err)
	return balanceResp.Balance
}

func (suite *CrosschainTestSuite) getSentToExternalTxIdByEvents(logs sdk.ABCIMessageLogs) (uint64, bool, error) {
	for _, eventLog := range logs {
		for _, event := range eventLog.Events {
			if event.Type != crosschaintypes.EventTypeSendToExternal {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != crosschaintypes.AttributeKeyOutgoingTxID {
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

func (suite *CrosschainTestSuite) addBridgeTokenClaim() {
	suite.T().Helper()
	suite.T().Logf("\n####################      Add bridge token claim      ####################\n")
	bridgeToken, err := crosschaintypes.NewQueryClient(suite.grpcClient).TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{ChainName: suite.chainName, Token: suite.BridgeToken.Token})
	expectDenom := types.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: fmt.Sprintf("%s%s", suite.chainName, suite.BridgeToken.Token),
	}.IBCDenom()
	if err == nil && bridgeToken.Denom == expectDenom {
		suite.T().Logf("bridge token already exists!tokenContract:[%v], denom:[%v], channelIbc:[%v]", suite.BridgeToken.Token, bridgeToken.Denom, bridgeToken.ChannelIbc)
		return
	}
	fxOriginatedTokenClaimMsg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Name:           "Pundix Pruse",
		Symbol:         "PURSE",
		Decimals:       18,
		BridgerAddress: suite.AdminAddress().String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	suite.BroadcastTx(fxOriginatedTokenClaimMsg)
	suite.T().Logf("\n")
}

func (suite *CrosschainTestSuite) signPendingValsetRequest() {
	suite.T().Helper()
	defer func() {
		suite.T().Logf("sign pending valset request defer ....\n")
		if err := recover(); err != nil {
			suite.T().Fatal(err)
		}
	}()
	gravityId := suite.queryGravityId()
	requestParams := &crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
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
			checkpoint, _ := valset.GetCheckpoint(gravityId)
			suite.T().Logf("need confirm valset: nonce:%v ethAddress:%v\n", valset.Nonce, suite.ethAddress().Hex())
			signature, err := crosschaintypes.NewEthereumSignature(checkpoint, suite.ethPrivKey)
			if err != nil {
				suite.T().Log(err)
				continue
			}
			suite.BroadcastTx(
				&crosschaintypes.MsgOracleSetConfirm{
					Nonce:           valset.Nonce,
					BridgerAddress:  suite.AdminAddress().String(),
					ExternalAddress: suite.ethAddress().Hex(),
					Signature:       hex.EncodeToString(signature),
					ChainName:       suite.chainName,
				})
		}
	}
}

func (suite *CrosschainTestSuite) queryGravityId() string {
	suite.T().Helper()
	chainParamsResp, err := crosschaintypes.NewQueryClient(suite.grpcClient).Params(suite.ctx, &crosschaintypes.QueryParamsRequest{ChainName: suite.chainName})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.T().Logf("chain params result:%+v\n", chainParamsResp.Params)
	return chainParamsResp.Params.GravityId
}

func (suite *CrosschainTestSuite) confirmBatch() {
	suite.T().Helper()

	gravityId := suite.queryGravityId()
	orchestrator := suite.AdminAddress()
	for {
		lastPendingBatchRequestResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).LastPendingBatchRequestByAddr(suite.ctx,
			&crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{BridgerAddress: orchestrator.String(), ChainName: suite.chainName})
		if err != nil {
			suite.T().Fatal(err)
		}
		outgoingTxBatch := lastPendingBatchRequestResponse.Batch
		if outgoingTxBatch == nil {
			break
		}
		checkpoint, err := outgoingTxBatch.GetCheckpoint(gravityId)
		if err != nil {
			suite.T().Fatal(err)
		}
		signatureBytes, err := crosschaintypes.NewEthereumSignature(checkpoint, suite.ethPrivKey)
		if err != nil {
			suite.T().Fatal(err)
		}

		err = crosschaintypes.ValidateEthereumSignature(checkpoint, signatureBytes, suite.ethAddress().Hex())
		if err != nil {
			suite.T().Fatal(err)
		}
		suite.BroadcastTx([]sdk.Msg{
			&crosschaintypes.MsgConfirmBatch{
				Nonce:           outgoingTxBatch.BatchNonce,
				TokenContract:   outgoingTxBatch.TokenContract,
				BridgerAddress:  suite.AdminAddress().String(),
				ExternalAddress: suite.ethAddress().Hex(),
				Signature:       hex.EncodeToString(signatureBytes),
				ChainName:       suite.chainName,
			},
		}...)
		suite.T().Logf("\n")
		time.Sleep(2 * time.Second)

		suite.BroadcastTx([]sdk.Msg{
			&crosschaintypes.MsgSendToExternalClaim{
				EventNonce:     suite.queryFxLastEventNonce(),
				BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
				BatchNonce:     outgoingTxBatch.BatchNonce,
				TokenContract:  outgoingTxBatch.TokenContract,
				BridgerAddress: suite.AdminAddress().String(),
				ChainName:      suite.chainName,
			},
		}...)
	}
}

func (suite *CrosschainTestSuite) batchRequest() {
	suite.T().Helper()

	batchFeeResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).BatchFees(suite.ctx, &crosschaintypes.QueryBatchFeeRequest{ChainName: suite.chainName})
	if err != nil {
		suite.T().Fatal(err)
	}
	orchestrator := suite.AdminAddress()
	feeReceive := suite.ethAddress().String()
	msgList := make([]sdk.Msg, 0, len(batchFeeResponse.BatchFees))
	for _, batchToken := range batchFeeResponse.BatchFees {
		if batchToken.TotalTxs >= 5 {
			denomResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
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

			msgList = append(msgList, &crosschaintypes.MsgRequestBatch{
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

func (suite *CrosschainTestSuite) fxToExternal(count int) {
	suite.T().Helper()
	suite.T().Logf("\n####################      FX to External      ####################\n")
	msgList := make([]sdk.Msg, 0, count)
	denom := types.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: fmt.Sprintf("%s%s", suite.chainName, suite.BridgeToken.Token),
	}.IBCDenom()
	for i := 0; i < count; i++ {
		msgList = append(msgList, &crosschaintypes.MsgSendToExternal{
			Sender:    suite.AdminAddress().String(),
			Dest:      suite.ethAddress().Hex(),
			Amount:    sdk.NewCoin(denom, sdk.NewInt(111111)),
			BridgeFee: sdk.NewCoin(denom, sdk.NewInt(1000)),
			ChainName: suite.chainName,
		})
	}
	suite.BroadcastTx(msgList...)
}

func (suite *CrosschainTestSuite) showAllBalance(address sdk.AccAddress) {
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

func (suite *CrosschainTestSuite) externalToFx() {
	suite.T().Helper()
	suite.T().Logf("\n####################      External to FX      ####################\n")
	suite.BroadcastTx(&crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Amount:         sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(100000000000000)),
		Sender:         suite.ethAddress().Hex(),
		Receiver:       suite.AdminAddress().String(),
		TargetIbc:      "",
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	})
	suite.T().Logf("\n")
}

func (suite *CrosschainTestSuite) externalToFxAndIbcTransfer() {
	suite.T().Helper()
	suite.T().Logf("\n####################      External to FX to Pundix      ####################\n")

	suite.BroadcastTx(&crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserver().ExternalBlockHeight + 1,
		TokenContract:  suite.BridgeToken.Token,
		Amount:         sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(100000000000000)),
		Sender:         suite.ethAddress().Hex(),
		Receiver:       suite.AdminAddress().String(),
		TargetIbc:      hex.EncodeToString([]byte("0x/transfer/channel-0")),
		BridgerAddress: suite.AdminAddress().String(),
		ChainName:      suite.chainName,
	})
	suite.T().Logf("\n")
}

func (suite *CrosschainTestSuite) setOrchestratorAddress() {
	suite.T().Helper()

	fxAddress := suite.AdminAddress()

	if !gethcommon.IsHexAddress(suite.ethAddress().Hex()) {
		suite.T().Fatal("eth address is invalid")
	}
	queryOrchestratorResponse, err := crosschaintypes.NewQueryClient(suite.grpcClient).GetOracleByBridgerAddr(suite.ctx, &crosschaintypes.QueryOracleByBridgerAddrRequest{
		BridgerAddress: fxAddress.String(),
		ChainName:      suite.chainName,
	})
	if queryOrchestratorResponse != nil && queryOrchestratorResponse.GetOracle() != nil {
		oracle := queryOrchestratorResponse.GetOracle()
		suite.T().Logf("already set orchestrator address! oracle:[%v], orchestrator:[%v], externalAddress:[%v]\n", oracle.OracleAddress, oracle.BridgerAddress, oracle.ExternalAddress)
		return
	}

	if err != nil {
		if !strings.Contains(err.Error(), "No Orchestrator: invalid") {
			suite.T().Fatal(err)
		}
		suite.T().Logf("not found validator!!error msg:%v\n", err.Error())
	}
	chainParams, err := crosschaintypes.NewQueryClient(suite.grpcClient).Params(suite.ctx, &crosschaintypes.QueryParamsRequest{ChainName: suite.chainName})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.BroadcastTx(&crosschaintypes.MsgBondedOracle{
		OracleAddress:   fxAddress.String(),
		BridgerAddress:  fxAddress.String(),
		ExternalAddress: suite.ethAddress().Hex(),
		DelegateAmount:  chainParams.Params.DelegateThreshold,
		ChainName:       suite.chainName,
	})
	suite.T().Logf("\n")
}
