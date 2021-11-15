package fxcore

import (
	"github.com/functionx/fx-core/app"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() app.EncodingConfig {
	encodingConfig := app.MakeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//codec.RegisterCrypto(encodingConfig.Amino)
	crosschaintypes.InitMsgValidatorBasicRouter()
	//types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
