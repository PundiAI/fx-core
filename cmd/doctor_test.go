package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

func Test_doctorCmd(t *testing.T) {
	cmd := doctorCmd()
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
}

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

func Test_checkGenesis(t *testing.T) {
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
			want: fxtypes.MainnetChainId,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i)
			},
		},
		{
			name: "testnet",
			args: args{genesisFile: "../public/testnet/genesis.json"},
			want: fxtypes.TestnetChainId,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i)
			},
		},
		{
			name: "not exist",
			args: args{genesisFile: "../public/invalid/genesis.json"},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkGenesis(tt.args.genesisFile)
			if !tt.wantErr(t, err, fmt.Sprintf("checkGenesis(%v)", tt.args.genesisFile)) {
				return
			}
			assert.Equalf(t, tt.want, got, "checkGenesis(%v)", tt.args.genesisFile)
		})
	}
}

func Test_checkVersionCompatibility(t *testing.T) {
	type args struct {
		version       string
		outputVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "v1",
			args: args{
				version:       "fxv1",
				outputVersion: "release/v1.0.1",
			},
			want: false,
		},
		{
			name: "v5",
			args: args{
				version:       "v5.0.x",
				outputVersion: "release/v5.0.1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, checkVersionCompatibility(tt.args.version, tt.args.outputVersion), "checkVersionCompatibility(%v, %v)", tt.args.version, tt.args.outputVersion)
		})
	}
}
