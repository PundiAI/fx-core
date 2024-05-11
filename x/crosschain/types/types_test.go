package types

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOracleSet_Checkpoint(t *testing.T) {
	bridgeValidators := BridgeValidators{
		{
			Power:           6667,
			ExternalAddress: "0x4D449a236D2Ff82C7De7E394c713517Ab2956995",
		},
	}
	oracleSet := NewOracleSet(1, 100, bridgeValidators)

	ourHash, err := oracleSet.GetCheckpoint("gravity-id")
	require.NoError(t, err)

	expected := "0xf9b58b26674421b9eca385c26ca3953405519edb34de3a971faf8e7bb0a3bb91"[2:]
	assert.Equal(t, expected, hex.EncodeToString(ourHash))
}

func TestOutgoingTxBatch_Checkpoint(t *testing.T) {
	senderAddr, err := sdk.AccAddressFromHexUnsafe("9738AB195F6FC972A9014ADB77E4EB6D7C32FDA8")
	require.NoError(t, err)
	erc20Addr := common.HexToAddress("0x8c15Ef5b4B21951d50E53E4fbdA8298FFAD25057")

	outgoingTxBatch := OutgoingTxBatch{
		BatchNonce:   1,
		BatchTimeout: 100,
		Transactions: []*OutgoingTransferTx{
			{
				Id:          0x1,
				Sender:      senderAddr.String(),
				DestAddress: "0x4D449a236D2Ff82C7De7E394c713517Ab2956995",
				Token:       NewERC20Token(sdkmath.NewInt(0x1), erc20Addr.String()),
				Fee:         NewERC20Token(sdkmath.NewInt(0x1), erc20Addr.String()),
			},
		},
		TokenContract: erc20Addr.String(),
		Block:         10,
		FeeReceive:    "0x9738Ab195f6fC972A9014ADb77e4eB6D7c32fDA8",
	}

	checkpoint, err := outgoingTxBatch.GetCheckpoint("gravity-id")
	require.NoError(t, err)

	expected := "0x4306502ddd8b91ec4f5db7485271c9508df1b081d0d7e40eefa202547afbacf4"[2:]
	assert.Equal(t, expected, hex.EncodeToString(checkpoint))
}

func TestOutgoingBridgeCall_Checkpoint(t *testing.T) {
	outgoingBridgeCall := OutgoingBridgeCall{
		Sender:   "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
		Receiver: "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
		Tokens: []ERC20Token{
			{
				Contract: "0x1429859428C0aBc9C2C47C8Ee9FBaf82cFA0F20f",
				Amount:   sdkmath.NewInt(1000),
			},
		},
		To:          "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
		Data:        "38ed17390000000000000000000000000000000000000000000000002411f07601e50822000000000000000000000000000000000000000000000000ffbc5effd7093f8200000000000000000000000000000000000000000000000000000000000000a00000000000000000000000009527be6aafdb188f57ba2827e3acca484b4e4d9b00000000000000000000000000000000000000000000000000000000663f6553000000000000000000000000000000000000000000000000000000000000000200000000000000000000000080b5a32e4f032b2a058b4f29ec95eefeeb87adcd00000000000000000000000010aeaec2072f36de9b55fa20d6de7a44b7195697",
		Memo:        "8803dbee0000000000000000000000000000000000000000000000008b9c699b98fca89b0000000000000000000000000000000000000000000000004252f728ec196cb100000000000000000000000000000000000000000000000000000000000000a00000000000000000000000006bd6b0d09c9b90553449b7c86ea46cfa7989baaa00000000000000000000000000000000000000000000000000000000663f6529000000000000000000000000000000000000000000000000000000000000000200000000000000000000000080b5a32e4f032b2a058b4f29ec95eefeeb87adcd000000000000000000000000c8b4d3e67238e38b20d38908646ff6f4f48de5ec",
		Nonce:       1,
		Timeout:     1067,
		BlockHeight: 0,
	}
	checkpoint, err := outgoingBridgeCall.GetCheckpoint("eth-fxcore")
	require.NoError(t, err)

	expected := "0xe7bbfb4bd9c8baf2f0c9c91868cc897a1f042797e5fa7e3fdf918180202ddfb9"[2:]
	assert.Equal(t, expected, hex.EncodeToString(checkpoint))
}

