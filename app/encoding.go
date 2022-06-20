package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	cryptocodec "github.com/evmos/ethermint/crypto/codec"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	ethermint "github.com/evmos/ethermint/types"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
	encodingConfig := makeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	crosschaintypes.InitMsgValidatorBasicRouter()
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
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
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	keyring.RegisterLegacyAminoCodec(amino)

	RegisterCryptoEthSecp256k1(amino)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ethermint.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

func RegisterCryptoEthSecp256k1(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&ethsecp256k1.PubKey{},
		ethsecp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&ethsecp256k1.PrivKey{},
		ethsecp256k1.PrivKeyName, nil)

	// NOTE: update SDK's amino codec to include the ethsecp256k1 keys.
	// DO NOT REMOVE unless deprecated on the SDK.
	legacy.Cdc = cdc
	keys.KeysCdc = cdc
}
