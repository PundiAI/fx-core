package crosschain

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/x/ibc/applications/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gethcommon "github.com/ethereum/go-ethereum/common"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

// TestOrchestratorChain init operator.
// 1. fxcored tx crosschain init-crosschain-params bsc 10000000000000000000000FX --title="Init Bsc chain params", --desc="about bsc chain description" --gravity-id="bsc" --oracles="fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2" --from fx1 -y --gas=auto --gas-adjustment=1.3
// 2. fxcored tx gov vote 1 yes --from fx1 -y --gas=auto --gas-adjustment=1.3
func TestOrchestratorChain(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	client := NewClient(t)

	go signPendingValsetRequest(client)
	setOrchestratorAddress(client)

	addBridgeTokenClaim(client)

	externalToFx(client)

	externalToFxAndIbcTransfer(client)

	showAllBalance(client, client.FxAddress())

	fxToExternal(client, 5)

	batchRequest(client)

	confirmBatch(client)

	sendToExternalAndCancel(client)
}

func sendToExternalAndCancel(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      FX to ETH      ####################\n")
	//denom := fmt.Sprintf("%s%s", "eth", ethTokenContract)
	sendToEthAmount, _ := sdk.NewIntFromString("20000000000000000000")
	sendToEthFee, _ := sdk.NewIntFromString("30000000000000000000")
	c.BroadcastTx(&[]sdk.Msg{&crosschaintypes.MsgSendToFxClaim{
		EventNonce:    c.QueryFxLastEventNonce(),
		BlockHeight:   c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract: purseTokenContract,
		Amount:        sendToEthAmount.Add(sendToEthFee),
		Sender:        c.ethAddress.Hex(),
		Receiver:      c.FxAddress().String(),
		TargetIbc:     "",
		Orchestrator:  c.FxAddress().String(),
		ChainName:     c.chainName,
	}})

	fxAddress := c.FxAddress()
	sendToEthBeforeBalance := getBalanceByAddress(c, fxAddress, purseDenom)
	c.t.Logf("send-to-eth before balance:[%v    %v]", sendToEthBeforeBalance.Amount.String(), sendToEthBeforeBalance.Denom)

	sendToEthHash := c.BroadcastTx(&[]sdk.Msg{&crosschaintypes.MsgSendToExternal{
		Sender:    c.FxAddress().String(),
		Dest:      c.ethAddress.Hex(),
		Amount:    sdk.NewCoin(purseDenom, sendToEthAmount),
		BridgeFee: sdk.NewCoin(purseDenom, sendToEthFee),
		ChainName: c.chainName,
	}})

	sendToEthAfterBalance := getBalanceByAddress(c, fxAddress, purseDenom)
	c.t.Logf("send-to-eth after balance:[%v    %v]", sendToEthAfterBalance.Amount.String(), sendToEthAfterBalance.Denom)

	time.Sleep(3 * time.Second)

	txResponse, err := c.TxClient.GetTx(c.ctx, &tx.GetTxRequest{Hash: sendToEthHash})
	require.NoError(c.t, err)
	txId, found, err := getSentToEthTxIdByEvents(txResponse.TxResponse.Logs)
	require.NoError(c.t, err)
	require.True(c.t, found)
	require.Greater(c.t, txId, uint64(0))
	c.t.Logf("send-to-eth txId:[%d]", txId)
	_ = c.BroadcastTx(&[]sdk.Msg{&crosschaintypes.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        c.FxAddress().String(),
		ChainName:     c.chainName,
	}})

	cancelSendToEthAfterBalance := getBalanceByAddress(c, fxAddress, purseDenom)
	c.t.Logf("cancel-send-to-eth after balance:[%v    %v]", cancelSendToEthAfterBalance.Amount.String(), cancelSendToEthAfterBalance.Denom)
	require.True(c.t, sendToEthBeforeBalance.Equal(cancelSendToEthAfterBalance))
}

