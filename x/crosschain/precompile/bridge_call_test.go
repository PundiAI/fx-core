package precompile_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func TestBridgeCallABI(t *testing.T) {
	bridgeCall := precompile.NewBridgeCallMethod(nil)

	require.Equal(t, 8, len(bridgeCall.Method.Inputs))
	require.Equal(t, 1, len(bridgeCall.Method.Outputs))
}

func TestContract_BridgeCall_Input(t *testing.T) {
	bridgeCall := precompile.NewBridgeCallMethod(nil)

	assert.Equal(t, `bridgeCall(string,address,address[],uint256[],address,bytes,uint256,bytes)`, bridgeCall.Method.Sig)
	assert.Equal(t, "payable", bridgeCall.Method.StateMutability)
	assert.Equal(t, 8, len(bridgeCall.Method.Inputs))

	inputs := bridgeCall.Method.Inputs
	type Args struct {
		DstChain string
		Refund   common.Address
		Tokens   []common.Address
		Amounts  []*big.Int
		To       common.Address
		Data     []byte
		Value    *big.Int
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
		To:    helpers.GenHexAddress(),
		Data:  []byte{1},
		Value: big.NewInt(1),
		Memo:  []byte{1},
	}
	inputData, err := inputs.Pack(
		args.DstChain,
		args.Refund,
		args.Tokens,
		args.Amounts,
		args.To,
		args.Data,
		args.Value,
		args.Memo,
	)
	assert.NoError(t, err)
	assert.True(t, len(inputData) > 0)

	inputValue, err := inputs.Unpack(inputData)
	assert.NoError(t, err)
	assert.NotNil(t, inputValue)

	args2 := Args{}
	err = inputs.Copy(&args2, inputValue)
	assert.NoError(t, err)

	assert.EqualValues(t, args, args2)
}

func TestContract_BridgeCall_Output(t *testing.T) {
	bridgeCall := precompile.NewBridgeCallMethod(nil)
	assert.Equal(t, 1, len(bridgeCall.Method.Outputs))

	outputs := bridgeCall.Method.Outputs
	eventNonce := big.NewInt(1)
	outputData, err := outputs.Pack(eventNonce)
	assert.NoError(t, err)
	assert.True(t, len(outputData) > 0)

	outputValue, err := outputs.Unpack(outputData)
	assert.NoError(t, err)
	assert.NotNil(t, outputValue)

	assert.Equal(t, eventNonce, outputValue[0])
}

func TestContract_BridgeCall_Event(t *testing.T) {
	bridgeCall := precompile.NewBridgeCallMethod(nil)

	assert.Equal(t, `BridgeCallEvent(address,address,address,address,uint256,uint256,string,address[],uint256[],bytes,bytes)`, bridgeCall.Event.Sig)
	assert.Equal(t, "0x4a9b24da6150ef33e7c41038842b7c94fe89a4fff22dccb2c3fd79f0176062c6", bridgeCall.Event.ID.String())
	assert.Equal(t, 11, len(bridgeCall.Event.Inputs))
	assert.Equal(t, 8, len(bridgeCall.Event.Inputs.NonIndexed()))
	for i := 0; i < 3; i++ {
		assert.Equal(t, true, bridgeCall.Event.Inputs[i].Indexed)
	}
	inputs := bridgeCall.Event.Inputs

	args := contract.ICrossChainBridgeCallEvent{
		TxOrigin:   helpers.GenHexAddress(),
		Value:      big.NewInt(1),
		EventNonce: big.NewInt(1),
		DstChain:   "eth",
		Tokens: []common.Address{
			helpers.GenHexAddress(),
		},
		Amounts: []*big.Int{
			big.NewInt(1),
		},
		Data: []byte{1},
		Memo: []byte{1},
	}
	inputData, err := inputs.NonIndexed().Pack(
		args.TxOrigin,
		args.Value,
		args.EventNonce,
		args.DstChain,
		args.Tokens,
		args.Amounts,
		args.Data,
		args.Memo,
	)
	assert.NoError(t, err)
	assert.True(t, len(inputData) > 0)

	inputValue, err := inputs.Unpack(inputData)
	assert.NoError(t, err)
	assert.NotNil(t, inputValue)

	var args2 contract.ICrossChainBridgeCallEvent
	err = inputs.Copy(&args2, inputValue)
	assert.NoError(t, err)
	assert.EqualValues(t, args, args2)
}

func TestContract_BridgeCall_NewBridgeCallEvent(t *testing.T) {
	bridgeCall := precompile.NewBridgeCallMethod(nil)

	sender := common.BytesToAddress([]byte{0x1})
	origin := common.BytesToAddress([]byte{0x2})
	nonce := big.NewInt(100)
	args := &types.BridgeCallArgs{
		DstChain: "eth",
		Refund:   common.BytesToAddress([]byte{0x3}),
		Tokens:   []common.Address{common.BytesToAddress([]byte{0x4}), common.BytesToAddress([]byte{0x5})},
		Amounts:  []*big.Int{big.NewInt(123), big.NewInt(456)},
		To:       common.BytesToAddress([]byte{0x4}),
		Data:     []byte{0x1, 0x2, 0x3},
		Value:    big.NewInt(100),
		Memo:     []byte{0x1, 0x2, 0x3},
	}
	dataNew, topicNew, err := bridgeCall.NewBridgeCallEvent(args, sender, origin, nonce)
	require.NoError(t, err)
	expectData := "0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000000365746800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000007b00000000000000000000000000000000000000000000000000000000000001c80000000000000000000000000000000000000000000000000000000000000003010203000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030102030000000000000000000000000000000000000000000000000000000000"
	require.EqualValues(t, expectData, hex.EncodeToString(dataNew))
	expectTopic := []common.Hash{
		common.HexToHash("4a9b24da6150ef33e7c41038842b7c94fe89a4fff22dccb2c3fd79f0176062c6"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
		common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004"),
	}
	assert.EqualValues(t, expectTopic, topicNew)
}
