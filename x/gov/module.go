package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/functionx/fx-core/v3/x/gov/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule         = AppModule{}
	_ module.EndBlockAppModule = AppModule{}
)

// AppModule implements an application module for the gov module.
type AppModule struct {
	gov.AppModule
	keeper keeper.Keeper
	ak     govtypes.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper) AppModule {
	return AppModule{
		AppModule: gov.NewAppModule(cdc, keeper.Keeper, ak, bk),
		keeper:    keeper,
		ak:        ak,
	}
}

// Deprecated: Route returns the message routing key
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	msgServer := keeper.NewMsgServerImpl(govkeeper.NewMsgServerImpl(am.keeper.Keeper), am.keeper)
	govv1betal.RegisterMsgServer(cfg.MsgServer(), govkeeper.NewLegacyMsgServerImpl(am.ak.GetModuleAddress(govtypes.ModuleName).String(), msgServer))
	govv1.RegisterMsgServer(cfg.MsgServer(), msgServer)

	legacyQueryServer := govkeeper.NewLegacyQueryServer(am.keeper.Keeper)
	govv1betal.RegisterQueryServer(cfg.QueryServer(), legacyQueryServer)
	govv1.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := govkeeper.NewMigrator(am.keeper.Keeper)
	err := cfg.RegisterMigration(govtypes.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(govtypes.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(err)
	}
}

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.keeper.EndBlocker(ctx)
	return []abci.ValidatorUpdate{}
}
