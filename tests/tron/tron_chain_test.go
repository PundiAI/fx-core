package test_tron

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

	trontypes "github.com/functionx/fx-core/x/tron/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

func TestOrchestratorChain(t *testing.T) {
	if testing.Short() {
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
	c.t.Logf("\n####################      FX to External      ####################\n")
	sendToExternalAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000))
	sendToExternalFee := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(2000))

	c.BroadcastTx([]sdk.Msg{&types.MsgSendToFxClaim{
		EventNonce:     c.QueryFxLastEventNonce(),
		BlockHeight:    c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract:  tusdTokenContract,
		Amount:         sendToExternalAmount.Add(sendToExternalFee),
		Sender:         c.externalAddress.String(),
		Receiver:       c.FxAddress().String(),
		TargetIbc:      "",
		BridgerAddress: c.FxAddress().String(),
		ChainName:      c.chainName,
	}})

	fxAddress := c.FxAddress()

	sendToExternalBeforeBalance := getBalanceByAddress(c, fxAddress, fxTusdTokenDenom)
	c.t.Logf("send-to-External before balance:[%v    %v]", sendToExternalBeforeBalance.Amount.String(), sendToExternalBeforeBalance.Denom)

	sendToExternalHash := c.BroadcastTx([]sdk.Msg{&types.MsgSendToExternal{
		Sender:    c.FxAddress().String(),
		Dest:      c.externalAddress.String(),
		Amount:    sdk.NewCoin(fxTusdTokenDenom, sendToExternalAmount),
		BridgeFee: sdk.NewCoin(fxTusdTokenDenom, sendToExternalFee),
		ChainName: c.chainName,
	}})

	sendToExternalAfterBalance := getBalanceByAddress(c, fxAddress, fxTusdTokenDenom)
	c.t.Logf("send-to-External after balance:[%v    %v]", sendToExternalAfterBalance.Amount.String(), sendToExternalAfterBalance.Denom)
	differentAmount := sendToExternalBeforeBalance.Amount.Sub(sendToExternalAfterBalance.Amount)
	require.True(c.t, sendToExternalAmount.Add(sendToExternalFee).Equal(differentAmount), "beforeBalance - afterBalance != sendToExternalFeeAmount+sendToExternalFee",
		sendToExternalBeforeBalance.Amount.String(),
		sendToExternalAfterBalance.Amount.String(),
		sendToExternalAmount.Add(sendToExternalFee).String())

	time.Sleep(3 * time.Second)

	txResponse, err := c.TxClient.GetTx(c.ctx, &tx.GetTxRequest{Hash: sendToExternalHash})
	require.NoError(c.t, err)
	txId, found, err := getSentToExternalTxIdByEvents(txResponse.TxResponse.Logs)
	require.NoError(c.t, err)
	require.True(c.t, found)
	require.Greater(c.t, txId, uint64(0))
	c.t.Logf("send-to-External txId:[%d]", txId)

	_ = c.BroadcastTx([]sdk.Msg{&types.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        c.FxAddress().String(),
		ChainName:     c.chainName,
	}})

	cancelSendToExternalAfterBalance := getBalanceByAddress(c, fxAddress, fxTusdTokenDenom)
	c.t.Logf("cancel-send-to-External after balance:[%v    %v]", cancelSendToExternalAfterBalance.Amount.String(), cancelSendToExternalAfterBalance.Denom)
	require.True(c.t, sendToExternalBeforeBalance.Equal(cancelSendToExternalAfterBalance), sendToExternalBeforeBalance.String(), cancelSendToExternalAfterBalance.String())
}

func getBalanceByAddress(c *Client, accAddr sdk.AccAddress, denom string) *sdk.Coin {
	balanceResp, err := c.bankQueryClient.Balance(c.ctx, banktypes.NewQueryBalanceRequest(accAddr, denom))
	require.NoError(c.t, err)
	return balanceResp.Balance
}

