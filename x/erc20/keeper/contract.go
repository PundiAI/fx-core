package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/contracts"
)

func (k Keeper) UpgradeSystemContract(ctx sdk.Context) error {
	ctx.Logger().Info("upgrade system contract", "height", ctx.BlockHeight())
	for _, contract := range contracts.GetUpgradeContracts(ctx.BlockHeight()) {
		return k.evmKeeper.CreateContractWithCode(ctx, contract.Address, contract.Code)
	}
	return nil
}
