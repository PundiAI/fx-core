package v8

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

var (
	baseCoin          = "usdt"
	channelId         = "channel-0"
	baseToBridgeCoins = sdk.NewCoins(
		sdk.NewCoin("polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F", sdkmath.NewInt(9309873)),
		sdk.NewCoin("tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", sdkmath.NewInt(105700000)),
	)
)

func unwrapEscrowBalance(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	if ctx.ChainID() != fxtypes.MainnetChainId {
		return nil
	}
	escrowAddress := getEscrowAddr()
	for _, coin := range baseToBridgeCoins {
		baseCoinBalance := bankKeeper.GetBalance(ctx, escrowAddress, baseCoin)
		if baseCoinBalance.IsZero() || baseCoinBalance.Amount.LT(coin.Amount) {
			ctx.Logger().With("wrap escrow balance").Error("base coin balance is not enough")
			continue
		}

		baseCoins := sdk.NewCoins(sdk.NewCoin(baseCoin, coin.Amount))
		// 1. send base coin from escrow to crosschain module
		if err := bankKeeper.SendCoinsFromAccountToModule(ctx, escrowAddress, crosschaintypes.ModuleName, baseCoins); err != nil {
			return err
		}

		// 2. burn base coin
		if err := bankKeeper.BurnCoins(ctx, crosschaintypes.ModuleName, baseCoins); err != nil {
			return err
		}

		// 3. send bridge coin from crosschain module to escrow
		if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, crosschaintypes.ModuleName, escrowAddress, sdk.NewCoins(coin)); err != nil {
			return err
		}
	}

	return nil
}

func UnwrapEscrowBalanceCheckBefore(t *testing.T, ctx sdk.Context, keeper bankkeeper.Keeper) {
	t.Helper()
	balances := keeper.GetAllBalances(ctx, getEscrowAddr())
	for _, coin := range baseToBridgeCoins {
		require.True(t, balances.AmountOf(coin.Denom).IsZero(), "balance should be zero", coin.String())
	}
}

func UnwrapEscrowBalanceCheckAfter(t *testing.T, ctx sdk.Context, keeper bankkeeper.Keeper) {
	t.Helper()
	balances := keeper.GetAllBalances(ctx, getEscrowAddr())
	for _, coin := range baseToBridgeCoins {
		find, balanceCoin := balances.Find(coin.Denom)
		require.True(t, find, "balance should be found", balances.String(), coin.String())
		require.Equal(t, coin.String(), balanceCoin.String(), "balance should be equal", coin.String(), balanceCoin.String())
	}
}

func getEscrowAddr() sdk.AccAddress {
	return transfertypes.GetEscrowAddress(transfertypes.PortID, channelId)
}
