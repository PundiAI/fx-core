package config

import (
	"github.com/tendermint/tendermint/config"
)

// DefaultTendermintConfig returns tendermint default configuration.
func DefaultTendermintConfig() *config.Config {
	defaultConfig := config.DefaultConfig()
	defaultConfig.Instrumentation.Namespace = "tendermint"
	return defaultConfig
}
