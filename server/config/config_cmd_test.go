package config_test

import (
	"bytes"
	"fmt"
	"testing"

	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/server/config"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func Test_Output(t *testing.T) {
	tests := []struct {
		name    string
		content interface{}
	}{
		{
			name:    "app.toml Output grpc.enable",
			content: true,
		},
		{
			name:    "app.toml Output bypass-min-fee.msg-types empty",
			content: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientCtx := client.Context{
				Output:       new(bytes.Buffer),
				OutputFormat: "json",
			}
			assert.NoError(t, config.Output(clientCtx, tt.content))
			assert.Equal(t, clientCtx.Output.(*bytes.Buffer).String(), fmt.Sprintf("%v\n", tt.content))
		})
	}
}

func Test_ConfigTomlConfig_Output(t *testing.T) {
	cfg := tmcfg.DefaultConfig()
	cfg.BaseConfig.Moniker = "anonymous"
	c := config.TmConfigToml{Config: cfg}
	buf := new(bytes.Buffer)
	clientCtx := client.Context{
		Output:       buf,
		OutputFormat: "json",
	}
	assert.NoError(t, c.Output(clientCtx))

	helpers.AssertJsonFile(t, "./data/config.json", buf.Bytes())
}

func Test_AppTomlConfig_Output(t *testing.T) {
	_, v := config.AppConfig(fxtypes.GetDefGasPrice())
	cfg := v.(config.Config)
	c := config.AppToml{Config: &cfg}
	buf := new(bytes.Buffer)
	clientCtx := client.Context{
		Output:       buf,
		OutputFormat: "json",
	}
	assert.NoError(t, c.Output(clientCtx))

	helpers.AssertJsonFile(t, "./data/app.json", buf.Bytes())
}