func getBalanceByAddress(c *Client, accAddr sdk.AccAddress, denom string) *sdk.Coin {
	balanceResp, err := c.bankQueryClient.Balance(c.ctx, banktypes.NewQueryBalanceRequest(accAddr, denom))
	require.NoError(c.t, err)
	return balanceResp.Balance
}

//
func getSentToEthTxIdByEvents(logs sdk.ABCIMessageLogs) (uint64, bool, error) {
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

func addBridgeTokenClaim(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      Add bridge token claim      ####################\n")
	bridgeToken, err := c.crosschainQueryClient.TokenToDenom(c.ctx, &crosschaintypes.QueryTokenToDenomRequest{ChainName: c.chainName, Token: purseTokenContract})
	expectDenom := types.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: fmt.Sprintf("%s%s", c.chainName, purseTokenContract),
	}.IBCDenom()
	if err == nil && bridgeToken.Denom == expectDenom {
		c.t.Logf("bridge token already exists!tokenContract:[%v], denom:[%v], channelIbc:[%v]", purseTokenContract, bridgeToken.Denom, bridgeToken.ChannelIbc)
		return
	}
	fxOriginatedTokenClaimMsg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:    c.QueryFxLastEventNonce(),
		BlockHeight:   c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract: purseTokenContract,
		Name:          "Pundix Pruse",
		Symbol:        purseTokenSymbol,
		Decimals:      18,
		Orchestrator:  c.FxAddress().String(),
		ChannelIbc:    hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:     c.chainName,
	}
	c.BroadcastTx(&[]sdk.Msg{fxOriginatedTokenClaimMsg})
	c.t.Logf("\n")
}

func signPendingValsetRequest(c *Client) {
	c.t.Helper()
	defer func() {
		c.t.Logf("sign pending valset request defer ....\n")
		if err := recover(); err != nil {
			c.t.Fatal(err)
		}
	}()
	gravityId := queryGravityId(c)
	requestParams := &crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
		OrchestratorAddress: c.FxAddress().String(),
		ChainName:           c.chainName,
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		queryResponse, err := c.crosschainQueryClient.LastPendingOracleSetRequestByAddr(c.ctx, requestParams)
		if err != nil {
			c.t.Logf("query last pending valset request is err!params:%+v, errors:%v\n", requestParams, err)
			continue
		}
		valsets := queryResponse.OracleSets
		if len(valsets) <= 0 {
			continue
		}
		for _, valset := range valsets {
			checkpoint := valset.GetCheckpoint(gravityId)
			c.t.Logf("need confirm valset: nonce:%v EthAddress:%v\n", valset.Nonce, c.ethAddress.Hex())
			signature, err := crosschaintypes.NewEthereumSignature(checkpoint, c.ethPrivKey)
			if err != nil {
				c.t.Log(err)
				continue
			}
			c.BroadcastTx(&[]sdk.Msg{
				&crosschaintypes.MsgOracleSetConfirm{
					Nonce:               valset.Nonce,
					OrchestratorAddress: c.FxAddress().String(),
					ExternalAddress:     c.ethAddress.Hex(),
					Signature:           hex.EncodeToString(signature),
					ChainName:           c.chainName,
				},
			})
		}
	}
}

var (
	chainParams *crosschaintypes.Params
)

func queryGravityId(c *Client) string {
	c.t.Helper()
	once := &sync.Once{}
	once.Do(func() {
		chainParamsResp, err := c.crosschainQueryClient.Params(c.ctx, &crosschaintypes.QueryParamsRequest{ChainName: c.chainName})
		if err != nil {
			c.t.Fatal(err)
		}
		chainParams = &chainParamsResp.Params
		c.t.Logf("chain params result:%+v\n", chainParams)
	})
	return chainParams.GravityId
}

