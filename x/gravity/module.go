package gravity

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/x/gravity/keeper"
	"github.com/functionx/fx-core/x/gravity/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic object for module implementation
type AppModuleBasic struct{}

// Name implements app module basic
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec implements app module basic
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis implements app module basic
func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	return []byte("{}")
}

// ValidateGenesis implements app module basic
func (AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
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
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule object for module implementation
type AppModule struct {
	AppModuleBasic
	keeper  keeper.Keeper
	migrate keeper.Migrator
}

// NewAppModule creates a new AppModule Object
func NewAppModule(keeper keeper.Keeper, migrate keeper.Migrator) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		migrate:        migrate,
	}
}

// RegisterInvariants implements app module
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Name implements app module
func (AppModule) Name() string {
	return types.ModuleName
}

// Route implements app module
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute implements app module
func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the distribution module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		return nil, nil
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	if err := cfg.RegisterMigration(types.ModuleName, 1, am.migrate.Migrate1to2); err != nil {
		panic(err)
	}
}

// InitGenesis initializes the genesis state for this module and implements app module.
func (am AppModule) InitGenesis(_ sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the current genesis state to a json.RawMessage
func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage {
	return []byte("{}")
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 {
	return 2
}

// BeginBlock implements app module
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock implements app module
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
