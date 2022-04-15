package keeper_test

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	fxtypes "github.com/functionx/fx-core/types"
	migratekeeper "github.com/functionx/fx-core/x/migrate/keeper"
)

func TestMigrateStakingHandler_Delegate(t *testing.T) {
	app, validators, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	val1, _, _ := validators[0], validators[1], validators[2]
	validator1 := GetValidator(t, app, val1)[0]
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	_, err := app.StakingKeeper.Delegate(ctx, alice, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, validator1, true)
	require.NoError(t, err)

	delegation, found := app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.True(t, found)
	shares := delegation.Shares

	_, found = app.StakingKeeper.GetDelegation(ctx, bob, val1.Address.Bytes())
	require.False(t, found)

	ctx = commitBlock(t, ctx, app)
	ctx = commitBlock(t, ctx, app)

	rewards1, err := GetDelegateRewards(ctx, app, alice, sdk.ValAddress(val1.Address))
	require.NoError(t, err)
	rewards2, err := GetDelegateRewards(ctx, app, bob, sdk.ValAddress(val1.Address))
	require.Equal(t, "delegation does not exist", err.Error())

	migrateKeeper := app.MigrateKeeper
	m := migratekeeper.NewDistrStakingMigrate(app.GetKey(distritypes.StoreKey), app.GetKey(stakingtypes.StoreKey), app.StakingKeeper)
	err = m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	delegation, found = app.StakingKeeper.GetDelegation(ctx, bob, val1.Address.Bytes())
	require.True(t, found)
	require.Equal(t, shares, delegation.Shares)

	_, found = app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.False(t, found)

	rewards3, err := GetDelegateRewards(ctx, app, alice, sdk.ValAddress(val1.Address))
	require.Equal(t, "delegation does not exist", err.Error())
	require.Equal(t, rewards2, rewards3)
	rewards4, err := GetDelegateRewards(ctx, app, bob, sdk.ValAddress(val1.Address))
	require.NoError(t, err)
	require.Equal(t, rewards1, rewards4)
}

func TestMigrateStakingHandler_Unbonding(t *testing.T) {
	app, validators, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	val1, _, _ := validators[0], validators[1], validators[2]
	validator1 := GetValidator(t, app, val1)[0]
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	//delegate
	delegateAmount := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000))
	_, err := app.StakingKeeper.Delegate(ctx, alice, delegateAmount, stakingtypes.Unbonded, validator1, true)
	require.NoError(t, err)

	_, found := app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.True(t, found)

	_, found = app.StakingKeeper.GetDelegation(ctx, bob, val1.Address.Bytes())
	require.False(t, found)

	//undelegate
	completionTime, err := app.StakingKeeper.Undelegate(ctx, alice, val1.Address.Bytes(), sdk.NewDec(1))
	require.NoError(t, err)

	delegation2, found := app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.True(t, found)

	unbondingDelegations := app.StakingKeeper.GetAllUnbondingDelegations(ctx, alice)
	require.Equal(t, 1, len(unbondingDelegations))
	require.Equal(t, unbondingDelegations[0].Entries[0].CompletionTime, completionTime)
	require.Equal(t, unbondingDelegations[0].DelegatorAddress, alice.String())

	slice := app.StakingKeeper.GetUBDQueueTimeSlice(ctx, completionTime)
	require.Equal(t, 1, len(slice))
	require.Equal(t, alice.String(), slice[0].DelegatorAddress)

	migrateKeeper := app.MigrateKeeper
	m := migratekeeper.NewDistrStakingMigrate(app.GetKey(distritypes.StoreKey), app.GetKey(stakingtypes.StoreKey), app.StakingKeeper)
	err = m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	_, found = app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.False(t, found)

	delegation3, found := app.StakingKeeper.GetDelegation(ctx, bob, val1.Address.Bytes())
	require.True(t, found)
	require.Equal(t, delegation2.Shares, delegation3.Shares)

	unbondingDelegations = app.StakingKeeper.GetAllUnbondingDelegations(ctx, alice)
	require.Equal(t, 0, len(unbondingDelegations))

	unbondingDelegations = app.StakingKeeper.GetAllUnbondingDelegations(ctx, bob)
	require.Equal(t, 1, len(unbondingDelegations))
	require.Equal(t, unbondingDelegations[0].Entries[0].CompletionTime, completionTime)
	require.Equal(t, unbondingDelegations[0].DelegatorAddress, bob.String())

	slice = app.StakingKeeper.GetUBDQueueTimeSlice(ctx, completionTime)
	require.Equal(t, 1, len(slice))
	require.Equal(t, bob.String(), slice[0].DelegatorAddress)

	bobBalanceV1 := app.BankKeeper.GetBalance(ctx, bob, fxtypes.MintDenom)
	require.True(t, bobBalanceV1.Amount.GT(sdk.NewInt(0)))

	ctx = commitUnbonding(t, ctx, app)

	ctx = commitBlock(t, ctx, app)

	bobBalanceV2 := app.BankKeeper.GetBalance(ctx, bob, fxtypes.MintDenom)
	require.Equal(t, bobBalanceV2.Sub(bobBalanceV1).Amount, delegateAmount.Quo(sdk.NewInt(10)))

}

