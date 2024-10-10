package gov

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v8/x/gov/client/cli"
	"github.com/functionx/fx-core/v8/x/gov/keeper"
	"github.com/functionx/fx-core/v8/x/gov/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.HasServices    = AppModule{}
	_ module.HasInvariants  = AppModule{}
	_ module.HasABCIGenesis = AppModule{}

	_ appmodule.AppModule     = AppModule{}
	_ appmodule.HasEndBlocker = AppModule{}
)

// AppModuleBasic defines the basic application module used by the gov module.
type AppModuleBasic struct {
	gov.AppModuleBasic
}

// NewAppModuleBasic creates a new AppModuleBasic object
func NewAppModuleBasic(legacyProposalHandlers []govclient.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		AppModuleBasic: gov.NewAppModuleBasic(legacyProposalHandlers),
	}
}

// GetQueryCmd returns the root query command for the gov module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterLegacyAminoCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	govv1betal.RegisterLegacyAminoCodec(cdc)
	govv1.RegisterLegacyAminoCodec(cdc)
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the gov module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := govv1.RegisterQueryHandlerClient(context.Background(), mux, govv1.NewQueryClient(clientCtx)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", govtypes.ModuleName, err.Error()))
	}
	if err := govv1betal.RegisterQueryHandlerClient(context.Background(), mux, govv1betal.NewQueryClient(clientCtx)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", govtypes.ModuleName, err.Error()))
	}
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(fmt.Sprintf("failed to %s register grpc gateway routes: %s", govtypes.ModuleName, err.Error()))
	}
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	govv1.RegisterInterfaces(registry)
	govv1betal.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)
}

// AppModule implements an application module for the gov module.
type AppModule struct {
	gov.AppModule
	keeper         *keeper.Keeper
	ak             govtypes.AccountKeeper
	bk             govtypes.BankKeeper
	cdc            codec.Codec
	legacySubspace govtypes.ParamSubspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, ss govtypes.ParamSubspace) AppModule {
	return AppModule{
		AppModule:      gov.NewAppModule(cdc, keeper.Keeper, ak, bk, ss),
		keeper:         keeper,
		ak:             ak,
		bk:             bk,
		cdc:            cdc,
		legacySubspace: ss,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	govMsgServer := govkeeper.NewMsgServerImpl(am.keeper.Keeper)
	msgServer := keeper.NewMsgServerImpl(govMsgServer, am.keeper)

	legacyMsgServer := govkeeper.NewLegacyMsgServerImpl(am.ak.GetModuleAddress(govtypes.ModuleName).String(), msgServer)
	govv1betal.RegisterMsgServer(cfg.MsgServer(), legacyMsgServer)
	govv1.RegisterMsgServer(cfg.MsgServer(), msgServer)
	types.RegisterMsgServer(cfg.MsgServer(), msgServer)

	legacyQueryServer := govkeeper.NewLegacyQueryServer(am.keeper.Keeper)
	govv1betal.RegisterQueryServer(cfg.QueryServer(), legacyQueryServer)
	govv1.RegisterQueryServer(cfg.QueryServer(), govkeeper.NewQueryServer(am.keeper.Keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the gov module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	am.AppModule.InitGenesis(ctx, cdc, data)

	// init fx gov params
	if err := am.keeper.InitCustomParams(ctx); err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) error {
	return EndBlocker(sdk.UnwrapSDKContext(ctx), am.keeper)
}
