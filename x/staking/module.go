package staking

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v4/x/staking/keeper"
	"github.com/functionx/fx-core/v4/x/staking/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	staking.AppModuleBasic
}

// DefaultGenesis returns default genesis state as raw bytes for the staking
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", stakingtypes.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

type AppModule struct {
	staking.AppModule
	AppModuleBasic AppModuleBasic
	Keeper         keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper) AppModule {
	stakingAppModule := staking.NewAppModule(cdc, keeper.Keeper, ak, bk)
	return AppModule{
		AppModuleBasic: AppModuleBasic{AppModuleBasic: stakingAppModule.AppModuleBasic},
		AppModule:      stakingAppModule,
		Keeper:         keeper,
	}
}

// InitGenesis performs genesis initialization for the staking module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	return am.Keeper.InitGenesis(ctx, &genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(am.Keeper.ExportGenesis(ctx))
}

// DefaultGenesis Override AppModule.DefaultGenesis by AppModuleBasic.DefaultGenesis
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return am.AppModuleBasic.DefaultGenesis(cdc)
}

// ValidateGenesis Override AppModule.ValidateGenesis by AppModuleBasic.ValidateGenesis
func (am AppModule) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return am.AppModuleBasic.ValidateGenesis(cdc, config, bz)
}
