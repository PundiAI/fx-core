package crosschain_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestContract_BridgeCall_Input(t *testing.T) {
	bridgeCallABI := crosschain.NewBridgeCallABI()

	assert.Equal(t, `bridgeCall(string,address,address[],uint256[],address,bytes,uint256,uint256,bytes)`, bridgeCallABI.Method.Sig)
	assert.Equal(t, "payable", bridgeCallABI.Method.StateMutability)
	require.Len(t, bridgeCallABI.Method.Inputs, 9)
	require.Len(t, bridgeCallABI.Method.Outputs, 1)

	inputs := bridgeCallABI.Method.Inputs
	type Args struct {
		DstChain string
		Refund   common.Address
		Tokens   []common.Address
		Amounts  []*big.Int
		To       common.Address
		QuoteId  *big.Int
		GasLimit *big.Int
		Data     []byte
		Memo     []byte
	}
	args := Args{
		DstChain: "eth",
		Refund:   helpers.GenHexAddress(),
		Tokens: []common.Address{
			helpers.GenHexAddress(),
		},
		Amounts: []*big.Int{
			big.NewInt(1),
		},
		To:       helpers.GenHexAddress(),
		QuoteId:  big.NewInt(1),
		GasLimit: big.NewInt(1),
		Data:     []byte{1},
		Memo:     []byte{1},
	}
	inputData, err := inputs.Pack(
		args.DstChain,
		args.Refund,
		args.Tokens,
		args.Amounts,
		args.To,
		args.Data,
		args.QuoteId,
		args.GasLimit,
		args.Memo,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, inputData)

	inputValue, err := inputs.Unpack(inputData)
	require.NoError(t, err)
	assert.NotNil(t, inputValue)

	args2 := Args{}
	err = inputs.Copy(&args2, inputValue)
	require.NoError(t, err)

	assert.EqualValues(t, args, args2)
}

func TestContract_BridgeCall_Output(t *testing.T) {
	bridgeCallABI := crosschain.NewBridgeCallABI()
	assert.Len(t, bridgeCallABI.Method.Outputs, 1)

	outputs := bridgeCallABI.Method.Outputs
	eventNonce := big.NewInt(1)
	outputData, err := outputs.Pack(eventNonce)
	require.NoError(t, err)
	assert.NotEmpty(t, outputData)

	outputValue, err := outputs.Unpack(outputData)
	require.NoError(t, err)
	assert.NotNil(t, outputValue)

	assert.Equal(t, eventNonce, outputValue[0])
}

func TestContract_BridgeCall_Event(t *testing.T) {
	bridgeCallABI := crosschain.NewBridgeCallABI()

	assert.Equal(t, `BridgeCallEvent(address,address,address,address,uint256,string,address[],uint256[],bytes,uint256,uint256,bytes)`, bridgeCallABI.Event.Sig)
	assert.Equal(t, "0xcaa0e5b7ba998f542b3804184a5d30836451c57f6d1f031c466272e188f4a70f", bridgeCallABI.Event.ID.String())
	assert.Len(t, bridgeCallABI.Event.Inputs, 12)
	assert.Len(t, bridgeCallABI.Event.Inputs.NonIndexed(), 9)
	for i := 0; i < 3; i++ {
		assert.True(t, bridgeCallABI.Event.Inputs[i].Indexed)
	}
	inputs := bridgeCallABI.Event.Inputs

	args := contract.ICrosschainBridgeCallEvent{
		TxOrigin:   helpers.GenHexAddress(),
		EventNonce: big.NewInt(1),
		DstChain:   "eth",
		Tokens: []common.Address{
			helpers.GenHexAddress(),
		},
		Amounts: []*big.Int{
			big.NewInt(1),
		},
		Data:     []byte{1},
		QuoteId:  big.NewInt(1),
		GasLimit: big.NewInt(1),
		Memo:     []byte{1},
	}
	inputData, err := inputs.NonIndexed().Pack(
		args.TxOrigin,
		args.EventNonce,
		args.DstChain,
		args.Tokens,
		args.Amounts,
		args.Data,
		args.QuoteId,
		args.GasLimit,
		args.Memo,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, inputData)

	inputValue, err := inputs.Unpack(inputData)
	require.NoError(t, err)
	assert.NotNil(t, inputValue)

	var args2 contract.ICrosschainBridgeCallEvent
	err = inputs.Copy(&args2, inputValue)
	require.NoError(t, err)
	assert.EqualValues(t, args, args2)
}

