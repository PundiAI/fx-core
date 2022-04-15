package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/app"
	fxtypes "github.com/functionx/fx-core/types"
	migratekeeper "github.com/functionx/fx-core/x/migrate/keeper"
)

func TestMigrateBankFunc(t *testing.T) {
	fxcore, _, delegateAddressArr := initTest(t)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	b1 := fxcore.BankKeeper.GetAllBalances(ctx, alice)
	require.False(t, b1.Empty())
	b2 := fxcore.BankKeeper.GetAllBalances(ctx, bob)
	require.False(t, b1.Empty())

	migrateKeeper := fxcore.MigrateKeeper
	m := migratekeeper.NewBankMigrate(fxcore.BankKeeper)
	err := m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	bb1 := fxcore.BankKeeper.GetAllBalances(ctx, alice)
	require.True(t, bb1.Empty())
	bb2 := fxcore.BankKeeper.GetAllBalances(ctx, bob)
	require.Equal(t, b1, bb2.Sub(b2))
}

func initTest(t *testing.T) (*app.App, []*tmtypes.Validator, []sdk.AccAddress) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(3,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})

	//update staking unbonding time
	stakingParams := fxcore.StakingKeeper.GetParams(ctx)
	stakingParams.UnbondingTime = 5 * time.Minute
	fxcore.StakingKeeper.SetParams(ctx, stakingParams)

	delegateAddressArr := app.AddTestAddrsIncremental(fxcore, ctx, 4, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000)))
	return fxcore, validator.Validators, delegateAddressArr
}

func GetValidator(t *testing.T, app *app.App, vals ...*tmtypes.Validator) []stakingtypes.Validator {
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	validators := make([]stakingtypes.Validator, 0, len(vals))
	for _, val := range vals {
		validator, found := app.StakingKeeper.GetValidator(ctx, val.Address.Bytes())
		require.True(t, found)
		validators = append(validators, validator)
	}
	return validators
}