//
func getSentToExternalTxIdByEvents(logs sdk.ABCIMessageLogs) (uint64, bool, error) {
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

func addBridgeTokenClaim(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      Add bridge token claim      ####################\n")
	bridgeToken, err := c.crosschainQueryClient.TokenToDenom(c.ctx, &types.QueryTokenToDenomRequest{ChainName: c.chainName, Token: tusdTokenContract})

	if err != nil && !strings.Contains(err.Error(), "bridge token is not exist") {
		c.t.Fatal(err)
	}
	if err == nil && bridgeToken.Denom == fmt.Sprintf("%v%v", c.chainName, tusdTokenContract) {
		c.t.Logf("bridge token already exists!tokenContract:[%v], denom:[%v], channelIbc:[%v]", tusdTokenContract, bridgeToken.Denom, bridgeToken.ChannelIbc)
		return
	}
	fxOriginatedTokenClaimMsg := &types.MsgBridgeTokenClaim{
		EventNonce:     c.QueryFxLastEventNonce(),
		BlockHeight:    c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract:  tusdTokenContract,
		Name:           tusdTokenName,
		Symbol:         tusdTokenSymbol,
		Decimals:       18,
		BridgerAddress: c.FxAddress().String(),
		ChannelIbc:     "",
		ChainName:      c.chainName,
	}
	c.BroadcastTx([]sdk.Msg{fxOriginatedTokenClaimMsg})
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
	requestParams := &types.QueryLastPendingOracleSetRequestByAddrRequest{
		BridgerAddress: c.FxAddress().String(),
		ChainName:      c.chainName,
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
			checkpoint, err := trontypes.GetCheckpointOracleSet(valset, gravityId)
			require.NoError(c.t, err)
			c.t.Logf("need confirm valset: nonce:%v ExternalAddress:%v\n", valset.Nonce, c.externalAddress.Hex())
			signature, err := trontypes.NewTronSignature(checkpoint, c.externalPrivKey)
			if err != nil {
				c.t.Log(err)
				continue
			}
			c.BroadcastTx([]sdk.Msg{
				&types.MsgOracleSetConfirm{
					Nonce:           valset.Nonce,
					BridgerAddress:  c.FxAddress().String(),
					ExternalAddress: c.externalAddress.String(),
					Signature:       hex.EncodeToString(signature),
					ChainName:       c.chainName,
				},
			})
		}
	}
}

var (
	chainParams *types.Params
)

func queryGravityId(c *Client) string {
	c.t.Helper()
	once := &sync.Once{}
	once.Do(func() {
		chainParamsResp, err := c.crosschainQueryClient.Params(c.ctx, &types.QueryParamsRequest{ChainName: c.chainName})
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
			&types.QueryLastPendingBatchRequestByAddrRequest{BridgerAddress: orchestrator.String(), ChainName: c.chainName})
		if err != nil {
			c.t.Fatal(err)
		}
		outgoingTxBatch := lastPendingBatchRequestResponse.Batch
		if outgoingTxBatch == nil {
			break
		}
		checkpoint, err := trontypes.GetCheckpointConfirmBatch(outgoingTxBatch, gravityId)
		if err != nil {
			c.t.Fatal(err)
		}
		signatureBytes, err := trontypes.NewTronSignature(checkpoint, c.externalPrivKey)
		if err != nil {
			c.t.Fatal(err)
		}

		err = trontypes.ValidateTronSignature(checkpoint, signatureBytes, c.externalAddress.String())
		if err != nil {
			c.t.Fatal(err)
		}
		c.BroadcastTx([]sdk.Msg{
			&types.MsgConfirmBatch{
				Nonce:           outgoingTxBatch.BatchNonce,
				TokenContract:   outgoingTxBatch.TokenContract,
				BridgerAddress:  c.FxAddress().String(),
				ExternalAddress: c.externalAddress.String(),
				Signature:       hex.EncodeToString(signatureBytes),
				ChainName:       c.chainName,
			},
		})
		c.t.Logf("\n")
		time.Sleep(2 * time.Second)

		c.BroadcastTx([]sdk.Msg{
			&types.MsgSendToExternalClaim{
				EventNonce:     c.QueryFxLastEventNonce(),
				BlockHeight:    c.QueryObserver().ExternalBlockHeight + 1,
				BatchNonce:     outgoingTxBatch.BatchNonce,
				TokenContract:  outgoingTxBatch.TokenContract,
				BridgerAddress: c.FxAddress().String(),
				ChainName:      c.chainName,
			},
		})
	}
}

