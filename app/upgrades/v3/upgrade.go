package v3

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	evmkeeper "github.com/functionx/fx-core/v3/x/evm/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// cache context
		cacheCtx, commit := ctx.CacheContext()

		// update wfx logic code
		updateWFXLogicCode(cacheCtx, keepers.EvmKeeper)

		// update metadata alias null
		updateMetadataAliasNull(cacheCtx, keepers.BankKeeper)

		// migrate evm param RejectUnprotectedTx to AllowUnprotectedTxs
		migrateRejectUnprotectedTx(cacheCtx, keepers.LegacyAmino, keepers.GetKey(paramstypes.StoreKey))

		// delete erc20 expiration ibc transfer hash
		deleteExpirationIBCTransferHash(ctx, keepers.Erc20Keeper, keepers.IBCKeeper)

		// run migrations
		toVM := runMigrations(cacheCtx, fromVM, mm, configurator)

		// init avalanche oracles
		initAvalancheOracles(cacheCtx, keepers.AvalancheKeeper)

		// update bsc oracles
		updateBSCOracles(cacheCtx, keepers.BscKeeper)

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
	// todo need add oracles
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
	// todo need add oracles
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{}
	} else if chainId == fxtypes.TestnetChainId {
		oracles = []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
	if len(oracles) <= 0 {
		return
	}
	bscKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}

func migrateRejectUnprotectedTx(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) {
	paramStoreKeyRejectUnprotectedTx := []byte("RejectUnprotectedTx")

	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(evmtypes.ModuleName), '/'))
	bzR := paramsStore.Get(paramStoreKeyRejectUnprotectedTx)

	var rejectUnprotectedTx bool
	if err := legacyAmino.UnmarshalJSON(bzR, &rejectUnprotectedTx); err != nil {
		panic(err.Error())
	}

	allowUnprotectedTxs := !rejectUnprotectedTx
	bzA, err := legacyAmino.MarshalJSON(allowUnprotectedTxs)
	if err != nil {
		panic(err.Error())
	}

	ctx.Logger().Info("migrate params", "module", evmtypes.ModuleName,
		"from", fmt.Sprintf("%s:%v", paramStoreKeyRejectUnprotectedTx, rejectUnprotectedTx),
		"to", fmt.Sprintf("%s:%v", evmtypes.ParamStoreKeyAllowUnprotectedTxs, allowUnprotectedTxs))

	paramsStore.Delete(paramStoreKeyRejectUnprotectedTx)
	paramsStore.Set(evmtypes.ParamStoreKeyAllowUnprotectedTxs, bzA)
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

func updateWFXLogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	wfx := fxtypes.GetWFX()
	err := k.UpdateContractCode(ctx, wfx.Address, wfx.Code)
	if err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
}

func updateMetadataAliasNull(ctx sdk.Context, bk bankkeeper.Keeper) {
	logger := ctx.Logger()
	bk.IterateAllDenomMetaData(ctx, func(md banktypes.Metadata) bool {
		if len(md.DenomUnits) != 2 || len(md.DenomUnits[1].Aliases) != 1 || md.DenomUnits[1].Aliases[0] != "null" {
			return false
		}
		logger.Info("update metadata alias null", "denom", md.Base)
		md.DenomUnits[1].Aliases = []string{}
		bk.SetDenomMetaData(ctx, md)
		return false
	})
}

func deleteExpirationIBCTransferHash(ctx sdk.Context, erc20Keeper erc20keeper.Keeper, ibcKeeper *ibckeeper.Keeper) {
	logger := ctx.Logger()
	counts := make(map[string]uint64, 10)

	iter := erc20Keeper.IBCTransferHashIterator(ctx)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := bytes.TrimPrefix(iter.Key(), erc20types.KeyPrefixIBCTransfer)
		split := strings.Split(string(key), "/")
		if len(split) != 3 {
			panic(fmt.Sprintf("invalid key: %s", string(key)))
		}
		port := split[0]
		channel := split[1]
		sequence, err := strconv.ParseUint(split[2], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("parse sequence %s error %s", split[2], err.Error()))
		}

		found := ibcKeeper.ChannelKeeper.HasPacketCommitment(ctx, port, channel, sequence)
		if found {
			continue
		}
		erc20Keeper.DeleteIBCTransferHash(ctx, port, channel, sequence)

		// statistics count
		counts[fmt.Sprintf("%s/%s", port, channel)] += 1
	}
	for portChannel, count := range counts {
		logger.Info("delete expiration ibc transfer hash", "port/channel", portChannel, "count", strconv.FormatUint(count, 10))
	}
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
