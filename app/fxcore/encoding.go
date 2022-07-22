package fxcore

import (
	"github.com/functionx/fx-core/app"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() app.EncodingConfig {
	encodingConfig := app.MakeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
