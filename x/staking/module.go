package staking

import (
	"context"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v8/x/staking/keeper"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

var (
	_ module.AppModuleBasic  = AppModuleBasic{}
	_ module.HasServices     = AppModule{}
	_ module.HasInvariants   = AppModule{}
	_ module.HasABCIGenesis  = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}

	_ appmodule.AppModule = AppModule{}
)

type AppModuleBasic struct {
	staking.AppModuleBasic
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	stakingtypes.RegisterLegacyAminoCodec(cdc)
	fxstakingtypes.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	stakingtypes.RegisterInterfaces(registry)
	fxstakingtypes.RegisterInterfaces(registry)
}

type AppModule struct {
	staking.AppModule
	AppModuleBasic AppModuleBasic
	Keeper         *keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, ls exported.Subspace) AppModule {
	stakingAppModule := staking.NewAppModule(cdc, keeper.Keeper, ak, bk, ls)
	return AppModule{
		AppModuleBasic: AppModuleBasic{AppModuleBasic: stakingAppModule.AppModuleBasic},
		AppModule:      stakingAppModule,
		Keeper:         keeper,
	}
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return am.Keeper.EndBlocker(ctx)
}
