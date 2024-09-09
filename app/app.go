package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
	evmante "github.com/evmos/ethermint/app/ante"
	srvflags "github.com/evmos/ethermint/server/flags"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	fxante "github.com/functionx/fx-core/v8/ante"
	"github.com/functionx/fx-core/v8/app/keepers"
	_ "github.com/functionx/fx-core/v8/docs/statik"
	fxcfg "github.com/functionx/fx-core/v8/server/config"
	fxauth "github.com/functionx/fx-core/v8/server/grpc/auth"
	gaspricev1 "github.com/functionx/fx-core/v8/server/grpc/gasprice/legacy/v1"
	gaspricev2 "github.com/functionx/fx-core/v8/server/grpc/gasprice/legacy/v2"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain"
	"github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/gravity"
	gravitykeeper "github.com/functionx/fx-core/v8/x/gravity/keeper"
	gravitytypes "github.com/functionx/fx-core/v8/x/gravity/types"
)

var _ servertypes.Application = (*App)(nil)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp
	*keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry
	mm                *module.Manager
	configurator      module.Configurator

	pendingTxListeners []evmante.PendingTxListener
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
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(
		fxtypes.Name,
		logger,
		db,
		txConfig.TxDecoder(),
		baseAppOptions...)

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	myApp := &App{
		BaseApp:           bApp,
		txConfig:          txConfig,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
	}

	// Setup keepers
	myApp.AppKeepers = keepers.NewAppKeeper(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		BlockedAccountAddrs(),
		skipUpgradeHeights,
		homePath,
		invCheckPeriod,
		appOpts,
	)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, myApp.AppKeepers.GetKVStoreKey()); err != nil {
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
	myApp.mm.SetOrderExportGenesis(orderInitBlockers()...)

	myApp.mm.RegisterInvariants(myApp.CrisisKeeper)
	myApp.configurator = module.NewConfigurator(myApp.appCodec, myApp.MsgServiceRouter(), myApp.GRPCQueryRouter())
	myApp.RegisterServices(myApp.configurator)

	autocliv1.RegisterQueryServer(myApp.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(myApp.mm.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(myApp.GRPCQueryRouter(), reflectionSvc)

	// initialize stores
	myApp.MountKVStores(myApp.GetKVStoreKey())
	myApp.MountTransientStores(myApp.GetTransientStoreKey())
	myApp.MountMemoryStores(myApp.GetMemoryStoreKey())

	myApp.setAnteHandler(appOpts)
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

func (app *App) setAnteHandler(appOpts servertypes.AppOptions) {
	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	BypassMinFeeMsgTypes := cast.ToStringSlice(appOpts.Get(fxcfg.BypassMinFeeMsgTypesKey))
	MaxBypassMinFeeMsgGasUsage := cast.ToUint64(appOpts.Get(fxcfg.BypassMinFeeMsgMaxGasUsageKey))
	anteOptions := fxante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		IbcKeeper:       app.IBCKeeper,
		GovKeeper:       app.GovKeeper,
		SignModeHandler: app.txConfig.SignModeHandler(),
		SigGasConsumer:  fxante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  maxGasWanted,
		TxFeeChecker:    fxante.NewCheckTxFeees(BypassMinFeeMsgTypes, MaxBypassMinFeeMsgGasUsage).Check,
		DisabledAuthzMsgs: []string{
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&vestingtypes.MsgCreateVestingAccount{}),
		},
		PendingTxListener: app.onPendingTx,
	}

	if err := anteOptions.Validate(); err != nil {
		panic(fmt.Errorf("failed to ante options validate: %w", err))
	}

	app.SetAnteHandler(fxante.NewAnteHandler(anteOptions))
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

// Configurator returns the configurator of the app
func (app *App) Configurator() module.Configurator {
	return app.configurator
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// InterfaceRegistry returns InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns App's TxConfig
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *App) DefaultGenesis() map[string]json.RawMessage {
	return NewDefAppGenesisByDenom(fxtypes.DefaultDenom, app.appCodec)
}

// RegisterServices registers all module services
func (app *App) RegisterServices(cfg module.Configurator) {
	app.mm.RegisterServices(cfg)

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

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Deprecated: Register gas price queries routes from grpc-gateway.
	gaspricev1.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	gaspricev2.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Deprecated: Register gravity queries routes from grpc-gateway.
	gravity.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

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

// RegisterNodeService registers the node gRPC service on the provided application gRPC query router.
func (app *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// RegisterPendingTxListener is used by json-rpc server to listen to pending transactions callback.
func (app *App) RegisterPendingTxListener(listener evmante.PendingTxListener) {
	app.pendingTxListeners = append(app.pendingTxListeners, listener)
}

func (app *App) onPendingTx(hash common.Hash) {
	for _, listener := range app.pendingTxListeners {
		listener(hash)
	}
}

// << only for test, do not use in production >>

// SimulationManager NOTE: This is solely to be used for testing purposes.
func (app *App) SimulationManager() *module.SimulationManager {
	return &module.SimulationManager{}
}

// LegacyAmino NOTE: This is solely to be used for testing purposes.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec NOTE: This is solely to be used for testing purposes.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// GetModules NOTE: This is solely to be used for testing purposes.
func (app *App) GetModules() map[string]module.AppModule {
	modules := make(map[string]module.AppModule, len(app.mm.Modules))
	for name, mod := range app.mm.Modules {
		modules[name] = mod.(module.AppModule)
	}
	return modules
}

// GetOrderBeginBlockersModules NOTE: This is solely to be used for testing purposes..
func (app *App) GetOrderBeginBlockersModules() []string {
	return app.mm.OrderBeginBlockers
}

// GetOrderEndBlockersModules NOTE: This is solely to be used for testing purposes.
func (app *App) GetOrderEndBlockersModules() []string {
	return app.mm.OrderEndBlockers
}

// GetOrderInitGenesisModules NOTE: This is solely to be used for testing purposes.
func (app *App) GetOrderInitGenesisModules() []string {
	return app.mm.OrderInitGenesis
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(string) interface{} {
	return nil
}
