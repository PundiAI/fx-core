package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[avalanchetypes.ModuleName] = mm.Modules[avalanchetypes.ModuleName].ConsensusVersion()
		fromVM[ethtypes.ModuleName] = mm.Modules[ethtypes.ModuleName].ConsensusVersion()
		initAvalancheOracles(ctx, keepers.AvalancheKeeper)
		updateBSCOracles(ctx, keepers.BscKeeper)
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func initAvalancheOracles(ctx sdk.Context, avalancheKeeper crosschainkeeper.Keeper) {
	var oracles []string
	chainId := ctx.ChainID()
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{}
	} else if chainId == fxtypes.TestnetChainId {
		oracles = []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
	avalancheKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}

func updateBSCOracles(ctx sdk.Context, bscKeeper crosschainkeeper.Keeper) {
	var oracles []string
	chainId := ctx.ChainID()
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{}
	} else if chainId == fxtypes.TestnetChainId {
		oracles = []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
	bscKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}
