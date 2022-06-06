package v2

import (
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/x/erc20/types"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(
	mm *module.Manager, configurator module.Configurator,
	bankStoreKey *sdk.KVStoreKey, bankKeeper bankKeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
	evmKeeper *evmkeeper.Keeper, erc20Keeper erc20keeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// update FX metadata
		UpdateFXMetadata(ctx, bankKeeper, bankStoreKey)

		// set max expected block time parameter. Replace the default with your expected value
		// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
		ibcKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

		for n, m := range mm.Modules {
			if initGenesis[n] {
				continue
			}
			if v, ok := runMigrates[n]; ok {
				fromVM[n] = v
				continue
			}
			fromVM[n] = m.ConsensusVersion()
		}

		ctx.Logger().Info("start to run module v2 migrations...")

		toVersion, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// init logic contract
		for _, contract := range fxtypes.GetInitContracts() {
			if len(contract.Code) <= 0 || contract.Address == common.HexToAddress(fxtypes.EmptyEvmAddress) {
				return nil, errors.New("invalid contract")
			}
			if err := evmKeeper.CreateContractWithCode(ctx, contract.Address, contract.Code); err != nil {
				return nil, err
			}
		}

		// register coin
		for _, metadata := range fxtypes.GetMetadata() {
			ctx.Logger().Info("add metadata", "coin", metadata.String())
			pair, err := erc20Keeper.RegisterCoin(ctx, metadata)
			if err != nil {
				return nil, err
			}
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				erc20types.EventTypeRegisterCoin,
				sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
			))
		}

		return toVersion, nil
	}
}

func UpdateFXMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper, key *types.KVStoreKey) {
	//delete fx
	deleteMetadata(ctx, key, strings.ToLower(fxtypes.DefaultDenom))
	//set FX
	md := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if err := md.Validate(); err != nil {
		panic("invalid FX metadata")
	}
	bankKeeper.SetDenomMetaData(ctx, md)
}

func deleteMetadata(ctx sdk.Context, key *types.KVStoreKey, base ...string) {
	store := ctx.KVStore(key)
	for _, b := range base {
		store.Delete(banktypes.DenomMetadataKey(b))
	}
}