func confirmBatch(c *Client) {
	c.t.Helper()

	gravityId := queryGravityId(c)
	orchestrator := c.FxAddress()
	for {
		lastPendingBatchRequestResponse, err := c.crosschainQueryClient.LastPendingBatchRequestByAddr(c.ctx,
			&crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{OrchestratorAddress: orchestrator.String(), ChainName: c.chainName})
		if err != nil {
			c.t.Fatal(err)
		}
		outgoingTxBatch := lastPendingBatchRequestResponse.Batch
		if outgoingTxBatch == nil {
			break
		}
		checkpoint, err := outgoingTxBatch.GetCheckpoint(gravityId)
		if err != nil {
			c.t.Fatal(err)
		}
		signatureBytes, err := crosschaintypes.NewEthereumSignature(checkpoint, c.ethPrivKey)
		if err != nil {
			c.t.Fatal(err)
		}

		err = crosschaintypes.ValidateEthereumSignature(checkpoint, signatureBytes, c.ethAddress.Hex())
		if err != nil {
			c.t.Fatal(err)
		}
		c.BroadcastTx(&[]sdk.Msg{
			&crosschaintypes.MsgConfirmBatch{
				Nonce:               outgoingTxBatch.BatchNonce,
				TokenContract:       outgoingTxBatch.TokenContract,
				OrchestratorAddress: c.FxAddress().String(),
				ExternalAddress:     c.ethAddress.Hex(),
				Signature:           hex.EncodeToString(signatureBytes),
				ChainName:           c.chainName,
			},
		})
		c.t.Logf("\n")
		time.Sleep(2 * time.Second)

		c.BroadcastTx(&[]sdk.Msg{
			&crosschaintypes.MsgSendToExternalClaim{
				EventNonce:    c.QueryFxLastEventNonce(),
				BlockHeight:   c.QueryObserver().ExternalBlockHeight + 1,
				BatchNonce:    outgoingTxBatch.BatchNonce,
				TokenContract: outgoingTxBatch.TokenContract,
				Orchestrator:  c.FxAddress().String(),
				ChainName:     c.chainName,
			},
		})
	}
}

//

func batchRequest(c *Client) {
	c.t.Helper()

	batchFeeResponse, err := c.crosschainQueryClient.BatchFees(c.ctx, &crosschaintypes.QueryBatchFeeRequest{ChainName: c.chainName})
	if err != nil {
		c.t.Fatal(err)
	}
	orchestrator := c.FxAddress()
	feeReceive := c.ethAddress.String()
	msgList := make([]sdk.Msg, 0, len(batchFeeResponse.BatchFees))
	for _, batchToken := range batchFeeResponse.BatchFees {
		if batchToken.TotalTxs >= 5 {
			denomResponse, err := c.crosschainQueryClient.TokenToDenom(c.ctx, &crosschaintypes.QueryTokenToDenomRequest{
				Token:     batchToken.TokenContract,
				ChainName: c.chainName,
			})
			if err != nil {
				c.t.Fatal(err)
			}
			if strings.HasPrefix(denomResponse.Denom, batchToken.TokenContract) {
				c.t.Logf("warn!!! not found token contract, tokenContract:[%v], erc20ToDenom response:[%v]\n", batchToken.TokenContract, denomResponse.Denom)
				continue
			}

			msgList = append(msgList, &crosschaintypes.MsgRequestBatch{
				Sender:     orchestrator.String(),
				Denom:      denomResponse.Denom,
				MinimumFee: batchToken.TotalFees,
				FeeReceive: feeReceive,
				ChainName:  c.chainName,
			})
		}
	}
	if len(msgList) <= 0 {
		return
	}
	c.BroadcastTx(&msgList)
	c.t.Logf("\n")
}

func fxToExternal(c *Client, count int) {
	c.t.Helper()
	c.t.Logf("\n####################      FX to External      ####################\n")
	msgList := make([]sdk.Msg, 0, count)
	denom := types.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: fmt.Sprintf("%s%s", c.chainName, purseTokenContract),
	}.IBCDenom()
	for i := 0; i < count; i++ {
		msgList = append(msgList, &crosschaintypes.MsgSendToExternal{
			Sender:    c.FxAddress().String(),
			Dest:      c.ethAddress.Hex(),
			Amount:    sdk.NewCoin(denom, sdk.NewInt(111111)),
			BridgeFee: sdk.NewCoin(denom, sdk.NewInt(1000)),
			ChainName: c.chainName,
		})
	}
	c.BroadcastTx(&msgList)
}

