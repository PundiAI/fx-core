package tron

import (
	"encoding/json"
	"fmt"

	"github.com/functionx/fx-core/x/tron/keeper"

	crosschainv54 "github.com/functionx/fx-core/x/crosschain/legacy/v045"

	"github.com/functionx/fx-core/x/crosschain"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/tron/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic object for module implementation
type AppModuleBasic struct{}

// Name implements app module basic
func (AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec implements app module basic
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// DefaultGenesis implements app module basic
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis implements app module basic
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, data json.RawMessage) error {
	var state crosschaintypes.GenesisState
	if err := cdc.UnmarshalJSON(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return state.ValidateBasic()
}

// RegisterRESTRoutes implements app module basic
func (AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {}

// GetQueryCmd implements app module basic
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// GetTxCmd implements app module basic
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the distribution module.
// also implements app modeul basic
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

// RegisterInterfaces implements app bmodule basic
func (AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {
	crosschaintypes.RegisterValidatorBasic(types.ModuleName, types.TronMsgValidate{})
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule object for module implementation
type AppModule struct {
	AppModuleBasic
	keeper        keeper.Keeper
	stakingKeeper crosschainv54.StakingKeeper
	paramsKey     sdk.StoreKey
}

// NewAppModule creates a new AppModule Object
func NewAppModule(keeper keeper.Keeper, stakingKeeper crosschainv54.StakingKeeper, paramsKey sdk.StoreKey) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		stakingKeeper:  stakingKeeper,
		paramsKey:      paramsKey,
	}
}

// RegisterInvariants implements app module
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route implements app module
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute implements app module
func (am AppModule) QuerierRoute() string { return types.QuerierRoute }

// LegacyQuerierHandler returns the distribution module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return crosschainkeeper.NewQuerier(am.keeper.Keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	m := crosschainkeeper.NewMigrator(am.keeper.Keeper, am.stakingKeeper, am.paramsKey)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(err)
	}
}

// InitGenesis initializes the genesis state for this module and implements app module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState crosschaintypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	crosschain.InitGenesis(ctx, am.keeper.Keeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the current genesis state to a json.RawMessage
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	state := crosschain.ExportGenesis(ctx, am.keeper.Keeper)
	return cdc.MustMarshalJSON(state)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 {
	return 2
}

// BeginBlock implements app module
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock implements app module
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	crosschain.EndBlocker(ctx, am.keeper.Keeper)
	return []abci.ValidatorUpdate{}
}
