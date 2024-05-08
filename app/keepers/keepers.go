package keepers

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
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
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	srvflags "github.com/evmos/ethermint/server/flags"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/spf13/cast"

	fxtypes "github.com/functionx/fx-core/v7/types"
	arbitrumtypes "github.com/functionx/fx-core/v7/x/arbitrum/types"
	fxauthzkeeper "github.com/functionx/fx-core/v7/x/authz/keeper"
	avalanchetypes "github.com/functionx/fx-core/v7/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	"github.com/functionx/fx-core/v7/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20"
	erc20keeper "github.com/functionx/fx-core/v7/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	fxevmkeeper "github.com/functionx/fx-core/v7/x/evm/keeper"
	fxgovkeeper "github.com/functionx/fx-core/v7/x/gov/keeper"
	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
	gravitytypes "github.com/functionx/fx-core/v7/x/gravity/types"
	fxtransfer "github.com/functionx/fx-core/v7/x/ibc/applications/transfer"
	fxtransferkeeper "github.com/functionx/fx-core/v7/x/ibc/applications/transfer/keeper"
	fxibctransfertypes "github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/v7/x/ibc/ibcrouter"
	layer2types "github.com/functionx/fx-core/v7/x/layer2/types"
	migratekeeper "github.com/functionx/fx-core/v7/x/migrate/keeper"
	migratetypes "github.com/functionx/fx-core/v7/x/migrate/types"
	optimismtypes "github.com/functionx/fx-core/v7/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v7/x/polygon/types"
	fxslashingkeeper "github.com/functionx/fx-core/v7/x/slashing/keeper"
	fxstakingkeeper "github.com/functionx/fx-core/v7/x/staking/keeper"
	tronkeeper "github.com/functionx/fx-core/v7/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

type CrossChainKeepers struct {
	BscKeeper       crosschainkeeper.Keeper
	PolygonKeeper   crosschainkeeper.Keeper
	AvalancheKeeper crosschainkeeper.Keeper
	EthKeeper       crosschainkeeper.Keeper
	TronKeeper      tronkeeper.Keeper
	ArbitrumKeeper  crosschainkeeper.Keeper
	OptimismKeeper  crosschainkeeper.Keeper
	Layer2Keeper    crosschainkeeper.Keeper
}

