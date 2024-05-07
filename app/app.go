package app

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
	srvflags "github.com/evmos/ethermint/server/flags"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	fxante "github.com/functionx/fx-core/v7/ante"
	"github.com/functionx/fx-core/v7/app/keepers"
	_ "github.com/functionx/fx-core/v7/docs/statik"
	fxcfg "github.com/functionx/fx-core/v7/server/config"
	fxauth "github.com/functionx/fx-core/v7/server/grpc/auth"
	gaspricev1 "github.com/functionx/fx-core/v7/server/grpc/gasprice/legacy/v1"
	gaspricev2 "github.com/functionx/fx-core/v7/server/grpc/gasprice/legacy/v2"
	fxrest "github.com/functionx/fx-core/v7/server/rest"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain"
	"github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/gravity"
	gravitykeeper "github.com/functionx/fx-core/v7/x/gravity/keeper"
	gravitytypes "github.com/functionx/fx-core/v7/x/gravity/types"
)

var _ servertypes.Application = (*App)(nil)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp
	*keepers.AppKeepers
	interfaceRegistry types.InterfaceRegistry
	mm                *module.Manager
	configurator      module.Configurator
}

func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(
		fxtypes.Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	myApp := &App{
		BaseApp:           bApp,
		interfaceRegistry: interfaceRegistry,
	}
	// Setup keepers
	myApp.AppKeepers = keepers.NewAppKeeper(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		myApp.BlockedModuleAccountAddrs(),
		skipUpgradeHeights,
		homePath,
		invCheckPeriod,
		appOpts,
	)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, myApp.AppKeepers.GetKVStoreKey()); err != nil {
		fmt.Printf("failed to load state streaming: %s", err)
		os.Exit(1)
	}

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	myApp.mm = module.NewManager(appModules(myApp, encodingConfig, skipGenesisInvariants)...)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	myApp.mm.SetOrderBeginBlockers(orderBeginBlockers()...)

	myApp.mm.SetOrderEndBlockers(orderEndBlockers()...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	myApp.mm.SetOrderInitGenesis(orderInitBlockers()...)

	myApp.mm.RegisterInvariants(&myApp.CrisisKeeper)
	myApp.mm.RegisterRoutes(myApp.Router(), myApp.QueryRouter(), encodingConfig.Amino)

	myApp.configurator = module.NewConfigurator(myApp.AppCodec(), myApp.MsgServiceRouter(), myApp.GRPCQueryRouter())
	myApp.RegisterServices(myApp.configurator)

	// initialize stores
	myApp.MountKVStores(myApp.GetKVStoreKey())
	myApp.MountTransientStores(myApp.GetTransientStoreKey())
	myApp.MountMemoryStores(myApp.GetMemoryStoreKey())

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	anteOptions := fxante.HandlerOptions{
		AccountKeeper:              myApp.AccountKeeper,
		BankKeeper:                 myApp.BankKeeper,
		EvmKeeper:                  myApp.EvmKeeper,
		FeeMarketKeeper:            myApp.FeeMarketKeeper,
		IbcKeeper:                  myApp.IBCKeeper,
		SignModeHandler:            encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:             fxante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:             maxGasWanted,
		BypassMinFeeMsgTypes:       cast.ToStringSlice(appOpts.Get(fxcfg.BypassMinFeeMsgTypesKey)),
		MaxBypassMinFeeMsgGasUsage: cast.ToUint64(appOpts.Get(fxcfg.BypassMinFeeMsgMaxGasUsageKey)),
	}

	if err := anteOptions.Validate(); err != nil {
		panic(fmt.Errorf("failed to ante options validate: %w", err))
	}

	myApp.SetAnteHandler(fxante.NewAnteHandler(anteOptions))
	myApp.SetInitChainer(myApp.InitChainer)
	myApp.SetBeginBlocker(myApp.BeginBlocker)
	myApp.SetEndBlocker(myApp.EndBlocker)

	myApp.setupUpgradeHandlers()
	myApp.setupUpgradeStoreLoaders()

	myApp.setPostHandler()

	if loadLatest {
		if err := myApp.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return myApp
}

func (app *App) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}
	app.SetPostHandler(postHandler)
}

// Name returns the name of the App
func (app *App) Name() string {
	return app.BaseApp.Name()
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.AppCodec(), genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *App) BlockedModuleAccountAddrs() map[string]bool {
	return app.ModuleAccountAddrs()
}

// InterfaceRegistry returns InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return &module.SimulationManager{}
}

func (app *App) RegisterServices(cfg module.Configurator) {
	for _, m := range app.mm.Modules {
		m.RegisterServices(cfg)
	}
	// Deprecated
	gaspricev1.RegisterQueryServer(cfg.QueryServer(), gaspricev1.Querier{})
	gaspricev2.RegisterQueryServer(cfg.QueryServer(), gaspricev2.Querier{})

	// Deprecated
	gravitytypes.RegisterQueryServer(cfg.QueryServer(), gravitykeeper.NewQueryServerImpl(app.EthKeeper))
	gravitytypes.RegisterMsgServer(cfg.MsgServer(), gravitykeeper.NewMsgServerImpl(app.EthKeeper))

	crosschaintypes.RegisterQueryServer(cfg.QueryServer(), app.CrosschainRouterKeeper)
	crosschaintypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerRouterImpl(app.CrosschainRouterKeeper))

	fxauth.RegisterQueryServer(cfg.QueryServer(), fxauth.Querier{})
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Deprecated: Register gas price queries routes from grpc-gateway.
	gaspricev1.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	gaspricev2.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Deprecated: Register gravity queries routes from grpc-gateway.
	gravity.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Deprecated: cosmos-sdk legacy rest.
	fxrest.RegisterRPCRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterTxRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterAuthRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterBankRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterEvidenceRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterMintRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterDistributeRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterSlashingRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterGovRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterStakingRESTRoutes(clientCtx, apiSvr.Router)
	fxrest.RegisterUpgradeRESTRoutes(clientCtx, apiSvr.Router)

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	// query that exposes operator configuration, most notably the operator's configured minimum gas price
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register crosschain queries routes from grpc-gateway.
	crosschain.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register fx auth queries routes from grpc-gateway.
	fxauth.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		statikFS, err := fs.New()
		if err != nil {
			panic(err)
		}
		staticServer := http.FileServer(statikFS)
		apiSvr.Router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

// RegisterNodeService registers the node gRPC service on the provided
// application gRPC query router.
func (app *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// TestingApp functions

// GetModules implements the TestingApp interface.
func (app *App) GetModules() map[string]module.AppModule {
	return app.mm.Modules
}

// GetOrderBeginBlockersModules implements the TestingApp interface.
func (app *App) GetOrderBeginBlockersModules() []string {
	return app.mm.OrderBeginBlockers
}

// GetOrderEndBlockersModules implements the TestingApp interface.
func (app *App) GetOrderEndBlockersModules() []string {
	return app.mm.OrderEndBlockers
}

// GetOrderInitGenesisModules implements the TestingApp interface.
func (app *App) GetOrderInitGenesisModules() []string {
	return app.mm.OrderInitGenesis
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(string) interface{} {
	return nil
}
