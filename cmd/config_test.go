package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/stretchr/testify/assert"
	tmcfg "github.com/tendermint/tendermint/config"

	fxcfg "github.com/functionx/fx-core/v3/server/config"
)

func Test_updateCfgCmd(t *testing.T) {
	tempDir := t.TempDir()
	defer assert.NoError(t, os.RemoveAll(tempDir))
	assert.NoError(t, os.MkdirAll(filepath.Join(tempDir, "config"), 0o700))

	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"version"})
	assert.NoError(t, svrcmd.Execute(rootCmd, tempDir))

	publicDir, err := os.ReadDir("../public")
	assert.NoError(t, err)
	for _, entry := range publicDir {
		appConfig, err := os.ReadFile(fmt.Sprintf("../public/%s/app.toml", entry.Name()))
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(filepath.Join(tempDir, "config/app.toml"), appConfig, 0o700))

		tmConfig, err := os.ReadFile(fmt.Sprintf("../public/%s/config.toml", entry.Name()))
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(filepath.Join(tempDir, "config/config.toml"), tmConfig, 0o700))

		rootCmd.SetArgs([]string{"config", "update"})
		assert.NoError(t, rootCmd.Execute())

		appConfigAfter, err := os.ReadFile(filepath.Join(tempDir, "config/app.toml"))
		assert.NoError(t, err)
		assert.Equal(t, string(appConfig), string(appConfigAfter))

		tmConfigAfter, err := os.ReadFile(filepath.Join(tempDir, "config/config.toml"))
		assert.NoError(t, err)
		assert.Equal(t, string(tmConfig), string(tmConfigAfter))
	}
}

func TestPublicTmConfig(t *testing.T) {
	tempDir := t.TempDir()
	defer assert.NoError(t, os.RemoveAll(tempDir))
	assert.NoError(t, os.MkdirAll(tempDir, 0o700))

	serverCtx := server.NewDefaultContext()
	fileName := fmt.Sprintf("%s/config.toml", t.TempDir())
	serverCtx.Config.BaseConfig.Moniker = "your-moniker"
	serverCtx.Config.Consensus.TimeoutCommit = 5 * time.Second
	serverCtx.Config.Instrumentation.Prometheus = true
	serverCtx.Config.P2P.AddrBookStrict = false
	serverCtx.Config.P2P.MaxNumOutboundPeers = 30
	serverCtx.Config.P2P.Seeds = "c5877d9d243af1a504caf5b7f7a9c915b3ae94ae@fxcore-mainnet-seed-node-1.functionx.io:26656,b289311ece065c813287e3a25835bb6378999aa5@fxcore-mainnet-seed-node-2.functionx.io:26656"
	serverCtx.Config.P2P.PersistentPeers = "96f04dffc25ffcce11e179581d2a3ab6cb5535d5@fxcore-mainnet-node-1.functionx.io:26656,836ded83bac83a4ac8511826fa1ad4ca2238f960@fxcore-mainnet-node-2.functionx.io:26656,7c7a260eeefda37eac896ae423e78cf345a2ef70@fxcore-mainnet-node-3.functionx.io:26656,0fee38117655b6961319950d6beb929fb194217c@fxcore-mainnet-node-4.functionx.io:26656,6e8818051a2ca9b8be67a6f2ba48c33d8c489d5c@fxcore-mainnet-node-5.functionx.io:26656"

	tmcfg.WriteConfigFile(fileName, serverCtx.Config)
	defConfig, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	mainnetConfig, err := os.ReadFile("../public/mainnet/config.toml")
	assert.NoError(t, err)
	assert.Equal(t, string(defConfig), string(mainnetConfig))

	serverCtx.Config.P2P.Seeds = "e922b34e660976a64d6024bde495666752141992@dhobyghaut-seed-node-1.functionx.io:26656"
	serverCtx.Config.P2P.PersistentPeers = "a817685c010402703820be2b5a90d9e07bc5c2d3@dhobyghaut-node-1.functionx.io:26656,d22e741b4e8e2586dbe38fd348d3de8dfbb889a0@dhobyghaut-node-2.functionx.io:26656,c1a985c7e4c0b5ce6d343d87e070a63b24a76594@dhobyghaut-node-3.functionx.io:26656,cc267dac09a38b67b3bda0033f62678cb54bf843@dhobyghaut-node-4.functionx.io:26656,0ea7e81071d4004a1fbbe304477d8ca3183a5282@dhobyghaut-node-5.functionx.io:26656"
	tmcfg.WriteConfigFile(fileName, serverCtx.Config)
	defConfig, err = os.ReadFile(fileName)
	assert.NoError(t, err)
	testnetConfig, err := os.ReadFile("../public/testnet/config.toml")
	assert.NoError(t, err)
	assert.Equal(t, string(defConfig), string(testnetConfig))
}

func TestPublicAppConfig(t *testing.T) {
	tempDir := t.TempDir()
	defer assert.NoError(t, os.RemoveAll(tempDir))
	assert.NoError(t, os.MkdirAll(tempDir, 0o700))

	config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
	appConfig := fxcfg.DefaultConfig()
	appConfig.Pruning = pruningtypes.PruningOptionCustom
	appConfig.PruningKeepRecent = "20000"
	appConfig.PruningInterval = "10"
	appConfig.MinGasPrices = "4000000000000FX"
	appConfig.IAVLDisableFastNode = true
	appConfig.Telemetry.EnableServiceLabel = true
	appConfig.Telemetry.Enabled = true
	appConfig.Telemetry.PrometheusRetentionTime = 60
	appConfig.API.Enable = true
	appConfig.API.Swagger = true
	appConfig.EVM.MaxTxGasWanted = 500000

	fileName := fmt.Sprintf("%s/app.toml", t.TempDir())
	config.WriteConfigFile(fileName, appConfig)
	defAppConfig, err := os.ReadFile(fileName)
	assert.NoError(t, err)

	mainnetConfig, err := os.ReadFile("../public/mainnet/app.toml")
	assert.NoError(t, err)
	assert.Equal(t, string(defAppConfig), string(mainnetConfig))

	testnetConfig, err := os.ReadFile("../public/testnet/app.toml")
	assert.NoError(t, err)
	assert.Equal(t, string(defAppConfig), string(testnetConfig))
}
