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
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfer "github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/functionx/fx-core/v3/x/avalanche"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	"github.com/functionx/fx-core/v3/x/bsc"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/erc20"
	erc20client "github.com/functionx/fx-core/v3/x/erc20/client"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	"github.com/functionx/fx-core/v3/x/eth"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	fxevm "github.com/functionx/fx-core/v3/x/evm"
	fxgov "github.com/functionx/fx-core/v3/x/gov"
	"github.com/functionx/fx-core/v3/x/gravity"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	fxibctransfer "github.com/functionx/fx-core/v3/x/ibc/applications/transfer"
	fxibctransfertypes "github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/v3/x/migrate"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	"github.com/functionx/fx-core/v3/x/polygon"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	"github.com/functionx/fx-core/v3/x/tron"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func init() {
	crosschaintypes.RegisterValidateBasic(bsctypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(polygontypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(avalanchetypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(ethtypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(trontypes.ModuleName, trontypes.TronMsgValidate{})
}

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	// used for secure addition and subtraction of balance using module account
	gravitytypes.ModuleName:   nil,
	bsctypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	polygontypes.ModuleName:   {authtypes.Minter, authtypes.Burner},
	avalanchetypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	ethtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	trontypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	evmtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	erc20types.ModuleName:     {authtypes.Minter, authtypes.Burner},
}

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
		erc20client.RegisterCoinProposalHandler,
		erc20client.RegisterERC20ProposalHandler,
		erc20client.ToggleTokenConversionProposalHandler,
		erc20client.UpdateDenomAliasProposalHandler,
	}...),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	fxibctransfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	// this line is used by starport scaffolding # stargate/app/moduleBasic
	gravity.AppModuleBasic{},
	bsc.AppModuleBasic{},
	polygon.AppModuleBasic{},
	avalanche.AppModuleBasic{},
	eth.AppModuleBasic{},
	tron.AppModuleBasic{},
	fxevm.AppModuleBasic{},
	feemarket.AppModuleBasic{},
	erc20.AppModuleBasic{},
	migrate.AppModuleBasic{},
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
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		fxgov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),

		gravity.NewAppModule(app.EthKeeper, app.GravityMigrator),
		bsc.NewAppModule(app.BscKeeper),
		polygon.NewAppModule(app.PolygonKeeper),
		avalanche.NewAppModule(app.AvalancheKeeper),
		eth.NewAppModule(app.EthKeeper),
		tron.NewAppModule(app.TronKeeper),
		fxevm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.LegacyAmino(), app.GetKey(paramstypes.StoreKey)),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20.NewAppModule(app.Erc20Keeper, app.IBCKeeper.ChannelKeeper),
		migrate.NewAppModule(app.MigrateKeeper),
		fxibctransfer.NewAppModule(app.FxTransferKeeper),
		ibctransfer.NewAppModule(app.IBCTransferKeeper),
	}
}

// orderBeginBlockers Tell the app's module manager how to set the order of
// BeginBlockers, which are run at the beginning of every block.
func orderBeginBlockers() []string {
	return []string{
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibchost.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		capabilitytypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,

		gravitytypes.ModuleName,
		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,

		erc20types.ModuleName,
		migratetypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibchost.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		gravitytypes.ModuleName,
		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,

		erc20types.ModuleName,
		migratetypes.ModuleName,
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
		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibchost.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		evmtypes.ModuleName,
		feemarkettypes.ModuleName,

		gravitytypes.ModuleName,
		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,

		erc20types.ModuleName,
		migratetypes.ModuleName,
	}
}
