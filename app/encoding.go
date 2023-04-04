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
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethermintcodec "github.com/evmos/ethermint/crypto/codec"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	etherminttypes "github.com/evmos/ethermint/types"

	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
	encodingConfig := makeEncodingConfig()
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	registerCryptoEthSecp256k1(encodingConfig.Amino)
	ethermintcodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	etherminttypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	crosschaintypes.RegisterLegacyAminoCodec(encodingConfig.Amino)
	crosschaintypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func makeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	encodingConfig := EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             cdc,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	keyring.RegisterLegacyAminoCodec(amino)
	return encodingConfig
}

func registerCryptoEthSecp256k1(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&ethsecp256k1.PubKey{},
		ethsecp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&ethsecp256k1.PrivKey{},
		ethsecp256k1.PrivKeyName, nil)

	// NOTE: update SDK's amino codec to include the ethsecp256k1 keys.
	// DO NOT REMOVE unless deprecated on the SDK.
	legacy.Cdc = cdc
	keys.KeysCdc = cdc
}

func init() {
	// gov v1beta1
	distrtypes.RegisterLegacyAminoCodec(govv1beta1.ModuleCdc.LegacyAmino)
	paramproposal.RegisterLegacyAminoCodec(govv1beta1.ModuleCdc.LegacyAmino)
	upgradetypes.RegisterLegacyAminoCodec(govv1beta1.ModuleCdc.LegacyAmino)
	crosschaintypes.RegisterLegacyAminoCodec(govv1beta1.ModuleCdc.LegacyAmino)
	erc20types.RegisterLegacyAminoCodec(govv1beta1.ModuleCdc.LegacyAmino)

	// gov v1
	crosschaintypes.RegisterLegacyAminoCodec(govv1.ModuleCdc.LegacyAmino)
	erc20types.RegisterLegacyAminoCodec(govv1.ModuleCdc.LegacyAmino)
	upgradetypes.RegisterLegacyAminoCodec(govv1.ModuleCdc.LegacyAmino)
}
