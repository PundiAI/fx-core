package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

func (k Keeper) SlashOracle(ctx sdk.Context, oracleAddress string, params types.Params) {
	oracleAddr, err := sdk.AccAddressFromBech32(oracleAddress)
	if err != nil {
		panic(err)
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		panic(fmt.Sprintf("not found oracle:%s", oracleAddr.String()))
	}
	if oracle.Jailed {
		return
	}
	slashAmount := oracle.DepositAmount.Amount.ToDec().Mul(params.SlashFraction).TruncateInt()
	tokensToBurn := sdk.MinInt(slashAmount, oracle.DepositAmount.Amount)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt())
	slashCoin := sdk.NewCoin(oracle.DepositAmount.Denom, tokensToBurn)
	if slashCoin.IsPositive() {
		oracle.DepositAmount = oracle.DepositAmount.Sub(slashCoin)
		if err = k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(slashCoin)); err != nil {
			panic(err)
		}
	}

	oracle.Jailed = true
	oracle.JailedHeight = ctx.BlockHeight()
	k.SetOracle(ctx, oracle)
	k.SetLastOracleSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))

}
