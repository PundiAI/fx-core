package keepers

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	ibctransfer "github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	srvflags "github.com/evmos/ethermint/server/flags"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/spf13/cast"

	fxtypes "github.com/functionx/fx-core/v3/types"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/erc20"
	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	fxevmkeeper "github.com/functionx/fx-core/v3/x/evm/keeper"
	fxgovkeeper "github.com/functionx/fx-core/v3/x/gov/keeper"
	gravitykeeper "github.com/functionx/fx-core/v3/x/gravity/keeper"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	fxtransfer "github.com/functionx/fx-core/v3/x/ibc/applications/transfer"
	fxtransferkeeper "github.com/functionx/fx-core/v3/x/ibc/applications/transfer/keeper"
	"github.com/functionx/fx-core/v3/x/ibc/ibcrouter"
	migratekeeper "github.com/functionx/fx-core/v3/x/migrate/keeper"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	tronkeeper "github.com/functionx/fx-core/v3/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
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
	AuthzKeeper       authzkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	GravityMigrator  gravitykeeper.Migrator
	CrosschainKeeper crosschainkeeper.RouterKeeper
	BscKeeper        crosschainkeeper.Keeper
	PolygonKeeper    crosschainkeeper.Keeper
	AvalancheKeeper  crosschainkeeper.Keeper
	EthKeeper        crosschainkeeper.Keeper
	TronKeeper       tronkeeper.Keeper

	EvmKeeper       *fxevmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	Erc20Keeper     erc20keeper.Keeper
	MigrateKeeper   migratekeeper.Keeper

	TransferModule ibctransfer.AppModule

	LegacyAmino *codec.LegacyAmino
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
	appKeepers := &AppKeepers{LegacyAmino: legacyAmino}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(
		appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
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
		authtypes.ProtoBaseAccount, maccPerms,
	)
	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
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
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	)
	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
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
		&stakingKeeper,
		authtypes.FeeCollectorName,
		blockedAddress,
	)
	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[slashingtypes.StoreKey],
		&stakingKeeper,
		appKeepers.GetSubspace(slashingtypes.ModuleName),
	)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
	)

	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp)

	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		stakingKeeper,
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
		appKeepers.GetSubspace(feemarkettypes.ModuleName),
		appKeepers.keys[feemarkettypes.StoreKey],
		appKeepers.tkeys[feemarkettypes.TransientKey],
	)

	appKeepers.EvmKeeper = fxevmkeeper.NewKeeper(
		evmkeeper.NewKeeper(
			appCodec,
			appKeepers.keys[evmtypes.StoreKey],
			appKeepers.tkeys[evmtypes.TransientKey],
			appKeepers.GetSubspace(evmtypes.ModuleName),
			appKeepers.AccountKeeper,
			appKeepers.BankKeeper,
			stakingKeeper,
			appKeepers.FeeMarketKeeper,
			cast.ToString(appOpts.Get(srvflags.EVMTracer)),
		),
		appKeepers.AccountKeeper,
	)

	erc20Keeper := erc20keeper.NewKeeper(
		appKeepers.keys[erc20types.StoreKey],
		appCodec,
		appKeepers.GetSubspace(erc20types.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.EvmKeeper,
		appKeepers.IBCTransferKeeper,
	)

	appKeepers.GravityMigrator = gravitykeeper.NewMigrator(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.keys[gravitytypes.StoreKey],
		appKeepers.keys[ethtypes.StoreKey],
		stakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	// init cross chain module
	appKeepers.BscKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		bsctypes.ModuleName,
		appKeepers.keys[bsctypes.StoreKey],
		appKeepers.GetSubspace(bsctypes.ModuleName),
		stakingKeeper,
		stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		erc20Keeper)

	appKeepers.PolygonKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		polygontypes.ModuleName,
		appKeepers.keys[polygontypes.StoreKey],
		appKeepers.GetSubspace(polygontypes.ModuleName),
		stakingKeeper,
		stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		erc20Keeper)

	appKeepers.AvalancheKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		avalanchetypes.ModuleName,
		appKeepers.keys[avalanchetypes.StoreKey],
		appKeepers.GetSubspace(avalanchetypes.ModuleName),
		stakingKeeper,
		stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		erc20Keeper)

	appKeepers.EthKeeper = crosschainkeeper.NewKeeper(
		appCodec,
		ethtypes.ModuleName,
		appKeepers.keys[ethtypes.StoreKey],
		appKeepers.GetSubspace(ethtypes.ModuleName),
		stakingKeeper,
		stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		erc20Keeper)

	appKeepers.TronKeeper = tronkeeper.NewKeeper(crosschainkeeper.NewKeeper(
		appCodec,
		trontypes.ModuleName,
		appKeepers.keys[trontypes.StoreKey],
		appKeepers.GetSubspace(trontypes.ModuleName),
		stakingKeeper,
		stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrkeeper.NewMsgServerImpl(appKeepers.DistrKeeper),
		appKeepers.BankKeeper,
		appKeepers.IBCTransferKeeper,
		erc20Keeper))

	// add cross-chain router
	crosschainRouter := crosschainkeeper.NewRouter()
	crosschainRouter.
		AddRoute(bsctypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.BscKeeper)).
		AddRoute(polygontypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.PolygonKeeper)).
		AddRoute(avalanchetypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.AvalancheKeeper)).
		AddRoute(ethtypes.ModuleName, crosschainkeeper.NewModuleHandler(appKeepers.EthKeeper)).
		AddRoute(trontypes.ModuleName, tronkeeper.NewModuleHandler(appKeepers.TronKeeper))

	appKeepers.CrosschainKeeper = crosschainkeeper.NewRouterKeeper(crosschainRouter)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(appKeepers.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(crosschaintypes.RouterKey, crosschain.NewChainProposalHandler(appKeepers.CrosschainKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(erc20Keeper))

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&stakingKeeper,
		govRouter,
	)
	appKeepers.GovKeeper = fxgovkeeper.NewKeeper(
		appKeepers.BankKeeper,
		&stakingKeeper,
		appKeepers.keys[govtypes.StoreKey],
		govKeeper,
	)

	transferRouter := fxtypes.NewRouter().
		AddRoute(gravitytypes.ModuleName, appKeepers.EthKeeper). // legacy router
		AddRoute(ethtypes.ModuleName, appKeepers.EthKeeper).
		AddRoute(bsctypes.ModuleName, appKeepers.BscKeeper).
		AddRoute(polygontypes.ModuleName, appKeepers.PolygonKeeper).
		AddRoute(trontypes.ModuleName, appKeepers.TronKeeper).
		AddRoute(avalanchetypes.ModuleName, appKeepers.AvalancheKeeper)
	appKeepers.Erc20Keeper = erc20Keeper.SetRouter(*transferRouter)

	appKeepers.EvmKeeper.SetHooks(appKeepers.Erc20Keeper.EVMHooks())

	ibcTransferRouter := transferRouter.
		AddRoute(erc20types.ModuleName, appKeepers.Erc20Keeper)
	appKeepers.FxTransferKeeper = appKeepers.FxTransferKeeper.SetRouter(*ibcTransferRouter)
	appKeepers.FxTransferKeeper = appKeepers.FxTransferKeeper.SetRefundHook(appKeepers.Erc20Keeper)

	ibcTransferModule := ibctransfer.NewIBCModule(appKeepers.IBCTransferKeeper)
	transferIBCModule := fxtransfer.NewIBCMiddleware(appKeepers.FxTransferKeeper, ibcTransferModule)
	ibcRouterMiddleware := ibcrouter.NewIBCMiddleware(transferIBCModule, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCTransferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibcRouterMiddleware)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

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
		stakingKeeper,
		appKeepers.SlashingKeeper,
	)

	return appKeepers
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	paramsKeeper.Subspace(bsctypes.ModuleName)
	paramsKeeper.Subspace(polygontypes.ModuleName)
	paramsKeeper.Subspace(avalanchetypes.ModuleName)
	paramsKeeper.Subspace(ethtypes.ModuleName)
	paramsKeeper.Subspace(trontypes.ModuleName)

	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(erc20types.ModuleName)
	return paramsKeeper
}
