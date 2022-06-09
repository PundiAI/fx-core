package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOutgoingTxBatch_GetFees(t *testing.T) {
	type fields struct {
		Transactions []*OutgoingTransferTx
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Int
	}{
		{
			name: "test 1",
			fields: fields{Transactions: []*OutgoingTransferTx{
				{
					Fee: NewERC20Token(sdk.NewInt(0), ""),
				},
				{
					Fee: NewERC20Token(sdk.NewInt(1), ""),
				},
				{
					Fee: NewERC20Token(sdk.NewInt(1), ""),
				},
			}},
			want: sdk.NewInt(2),
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
			Fee: NewERC20Token(sdk.NewInt(1), ""),
		}}}
		fees := outgoingTxBatch.GetFees()
		assert.Equal(b, fees, sdk.NewInt(1))
	}
}
