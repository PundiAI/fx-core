package public

import (
	"embed"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	GenesisFileName   = "genesis.json"
	ConfigFileName    = "config.toml"
	ConfigAppFileName = "app.toml"
)

//go:embed *
var config embed.FS

func ImportGenesisConfig(cmd *cobra.Command) error {
	if !strings.EqualFold(cmd.Use, "start") {
		return nil
	}
	home, err := cmd.Flags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}
	network, err := cmd.Flags().GetString("network")
	if err != nil {
		return err
	}
	if len(network) <= 0 || len(home) <= 0 {
		return nil
	}
	logger, err := ctxLogger(cmd)
	if err != nil {
		return err
	}
	configPath := filepath.Join(home, "config")
	genesisPath := filepath.Join(configPath, GenesisFileName)

	logger.Debug("path", "config", configPath)
	if tmos.FileExists(genesisPath) || network == "" {
		if tmos.FileExists(genesisPath) {
			logger.Info("genesis file exist", GenesisFileName, genesisPath)
		}
		return nil
	}
	if !tmos.FileExists(configPath) {
		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			return err
		}
	}
	logger.Info("Starting by network", "network", network, "home", home)
	//config.toml
	logger.Info("write config file", "file", filepath.Join(configPath, ConfigFileName))
	configTomlData, err := config.ReadFile(ConfigFileName)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(configPath, ConfigFileName), configTomlData, 0666); err != nil {
		return err
	}
	//app.toml
	logger.Info("write app file", "file", filepath.Join(configPath, ConfigAppFileName))
	configAppData, err := config.ReadFile(ConfigAppFileName)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(configPath, ConfigAppFileName), configAppData, 0666); err != nil {
		return err
	}
	//genesis.json
	logger.Info("write genesis file", "file", filepath.Join(configPath, GenesisFileName))
	genesisData, err := config.ReadFile(GenesisFileName)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(configPath, GenesisFileName), genesisData, 0666); err != nil {
		return err
	}
	return nil
}

func ctxLogger(cmd *cobra.Command) (server.ZeroLogWrapper, error) {
	logFormat, err := cmd.Flags().GetString(flags.FlagLogFormat)
	if err != nil {
		return server.ZeroLogWrapper{}, err
	}
	logLevel, err := cmd.Flags().GetString(flags.FlagLogLevel)
	if err != nil {
		return server.ZeroLogWrapper{}, err
	}
	var logWriter io.Writer
	if strings.ToLower(logFormat) == tmcfg.LogFormatPlain {
		logWriter = zerolog.ConsoleWriter{Out: os.Stderr}
	} else {
		logWriter = os.Stderr
	}
	logLvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return server.ZeroLogWrapper{}, fmt.Errorf("failed to parse log level (%s): %w", logLevel, err)
	}

	return server.ZeroLogWrapper{Logger: zerolog.New(logWriter).Level(logLvl).With().Timestamp().Logger()}, nil
}
