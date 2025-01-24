package crosschain_test

import (
	"encoding/hex"
	"math/big"
	"strings"
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

func (suite *CrosschainPrecompileTestSuite) TestContract_BridgeCall_NativeCoin() {
	symbol := helpers.NewRandSymbol()
	suite.AddBridgeToken(symbol, true)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	erc20Contract := suite.GetERC20Token(baseDenom).GetERC20Contract()
	suite.erc20TokenSuite.WithContract(erc20Contract)

	amount := sdkmath.NewInt(100)
	suite.AddNativeCoinToEVM(baseDenom, amount)

	feeAmount := big.NewInt(1)
	transferAmount := big.NewInt(2)

	approveAmount := transferAmount
	if !suite.IsCallPrecompile() {
		approveAmount = big.NewInt(0).Add(transferAmount, feeAmount)
	}
	suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, approveAmount)

	txResponse := suite.BridgeCall(suite.Ctx, nil,
		suite.signer.Address(), suite.NewBridgeCallArgs(erc20Contract, transferAmount))
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 3)

	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(big.NewInt(97), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeFeeColBalance := suite.erc20TokenSuite.BalanceOf(suite.Ctx,
		common.BytesToAddress(authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName)))
	suite.Equal(feeAmount, bridgeFeeColBalance)

	bridgeToken := suite.GetBridgeToken(baseDenom)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), bridgeCoin)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_BridgeCall_NativeERC20() {
	symbol := helpers.NewRandSymbol()

	erc20TokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, symbol)
	suite.AddBridgeToken(erc20TokenAddr.String(), false)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	amount := sdkmath.NewInt(100)
	suite.AddNativeERC20ToEVM(baseDenom, amount)

	feeAmount := big.NewInt(1)
	transferAmount := big.NewInt(2)

	approveAmount := transferAmount
	if !suite.IsCallPrecompile() {
		approveAmount = big.NewInt(0).Add(transferAmount, feeAmount)
	}
	suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, approveAmount)

	txResponse := suite.BridgeCall(suite.Ctx, nil,
		suite.signer.Address(), suite.NewBridgeCallArgs(erc20TokenAddr, transferAmount),
	)
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 3)

	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(big.NewInt(97), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(0))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeFeeColBalance := suite.erc20TokenSuite.BalanceOf(suite.Ctx,
		common.BytesToAddress(authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName)))
	suite.Equal(feeAmount, bridgeFeeColBalance)

	bridgeToken := suite.GetBridgeToken(baseDenom)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(2))
	suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_BridgeCall_IBCToken() {
	symbol := helpers.NewRandSymbol()

	suite.AddBridgeToken(symbol, true, true)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	erc20Contract := suite.GetERC20Token(baseDenom).GetERC20Contract()
	suite.erc20TokenSuite.WithContract(erc20Contract)

	amount := sdkmath.NewInt(100)
	suite.AddNativeCoinToEVM(baseDenom, amount, true)

	feeAmount := big.NewInt(1)
	transferAmount := big.NewInt(2)

	approveAmount := transferAmount
	if !suite.IsCallPrecompile() {
		approveAmount = big.NewInt(0).Add(transferAmount, feeAmount)
	}
	suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, approveAmount)

	txResponse := suite.BridgeCall(suite.Ctx, nil,
		suite.signer.Address(), suite.NewBridgeCallArgs(erc20Contract, transferAmount),
	)
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 3)

	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(big.NewInt(97), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeFeeColBalance := suite.erc20TokenSuite.BalanceOf(suite.Ctx,
		common.BytesToAddress(authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName)))
	suite.Equal(feeAmount, bridgeFeeColBalance)

	bridgeToken := suite.GetBridgeToken(baseDenom)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(2))
	suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)
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
