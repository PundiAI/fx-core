package precompile_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/precompile"
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
