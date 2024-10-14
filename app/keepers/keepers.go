package keepers

import (
	"os"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
	srvflags "github.com/evmos/ethermint/server/flags"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	v0evmtypes "github.com/evmos/ethermint/x/evm/migrations/v0/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/spf13/cast"

	fxtypes "github.com/functionx/fx-core/v8/types"
	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v8/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschainprecompile "github.com/functionx/fx-core/v8/x/crosschain/precompile"
	erc20keeper "github.com/functionx/fx-core/v8/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	fxevmkeeper "github.com/functionx/fx-core/v8/x/evm/keeper"
	fxgovkeeper "github.com/functionx/fx-core/v8/x/gov/keeper"
	fxgovtypes "github.com/functionx/fx-core/v8/x/gov/types"
	ibcmiddleware "github.com/functionx/fx-core/v8/x/ibc/middleware"
	ibcmiddlewarekeeper "github.com/functionx/fx-core/v8/x/ibc/middleware/keeper"
	layer2types "github.com/functionx/fx-core/v8/x/layer2/types"
	migratekeeper "github.com/functionx/fx-core/v8/x/migrate/keeper"
	migratetypes "github.com/functionx/fx-core/v8/x/migrate/types"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v8/x/polygon/types"
	fxstakingkeeper "github.com/functionx/fx-core/v8/x/staking/keeper"
	stakingprecompile "github.com/functionx/fx-core/v8/x/staking/precompile"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

type CrossChainKeepers struct {
	BscKeeper       crosschainkeeper.Keeper
	PolygonKeeper   crosschainkeeper.Keeper
	AvalancheKeeper crosschainkeeper.Keeper
	EthKeeper       crosschainkeeper.Keeper
	TronKeeper      crosschainkeeper.Keeper
	ArbitrumKeeper  crosschainkeeper.Keeper
	OptimismKeeper  crosschainkeeper.Keeper
	Layer2Keeper    crosschainkeeper.Keeper
}

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey
	objKeys map[string]*storetypes.ObjectStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *fxstakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *fxgovkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCKeeper           *ibckeeper.Keeper
	EvidenceKeeper      evidencekeeper.Keeper
	IBCMiddlewareKeeper ibcmiddlewarekeeper.Keeper
	IBCTransferKeeper   ibctransferkeeper.Keeper
	FeeGrantKeeper      feegrantkeeper.Keeper
	AuthzKeeper         authzkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	CrosschainRouterKeeper crosschainkeeper.RouterKeeper
	CrossChainKeepers

	EvmKeeper       *fxevmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	Erc20Keeper     erc20keeper.Keeper
	MigrateKeeper   migratekeeper.Keeper
}

func NewAppKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	logger log.Logger,
	appOpts servertypes.AppOptions,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	if err := bApp.RegisterStreamingServices(appOpts, appKeepers.keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	authAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]),
		authAddr,
		runtime.EventService{},
	)
	bApp.SetParamStore(appKeepers.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.CapabilityKeeper.Seal()

	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authAddr,
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[banktypes.StoreKey]),
		appKeepers.objKeys[banktypes.ObjectStoreKey],
		appKeepers.AccountKeeper,
		blockedAddress,
		authAddr,
		logger,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]),
		appKeepers.AccountKeeper,
	)

	appKeepers.StakingKeeper = fxstakingkeeper.NewKeeper(
		stakingkeeper.NewKeeper(
			appCodec,
			runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]),
			appKeepers.AccountKeeper,
			appKeepers.BankKeeper,
			authAddr,
			authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
			authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
		),
		appKeepers.keys[stakingtypes.StoreKey],
		appCodec,
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[minttypes.StoreKey]),
		appKeepers.StakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authAddr,
	)
	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authAddr,
	)
	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]),
		appKeepers.StakingKeeper,
		authAddr,
	)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[crisistypes.StoreKey]),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authAddr,
		appKeepers.AccountKeeper.AddressCodec(),
	)

	// register the staking hooks
	// NOTE: appKeepers.StakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper.Keeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]),
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(appKeepers.keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		bApp,
		authAddr,
	)
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
		authAddr,
	)
	appKeepers.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
		authAddr,
	)
	appKeepers.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		appKeepers.keys[feemarkettypes.StoreKey],
		// appKeepers.tkeys[feemarkettypes.TransientKey],
		appKeepers.GetSubspace(feemarkettypes.ModuleName),
	)

	// cross chain precompile
	precompileRouter := crosschainprecompile.NewRouter()
	evmKeeper := evmkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evmtypes.StoreKey],
		appKeepers.objKeys[evmtypes.ObjectStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.FeeMarketKeeper,
		cast.ToString(appOpts.Get(srvflags.EVMTracer)),
		appKeepers.GetSubspace(evmtypes.ModuleName),
		[]evmkeeper.CustomContractFn{
			func(_ sdk.Context, _ ethparams.Rules) vm.PrecompiledContract {
				return crosschainprecompile.NewPrecompiledContract(appKeepers.BankKeeper, appKeepers.Erc20Keeper,
					appKeepers.GovKeeper, precompileRouter)
			},
			func(_ sdk.Context, _ ethparams.Rules) vm.PrecompiledContract {
				return stakingprecompile.NewPrecompiledContract(appKeepers.BankKeeper, appKeepers.StakingKeeper,
					appKeepers.DistrKeeper, fxtypes.DefaultDenom, appKeepers.GovKeeper, appKeepers.SlashingKeeper)
			},
		},
	)
	evmKeeper.WithChainIDString(fxtypes.ChainIdWithEIP155(cast.ToString(appOpts.Get(flags.FlagChainID))))
	appKeepers.EvmKeeper = fxevmkeeper.NewKeeper(
		evmKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.Erc20Keeper = erc20keeper.NewKeeper(
		appKeepers.keys[erc20types.StoreKey],
		appCodec,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		appKeepers.IBCTransferKeeper,
		authAddr,
	)

	// init cross chain module
	appKeepers.BscKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		bsctypes.ModuleName,
		appKeepers.keys[bsctypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.PolygonKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		polygontypes.ModuleName,
		appKeepers.keys[polygontypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.AvalancheKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		avalanchetypes.ModuleName,
		appKeepers.keys[avalanchetypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.EthKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		ethtypes.ModuleName,
		appKeepers.keys[ethtypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.ArbitrumKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		arbitrumtypes.ModuleName,
		appKeepers.keys[arbitrumtypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.OptimismKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		optimismtypes.ModuleName,
		appKeepers.keys[optimismtypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.Layer2Keeper = crosschainkeeper.NewKeeper(
		appCodec,
		layer2types.ModuleName,
		appKeepers.keys[layer2types.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)
	appKeepers.TronKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		trontypes.ModuleName,
		appKeepers.keys[trontypes.StoreKey],
		appKeepers.StakingKeeper,
		stakingkeeper.NewMsgServerImpl(appKeepers.StakingKeeper.Keeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		appKeepers.Erc20Keeper,
		appKeepers.AccountKeeper,
		appKeepers.EvmKeeper,
		appKeepers.EvmKeeper,
		authAddr,
	)

	// add cross-chain router
	crosschainRouter := crosschainkeeper.NewRouter()
	crosschainRouter.
		AddRoute(bsctypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.BscKeeper)).
		AddRoute(polygontypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.PolygonKeeper)).
		AddRoute(avalanchetypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.AvalancheKeeper)).
		AddRoute(ethtypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.EthKeeper)).
		AddRoute(arbitrumtypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.ArbitrumKeeper)).
		AddRoute(optimismtypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.OptimismKeeper)).
		AddRoute(layer2types.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.Layer2Keeper)).
		AddRoute(trontypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.TronKeeper))

	appKeepers.CrosschainRouterKeeper = crosschainkeeper.NewRouterKeeper(crosschainRouter)

	// cross chain precompile
	precompileRouter.
		AddRoute(bsctypes.ModuleName, appKeepers.BscKeeper).
		AddRoute(polygontypes.ModuleName, appKeepers.PolygonKeeper).
		AddRoute(avalanchetypes.ModuleName, appKeepers.AvalancheKeeper).
		AddRoute(ethtypes.ModuleName, appKeepers.EthKeeper).
		AddRoute(arbitrumtypes.ModuleName, appKeepers.ArbitrumKeeper).
		AddRoute(optimismtypes.ModuleName, appKeepers.OptimismKeeper).
		AddRoute(layer2types.ModuleName, appKeepers.Layer2Keeper).
		AddRoute(trontypes.ModuleName, appKeepers.TronKeeper)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper))

	// provider depends on gov, so gov must be registered first
	govConfig := fxgovtypes.DefaultConfig()
	// set the MaxMetadataLen for proposals to the same value as it was pre-sdk v0.47.x
	govConfig.MaxMetadataLen = 10200
	_govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.DistrKeeper,
		bApp.MsgServiceRouter(),
		govConfig.Config,
		authAddr,
	)
	_govKeeper.SetLegacyRouter(govRouter)

	appKeepers.GovKeeper = fxgovkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.keys,
		_govKeeper,
		appCodec,
		authAddr,
	)
	appKeepers.IBCMiddlewareKeeper = ibcmiddlewarekeeper.NewKeeper(appCodec, appKeepers.EvmKeeper,
		appKeepers.EthKeeper) // TODO: replace by crosschain keeper
	ibcTransferModule := ibctransfer.NewIBCModule(appKeepers.IBCTransferKeeper)
	transferIBCModule := ibcmiddleware.NewIBCMiddleware(appKeepers.IBCMiddlewareKeeper, appKeepers.IBCKeeper.ChannelKeeper, ibcTransferModule)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	appKeepers.MigrateKeeper = migratekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[migratetypes.StoreKey],
		appKeepers.AccountKeeper,
	)
	appKeepers.MigrateKeeper = appKeepers.MigrateKeeper.SetMigrateI(
		migratekeeper.NewBankMigrate(appKeepers.BankKeeper),
		migratekeeper.NewDistrStakingMigrate(appKeepers.keys[distrtypes.StoreKey], appKeepers.keys[stakingtypes.StoreKey], appKeepers.StakingKeeper),
		migratekeeper.NewGovMigrate(appKeepers.GovKeeper, appKeepers.AccountKeeper),
	)

	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]),
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
		appKeepers.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)

	return appKeepers
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, ok := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	if !ok {
		panic("couldn't load subspace for module: " + moduleName)
	}
	return subspace
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())
	paramsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())
	paramsKeeper.Subspace(stakingtypes.ModuleName).WithKeyTable(stakingtypes.ParamKeyTable())
	paramsKeeper.Subspace(minttypes.ModuleName).WithKeyTable(minttypes.ParamKeyTable())
	paramsKeeper.Subspace(distrtypes.ModuleName).WithKeyTable(distrtypes.ParamKeyTable())
	paramsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable())
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())

	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	ibcKeyTable := ibcclienttypes.ParamKeyTable()
	ibcKeyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(ibcKeyTable)

	paramsKeeper.Subspace(evmtypes.ModuleName).WithKeyTable(v0evmtypes.ParamKeyTable())
	paramsKeeper.Subspace(feemarkettypes.ModuleName).WithKeyTable(feemarkettypes.ParamKeyTable())

	paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
	return paramsKeeper
}
