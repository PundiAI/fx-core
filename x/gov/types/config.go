package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Config is a config struct used for intialising the gov module to avoid using globals.
type Config struct {
	// MaxTitleLen defines the maximum proposal title length.
	MaxTitleLen uint64
	// MaxSummaryLen defines the maximum proposal summary length.
	MaxSummaryLen uint64
	// Config cosmos gov config
	govtypes.Config
}

// DefaultConfig returns the default config for gov.
func DefaultConfig() Config {
	return Config{
		MaxTitleLen:   uint64(v1beta1.MaxTitleLength),
		MaxSummaryLen: uint64(v1beta1.MaxDescriptionLength),
		Config:        govtypes.DefaultConfig(),
	}
}
