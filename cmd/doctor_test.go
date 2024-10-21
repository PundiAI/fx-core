package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v8/types"
)

func Test_doctorCmd(t *testing.T) {
	cmd := doctorCmd()
	cmd.SetArgs([]string{})
	require.NoError(t, cmd.Execute())
}

func Test_getGenesisSha256(t *testing.T) {
	tests := []struct {
		name        string
		genesisFile string
		want        string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "mainnet",
			genesisFile: "../public/mainnet/genesis.json",
			want:        fxtypes.MainnetGenesisHash,
			wantErr:     assert.NoError,
		},
		{
			name:        "testnet",
			genesisFile: "../public/testnet/genesis.json",
			want:        fxtypes.TestnetGenesisHash,
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := getGenesisDocAndSha256(tt.genesisFile)
			if !tt.wantErr(t, err, fmt.Sprintf("getGenesisDocAndSha256(%v)", tt.genesisFile)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getGenesisDocAndSha256(%v)", tt.genesisFile)
		})
	}
}

func Test_checkGenesis(t *testing.T) {
	tests := []struct {
		name        string
		genesisFile string
		want        string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "mainnet",
			genesisFile: "../public/mainnet/genesis.json",
			want:        fxtypes.MainnetChainId,
			wantErr:     assert.NoError,
		},
		{
			name:        "testnet",
			genesisFile: "../public/testnet/genesis.json",
			want:        fxtypes.TestnetChainId,
			wantErr:     assert.NoError,
		},
		{
			name:        "not exist",
			genesisFile: "../public/invalid/genesis.json",
			want:        "",
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkGenesis(tt.genesisFile)
			if !tt.wantErr(t, err, fmt.Sprintf("checkGenesis(%v)", tt.genesisFile)) {
				return
			}
			assert.Equalf(t, tt.want, got, "checkGenesis(%v)", tt.genesisFile)
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