func TestBridgeValidators_PowerDiff(t *testing.T) {
	specs := map[string]struct {
		start BridgeValidators
		diff  BridgeValidators
		exp   float64
	}{
		"no diff": {
			start: BridgeValidators{
				{Power: 1, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 2, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 3, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			diff: BridgeValidators{
				{Power: 1, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 2, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 3, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: 0.0,
		},
		"one": {
			start: BridgeValidators{
				{Power: 1073741823, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 1073741823, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 2147483646, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			diff: BridgeValidators{
				{Power: 858993459, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 858993459, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 2576980377, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: 0.2,
		},
		"real world": {
			start: BridgeValidators{
				{Power: 678509841, ExternalAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 685294939, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 671724742, ExternalAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, ExternalAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 617443955, ExternalAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 6785098, ExternalAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
				{Power: 291759231, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			diff: BridgeValidators{
				{Power: 642345266, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 678509841, ExternalAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, ExternalAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, ExternalAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 671724742, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 617443955, ExternalAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 291759231, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
				{Power: 6785098, ExternalAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
			},
			exp: 0.010000000011641532,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			assert.Equal(t, spec.exp, spec.start.PowerDiff(spec.diff))
		})
	}
}

func TestBridgeValidators_Sort(t *testing.T) {
	address1 := common.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()
	address2 := common.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()
	address3 := common.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()

	specs := map[string]struct {
		src BridgeValidators
		exp BridgeValidators
	}{
		"by power desc": {
			src: BridgeValidators{
				{Power: 1, ExternalAddress: address3},
				{Power: 2, ExternalAddress: address1},
				{Power: 3, ExternalAddress: address2},
			},
			exp: BridgeValidators{
				{Power: 3, ExternalAddress: address2},
				{Power: 2, ExternalAddress: address1},
				{Power: 1, ExternalAddress: address3},
			},
		},
		"by eth addr on same power": {
			src: BridgeValidators{
				{Power: 1, ExternalAddress: address2},
				{Power: 1, ExternalAddress: address1},
				{Power: 1, ExternalAddress: address3},
			},
			exp: BridgeValidators{
				{Power: 1, ExternalAddress: address1},
				{Power: 1, ExternalAddress: address2},
				{Power: 1, ExternalAddress: address3},
			},
		},
		// if you're thinking about changing this due to a change in the sorting algorithm
		// you MUST go change this in gravity_utils/types.rs as well. You will also break all
		// bridges in production when they try to migrate so use extreme caution!
		"real world": {
			src: BridgeValidators{
				{Power: 678509841, ExternalAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 685294939, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 671724742, ExternalAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, ExternalAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 617443955, ExternalAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 6785098, ExternalAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
				{Power: 291759231, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: BridgeValidators{
				{Power: 685294939, ExternalAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 678509841, ExternalAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, ExternalAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, ExternalAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 671724742, ExternalAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 617443955, ExternalAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 291759231, ExternalAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
				{Power: 6785098, ExternalAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			sort.Sort(spec.src)
			assert.Equal(t, spec.src, spec.exp)

			rand.Shuffle(len(spec.src), func(i, j int) {
				spec.src[i], spec.src[j] = spec.src[j], spec.src[i]
			})
			sort.Sort(spec.src)
			assert.Equal(t, spec.src, spec.exp)
		})
	}
}

func TestOutgoingTxBatch_GetFees(t *testing.T) {
	type fields struct {
		Transactions []*OutgoingTransferTx
	}
	tests := []struct {
		name   string
		fields fields
		want   sdkmath.Int
	}{
		{
			name: "test 1",
			fields: fields{Transactions: []*OutgoingTransferTx{
				{
					Fee: NewERC20Token(sdkmath.NewInt(0), ""),
				},
				{
					Fee: NewERC20Token(sdkmath.NewInt(1), ""),
				},
				{
					Fee: NewERC20Token(sdkmath.NewInt(1), ""),
				},
			}},
			want: sdkmath.NewInt(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := OutgoingTxBatch{
				Transactions: tt.fields.Transactions,
			}
			if got := m.GetFees(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFees() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		outgoingTxBatch := OutgoingTxBatch{Transactions: []*OutgoingTransferTx{{
			Fee: NewERC20Token(sdkmath.NewInt(1), ""),
		}}}
		fees := outgoingTxBatch.GetFees()
		assert.Equal(b, fees, sdkmath.NewInt(1))
	}
}
