package staking

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v8/x/staking/keeper"
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

type AppModule struct {
	staking.AppModule
	AppModuleBasic AppModuleBasic
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, ls exported.Subspace) AppModule {
	stakingAppModule := staking.NewAppModule(cdc, keeper.Keeper, ak, bk, ls)
	return AppModule{
		AppModuleBasic: AppModuleBasic{AppModuleBasic: stakingAppModule.AppModuleBasic},
		AppModule:      stakingAppModule,
	}
}
