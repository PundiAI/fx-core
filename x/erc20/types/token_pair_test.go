package types_test

import (
	"reflect"
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func TestTokenPair_GetID(t *testing.T) {
	type fields struct {
		Erc20Address  string
		Denom         string
		Enabled       bool
		ContractOwner types.Owner
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"valid",
			fields{
				Erc20Address: "0x0000000000000000000000000000000000000000",
				Denom:        "test",
			},
			tmhash.Sum([]byte("0x0000000000000000000000000000000000000000|test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := types.TokenPair{
				Erc20Address:  tt.fields.Erc20Address,
				Denom:         tt.fields.Denom,
				Enabled:       tt.fields.Enabled,
				ContractOwner: tt.fields.ContractOwner,
			}
			if got := tp.GetID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
