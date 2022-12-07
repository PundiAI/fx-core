package v3

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
	"path/filepath"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	fxevmkeeper "github.com/functionx/fx-core/v3/x/evm/keeper"
	evmlegacyv3 "github.com/functionx/fx-core/v3/x/evm/legacy/v3"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[avalanchetypes.ModuleName] = mm.Modules[avalanchetypes.ModuleName].ConsensusVersion()
		fromVM[ethtypes.ModuleName] = mm.Modules[ethtypes.ModuleName].ConsensusVersion()

		// cache context
		cacheCtx, commit := ctx.CacheContext()

		// init avalanche oracles
		initAvalancheOracles(cacheCtx, keepers.AvalancheKeeper)

		// update bsc oracles
		updateBSCOracles(cacheCtx, keepers.BscKeeper)

		// update wfx code
		updateWFXCode(cacheCtx, keepers.EvmKeeper)

		// migrate evm param RejectUnprotectedTx to AllowUnprotectedTxs
		migrateRejectUnprotectedTx(cacheCtx, keepers.LegacyAmino, keepers.GetKey(paramstypes.StoreKey))

		// run migrations
		toVM := runMigrations(cacheCtx, fromVM, mm, configurator)

		// register coin
		registerCoin(cacheCtx, keepers.Erc20Keeper)

		//commit upgrade
		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		return toVM, nil
	}
}

func initAvalancheOracles(ctx sdk.Context, avalancheKeeper crosschainkeeper.Keeper) {
	var oracles []string
	chainId := ctx.ChainID()
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{
			"", "",
		}
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

func migrateRejectUnprotectedTx(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) {
	err := evmlegacyv3.MigrateRejectUnprotectedTx(ctx, legacyAmino, paramsKey)
	if err != nil {
		panic(fmt.Sprintf("migrate param RejectUnprotectedTx error %s", err.Error()))
	}
}

func runMigrations(ctx sdk.Context, fromVM module.VersionMap, mm *module.Manager, mc module.Configurator) module.VersionMap {
	ctx.Logger().Info("start to run module v3 migrations...")
	toVM, err := mm.RunMigrations(ctx, mc, fromVM)
	if err != nil {
		panic(fmt.Sprintf("run migrations: %s", err.Error()))
	}
	return toVM
}

func registerCoin(ctx sdk.Context, k erc20keeper.Keeper) {
	for _, metadata := range GetMetadata(ctx.ChainID()) {
		ctx.Logger().Info("add metadata", "coin", metadata.String())
		pair, err := k.RegisterCoin(ctx, metadata)
		if err != nil {
			panic(fmt.Sprintf("register %s: %s", metadata.Base, err.Error()))
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			erc20types.EventTypeRegisterCoin,
			sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
		))
	}
}

// updateWFXCode update wfx code
func updateWFXCode(ctx sdk.Context, ek *fxevmkeeper.Keeper) {
	logger := ctx.Logger()
	wfx := fxtypes.GetWFX()
	codeHash := crypto.Keccak256Hash(wfx.Code)

	logger.Info("update wfx code", "address", wfx.Address.String(), "version", wfx.Version, "code-hash", codeHash.String())
	acc := ek.GetAccount(ctx, wfx.Address)
	if acc == nil {
		panic(fmt.Sprintf("account %s not found", wfx.Address.String()))
	}
	acc.CodeHash = codeHash.Bytes()
	ek.SetCode(ctx, acc.CodeHash, wfx.Code)
	if err := ek.SetAccount(ctx, wfx.Address, *acc); err != nil {
		panic(fmt.Sprintf("evm set account %s error %s", wfx.Address.String(), err.Error()))
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventUpdateContract,
		sdk.NewAttribute(AttributeKeyContract, wfx.Address.String()),
		sdk.NewAttribute(AttributeKeyVersion, wfx.Version),
	))
}

// PreUpgradeCmd called by cosmovisor
func PreUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "fxv3 pre-upgrade, called by cosmovisor, before migrations upgrade",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			serverCtx.Logger.Info("pre-upgrade", "action", "update app.toml and config.toml")

			rootDir := serverCtx.Config.RootDir
			fileName := filepath.Join(rootDir, "config", "config.toml")
			tmcfg.WriteConfigFile(fileName, serverCtx.Config)

			config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
			appConfig := fxcfg.DefaultConfig()
			if err := serverCtx.Viper.Unmarshal(appConfig); err != nil {
				return err
			}
			fileName = filepath.Join(rootDir, "config", "app.toml")
			config.WriteConfigFile(fileName, appConfig)

			clientCtx := client.GetClientContextFromCmd(cmd)
			return clientCtx.PrintString("fxv3 pre-upgrade success")
		},
	}
	return cmd
}
