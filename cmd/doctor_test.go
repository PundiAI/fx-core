package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

func Test_getGenesisSha256(t *testing.T) {
	type args struct {
		genesisFile string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "mainnet",
			args: args{genesisFile: "../public/mainnet/genesis.json"},
			want: fxtypes.MainnetGenesisHash,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i)
			},
		},
		{
			name: "testnet",
			args: args{genesisFile: "../public/testnet/genesis.json"},
			want: fxtypes.TestnetGenesisHash,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGenesisSha256(tt.args.genesisFile)
			if !tt.wantErr(t, err, fmt.Sprintf("getGenesisSha256(%v)", tt.args.genesisFile)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getGenesisSha256(%v)", tt.args.genesisFile)
		})
	}
}