func TestMigrateStakingHandler_Redelegate(t *testing.T) {
	app, validators, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	val1, val2, _ := validators[0], validators[1], validators[2]
	validator12 := GetValidator(t, app, val1, val2)
	validator1, _ := validator12[0], validator12[1]
	alice, bob, _, _ := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	//delegate
	_, err := app.StakingKeeper.Delegate(ctx, alice, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, validator1, true)
	require.NoError(t, err)

	_, found := app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.True(t, found)

	_, found = app.StakingKeeper.GetDelegation(ctx, bob, val1.Address.Bytes())
	require.False(t, found)

	//redelegate
	completionTime, err := app.StakingKeeper.BeginRedelegation(ctx, alice, val1.Address.Bytes(), val2.Address.Bytes(), sdk.NewDec(1))
	require.NoError(t, err)

	delegation2, found := app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.True(t, found)

	delegation3, found := app.StakingKeeper.GetDelegation(ctx, alice, val2.Address.Bytes())
	require.True(t, found)

	queue := app.StakingKeeper.GetRedelegationQueueTimeSlice(ctx, completionTime)
	require.Equal(t, 1, len(queue))
	require.Equal(t, queue[0].DelegatorAddress, alice.String())

	migrateKeeper := app.MigrateKeeper
	m := migratekeeper.NewDistrStakingMigrate(app.GetKey(distritypes.StoreKey), app.GetKey(stakingtypes.StoreKey), app.StakingKeeper)
	err = m.Validate(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, alice, bob)
	require.NoError(t, err)

	_, found = app.StakingKeeper.GetDelegation(ctx, alice, val1.Address.Bytes())
	require.False(t, found)

	delegation4, found := app.StakingKeeper.GetDelegation(ctx, bob, val1.Address.Bytes())
	require.True(t, found)
	require.Equal(t, delegation2.Shares, delegation4.Shares)

	delegation5, found := app.StakingKeeper.GetDelegation(ctx, bob, val2.Address.Bytes())
	require.True(t, found)
	require.Equal(t, delegation3.Shares, delegation5.Shares)

	queue = app.StakingKeeper.GetRedelegationQueueTimeSlice(ctx, completionTime)
	require.Equal(t, 1, len(queue))
	require.Equal(t, queue[0].DelegatorAddress, bob.String())
}

func GetDelegateRewards(ctx sdk.Context, fxcore *app.App, delegate sdk.AccAddress, validator sdk.ValAddress) (sdk.DecCoins, error) {
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, fxcore.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, fxcore.DistrKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	rewards, err := queryClient.DelegationRewards(context.Background(), &types.QueryDelegationRewardsRequest{
		DelegatorAddress: delegate.String(),
		ValidatorAddress: validator.String(),
	})
	if err != nil {
		return nil, err
	}
	return rewards.Rewards, nil
}

func commitUnbonding(t *testing.T, ctx sdk.Context, fxcore *app.App) sdk.Context {
	i := 0
	for i < 70 {
		ctx = commitBlock(t, ctx, fxcore)
		i++
	}
	return ctx
}
