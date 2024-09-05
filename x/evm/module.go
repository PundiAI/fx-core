package evm

import (
	"encoding/json"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/evmos/ethermint/x/evm"
	"github.com/evmos/ethermint/x/evm/simulation"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/evm/keeper"
	fxevmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.HasGenesis     = AppModule{}
)

// AppModuleBasic defines the basic application module used by the evm module.
type AppModuleBasic struct {
	evm.AppModuleBasic
}

// DefaultGenesis returns default genesis state as raw bytes for the evm
// module.
func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return b.AppModuleBasic.DefaultGenesis(cdc)
}

// ValidateGenesis is the validation check of the Genesis
func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, txConfig client.TxEncodingConfig, bz json.RawMessage) error {
	return b.AppModuleBasic.ValidateGenesis(cdc, txConfig, bz)
}

// RegisterLegacyAminoCodec registers the evm module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	evmtypes.RegisterLegacyAminoCodec(cdc)
	fxevmtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers interfaces and implementations of the evm module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	evmtypes.RegisterInterfaces(registry)
	fxevmtypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(c client.Context, serveMux *runtime.ServeMux) {
	b.AppModuleBasic.RegisterGRPCGatewayRoutes(c, serveMux)
}

// GetTxCmd returns the root tx command for the evm module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return b.AppModuleBasic.GetTxCmd()
}

// GetQueryCmd returns no root query command for the evm module.
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return b.AppModuleBasic.GetQueryCmd()
}

// AppModule implements an application module for the evm module.
type AppModule struct {
	AppModuleBasic
	evm.AppModule
	keeper         *keeper.Keeper
	legacySubspace evmtypes.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, accountKeeper evmtypes.AccountKeeper, legacySubspace evmtypes.Subspace) AppModule {
	return AppModule{
		AppModule:      evm.NewAppModule(k.Keeper, accountKeeper, legacySubspace),
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		legacySubspace: legacySubspace,
	}
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.AppModule.RegisterServices(cfg)
	fxevmtypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
}

// BeginBlock returns the begin block for the evm module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
	// not reset chain-id on the begin-block
}

// InitGenesis performs genesis initialization for the evm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	ctx = ctx.WithChainID(fxtypes.ChainIdWithEIP155())
	validatorUpdates := am.AppModule.InitGenesis(ctx, cdc, data)
	am.keeper.InitGenesis(ctx)
	return validatorUpdates
}

// ExportGenesis returns the exported genesis state as raw bytes for the evm
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return am.AppModule.ExportGenesis(ctx, cdc)
}

// GenerateGenesisState creates a randomized GenState of the evm module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	RandomizedGenState(simState)
}

// RandomizedGenState generates a random GenesisState for the EVM module
func RandomizedGenState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
	stateBytes, ok := simState.GenState[evmtypes.ModuleName]
	if !ok {
		panic("could not find evm genesis state")
	}

	var evmGenesis evmtypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(stateBytes, &evmGenesis)
	evmGenesis.Params.EvmDenom = fxtypes.DefaultDenom
	simState.GenState[evmtypes.ModuleName] = simState.Cdc.MustMarshalJSON(&evmGenesis)
}
