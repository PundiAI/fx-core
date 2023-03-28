package tests

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type CrosschainTestSuite struct {
	*TestSuite
	params          crosschaintypes.Params
	chainName       string
	oraclePrivKey   cryptotypes.PrivKey
	bridgerPrivKey  cryptotypes.PrivKey
	externalPrivKey *ecdsa.PrivateKey
	privKey         cryptotypes.PrivKey
}

func NewCrosschainWithTestSuite(chainName string, ts *TestSuite) CrosschainTestSuite {
	externalPrivKey, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(err.Error())
	}
	return CrosschainTestSuite{
		TestSuite:       ts,
		chainName:       chainName,
		oraclePrivKey:   helpers.NewPriKey(),
		bridgerPrivKey:  helpers.NewPriKey(),
		externalPrivKey: externalPrivKey,
		privKey:         helpers.NewEthPrivKey(),
	}
}

func (suite *CrosschainTestSuite) Init() {
	suite.Send(suite.OracleAddr(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.BridgerAddr(), suite.NewCoin(sdkmath.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.AccAddress(), suite.NewCoin(sdkmath.NewInt(1_000).MulRaw(1e18)))
	suite.params = suite.QueryParams()
}

func (suite *CrosschainTestSuite) OracleAddr() sdk.AccAddress {
	return suite.oraclePrivKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) ExternalAddr() string {
	if suite.chainName == trontypes.ModuleName {
		return tronaddress.PubkeyToAddress(suite.externalPrivKey.PublicKey).String()
	}
	return ethcrypto.PubkeyToAddress(suite.externalPrivKey.PublicKey).String()
}

func (suite *CrosschainTestSuite) BridgerAddr() sdk.AccAddress {
	return suite.bridgerPrivKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) AccAddress() sdk.AccAddress {
	return suite.privKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) HexAddress() gethcommon.Address {
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

func (suite *CrosschainTestSuite) QueryPendingUnbatchedTx(sender sdk.AccAddress) []*crosschaintypes.OutgoingTransferTx {
	pendingTx, err := suite.CrosschainQuery().GetPendingSendToExternal(suite.ctx, &crosschaintypes.QueryPendingSendToExternalRequest{
		ChainName:     suite.chainName,
		SenderAddress: sender.String(),
	})
	suite.NoError(err)
	return pendingTx.UnbatchedTransfers
}

func (suite *CrosschainTestSuite) queryFxLastEventNonce() uint64 {
	lastEventNonce, err := suite.CrosschainQuery().LastEventNonceByAddr(suite.ctx,
		&crosschaintypes.QueryLastEventNonceByAddrRequest{
			ChainName:      suite.chainName,
			BridgerAddress: suite.BridgerAddr().String(),
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

func (suite *CrosschainTestSuite) AddBridgeTokenClaim(name, symbol string, decimals uint64, token, channelIBCHex string) string {
	response, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.ErrorContains(err, "code = NotFound desc = bridge token")
	suite.Nil(response)

	suite.BroadcastTx(suite.bridgerPrivKey, &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: suite.BridgerAddr().String(),
		ChannelIbc:     channelIBCHex,
		ChainName:      suite.chainName,
	})

	response, err = suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	if len(channelIBCHex) > 0 {
		bridgeDenom := fmt.Sprintf("%s%s", suite.chainName, token)
		trace, err := fxtypes.GetIbcDenomTrace(bridgeDenom, channelIBCHex)
		suite.NoError(err)

		bridgeDenom = trace.IBCDenom()
		suite.Equal(bridgeDenom, response.Denom)
	} else if response.Denom != fxtypes.DefaultDenom {
		suite.Equal(fmt.Sprintf("%s%s", suite.chainName, token), response.Denom)
	}

	return response.Denom
}

func (suite *CrosschainTestSuite) GetBridgeTokens() (denoms []*crosschaintypes.BridgeToken) {
	response, err := suite.CrosschainQuery().BridgeTokens(suite.ctx, &crosschaintypes.QueryBridgeTokensRequest{
		ChainName: suite.chainName,
	})
	suite.NoError(err)
	return response.BridgeTokens
}

func (suite *CrosschainTestSuite) GetBridgeDenomByToken(token string) (denom string) {
	response, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	suite.NotEmpty(response.Denom)
	return response.Denom
}

func (suite *CrosschainTestSuite) BondedOracle() {
	response, err := suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Error(err, crosschaintypes.ErrNoFoundOracle)
	suite.Nil(response)

	txResponse := suite.BroadcastTx(suite.oraclePrivKey, &crosschaintypes.MsgBondedOracle{
		OracleAddress:    suite.OracleAddr().String(),
		BridgerAddress:   suite.BridgerAddr().String(),
		ExternalAddress:  suite.ExternalAddr(),
		ValidatorAddress: suite.GetFirstValAddr().String(),
		DelegateAmount:   suite.params.DelegateThreshold,
		ChainName:        suite.chainName,
	})

	response, err = suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)
	suite.Equal(crosschaintypes.Oracle{
		OracleAddress:     suite.OracleAddr().String(),
		BridgerAddress:    suite.BridgerAddr().String(),
		ExternalAddress:   suite.ExternalAddr(),
		DelegateAmount:    suite.params.DelegateThreshold.Amount,
		StartHeight:       txResponse.Height,
		Online:            true,
		DelegateValidator: suite.GetFirstValAddr().String(),
		SlashTimes:        0,
	}, *response.Oracle)
}

func (suite *CrosschainTestSuite) SendUpdateChainOraclesProposal() (*sdk.TxResponse, uint64) {
	msg := &crosschaintypes.MsgUpdateChainOracles{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Oracles:   []string{suite.OracleAddr().String()},
		ChainName: suite.chainName,
	}
	return suite.BroadcastProposalTx2([]sdk.Msg{msg}, "UpdateChainOraclesProposal", "UpdateChainOraclesProposal")
}

func (suite *CrosschainTestSuite) SendOracleSetConfirm() {
	queryResponse, err := suite.CrosschainQuery().LastPendingOracleSetRequestByAddr(suite.ctx,
		&crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)

	for _, oracleSet := range queryResponse.OracleSets {
		var signature []byte
		if suite.chainName == trontypes.ModuleName {
			checkpoint, err := trontypes.GetCheckpointOracleSet(oracleSet, suite.params.GravityId)
			suite.NoError(err)
			signature, err = trontypes.NewTronSignature(checkpoint, suite.externalPrivKey)
			suite.NoError(err)
			err = trontypes.ValidateTronSignature(checkpoint, signature, suite.ExternalAddr())
			suite.NoError(err)
		} else {
			checkpoint, err := oracleSet.GetCheckpoint(suite.params.GravityId)
			suite.NoError(err)
			signature, err = crosschaintypes.NewEthereumSignature(checkpoint, suite.externalPrivKey)
			suite.NoError(err)
			err = crosschaintypes.ValidateEthereumSignature(checkpoint, signature, suite.ExternalAddr())
			suite.NoError(err)
		}

		suite.BroadcastTx(suite.bridgerPrivKey, &crosschaintypes.MsgOracleSetConfirm{
			Nonce:           oracleSet.Nonce,
			BridgerAddress:  suite.BridgerAddr().String(),
			ExternalAddress: suite.ExternalAddr(),
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		})
	}
}

func (suite *CrosschainTestSuite) SendToFxClaim(token string, amount sdkmath.Int, targetIbc string) {
	sender := suite.HexAddress().Hex()
	if suite.chainName == trontypes.ModuleName {
		sender = trontypes.AddressFromHex(sender)
	}
	suite.BroadcastTx(suite.bridgerPrivKey, &crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Amount:         amount,
		Sender:         sender,
		Receiver:       suite.AccAddress().String(),
		TargetIbc:      hex.EncodeToString([]byte(targetIbc)),
		BridgerAddress: suite.BridgerAddr().String(),
		ChainName:      suite.chainName,
	})
	bridgeToken, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	if bridgeToken.Denom == fxtypes.DefaultDenom && len(targetIbc) <= 0 {
		balances := suite.QueryBalances(suite.AccAddress())
		suite.True(balances.IsAllGTE(sdk.NewCoins(sdk.NewCoin(bridgeToken.Denom, amount))))
	}
}

func (suite *CrosschainTestSuite) SendToExternal(count int, amount sdk.Coin) uint64 {
	msgList := make([]sdk.Msg, 0, count)
	for i := 0; i < count; i++ {
		dest := suite.HexAddress().Hex()
		if suite.chainName == trontypes.ModuleName {
			dest = trontypes.AddressFromHex(dest)
		}
		msgList = append(msgList, &crosschaintypes.MsgSendToExternal{
			Sender:    suite.AccAddress().String(),
			Dest:      dest,
			Amount:    amount.SubAmount(sdkmath.NewInt(1)),
			BridgeFee: sdk.NewCoin(amount.Denom, sdkmath.NewInt(1)),
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

func (suite *CrosschainTestSuite) SendToExternalAndCancel(coin sdk.Coin) {
	balBefore := suite.QueryBalances(suite.AccAddress())

	txId := suite.SendToExternal(1, coin)
	suite.Greater(txId, uint64(0))

	suite.SendCancelSendToExternal(txId)

	balAfter := suite.QueryBalances(suite.AccAddress())
	suite.Equal(balBefore.AmountOf(coin.Denom), balAfter.AmountOf(coin.Denom))
}

func (suite *CrosschainTestSuite) SendCancelSendToExternal(txId uint64) {
	suite.BroadcastTx(suite.privKey, &crosschaintypes.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        suite.AccAddress().String(),
		ChainName:     suite.chainName,
	})
}

func (suite *CrosschainTestSuite) SendIncreaseBridgeFee(txId uint64, bridgeFee sdk.Coin) {
	suite.BroadcastTx(suite.privKey, &crosschaintypes.MsgIncreaseBridgeFee{
		ChainName:     suite.chainName,
		TransactionId: txId,
		Sender:        suite.AccAddress().String(),
		AddBridgeFee:  bridgeFee,
	})
}

func (suite *CrosschainTestSuite) CheckIncreaseBridgeFee(sender sdk.AccAddress, txId uint64) {
	unbatchedTxs := suite.QueryPendingUnbatchedTx(sender)
	bridgeFee := sdk.ZeroInt()
	bridgeToken := ""
	for _, tx := range unbatchedTxs {
		if tx.Id != txId {
			continue
		}
		bridgeFee = tx.Fee.Amount
		bridgeToken = tx.Fee.Contract
	}
	suite.NotEmpty(bridgeToken)

	bridgeDenom := suite.GetBridgeDenomByToken(bridgeToken)

	addBridgeFee := sdkmath.NewInt(10)
	suite.SendIncreaseBridgeFee(txId, sdk.NewCoin(bridgeDenom, addBridgeFee))

	unbatchedTxs = suite.QueryPendingUnbatchedTx(sender)
	for _, tx := range unbatchedTxs {
		if tx.Id == txId {
			suite.Equal(tx.Fee.Amount, bridgeFee.Add(addBridgeFee))
			break
		}
	}
}

func (suite *CrosschainTestSuite) SendBatchRequest(minTxs uint64) {
	msgList := make([]sdk.Msg, 0)
	batchFeeResponse, err := suite.CrosschainQuery().BatchFees(suite.ctx, &crosschaintypes.QueryBatchFeeRequest{ChainName: suite.chainName})
	suite.NoError(err)
	suite.True(len(batchFeeResponse.BatchFees) >= 1)
	for _, batchToken := range batchFeeResponse.BatchFees {
		suite.Equal(batchToken.TotalTxs, minTxs)

		denomResponse, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
			Token:     batchToken.TokenContract,
			ChainName: suite.chainName,
		})
		suite.NoError(err)

		feeReceive := suite.HexAddress().String()
		if suite.chainName == trontypes.ModuleName {
			feeReceive = trontypes.AddressFromHex(feeReceive)
		}
		msgList = append(msgList, &crosschaintypes.MsgRequestBatch{
			Sender:     suite.BridgerAddr().String(),
			Denom:      denomResponse.Denom,
			MinimumFee: batchToken.TotalFees,
			FeeReceive: feeReceive,
			ChainName:  suite.chainName,
		})
	}
	suite.BroadcastTx(suite.bridgerPrivKey, msgList...)
}

func (suite *CrosschainTestSuite) SendConfirmBatch() {
	response, err := suite.CrosschainQuery().LastPendingBatchRequestByAddr(
		suite.ctx,
		&crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)
	suite.NotNil(response.Batch)

	outgoingTxBatch := response.Batch
	var signatureBytes []byte
	if suite.chainName == trontypes.ModuleName {
		checkpoint, err := trontypes.GetCheckpointConfirmBatch(outgoingTxBatch, suite.params.GravityId)
		suite.NoError(err)
		signatureBytes, err = trontypes.NewTronSignature(checkpoint, suite.externalPrivKey)
		suite.NoError(err)
		err = trontypes.ValidateTronSignature(checkpoint, signatureBytes, suite.ExternalAddr())
		suite.NoError(err)
	} else {
		checkpoint, err := outgoingTxBatch.GetCheckpoint(suite.params.GravityId)
		suite.NoError(err)
		signatureBytes, err = crosschaintypes.NewEthereumSignature(checkpoint, suite.externalPrivKey)
		suite.NoError(err)
		err = crosschaintypes.ValidateEthereumSignature(checkpoint, signatureBytes, suite.ExternalAddr())
		suite.NoError(err)
	}

	suite.BroadcastTx(suite.bridgerPrivKey,
		&crosschaintypes.MsgConfirmBatch{
			Nonce:           outgoingTxBatch.BatchNonce,
			TokenContract:   outgoingTxBatch.TokenContract,
			BridgerAddress:  suite.BridgerAddr().String(),
			ExternalAddress: suite.ExternalAddr(),
			Signature:       hex.EncodeToString(signatureBytes),
			ChainName:       suite.chainName,
		},
		&crosschaintypes.MsgSendToExternalClaim{
			EventNonce:     suite.queryFxLastEventNonce(),
			BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
			BatchNonce:     outgoingTxBatch.BatchNonce,
			TokenContract:  outgoingTxBatch.TokenContract,
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
}