type AppKeepers struct {
	legacyAmino *codec.LegacyAmino
	appCodec    codec.Codec

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    fxstakingkeeper.Keeper
	SlashingKeeper   fxslashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        fxgovkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper

	// IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCKeeper         *ibckeeper.Keeper
	EvidenceKeeper    evidencekeeper.Keeper
	FxTransferKeeper  fxtransferkeeper.Keeper
	IBCTransferKeeper ibctransferkeeper.Keeper
	FeeGrantKeeper    feegrantkeeper.Keeper
	AuthzKeeper       fxauthzkeeper.Keeper

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
	appOpts servertypes.AppOptions,
) *AppKeepers {
	appKeepers := &AppKeepers{
		legacyAmino: legacyAmino,
		appCodec:    appCodec,
	}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.generateKeys()

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(
		appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.CapabilityKeeper.Seal()

	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		appKeepers.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
		fxtypes.AddressPrefix,
	)
	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)
	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(banktypes.ModuleName),
		blockedAddress,
	)
	appKeepers.StakingKeeper = fxstakingkeeper.NewKeeper(stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	), appKeepers.keys[stakingtypes.StoreKey], appCodec, appKeepers.AccountKeeper)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		&appKeepers.StakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
	)
	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.GetSubspace(distrtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
	)
	appKeepers.SlashingKeeper = fxslashingkeeper.NewKeeper(
		slashingkeeper.NewKeeper(
			appCodec,
			appKeepers.keys[slashingtypes.StoreKey],
			&appKeepers.StakingKeeper,
			appKeepers.GetSubspace(slashingtypes.ModuleName),
		), appKeepers.keys[slashingtypes.StoreKey],
	)
	appKeepers.StakingKeeper = *appKeepers.StakingKeeper.SetSlashingKeeper(appKeepers.SlashingKeeper)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
	)

	// register the staking hooks
	// NOTE: appKeepers.StakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper.Keeper = *appKeepers.StakingKeeper.Keeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

	appKeepers.AuthzKeeper = fxauthzkeeper.NewKeeper(authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	), appKeepers.keys[authzkeeper.StoreKey], appCodec)

	// grant privileges
	appKeepers.StakingKeeper = *appKeepers.StakingKeeper.SetAuthzKeeper(appKeepers.AuthzKeeper)

	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		&appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
	)
	appKeepers.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)
	appKeepers.FxTransferKeeper = fxtransferkeeper.NewKeeper(
		appKeepers.IBCTransferKeeper,
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)
	appKeepers.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		appKeepers.keys[feemarkettypes.StoreKey],
		appKeepers.tkeys[feemarkettypes.TransientKey],
		appKeepers.GetSubspace(feemarkettypes.ModuleName),
	)
	appKeepers.EvmKeeper = fxevmkeeper.NewKeeper(
		evmkeeper.NewKeeper(
			appCodec,
			appKeepers.keys[evmtypes.StoreKey],
			appKeepers.tkeys[evmtypes.TransientKey],
			authtypes.NewModuleAddress(govtypes.ModuleName),
			appKeepers.AccountKeeper,
			appKeepers.BankKeeper,
			&appKeepers.StakingKeeper,
			appKeepers.FeeMarketKeeper,
			cast.ToString(appOpts.Get(srvflags.EVMTracer)),
			appKeepers.GetSubspace(evmtypes.ModuleName),
		),
		appKeepers.AccountKeeper,
	)

	appKeepers.Erc20Keeper = erc20keeper.NewKeeper(
		appKeepers.keys[erc20types.StoreKey],
		appCodec,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.EvmKeeper,
		appKeepers.IBCTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.TronKeeper = tronkeeper.NewKeeper(crosschainkeeper.NewKeeper(
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	))

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
		AddRoute(trontypes.ModuleName, tronkeeper.NewModuleHandler(appKeepers.TronKeeper))

	appKeepers.CrosschainRouterKeeper = crosschainkeeper.NewRouterKeeper(crosschainRouter)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(appKeepers.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(crosschaintypes.RouterKey, crosschain.NewCrosschainProposalHandler(appKeepers.CrosschainRouterKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(appKeepers.Erc20Keeper))

	govConfig := fxgovtypes.DefaultConfig()
	_govKeeper := govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		govRouter,
		bApp.MsgServiceRouter(),
		govConfig.Config,
	)
	_govKeeper = *_govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	appKeepers.GovKeeper = fxgovkeeper.NewKeeper(
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.keys[govtypes.StoreKey],
		_govKeeper,
		govConfig,
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	ibcTransferRouter := fxibctransfertypes.NewRouter().
		AddRoute(gravitytypes.ModuleName, appKeepers.EthKeeper).
		AddRoute(ethtypes.ModuleName, appKeepers.EthKeeper).
		AddRoute(bsctypes.ModuleName, appKeepers.BscKeeper).
		AddRoute(polygontypes.ModuleName, appKeepers.PolygonKeeper).
		AddRoute(trontypes.ModuleName, appKeepers.TronKeeper).
		AddRoute(avalanchetypes.ModuleName, appKeepers.AvalancheKeeper).
		AddRoute(arbitrumtypes.ModuleName, appKeepers.ArbitrumKeeper).
		AddRoute(optimismtypes.ModuleName, appKeepers.OptimismKeeper).
		AddRoute(layer2types.ModuleName, appKeepers.Layer2Keeper).
		AddRoute(erc20types.ModuleName, appKeepers.Erc20Keeper)
	appKeepers.FxTransferKeeper = appKeepers.FxTransferKeeper.SetRouter(*ibcTransferRouter)
	appKeepers.FxTransferKeeper = appKeepers.FxTransferKeeper.SetRefundHook(appKeepers.Erc20Keeper)
	appKeepers.FxTransferKeeper = appKeepers.FxTransferKeeper.SetErc20Keeper(appKeepers.Erc20Keeper)
	appKeepers.FxTransferKeeper = appKeepers.FxTransferKeeper.SetEvmKeeper(appKeepers.EvmKeeper)

	ibcTransferModule := ibctransfer.NewIBCModule(appKeepers.IBCTransferKeeper)
	transferIBCModule := fxtransfer.NewIBCMiddleware(appKeepers.FxTransferKeeper, ibcTransferModule)
	ibcRouterMiddleware := ibcrouter.NewIBCMiddleware(transferIBCModule, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCTransferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibcRouterMiddleware)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	appKeepers.MigrateKeeper = migratekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[migratetypes.StoreKey],
		appKeepers.AccountKeeper,
	)
	appKeepers.MigrateKeeper = appKeepers.MigrateKeeper.SetMigrateI(
		migratekeeper.NewBankMigrate(appKeepers.BankKeeper),
		migratekeeper.NewDistrStakingMigrate(appKeepers.keys[distrtypes.StoreKey], appKeepers.keys[stakingtypes.StoreKey], appKeepers.StakingKeeper),
		migratekeeper.NewGovMigrate(appKeepers.keys[govtypes.StoreKey], appKeepers.GovKeeper),
	)

	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)

	// set evm precompiled contracts
	appKeepers.EvmPrecompiled()

	return appKeepers
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// LegacyAmino NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (appKeepers *AppKeepers) LegacyAmino() *codec.LegacyAmino {
	return appKeepers.legacyAmino
}

// AppCodec NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (appKeepers *AppKeepers) AppCodec() codec.Codec {
	return appKeepers.appCodec
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)

	paramsKeeper.Subspace(bsctypes.ModuleName)
	paramsKeeper.Subspace(polygontypes.ModuleName)
	paramsKeeper.Subspace(avalanchetypes.ModuleName)
	paramsKeeper.Subspace(ethtypes.ModuleName)
	paramsKeeper.Subspace(trontypes.ModuleName)

	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName).WithKeyTable(evmtypes.ParamKeyTable()) // nolint: staticcheck
	paramsKeeper.Subspace(feemarkettypes.ModuleName).WithKeyTable(feemarkettypes.ParamKeyTable())
	paramsKeeper.Subspace(erc20types.ModuleName)
	return paramsKeeper
}