//

func batchRequest(c *Client) {
	c.t.Helper()

	batchFeeResponse, err := c.crosschainQueryClient.BatchFees(c.ctx, &types.QueryBatchFeeRequest{ChainName: c.chainName})
	if err != nil {
		c.t.Fatal(err)
	}
	orchestrator := c.FxAddress()
	feeReceive := c.externalAddress.String()
	msgList := make([]sdk.Msg, 0, len(batchFeeResponse.BatchFees))
	for _, batchToken := range batchFeeResponse.BatchFees {
		if batchToken.TotalTxs >= 5 {
			denomResponse, err := c.crosschainQueryClient.TokenToDenom(c.ctx, &types.QueryTokenToDenomRequest{
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

			msgList = append(msgList, &types.MsgRequestBatch{
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
	c.BroadcastTx(msgList)
	c.t.Logf("\n")
}

func fxToExternal(c *Client, count int) {
	c.t.Helper()
	c.t.Logf("\n####################      FX to External      ####################\n")
	sendToFxBeforeBalance, err := c.bankQueryClient.Balance(c.ctx, &banktypes.QueryBalanceRequest{
		Address: c.FxAddress().String(),
		Denom:   fxTusdTokenDenom,
	})
	require.NoError(c.t, err)
	sendToExternalAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(1900))
	sendToExternalFee := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(100))
	totalSendToExternalAmount := sdk.ZeroInt()
	msgList := make([]sdk.Msg, 0, count)
	for i := 0; i < count; i++ {
		msgList = append(msgList, &types.MsgSendToExternal{
			Sender:    c.FxAddress().String(),
			Dest:      c.externalAddress.String(),
			Amount:    sdk.NewCoin(fxTusdTokenDenom, sendToExternalAmount),
			BridgeFee: sdk.NewCoin(fxTusdTokenDenom, sendToExternalFee),
			ChainName: c.chainName,
		})
		totalSendToExternalAmount = totalSendToExternalAmount.Add(sendToExternalAmount).Add(sendToExternalFee)
	}
	c.BroadcastTx(msgList)
	sendToFxBeforeAfter, err := c.bankQueryClient.Balance(c.ctx, &banktypes.QueryBalanceRequest{
		Address: c.FxAddress().String(),
		Denom:   fxTusdTokenDenom,
	})
	require.NoError(c.t, err)
	differentAmount := sendToFxBeforeBalance.Balance.Amount.Sub(sendToFxBeforeAfter.Balance.Amount)
	require.True(c.t, totalSendToExternalAmount.Equal(differentAmount), "beforeBalance - afterBalance !=  totalSendToExternalAmount",
		sendToFxBeforeBalance.Balance.Amount.String(),
		sendToFxBeforeAfter.Balance.Amount.String(),
		totalSendToExternalAmount.String(),
	)
	c.t.Logf("\n")
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
	sendToFxBeforeBalance, err := c.bankQueryClient.Balance(c.ctx, &banktypes.QueryBalanceRequest{
		Address: c.FxAddress().String(),
		Denom:   fxTusdTokenDenom,
	})
	require.NoError(c.t, err)
	sendToFxAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000))
	c.BroadcastTx([]sdk.Msg{&types.MsgSendToFxClaim{
		EventNonce:     c.QueryFxLastEventNonce(),
		BlockHeight:    c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract:  tusdTokenContract,
		Amount:         sendToFxAmount,
		Sender:         c.externalAddress.String(),
		Receiver:       c.FxAddress().String(),
		TargetIbc:      "",
		BridgerAddress: c.FxAddress().String(),
		ChainName:      c.chainName,
	}})
	sendToFxBeforeAfter, err := c.bankQueryClient.Balance(c.ctx, &banktypes.QueryBalanceRequest{
		Address: c.FxAddress().String(),
		Denom:   fxTusdTokenDenom,
	})
	require.NoError(c.t, err)
	differentAmount := sendToFxBeforeAfter.Balance.Amount.Sub(sendToFxBeforeBalance.Balance.Amount)
	require.True(c.t, sendToFxAmount.Equal(differentAmount), "beforeBalance + sendToFxAmount != afterBalance",
		sendToFxBeforeBalance.Balance.Amount.String(),
		sendToFxAmount.String(),
		sendToFxBeforeAfter.Balance.Amount.String())
	c.t.Logf("\n")
}

func externalToFxAndIbcTransfer(c *Client) {
	c.t.Helper()
	c.t.Logf("\n####################      External to FX to Pundix      ####################\n")

	sendToFxBeforeBalance, err := c.bankQueryClient.Balance(c.ctx, &banktypes.QueryBalanceRequest{
		Address: c.FxAddress().String(),
		Denom:   fxTusdTokenDenom,
	})
	require.NoError(c.t, err)
	sendToFxAmount := sdk.NewIntWithDecimal(10, 18).Mul(sdk.NewInt(10000))
	c.BroadcastTx([]sdk.Msg{&types.MsgSendToFxClaim{
		EventNonce:     c.QueryFxLastEventNonce(),
		BlockHeight:    c.QueryObserver().ExternalBlockHeight + 1,
		TokenContract:  tusdTokenContract,
		Amount:         sendToFxAmount,
		Sender:         c.externalAddress.String(),
		Receiver:       c.FxAddress().String(),
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: c.FxAddress().String(),
		ChainName:      c.chainName,
	}})
	sendToFxBeforeAfter, err := c.bankQueryClient.Balance(c.ctx, &banktypes.QueryBalanceRequest{
		Address: c.FxAddress().String(),
		Denom:   fxTusdTokenDenom,
	})
	require.NoError(c.t, err)
	//differentAmount := sendToFxBeforeAfter.Balance.Amount.Sub(sendToFxBeforeBalance.Balance.Amount)
	require.True(c.t, sendToFxBeforeAfter.Balance.Amount.Equal(sendToFxBeforeBalance.Balance.Amount), "externalToFxAndIbcTransfer beforeBalance  != afterBalance",
		sendToFxBeforeBalance.Balance.Amount.String(),
		//sendToFxAmount.String(),
		sendToFxBeforeAfter.Balance.Amount.String())
	c.t.Logf("\n")
}

func setOrchestratorAddress(c *Client) {
	c.t.Helper()

	fxAddress := c.FxAddress()
	if err := trontypes.ValidateTronAddress(c.externalAddress.String()); err != nil {
		c.t.Fatal(err, "external address is invalid", c.externalAddress.String())
	}
	queryOracleResponse, err := c.crosschainQueryClient.GetOracleByAddr(c.ctx, &types.QueryOracleByAddrRequest{
		OracleAddress: fxAddress.String(),
		ChainName:     c.chainName,
	})
	if queryOracleResponse != nil && queryOracleResponse.GetOracle() != nil {
		oracle := queryOracleResponse.GetOracle()
		c.t.Logf("already set orchestrator address! oracle:[%v], orchestrator:[%v], externalAddress:[%v]\n", oracle.OracleAddress, oracle.BridgerAddress, oracle.ExternalAddress)
		return
	}

	if err != nil {
		if !strings.Contains(err.Error(), "No oracleAddr") {
			c.t.Fatal(err)
		}
		c.t.Logf("not found validator!!error msg:%v\n", err.Error())
	}
	chainParams, err := c.crosschainQueryClient.Params(c.ctx, &types.QueryParamsRequest{ChainName: c.chainName})
	if err != nil {
		c.t.Fatal(err)
	}
	c.BroadcastTx([]sdk.Msg{&types.MsgCreateOracleBridger{
		OracleAddress:   fxAddress.String(),
		BridgerAddress:  fxAddress.String(),
		ExternalAddress: c.externalAddress.String(),
		DelegateAmount:  chainParams.Params.DelegateThreshold,
		ChainName:       c.chainName,
	}})
	c.t.Logf("\n")
}
