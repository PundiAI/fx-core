package keeper

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/contracts"
	fxcoretypes "github.com/functionx/fx-core/types"
)

func (k Keeper) InitSystemContract(ctx sdk.Context) error {
	network := fxcoretypes.Network()
	upgrade := contracts.GetInitConfig(network)
	if upgrade == nil {
		return errors.New("empty system contract")
	}
	return k.upgradeSystemContract(ctx, upgrade)
}
func (k Keeper) UpgradeSystemContract(ctx sdk.Context) error {
	network := fxcoretypes.Network()
	blockHeight := ctx.BlockHeight()

	bc := contracts.GetUpgradeBlockConfig(network)

	if blockHeight == bc.TestUpgradeBlock {
		if err := k.upgradeSystemContract(ctx, contracts.GetTestConfig(network)); err != nil {
			return err
		}
	}
	return nil
}
func (k Keeper) upgradeSystemContract(ctx sdk.Context, upgrade *contracts.Upgrade) error {
	if upgrade == nil {
		ctx.Logger().Info("empty upgrade config", "height", ctx.BlockHeight())
		return nil
	}
	ctx.Logger().Info("upgrade system contract", "name", upgrade.Name, "height", ctx.BlockHeight())
	for _, cfg := range upgrade.Configs {
		if cfg.ContractAddr.Hex() == contracts.EmptyAddress {
			continue
		}
		ctx.Logger().Info("upgrade contract", "address", cfg.ContractAddr.Hex())

		if err := k.CreateContractWithCode(ctx, cfg.ContractAddr, cfg.Code); err != nil {
			return err
		}
	}
	return nil
}
