package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/functionx/fx-core/app/fxcore"
	migratekeeper "github.com/functionx/fx-core/x/migrate/keeper"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"testing"
)

func TestMigrateBankFunc(t *testing.T) {
	app, _, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	b1 := app.BankKeeper.GetAllBalances(ctx, alice)
	require.False(t, b1.Empty())
	b2 := app.BankKeeper.GetAllBalances(ctx, bob)
	require.False(t, b1.Empty())

	migrateKeeper := app.MigrateKeeper
	m := migratekeeper.NewBankMigrate(app.BankKeeper)
	err := m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	bb1 := app.BankKeeper.GetAllBalances(ctx, alice)
	require.True(t, bb1.Empty())
	bb2 := app.BankKeeper.GetAllBalances(ctx, bob)
	require.Equal(t, b1, bb2.Sub(b2))
}

func initTest(t *testing.T) (*fxcore.App, []*tmtypes.Validator, []sdk.AccAddress) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(3,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	delegateAddressArr := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(10000)))
	return app, validator.Validators, delegateAddressArr
}

func GetValidator(t *testing.T, app *fxcore.App, vals ...*tmtypes.Validator) []stakingtypes.Validator {
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	validators := make([]stakingtypes.Validator, 0, len(vals))
	for _, val := range vals {
		validator, found := app.StakingKeeper.GetValidator(ctx, val.Address.Bytes())
		require.True(t, found)
		validators = append(validators, validator)
	}
	return validators
}
