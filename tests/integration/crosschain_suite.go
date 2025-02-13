package integration

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	"github.com/pundiai/fx-core/v8/x/eth/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

type CrosschainSuite struct {
	*FxCoreSuite

	chainName string
	params    crosschaintypes.Params
	oracle    *helpers.Signer
	bridger   *helpers.Signer
	external  *ecdsa.PrivateKey

	signer *helpers.Signer

	contractAddr common.Address
	crosschain   *contract.ICrosschain
}

func NewCrosschainSuite(chainName string, suite *FxCoreSuite) *CrosschainSuite {
	externalPrivKey, err := ethcrypto.GenerateKey()
	suite.Require().NoError(err)
	contractAddr := common.HexToAddress(contract.CrosschainAddress)
	crosschain, err := contract.NewICrosschain(contractAddr, suite.ethCli)
	suite.Require().NoError(err)
	return &CrosschainSuite{
		FxCoreSuite:  suite,
		chainName:    chainName,
		oracle:       helpers.NewSigner(helpers.NewEthPrivKey()),
		bridger:      helpers.NewSigner(helpers.NewEthPrivKey()),
		external:     externalPrivKey,
		contractAddr: contractAddr,
		crosschain:   crosschain,
	}
}

func (suite *CrosschainSuite) Init() {
	suite.Send(suite.OracleAddr(), suite.NewCoin(sdkmath.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.BridgerAddr(), suite.NewCoin(sdkmath.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.signer.AccAddress(), suite.NewCoin(sdkmath.NewInt(1_000).MulRaw(1e18)))
	suite.params = suite.QueryParams()
}

func (suite *CrosschainSuite) OracleAddr() sdk.AccAddress {
	return suite.oracle.AccAddress()
}

func (suite *CrosschainSuite) ExternalAddr() string {
	address := ethcrypto.PubkeyToAddress(suite.external.PublicKey)
	return fxtypes.ExternalAddrToStr(suite.chainName, address.Bytes())
}

func (suite *CrosschainSuite) BridgerAddr() sdk.AccAddress {
	return suite.bridger.AccAddress()
}

func (suite *CrosschainSuite) HexAddressString() string {
	return fxtypes.ExternalAddrToStr(suite.chainName, suite.signer.Address().Bytes())
}

func (suite *CrosschainSuite) CrosschainQuery() crosschaintypes.QueryClient {
	return suite.grpcCli.CrosschainQuery()
}

func (suite *CrosschainSuite) QueryParams() crosschaintypes.Params {
	response, err := suite.CrosschainQuery().Params(suite.ctx,
		&crosschaintypes.QueryParamsRequest{ChainName: suite.chainName})
	suite.Require().NoError(err)
	return response.Params
}

func (suite *CrosschainSuite) queryFxLastEventNonce() uint64 {
	lastEventNonce, err := suite.CrosschainQuery().LastEventNonceByAddr(suite.ctx,
		&crosschaintypes.QueryLastEventNonceByAddrRequest{
			ChainName:      suite.chainName,
			BridgerAddress: suite.BridgerAddr().String(),
		},
	)
	suite.Require().NoError(err)
	return lastEventNonce.EventNonce + 1
}

func (suite *CrosschainSuite) queryObserverExternalBlockHeight() uint64 {
	response, err := suite.CrosschainQuery().LastObservedBlockHeight(suite.ctx,
		&crosschaintypes.QueryLastObservedBlockHeightRequest{
			ChainName: suite.chainName,
		},
	)
	suite.Require().NoError(err)
	return response.ExternalBlockHeight
}

func (suite *CrosschainSuite) AddBridgeTokenClaim(name, symbol string, decimals uint64, token string) string {
	response, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.ErrorContains(err, "code = NotFound desc = bridge token")
	suite.Nil(response)

	suite.BroadcastTx(suite.bridger, &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: suite.BridgerAddr().String(),
		ChainName:      suite.chainName,
	})

	response, err = suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.Require().NoError(err)
	if response.Denom != fxtypes.DefaultDenom {
		suite.Equal(crosschaintypes.NewBridgeDenom(suite.chainName, token), response.Denom)
	}

	return response.Denom
}

func (suite *CrosschainSuite) GetBridgeTokens() (denoms []erc20types.BridgeToken) {
	response, err := suite.CrosschainQuery().BridgeTokens(suite.ctx, &crosschaintypes.QueryBridgeTokensRequest{
		ChainName: suite.chainName,
	})
	suite.Require().NoError(err)
	return response.BridgeTokens
}

func (suite *CrosschainSuite) GetBridgeDenomByToken(token string) (denom string) {
	response, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(response.Denom)
	return response.Denom
}

func (suite *CrosschainSuite) GetBridgeTokenByDenom(denom string) (token string) {
	response, err := suite.CrosschainQuery().DenomToToken(suite.ctx, &crosschaintypes.QueryDenomToTokenRequest{
		ChainName: suite.chainName,
		Denom:     denom,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(response.Token)
	return response.Token
}

func (suite *CrosschainSuite) BondedOracle() {
	response, err := suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Require().Error(err, crosschaintypes.ErrNoFoundOracle)
	suite.Nil(response)

	txResponse := suite.BroadcastTx(suite.oracle, &crosschaintypes.MsgBondedOracle{
		OracleAddress:    suite.OracleAddr().String(),
		BridgerAddress:   suite.BridgerAddr().String(),
		ExternalAddress:  suite.ExternalAddr(),
		ValidatorAddress: suite.GetValAddr().String(),
		DelegateAmount:   suite.params.DelegateThreshold,
		ChainName:        suite.chainName,
	})

	response, err = suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Require().NoError(err)
	suite.Equal(crosschaintypes.Oracle{
		OracleAddress:     suite.OracleAddr().String(),
		BridgerAddress:    suite.BridgerAddr().String(),
		ExternalAddress:   suite.ExternalAddr(),
		DelegateAmount:    suite.params.DelegateThreshold.Amount,
		StartHeight:       txResponse.Height,
		Online:            true,
		DelegateValidator: suite.GetValAddr().String(),
		SlashTimes:        0,
	}, *response.Oracle)
}

func (suite *CrosschainSuite) SendUpdateChainOraclesProposal() (*sdk.TxResponse, uint64) {
	msg := &crosschaintypes.MsgUpdateChainOracles{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Oracles:   []string{suite.OracleAddr().String()},
		ChainName: suite.chainName,
	}
	return suite.BroadcastProposalTxV1(msg)
}

func (suite *CrosschainSuite) SendOracleSetConfirm() {
	queryResponse, err := suite.CrosschainQuery().LastPendingOracleSetRequestByAddr(suite.ctx,
		&crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Require().NoError(err)
	suite.NotEmpty(queryResponse.OracleSets)

	for _, oracleSet := range queryResponse.OracleSets {
		checkpoint, err := oracleSet.GetCheckpoint(suite.params.GravityId)
		suite.Require().NoError(err)

		var signature []byte
		if suite.chainName == trontypes.ModuleName {
			signature, err = trontypes.NewTronSignature(checkpoint, suite.external)
			suite.Require().NoError(err)
			err = trontypes.ValidateTronSignature(checkpoint, signature, suite.ExternalAddr())
			suite.Require().NoError(err)
		} else {
			signature, err = types.NewEthereumSignature(checkpoint, suite.external)
			suite.Require().NoError(err)
			err = types.ValidateEthereumSignature(checkpoint, signature, suite.ExternalAddr())
			suite.Require().NoError(err)
		}

		suite.BroadcastTx(suite.bridger, &crosschaintypes.MsgOracleSetConfirm{
			Nonce:           oracleSet.Nonce,
			BridgerAddress:  suite.BridgerAddr().String(),
			ExternalAddress: suite.ExternalAddr(),
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		})
	}
}

func (suite *CrosschainSuite) BridgeCallClaim(to string, tokens []string, amounts []sdkmath.Int) {
	suite.BroadcastTx(suite.bridger, &crosschaintypes.MsgBridgeCallClaim{
		ChainName:      suite.chainName,
		Sender:         suite.HexAddressString(),
		Refund:         suite.HexAddressString(),
		To:             to,
		QuoteId:        sdkmath.ZeroInt(),
		TokenContracts: tokens,
		Amounts:        amounts,
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		BridgerAddress: suite.BridgerAddr().String(),
		TxOrigin:       suite.HexAddressString(),
	})
	suite.ExecuteClaim()
}

func (suite *CrosschainSuite) SendToFxClaim(receiver sdk.AccAddress, token string, amount sdkmath.Int) {
	beforeBalance := suite.GetAllBalances(suite.signer.AccAddress())

	suite.BroadcastTx(suite.bridger, &crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Amount:         amount,
		Sender:         suite.HexAddressString(),
		Receiver:       receiver.String(),
		BridgerAddress: suite.BridgerAddr().String(),
		ChainName:      suite.chainName,
	})
	suite.ExecuteClaim()
	bridgeToken, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.Require().NoError(err)
	if bridgeToken.Denom == fxtypes.DefaultDenom {
		balances := suite.GetAllBalances(receiver)
		suite.Require().Equal(beforeBalance.Add(sdk.NewCoin(fxtypes.DefaultDenom, amount)), balances)
	} else {
		afterBalance := suite.GetAllBalances(suite.signer.AccAddress())
		suite.Equal(afterBalance, beforeBalance)
	}
}

func (suite *CrosschainSuite) SendToExternal(count int, amount sdk.Coin) (*sdk.TxResponse, uint64) {
	msgList := make([]sdk.Msg, 0, count)
	for i := 0; i < count; i++ {
		msgList = append(msgList, &crosschaintypes.MsgSendToExternal{
			Sender:    suite.signer.AccAddress().String(),
			Dest:      suite.HexAddressString(),
			Amount:    amount.SubAmount(sdkmath.NewInt(1)),
			BridgeFee: sdk.NewCoin(amount.Denom, sdkmath.NewInt(1)),
			ChainName: suite.chainName,
		})
	}
	txResponse := suite.BroadcastTx(suite.signer, msgList...)
	for _, eventLog := range txResponse.Logs {
		for _, event := range eventLog.Events {
			if event.Type != crosschaintypes.EventTypeSendToExternal {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != crosschaintypes.AttributeKeyOutgoingTxID &&
					attribute.Key != crosschaintypes.AttributeKeyPendingOutgoingTxID {
					continue
				}
				txId, err := strconv.ParseUint(attribute.Value, 10, 64)
				suite.Require().NoError(err)
				return txResponse, txId
			}
		}
	}
	return txResponse, 0
}

func (suite *CrosschainSuite) SendToExternalAndCancel(coin sdk.Coin) {
	balBefore := suite.GetAllBalances(suite.signer.AccAddress())

	_, txId := suite.SendToExternal(1, coin)
	suite.Greater(txId, uint64(0))

	balAfter := suite.GetAllBalances(suite.signer.AccAddress())
	suite.Equal(balBefore.AmountOf(coin.Denom), balAfter.AmountOf(coin.Denom))
}

func (suite *CrosschainSuite) SendConfirmBatch() {
	response, err := suite.CrosschainQuery().LastPendingBatchRequestByAddr(
		suite.ctx,
		&crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Require().NoError(err)
	suite.NotNil(response.GetBatches())

	for _, outgoingTxBatch := range response.GetBatches() {
		checkpoint, err := outgoingTxBatch.GetCheckpoint(suite.params.GravityId)
		suite.Require().NoError(err)
		signatureBytes := suite.SignatureCheckpoint(checkpoint)

		suite.BroadcastTx(suite.bridger,
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
}

func (suite *CrosschainSuite) SendToExternalAndConfirm(coin sdk.Coin) {
	suite.SendToExternal(1, coin)
	suite.SendConfirmBatch()
}

func (suite *CrosschainSuite) FormatAddress(address common.Address) string {
	return fxtypes.ExternalAddrToStr(suite.chainName, address.Bytes())
}

func (suite *CrosschainSuite) BridgeCallConfirm(nonce uint64, isSuccess bool) {
	bridgeCall := suite.QueryBridgeCallByNonce(nonce)
	checkpoint, err := bridgeCall.GetCheckpoint(suite.params.GravityId)
	suite.Require().NoError(err)
	signatureBytes := suite.SignatureCheckpoint(checkpoint)

	suite.BroadcastTx(suite.bridger,
		&crosschaintypes.MsgBridgeCallConfirm{
			Nonce:           nonce,
			BridgerAddress:  suite.BridgerAddr().String(),
			ExternalAddress: suite.ExternalAddr(),
			Signature:       hex.EncodeToString(signatureBytes),
			ChainName:       suite.chainName,
		},
	)
	suite.BroadcastTx(suite.bridger,
		&crosschaintypes.MsgBridgeCallResultClaim{
			ChainName:      suite.chainName,
			BridgerAddress: suite.BridgerAddr().String(),
			EventNonce:     suite.queryFxLastEventNonce(),
			BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
			Nonce:          nonce,
			TxOrigin:       suite.ExternalAddr(),
			Success:        isSuccess,
			Cause:          "",
		},
	)
	suite.ExecuteClaim()
}

func (suite *CrosschainSuite) SignatureCheckpoint(checkpoint []byte) []byte {
	var signatureBytes []byte
	var err error
	if suite.chainName == trontypes.ModuleName {
		signatureBytes, err = trontypes.NewTronSignature(checkpoint, suite.external)
		suite.Require().NoError(err)
		suite.Require().NoError(trontypes.ValidateTronSignature(checkpoint, signatureBytes, suite.ExternalAddr()))
	} else {
		signatureBytes, err = types.NewEthereumSignature(checkpoint, suite.external)
		suite.Require().NoError(err)
		suite.Require().NoError(types.ValidateEthereumSignature(checkpoint, signatureBytes, suite.ExternalAddr()))
	}
	return signatureBytes
}

func (suite *CrosschainSuite) QueryBridgeCallByNonce(nonce uint64) *crosschaintypes.OutgoingBridgeCall {
	response, err := suite.CrosschainQuery().BridgeCallByNonce(suite.ctx, &crosschaintypes.QueryBridgeCallByNonceRequest{
		ChainName: suite.chainName,
		Nonce:     nonce,
	})
	suite.Require().NoError(err)
	return response.GetBridgeCall()
}

func (suite *CrosschainSuite) ExecuteClaim() *ethtypes.Transaction {
	externalClaims := suite.PendingExecuteClaim()
	suite.Require().True(len(externalClaims) > 0)

	ethTx, err := suite.crosschain.ExecuteClaim(suite.TransactOpts(suite.signer),
		suite.chainName, new(big.Int).SetUint64(externalClaims[0].GetEventNonce()))
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	return ethTx
}

func (suite *CrosschainSuite) PendingExecuteClaim() []crosschaintypes.ExternalClaim {
	response, err := suite.CrosschainQuery().PendingExecuteClaim(suite.ctx, &crosschaintypes.QueryPendingExecuteClaimRequest{
		ChainName: suite.chainName,
	})
	suite.Require().NoError(err)
	externalClaims := make([]crosschaintypes.ExternalClaim, 0, len(response.Claims))
	for _, claim := range response.Claims {
		var externalClaim crosschaintypes.ExternalClaim
		err = suite.codec.UnpackAny(claim, &externalClaim)
		suite.Require().NoError(err)
		externalClaims = append(externalClaims, externalClaim)
	}
	return externalClaims
}

func (suite *CrosschainSuite) UpdateParams(opts ...func(params *crosschaintypes.Params)) (*sdk.TxResponse, uint64) {
	params := suite.QueryParams()
	for _, opt := range opts {
		opt(&params)
	}
	suite.params = params
	msg := &crosschaintypes.MsgUpdateParams{
		ChainName: suite.chainName,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    params,
	}
	return suite.BroadcastProposalTxV1(msg)
}

func (suite *CrosschainSuite) Crosschain(token common.Address, recipient string, amount, fee *big.Int, target string) *ethtypes.Transaction {
	erc20TokenSuite := NewERC20TokenSuite(suite.EthSuite, token, suite.signer)

	erc20TokenSuite.Approve(suite.contractAddr, big.NewInt(0).Add(amount, fee))

	beforeBalanceOf := erc20TokenSuite.BalanceOf(suite.signer.Address())

	ethTx, err := suite.crosschain.CrossChain(suite.TransactOpts(suite.signer), token, recipient, amount, fee, contract.MustStrToByte32(target), "")
	suite.Require().NoError(err)
	suite.WaitMined(ethTx)

	afterBalanceOf := erc20TokenSuite.BalanceOf(suite.signer.Address())
	suite.Require().True(new(big.Int).Sub(beforeBalanceOf, afterBalanceOf).Cmp(new(big.Int).Add(amount, fee)) == 0)
	return ethTx
}
