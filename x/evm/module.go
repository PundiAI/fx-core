package evm

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/evmos/ethermint/x/evm"
	"github.com/evmos/ethermint/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/keeper"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.BeginBlockAppModule = AppModule{}
	_ module.EndBlockAppModule   = AppModule{}
)

// AppModuleBasic defines the basic application module used by the evm module.
type AppModuleBasic struct {
	evm.AppModuleBasic
}

// DefaultGenesis returns default genesis state as raw bytes for the evm
// module.
func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := types.DefaultGenesisState()
	genesisState.Params.EvmDenom = fxtypes.DefaultDenom
	return cdc.MustMarshalJSON(genesisState)
}

// AppModule implements an application module for the evm module.
type AppModule struct {
	AppModuleBasic
	keeper         *keeper.Keeper
	accountKeeper  types.AccountKeeper
	legacyAmino    *codec.LegacyAmino
	paramsStoreKey storetypes.StoreKey
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, accountKeeper types.AccountKeeper, legacyAmino *codec.LegacyAmino, paramsStoreKey storetypes.StoreKey) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		accountKeeper:  accountKeeper,
		legacyAmino:    legacyAmino,
		paramsStoreKey: paramsStoreKey,
	}
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	migrator := keeper.NewMigrator(am.keeper, am.legacyAmino, am.paramsStoreKey)
	if err := cfg.RegisterMigration(types.ModuleName, 2, migrator.Migrate2to3); err != nil {
		panic(err)
	}
}

// Deprecated: Route returns the message routing key
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute implements app module
func (AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler returns no sdk.Querier
func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return nil
}

// BeginBlock returns the begin block for the evm module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	am.keeper.BeginBlock(ctx, req)
}

// EndBlock returns the end blocker for the evm module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.keeper.EndBlock(ctx, req)
}

// InitGenesis performs genesis initialization for the evm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, am.accountKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the evm
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	state := evm.ExportGenesis(ctx, am.keeper.Keeper, am.accountKeeper)
	return cdc.MustMarshalJSON(state)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 {
	return 3
}
