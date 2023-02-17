package staking

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/staking/keeper"
	"github.com/functionx/fx-core/v3/x/staking/types"
)

var (
	_ module.AppModule = AppModule{}
	// _ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModule implements an application module for the staking module.
type AppModule struct {
	staking.AppModule
	keeper        *keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    stakingtypes.BankKeeper
	evmKeeper     types.EvmKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak types.AccountKeeper, bk stakingtypes.BankKeeper, evmKeeper types.EvmKeeper) AppModule {
	return AppModule{
		AppModule:     staking.NewAppModule(cdc, keeper.Keeper, ak, bk),
		keeper:        keeper,
		accountKeeper: ak,
		bankKeeper:    bk,
		evmKeeper:     evmKeeper,
	}
}

// Name returns the staking module's name.
func (AppModule) Name() string {
	return stakingtypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(stakingtypes.RouterKey, NewHandler(am.keeper))
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := stakingkeeper.Querier{Keeper: am.keeper.Keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(am.keeper.Keeper)
	// TODO: if upgrade cosmos-sdk check is needed migrate
	_ = cfg.RegisterMigration(stakingtypes.ModuleName, 1, m.Migrate1to2)
}

// InitGenesis performs genesis initialization for the staking module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState stakingtypes.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)

	if err := am.evmKeeper.CreateContractWithCode(ctx, fxtypes.GetLPToken().Address, fxtypes.GetLPToken().Code); err != nil {
		panic(err)
	}
	CreateLPTokenModuleAccount(ctx, types.LPTokenOwnerModuleName, am.accountKeeper)
	return staking.InitGenesis(ctx, am.keeper.Keeper, am.accountKeeper, am.bankKeeper, &genesisState)
}

func CreateLPTokenModuleAccount(ctx sdk.Context, lpTokenModuleName string, ak types.AccountKeeper) {
	moduleAddress, permissions := ak.GetModuleAddressAndPermissions(lpTokenModuleName)
	if moduleAddress == nil || moduleAddress.Empty() {
		panic("LPTokenOwnerModuleName module account has not been set")
	}

	moduleAccount := ak.GetAccount(ctx, moduleAddress)
	if _, ok := moduleAccount.(authtypes.ModuleAccountI); ok {
		return
	}

	// create a new module account
	macc := authtypes.NewEmptyModuleAccount(lpTokenModuleName, permissions...)
	maccI := (ak.NewAccount(ctx, macc)).(authtypes.ModuleAccountI) // set the account number
	ak.SetModuleAccount(ctx, maccI)
}
