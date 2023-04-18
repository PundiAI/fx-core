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
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/functionx/fx-core/v4/x/arbitrum"
	arbitrumtypes "github.com/functionx/fx-core/v4/x/arbitrum/types"
	"github.com/functionx/fx-core/v4/x/avalanche"
	avalanchetypes "github.com/functionx/fx-core/v4/x/avalanche/types"
	"github.com/functionx/fx-core/v4/x/bsc"
	bsctypes "github.com/functionx/fx-core/v4/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v4/x/crosschain/types"
	"github.com/functionx/fx-core/v4/x/erc20"
	erc20client "github.com/functionx/fx-core/v4/x/erc20/client"
	erc20types "github.com/functionx/fx-core/v4/x/erc20/types"
	"github.com/functionx/fx-core/v4/x/eth"
	ethtypes "github.com/functionx/fx-core/v4/x/eth/types"
	fxevm "github.com/functionx/fx-core/v4/x/evm"
	fxgov "github.com/functionx/fx-core/v4/x/gov"
	fxibctransfer "github.com/functionx/fx-core/v4/x/ibc/applications/transfer"
	fxibctransfertypes "github.com/functionx/fx-core/v4/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/v4/x/migrate"
	migratetypes "github.com/functionx/fx-core/v4/x/migrate/types"
	"github.com/functionx/fx-core/v4/x/optimism"
	optimismtypes "github.com/functionx/fx-core/v4/x/optimism/types"
	"github.com/functionx/fx-core/v4/x/polygon"
	polygontypes "github.com/functionx/fx-core/v4/x/polygon/types"
	fxstaking "github.com/functionx/fx-core/v4/x/staking"
	"github.com/functionx/fx-core/v4/x/tron"
	trontypes "github.com/functionx/fx-core/v4/x/tron/types"
)

func init() {
	crosschaintypes.RegisterValidateBasic(bsctypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(polygontypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(avalanchetypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(ethtypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(trontypes.ModuleName, trontypes.TronMsgValidate{})
	crosschaintypes.RegisterValidateBasic(arbitrumtypes.ModuleName, crosschaintypes.MsgValidate{})
	crosschaintypes.RegisterValidateBasic(optimismtypes.ModuleName, crosschaintypes.MsgValidate{})
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
	bsctypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
	polygontypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
	avalanchetypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	ethtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
	trontypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
	arbitrumtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	optimismtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
	erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
}

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	fxstaking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	fxgov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		erc20client.RegisterCoinProposalHandler,
		erc20client.RegisterERC20ProposalHandler,
		erc20client.ToggleTokenConversionProposalHandler,
		erc20client.UpdateDenomAliasProposalHandler,
	}),
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
	bsc.AppModuleBasic{},
	polygon.AppModuleBasic{},
	avalanche.AppModuleBasic{},
	eth.AppModuleBasic{},
	tron.AppModuleBasic{},
	arbitrum.AppModule{},
	optimism.AppModule{},
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
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		fxstaking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),

		bsc.NewAppModule(app.BscKeeper, app.GetSubspace(bsctypes.ModuleName)),
		polygon.NewAppModule(app.PolygonKeeper, app.GetSubspace(polygontypes.ModuleName)),
		avalanche.NewAppModule(app.AvalancheKeeper, app.GetSubspace(avalanchetypes.ModuleName)),
		eth.NewAppModule(app.EthKeeper, app.GetSubspace(ethtypes.ModuleName)),
		tron.NewAppModule(app.TronKeeper, app.GetSubspace(trontypes.ModuleName)),
		arbitrum.NewAppModule(app.ArbitrumKeeper),
		optimism.NewAppModule(app.OptimismKeeper),
		fxevm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.LegacyAmino(), app.GetKey(paramstypes.StoreKey), app.GetSubspace(evmtypes.ModuleName)),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
		erc20.NewAppModule(app.Erc20Keeper, app.GetSubspace(erc20types.ModuleName)),
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
		upgradetypes.ModuleName,    // *
		capabilitytypes.ModuleName, // *
		minttypes.ModuleName,       // *
		distrtypes.ModuleName,      // *
		slashingtypes.ModuleName,   // *
		evidencetypes.ModuleName,   // *
		stakingtypes.ModuleName,    // *
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName, // *
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,

		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,
		arbitrumtypes.ModuleName,
		optimismtypes.ModuleName,

		feemarkettypes.ModuleName, // *
		evmtypes.ModuleName,
		erc20types.ModuleName,
		migratetypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,  // *
		govtypes.ModuleName,     // *
		stakingtypes.ModuleName, // *
		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibchost.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName, // *
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		crosschaintypes.ModuleName, // *
		bsctypes.ModuleName,        // *
		trontypes.ModuleName,       // *
		polygontypes.ModuleName,    // *
		avalanchetypes.ModuleName,  // *
		ethtypes.ModuleName,        // *
		arbitrumtypes.ModuleName,   // *
		optimismtypes.ModuleName,   // *

		evmtypes.ModuleName,
		feemarkettypes.ModuleName, // *
		erc20types.ModuleName,     // *
		migratetypes.ModuleName,
	}
}

func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		ibctransfertypes.ModuleName,
		fxibctransfertypes.CompatibleModuleName,
		ibchost.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		avalanchetypes.ModuleName,
		ethtypes.ModuleName,
		arbitrumtypes.ModuleName,
		optimismtypes.ModuleName,

		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		erc20types.ModuleName,
		migratetypes.ModuleName,
	}
}
