package v2_test

import (
	"fmt"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	v2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
)

func TestMigrateBank(t *testing.T) {
	app := helpers.Setup(false, false)
	ctx := app.NewContext(false, tmproto.Header{Height: int64(rand.Uint32())})

	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(int64(rand.Uint32()))))
	index := rand.Intn(100) + 1
	for i := 0; i < index; i++ {
		coins = coins.Add(sdk.Coin{
			Denom:  fmt.Sprintf("%s%s", ethtypes.ModuleName, helpers.GenerateAddress().Hex()),
			Amount: sdk.NewInt(int64(rand.Uint32())),
		})
	}
	coins = coins.Sort()

	balances0 := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(gravitytypes.ModuleName))
	assert.Equal(t, balances0.Len(), 0)

	balances1 := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(ethtypes.ModuleName))
	assert.Equal(t, balances1.String(), "378600525462891000000000000FX")

	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	assert.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, gravitytypes.ModuleName, coins)
	assert.NoError(t, err)

	balances2 := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(gravitytypes.ModuleName))
	assert.Equal(t, balances2.String(), coins.String())

	err = v2.MigrateBank(ctx, app.AccountKeeper, app.BankKeeper, ethtypes.ModuleName)
	assert.NoError(t, err)

	balances3 := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(ethtypes.ModuleName))
	coins = coins.Add(balances1...).Sort()
	assert.Equal(t, balances3.String(), coins.String())
}
