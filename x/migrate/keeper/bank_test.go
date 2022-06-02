package keeper_test

import (
	"github.com/functionx/fx-core/app/helpers"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"github.com/functionx/fx-core/crypto/ethsecp256k1"

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
	myApp, _, delegateAddressArr := initTest(t)
	ctx := myApp.BaseApp.NewContext(false, tmproto.Header{})
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	b1 := myApp.BankKeeper.GetAllBalances(ctx, alice)
	require.False(t, b1.Empty())
	b2 := myApp.BankKeeper.GetAllBalances(ctx, bob)
	require.False(t, b1.Empty())

	migrateKeeper := myApp.MigrateKeeper
	m := migratekeeper.NewBankMigrate(myApp.BankKeeper)
	err := m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	bb1 := myApp.BankKeeper.GetAllBalances(ctx, alice)
	require.True(t, bb1.Empty())
	bb2 := myApp.BankKeeper.GetAllBalances(ctx, bob)
	require.Equal(t, b1, bb2.Sub(b2))
}

func initTest(t *testing.T) (*app.App, []*tmtypes.Validator, []sdk.AccAddress) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := helpers.GenerateGenesisValidator(t, 3, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initBalances)))
	myApp := helpers.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := myApp.BaseApp.NewContext(false, tmproto.Header{})
	//378664 825 462891000000000000FX,
	//378664 525 462891000000000000FX

	//update staking unbonding time
	stakingParams := myApp.StakingKeeper.GetParams(ctx)
	stakingParams.UnbondingTime = 5 * time.Minute
	myApp.StakingKeeper.SetParams(ctx, stakingParams)

	delegateAddressArr := helpers.AddTestAddrsIncremental(myApp, ctx, 4, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000)))

	//Note: address not matching pubKey, use for test
	for i, addr := range delegateAddressArr {
		account := myApp.AccountKeeper.GetAccount(ctx, addr)
		if i%2 == 0 {
			err := account.SetPubKey(secp256k1.GenPrivKey().PubKey())
			require.NoError(t, err)
		} else {
			key, _ := ethsecp256k1.GenerateKey()
			err := account.SetPubKey(key.PubKey())
			require.NoError(t, err)
		}
		myApp.AccountKeeper.SetAccount(ctx, account)
	}

	return myApp, validator.Validators, delegateAddressArr
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
