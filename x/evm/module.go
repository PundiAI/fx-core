package evm

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/evmos/ethermint/x/evm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/keeper"
)

var (
	_ module.AppModule = AppModule{}
)

// ____________________________________________________________________________

// AppModule implements an application module for the evm module.
type AppModule struct {
	evm.AppModule
	keeper *keeper.Keeper
	ak     types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModule: evm.NewAppModule(k.Keeper, ak),
		keeper:    k,
		ak:        ak,
	}
}

// DefaultGenesis returns default genesis state as raw bytes for the evm
// module.
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := types.DefaultGenesisState()
	genesisState.Params.EvmDenom = fxtypes.DefaultDenom
	return cdc.MustMarshalJSON(genesisState)
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := evmkeeper.NewMigrator(*am.keeper.Keeper)
	err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(err)
	}
}

// Deprecated: Route returns the message routing key
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// BeginBlock returns the begin block for the evm module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	am.keeper.BeginBlock(ctx, req)
}

// InitGenesis performs genesis initialization for the evm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, am.ak, genesisState)
	return []abci.ValidatorUpdate{}
}
