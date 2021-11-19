package fxcore

import (
	"github.com/functionx/fx-core/app"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() app.EncodingConfig {
	encodingConfig := app.MakeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	crosschaintypes.InitMsgValidatorBasicRouter()
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
