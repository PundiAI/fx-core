package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) CreateBridgeAccount(ctx sdk.Context, address string) {
	accAddress := types.ExternalAddrToAccAddr(k.moduleName, address)
	if account := k.ak.GetAccount(ctx, accAddress); account != nil {
		return
	}
	k.ak.NewAccountWithAddress(ctx, accAddress)
}