func TestContract_BridgeCall_NewBridgeCallEvent(t *testing.T) {
	bridgeCallABI := crosschain.NewBridgeCallABI()

	sender := common.BytesToAddress([]byte{0x1})
	origin := common.BytesToAddress([]byte{0x2})
	nonce := big.NewInt(100)
	args := &contract.BridgeCallArgs{
		DstChain: "eth",
		Refund:   common.BytesToAddress([]byte{0x3}),
		Tokens:   []common.Address{common.BytesToAddress([]byte{0x4}), common.BytesToAddress([]byte{0x5})},
		Amounts:  []*big.Int{big.NewInt(123), big.NewInt(456)},
		To:       common.BytesToAddress([]byte{0x4}),
		Data:     []byte{0x1, 0x2, 0x3},
		QuoteId:  big.NewInt(100),
		GasLimit: big.NewInt(0),
		Memo:     []byte{0x1, 0x2, 0x3},
	}
	dataNew, topicNew, err := bridgeCallABI.NewBridgeCallEvent(args, sender, origin, nonce)
	require.NoError(t, err)
	expectData := "000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000260000000000000000000000000000000000000000000000000000000000000000365746800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000007b00000000000000000000000000000000000000000000000000000000000001c80000000000000000000000000000000000000000000000000000000000000003010203000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030102030000000000000000000000000000000000000000000000000000000000"
	require.EqualValues(t, expectData, hex.EncodeToString(dataNew))
	expectTopic := []common.Hash{
		common.HexToHash("0xcaa0e5b7ba998f542b3804184a5d30836451c57f6d1f031c466272e188f4a70f"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004"),
	}
	assert.EqualValues(t, expectTopic, topicNew)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_BridgeCall() {
	testCases := []struct {
		name                       string
		malleate                   func() *erc20types.ERC20Token
		feeAmount                  *big.Int
		transferAmount             *big.Int
		erc20ModuleAmount          sdkmath.Int // default base denom amount
		crosschainModuleAmount     sdkmath.Int // default bridge denom amount
		crosschainModuleBaseAmount sdkmath.Int
		chainNameAmount            sdkmath.Int // default bridge denom amount
	}{
		{
			name: "native coin",
			malleate: func() *erc20types.ERC20Token {
				bridgeToken := suite.AddBridgeToken(helpers.NewRandSymbol(), true)

				suite.Quote(bridgeToken.Denom)

				erc20Token := suite.GetERC20Token(bridgeToken.Denom)
				suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract())

				suite.AddNativeCoinToEVM(bridgeToken.Denom, sdkmath.NewInt(100))

				return erc20Token
			},
			feeAmount:                  big.NewInt(1),
			transferAmount:             big.NewInt(2),
			erc20ModuleAmount:          sdkmath.NewInt(98),
			crosschainModuleAmount:     sdkmath.NewInt(98),
			crosschainModuleBaseAmount: sdkmath.NewInt(0),
			chainNameAmount:            sdkmath.NewInt(0),
		},
		{
			name: "native erc20",
			malleate: func() *erc20types.ERC20Token {
				erc20TokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, helpers.NewRandSymbol())
				bridgeToken := suite.AddBridgeToken(erc20TokenAddr.String(), false)

				suite.Quote(bridgeToken.Denom)

				suite.AddNativeERC20ToEVM(bridgeToken.Denom, sdkmath.NewInt(100))

				return suite.GetERC20Token(bridgeToken.Denom)
			},
			feeAmount:                  big.NewInt(1),
			transferAmount:             big.NewInt(2),
			erc20ModuleAmount:          sdkmath.NewInt(0),
			crosschainModuleAmount:     sdkmath.NewInt(0),
			crosschainModuleBaseAmount: sdkmath.NewInt(2),
			chainNameAmount:            sdkmath.NewInt(2),
		},
		{
			name: "IBC Token",
			malleate: func() *erc20types.ERC20Token {
				bridgeToken := suite.AddBridgeToken(helpers.NewRandSymbol(), true, true)

				suite.Quote(bridgeToken.Denom)

				erc20Contract := suite.GetERC20Token(bridgeToken.Denom).GetERC20Contract()
				suite.erc20TokenSuite.WithContract(erc20Contract)

				suite.AddNativeCoinToEVM(bridgeToken.Denom, sdkmath.NewInt(100), true)

				return suite.GetERC20Token(bridgeToken.Denom)
			},
			feeAmount:                  big.NewInt(1),
			transferAmount:             big.NewInt(2),
			erc20ModuleAmount:          sdkmath.NewInt(98),
			crosschainModuleAmount:     sdkmath.NewInt(0),
			crosschainModuleBaseAmount: sdkmath.NewInt(2),
			chainNameAmount:            sdkmath.NewInt(2),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			erc20Token := tc.malleate()

			accountsBalances := suite.App.BankKeeper.GetAccountsBalances(suite.Ctx)

			// ===========> Test Precompile BridgeCall

			approveAmount := tc.transferAmount
			if !suite.IsCallPrecompile() {
				approveAmount = big.NewInt(0).Add(tc.transferAmount, tc.feeAmount)
			}
			suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, approveAmount)

			bridgeCallArgs := suite.NewBridgeCallArgs(erc20Token.GetERC20Contract(), tc.transferAmount)
			txResponse := suite.BridgeCall(suite.Ctx, nil,
				suite.signer.Address(), bridgeCallArgs)
			suite.NotNil(txResponse)
			suite.GreaterOrEqual(len(txResponse.Logs), 3)
			bridgeCallEvent, err := crosschain.NewBridgeCallABI().UnpackEvent(txResponse.Logs[len(txResponse.Logs)-1].ToEthereum())
			suite.Require().NoError(err)
			suite.Equal(big.NewInt(1), bridgeCallEvent.EventNonce)
			suite.Equal(suite.GetSender(), bridgeCallEvent.Sender)
			suite.Equal(suite.signer.Address(), bridgeCallEvent.TxOrigin)
			suite.Equal(bridgeCallArgs.Refund, bridgeCallEvent.Refund)
			suite.Equal(bridgeCallArgs.To, bridgeCallEvent.To)
			suite.Equal(bridgeCallArgs.DstChain, bridgeCallEvent.DstChain)
			suite.Equal(bridgeCallArgs.Tokens, bridgeCallEvent.Tokens)
			suite.Equal(bridgeCallArgs.Amounts, bridgeCallEvent.Amounts)
			suite.Equal(bridgeCallArgs.Data, bridgeCallEvent.Data)
			suite.Equal(bridgeCallArgs.QuoteId, bridgeCallEvent.QuoteId)
			suite.Equal(bridgeCallArgs.GasLimit.String(), bridgeCallEvent.GasLimit.String())
			suite.Equal(bridgeCallArgs.Memo, bridgeCallEvent.Memo)

			balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
			suite.Equal(big.NewInt(97), balance)

			bridgeFeeColAddr := common.BytesToAddress(authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName))
			bridgeFeeColBalance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, bridgeFeeColAddr)
			suite.Equal(tc.feeAmount, bridgeFeeColBalance)

			bridgeToken := suite.GetBridgeToken(erc20Token.Denom)
			bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), tc.crosschainModuleAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), bridgeCoin)

			baseCoin := sdk.NewCoin(bridgeToken.Denom, tc.crosschainModuleBaseAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), baseCoin)

			bridgeCoin = sdk.NewCoin(bridgeToken.BridgeDenom(), tc.chainNameAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)

			baseCoin = sdk.NewCoin(erc20Token.Denom, sdkmath.NewInt(0))
			suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), baseCoin)

			bridgeCoin = sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(0))
			suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), bridgeCoin)

			baseCoin = sdk.NewCoin(erc20Token.Denom, tc.erc20ModuleAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

			// ===========> Test MsgBridgeCallResultClaim

			keeper := suite.App.CrosschainKeepers.GetKeeper(suite.chainName)
			outgoingBridgeCall, found := keeper.GetOutgoingBridgeCallByNonce(suite.Ctx, bridgeCallEvent.EventNonce.Uint64())
			suite.True(found)

			quoteInfo, found := keeper.GetOutgoingBridgeCallQuoteInfo(suite.Ctx, bridgeCallEvent.EventNonce.Uint64())
			suite.True(found)

			bridgeCallResultClaim := &crosschaintypes.MsgBridgeCallResultClaim{
				ChainName:      suite.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,   // mock fx bridge contract event nonce
				BlockHeight:    100, // mock block height
				Nonce:          outgoingBridgeCall.Nonce,
				TxOrigin:       helpers.GenExternalAddr(suite.chainName),
				Success:        false,
				Cause:          "",
			}
			suite.executeClaim(bridgeCallResultClaim)

			_, found = keeper.GetOutgoingBridgeCallByNonce(suite.Ctx, bridgeCallEvent.EventNonce.Uint64())
			suite.False(found)

			_, found = keeper.GetOutgoingBridgeCallQuoteInfo(suite.Ctx, bridgeCallEvent.EventNonce.Uint64())
			suite.False(found)

			bridgeFeeColBalance = suite.erc20TokenSuite.BalanceOf(suite.Ctx, bridgeFeeColAddr)
			suite.Equal(big.NewInt(0).String(), bridgeFeeColBalance.String())

			balance = suite.erc20TokenSuite.BalanceOf(suite.Ctx, common.HexToAddress(quoteInfo.Oracle))
			suite.Equal(tc.feeAmount.String(), balance.String())

			if !bridgeCallResultClaim.Success {
				balance = suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
				suite.Equal(big.NewInt(99), balance)

				suite.Equal(accountsBalances, suite.App.BankKeeper.GetAccountsBalances(suite.Ctx))
			}
		})
	}
}

func (suite *CrosschainPrecompileTestSuite) TestContract_BridgeCall_OriginToken() {
	suite.AddBridgeToken(fxtypes.DefaultSymbol, false)

	suite.Quote(fxtypes.DefaultDenom)

	balance := suite.Balance(suite.signer.AccAddress())

	value := big.NewInt(2)
	if !suite.IsCallPrecompile() {
		value = big.NewInt(3) // add fee
	}
	txResponse := suite.BridgeCall(suite.Ctx, value,
		suite.signer.Address(), suite.NewBridgeCallArgs(common.Address{}, nil))
	suite.NotNil(txResponse)
	suite.Len(txResponse.Logs, 1)

	suite.AssertBalance(suite.signer.AccAddress(), balance.Sub(helpers.NewStakingCoin(3, 0))...)
	suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName), helpers.NewStakingCoin(1, 0))
	suite.AssertBalance(authtypes.NewModuleAddress(ethtypes.ModuleName), helpers.NewStakingCoin(2, 0))
}
