package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	cryptocodec "github.com/functionx/fx-core/crypto/codec"
	enccodec "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
	encodingConfig := makeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	crosschaintypes.InitMsgValidatorBasicRouter()
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func makeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	encodingConfig := EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterCrypto(amino)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
