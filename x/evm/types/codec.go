package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

// RegisterInterfaces registers the client interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&evmtypes.MsgEthereumTx{},
		&MsgCallContract{},
	)
	registry.RegisterImplementations(
		(*tx.TxExtensionOptionI)(nil),
		&evmtypes.ExtensionOptionsEthereumTx{},
	)
	registry.RegisterInterface(
		"ethermint.evm.v1.TxData",
		(*evmtypes.TxData)(nil),
		&evmtypes.DynamicFeeTx{},
		&evmtypes.AccessListTx{},
		&evmtypes.LegacyTx{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
