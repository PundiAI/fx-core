package tron

import (
	"encoding/json"
	"github.com/functionx/fx-core/x/crosschain"
	"github.com/functionx/fx-core/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/x/crosschain/keeper"
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
	return trontypes.ModuleName
}

// RegisterLegacyAminoCodec implements app module basic
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
}

// DefaultGenesis implements app module basic
func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	return nil
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
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {
}

// RegisterInterfaces implements app bmodule basic
func (b AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {
	types.RegisterValidatorBasic(trontypes.ModuleName, trontypes.MsgValidateBasic{})
}

//____________________________________________________________________________

// AppModule object for module implementation
type AppModule struct {
	AppModuleBasic
	keeper     keeper.Keeper
	bankKeeper bankkeeper.Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(keeper keeper.Keeper, bankKeeper bankkeeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		bankKeeper:     bankKeeper,
	}
}

// Name implements app module
func (AppModule) Name() string {
	return trontypes.ModuleName
}

// RegisterInvariants implements app module
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
	//  make some invariants in the gravity module to ensure that coins aren't being fraudlently minted etc...
}

// Route implements app module
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute implements app module
func (am AppModule) QuerierRoute() string {
	return trontypes.QuerierRoute
}

// LegacyQuerierHandler returns the distribution module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(_ module.Configurator) {}

// InitGenesis initializes the genesis state for this module and implements app module.
func (am AppModule) InitGenesis(_ sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the current genesis state to a json.RawMessage
func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage {
	return nil
}

// BeginBlock implements app module
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock implements app module
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	crosschain.EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