func showAllBalance(c *Client, address sdk.AccAddress) {
	c.t.Helper()
	c.t.Logf("\n####################      Query Address Balance      ####################\n")
	queryAllBalancesResponse, err := c.bankQueryClient.AllBalances(c.ctx, banktypes.NewQueryAllBalancesRequest(address, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: true,
	}))
	if err != nil {
		c.t.Fatal(err)
	}
	c.t.Logf("address: [%v] all balance\n", address.String())
	for _, balance := range queryAllBalancesResponse.Balances {
		c.t.Logf("denom:%v, amount:%v\n", balance.Denom, balance.Amount.String())
	}
	c.t.Logf("\n")
}

func externalToFx(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      External to FX      ####################\n")
	c.BroadcastTx(&[]sdk.Msg{&crosschaintypes.MsgSendToFxClaim{
		EventNonce:    c.QueryFxLastEventNonce(),
		BlockHeight:   c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract: purseTokenContract,
		Amount:        sdk.NewInt(11111111),
		Sender:        c.ethAddress.Hex(),
		Receiver:      c.FxAddress().String(),
		TargetIbc:     "",
		Orchestrator:  c.FxAddress().String(),
		ChainName:     c.chainName,
	}})
	c.t.Logf("\n")
}

func externalToFxAndIbcTransfer(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      External to FX to Pundix      ####################\n")

	c.BroadcastTx(&[]sdk.Msg{&crosschaintypes.MsgSendToFxClaim{
		EventNonce:    c.QueryFxLastEventNonce(),
		BlockHeight:   c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract: purseTokenContract,
		Amount:        sdk.NewInt(22222222),
		Sender:        c.ethAddress.Hex(),
		Receiver:      c.FxAddress().String(),
		TargetIbc:     hex.EncodeToString([]byte("px/transfer/channel-0")),
		Orchestrator:  c.FxAddress().String(),
		ChainName:     c.chainName,
	}})
	c.t.Logf("\n")
}

func setOrchestratorAddress(c *Client) {
	c.t.Helper()

	fxAddress := c.FxAddress()

	if !gethcommon.IsHexAddress(c.ethAddress.Hex()) {
		c.t.Fatal("eth address is invalid")
	}
	queryOrchestratorResponse, err := c.crosschainQueryClient.GetOracleByOrchestrator(c.ctx, &crosschaintypes.QueryOracleByOrchestratorRequest{
		OrchestratorAddress: fxAddress.String(),
		ChainName:           c.chainName,
	})
	if queryOrchestratorResponse != nil && queryOrchestratorResponse.GetOracle() != nil {
		oracle := queryOrchestratorResponse.GetOracle()
		c.t.Logf("already set orchestrator address! oracle:[%v], orchestrator:[%v], externalAddress:[%v]\n", oracle.OracleAddress, oracle.OrchestratorAddress, oracle.ExternalAddress)
		return
	}

	if err != nil {
		if !strings.Contains(err.Error(), "No Orchestrator: invalid: invalid request") {
			c.t.Fatal(err)
		}
		c.t.Logf("not found validator!!error msg:%v\n", err.Error())
	}
	chainParams, err := c.crosschainQueryClient.Params(c.ctx, &crosschaintypes.QueryParamsRequest{ChainName: c.chainName})
	if err != nil {
		c.t.Fatal(err)
	}
	c.BroadcastTx(&[]sdk.Msg{&crosschaintypes.MsgSetOrchestratorAddress{
		Oracle:          fxAddress.String(),
		Orchestrator:    fxAddress.String(),
		ExternalAddress: c.ethAddress.Hex(),
		Deposit:         chainParams.Params.DepositThreshold,
		ChainName:       c.chainName,
	}})
	c.t.Logf("\n")
}
