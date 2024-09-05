package app

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/functionx/fx-core/v8/x/arbitrum"
	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	"github.com/functionx/fx-core/v8/x/avalanche"
	avalanchetypes "github.com/functionx/fx-core/v8/x/avalanche/types"
	"github.com/functionx/fx-core/v8/x/bsc"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	crosschainclient "github.com/functionx/fx-core/v8/x/crosschain/client"
	"github.com/functionx/fx-core/v8/x/erc20"
	erc20client "github.com/functionx/fx-core/v8/x/erc20/client"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	"github.com/functionx/fx-core/v8/x/eth"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	fxevm "github.com/functionx/fx-core/v8/x/evm"
	fxgov "github.com/functionx/fx-core/v8/x/gov"
	fxibctransfer "github.com/functionx/fx-core/v8/x/ibc/applications/transfer"
	fxibctransfertypes "github.com/functionx/fx-core/v8/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/v8/x/layer2"
	layer2types "github.com/functionx/fx-core/v8/x/layer2/types"
	"github.com/functionx/fx-core/v8/x/migrate"
	migratetypes "github.com/functionx/fx-core/v8/x/migrate/types"
	"github.com/functionx/fx-core/v8/x/optimism"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
	"github.com/functionx/fx-core/v8/x/polygon"
	polygontypes "github.com/functionx/fx-core/v8/x/polygon/types"
	fxstaking "github.com/functionx/fx-core/v8/x/staking"
	"github.com/functionx/fx-core/v8/x/tron"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	bsctypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
	polygontypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
	avalanchetypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	ethtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
	trontypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
	arbitrumtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	optimismtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	layer2types.ModuleName:         {authtypes.Minter, authtypes.Burner},
	evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
	erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
}

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	fxstaking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	fxgov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		erc20client.LegacyRegisterCoinProposalHandler,
		erc20client.LegacyRegisterERC20ProposalHandler,
		erc20client.LegacyToggleTokenConversionProposalHandler,
		erc20client.LegacyUpdateDenomAliasProposalHandler,
		crosschainclient.LegacyUpdateChainOraclesProposalHandler,
	}),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	fxibctransfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	bsc.AppModuleBasic{},
	polygon.AppModuleBasic{},
	avalanche.AppModuleBasic{},
	eth.AppModuleBasic{},
	tron.AppModuleBasic{},
	arbitrum.AppModule{},
	optimism.AppModule{},
	layer2.AppModule{},
	fxevm.AppModuleBasic{},
	feemarket.AppModuleBasic{},
	erc20.AppModuleBasic{},
	migrate.AppModuleBasic{},
	consensus.AppModuleBasic{},
)

func appModules(
	app *App,
	encodingConfig EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Codec
	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		fxgov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper.Keeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		fxstaking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper.Keeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),

		bsc.NewAppModule(app.BscKeeper),
		polygon.NewAppModule(app.PolygonKeeper),
		avalanche.NewAppModule(app.AvalancheKeeper),
		eth.NewAppModule(app.EthKeeper),
		tron.NewAppModule(app.TronKeeper),
		arbitrum.NewAppModule(app.ArbitrumKeeper),
		optimism.NewAppModule(app.OptimismKeeper),
		layer2.NewAppModule(app.Layer2Keeper),
		fxevm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.GetSubspace(evmtypes.ModuleName)),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
		erc20.NewAppModule(app.Erc20Keeper),
		migrate.NewAppModule(app.MigrateKeeper),
		fxibctransfer.NewAppModule(app.FxTransferKeeper),
		ibctransfer.NewAppModule(app.IBCTransferKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
	}
}

// orderBeginBlockers Tell the app's module manager how to set the order of
// BeginBlockers, which are run at the beginning of every block.
func orderBeginBlockers() []string {
	return []string{
		// upgrades should be run first
		upgradetypes.ModuleName,    // begin
		capabilitytypes.ModuleName, // begin
		minttypes.ModuleName,       // begin
		distrtypes.ModuleName,      // begin
		slashingtypes.ModuleName,   // begin
		evidencetypes.ModuleName,   // begin
		stakingtypes.ModuleName,    // begin
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName, // begin
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,

		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,
		arbitrumtypes.ModuleName,
		optimismtypes.ModuleName,
		layer2types.ModuleName,

		feemarkettypes.ModuleName, // begin
		evmtypes.ModuleName,       // begin
		erc20types.ModuleName,
		migratetypes.ModuleName,

		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibcexported.ModuleName,

		consensustypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,  // end
		govtypes.ModuleName,     // end
		stakingtypes.ModuleName, // end
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName, // end
		authz.ModuleName,
		feegrant.ModuleName, // end
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		bsctypes.ModuleName,       // end
		trontypes.ModuleName,      // end
		polygontypes.ModuleName,   // end
		avalanchetypes.ModuleName, // end
		ethtypes.ModuleName,       // end
		arbitrumtypes.ModuleName,  // end
		optimismtypes.ModuleName,  // end
		layer2types.ModuleName,    // end

		evmtypes.ModuleName,       // end
		feemarkettypes.ModuleName, // end
		erc20types.ModuleName,
		migratetypes.ModuleName,

		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibcexported.ModuleName,

		consensustypes.ModuleName,
	}
}

func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,
		arbitrumtypes.ModuleName,
		optimismtypes.ModuleName,
		layer2types.ModuleName,

		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		erc20types.ModuleName,
		migratetypes.ModuleName,

		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibcexported.ModuleName,

		consensustypes.ModuleName,
	}
}

func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// BlockedAccountAddrs returns all the app's blocked module account addresses.
func BlockedAccountAddrs() map[string]bool {
	modAccAddrs := ModuleAccountAddrs()
	// allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	return modAccAddrs
}
