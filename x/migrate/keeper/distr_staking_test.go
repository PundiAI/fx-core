package keeper_test

import (
	"context"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/app"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	migratekeeper "github.com/pundiai/fx-core/v8/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateStakingDelegate() {
	suite.MintToken(suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Len(keys, 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Len(ethKeys, 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
	suite.Require().NoError(err)
	val1 := validators[0]

	// acc delegate
	_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, acc, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	// check acc delegate
	val1ValAddr := suite.ValStringToVal(val1.GetOperator())
	delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().NoError(err)
	shares := delegation.Shares

	// check eth acc delegate
	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val1ValAddr)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	// commit block
	suite.Ctx = commitBlock(suite.T(), suite.Ctx, suite.App)
	suite.Ctx = commitBlock(suite.T(), suite.Ctx, suite.App)

	// check reward
	rewards1, err := GetDelegateRewards(suite.Ctx, suite.App, acc, val1ValAddr.String())
	suite.Require().NoError(err)
	rewards2, err := GetDelegateRewards(suite.Ctx, suite.App, ethAcc.Bytes(), val1ValAddr.String())
	suite.Require().EqualError(err, "no delegation for (address, validator) tuple")

	// migrate
	m := migratekeeper.NewDistrStakingMigrate(suite.App.GetKey(distrtypes.StoreKey), suite.App.GetKey(stakingtypes.StoreKey), suite.App.StakingKeeper)
	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	// check eth acc delegate
	delegation, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val1ValAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(shares, delegation.Shares)

	// check acc delegate
	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	// check reward
	rewards3, err := GetDelegateRewards(suite.Ctx, suite.App, acc, val1ValAddr.String())
	suite.Require().EqualError(err, "no delegation for (address, validator) tuple")
	suite.Require().Equal(rewards2, rewards3)
	rewards4, err := GetDelegateRewards(suite.Ctx, suite.App, ethAcc.Bytes(), val1ValAddr.String())
	suite.Require().NoError(err)
	suite.Require().Equal(rewards1, rewards4)
}

func (suite *KeeperTestSuite) TestMigrateStakingUnbonding() {
	suite.MintToken(suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Len(keys, 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Len(ethKeys, 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
	suite.Require().NoError(err)
	val1 := validators[0]

	// delegate
	delegateAmount := sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000))
	_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, acc, delegateAmount, stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	val1ValAddr := suite.ValStringToVal(val1.GetOperator())
	del, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().NoError(err)

	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val1ValAddr)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	// undelegate
	completionTime, _, err := suite.App.StakingKeeper.Undelegate(suite.Ctx, acc, val1ValAddr, del.Shares.Quo(sdkmath.LegacyNewDec(10)))
	suite.Require().NoError(err)

	delegation2, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().NoError(err)

	unbondingDelegations, err := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, acc)
	suite.Require().NoError(err)
	suite.Require().Len(unbondingDelegations, 1)
	suite.Require().Equal(unbondingDelegations[0].Entries[0].CompletionTime, completionTime)
	suite.Require().Equal(unbondingDelegations[0].DelegatorAddress, acc.String())

	slice, err := suite.App.StakingKeeper.GetUBDQueueTimeSlice(suite.Ctx, completionTime)
	suite.Require().NoError(err)
	suite.Require().Len(slice, 1)
	suite.Require().Equal(acc.String(), slice[0].DelegatorAddress)

	m := migratekeeper.NewDistrStakingMigrate(suite.App.GetKey(distrtypes.StoreKey), suite.App.GetKey(stakingtypes.StoreKey), suite.App.StakingKeeper)
	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	delegation3, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val1ValAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(delegation2.Shares, delegation3.Shares)

	unbondingDelegations, err = suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, acc)
	suite.Require().NoError(err)
	suite.Require().Empty(unbondingDelegations)

	unbondingDelegations, err = suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, ethAcc.Bytes())
	suite.Require().NoError(err)
	suite.Require().Len(unbondingDelegations, 1)
	suite.Require().Equal(unbondingDelegations[0].Entries[0].CompletionTime, completionTime)
	suite.Require().Equal(unbondingDelegations[0].DelegatorAddress, sdk.AccAddress(ethAcc.Bytes()).String())

	slice, err = suite.App.StakingKeeper.GetUBDQueueTimeSlice(suite.Ctx, completionTime)
	suite.Require().NoError(err)
	suite.Require().Len(slice, 1)
	suite.Require().Equal(sdk.AccAddress(ethAcc.Bytes()).String(), slice[0].DelegatorAddress)

	ethAccBalanceV1 := suite.App.BankKeeper.GetBalance(suite.Ctx, ethAcc.Bytes(), fxtypes.DefaultDenom)
	suite.Require().True(ethAccBalanceV1.Amount.GT(sdkmath.NewInt(0)))

	suite.Ctx = commitUnbonding(suite.T(), suite.Ctx, suite.App)

	suite.Ctx = commitBlock(suite.T(), suite.Ctx, suite.App)

	ethAccBalanceV2 := suite.App.BankKeeper.GetBalance(suite.Ctx, ethAcc.Bytes(), fxtypes.DefaultDenom)
	suite.Require().Equal(ethAccBalanceV2.Sub(ethAccBalanceV1).Amount, delegateAmount.Quo(sdkmath.NewInt(10)))
}

func (suite *KeeperTestSuite) TestMigrateStakingRedelegate() {
	suite.MintToken(suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Len(keys, 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Len(ethKeys, 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
	suite.Require().NoError(err)
	val1, val2 := validators[0], validators[1]

	// delegate
	_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, acc, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	val1ValAddr := suite.ValStringToVal(val1.GetOperator())
	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().NoError(err)

	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val1ValAddr)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	// redelegate
	var completionTime time.Time
	entries, err := suite.App.StakingKeeper.MaxEntries(suite.Ctx)
	suite.Require().NoError(err)
	val2ValAddr := suite.ValStringToVal(val2.GetOperator())
	for i := 0; i < int(entries); i++ {
		completionTime, err = suite.App.StakingKeeper.BeginRedelegation(suite.Ctx, acc, val1ValAddr, val2ValAddr, sdkmath.LegacyNewDec(1))
		suite.Require().NoError(err)
	}

	delegation2, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().NoError(err)

	delegation3, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val2ValAddr)
	suite.Require().NoError(err)

	queue, err := suite.App.StakingKeeper.GetRedelegationQueueTimeSlice(suite.Ctx, completionTime)
	suite.Require().NoError(err)
	suite.Require().Len(queue, int(entries))
	suite.Require().Equal(queue[0].DelegatorAddress, acc.String())

	m := migratekeeper.NewDistrStakingMigrate(suite.App.GetKey(distrtypes.StoreKey), suite.App.GetKey(stakingtypes.StoreKey), suite.App.StakingKeeper)
	err = m.Validate(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, acc, val1ValAddr)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	delegation4, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val1ValAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(delegation2.Shares, delegation4.Shares)

	delegation5, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, ethAcc.Bytes(), val2ValAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(delegation3.Shares, delegation5.Shares)

	queue, err = suite.App.StakingKeeper.GetRedelegationQueueTimeSlice(suite.Ctx, completionTime)
	suite.Require().NoError(err)
	suite.Require().Len(queue, int(entries))
	suite.Require().Equal(queue[0].DelegatorAddress, sdk.AccAddress(ethAcc.Bytes()).String())
}

func GetDelegateRewards(ctx sdk.Context, app *app.App, delegate []byte, validator string) (sdk.DecCoins, error) {
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	distrtypes.RegisterQueryServer(queryHelper, distributionkeeper.NewQuerier(app.DistrKeeper))
	queryClient := distrtypes.NewQueryClient(queryHelper)
	rewards, err := queryClient.DelegationRewards(context.Background(), &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: sdk.AccAddress(delegate).String(),
		ValidatorAddress: validator,
	})
	if err != nil {
		return nil, err
	}
	return rewards.Rewards, nil
}

func commitUnbonding(t *testing.T, ctx sdk.Context, app *app.App) sdk.Context {
	t.Helper()
	i := 0
	for i < 70 {
		ctx = commitBlock(t, ctx, app)
		i++
	}
	return ctx
}

func commitBlock(t *testing.T, ctx sdk.Context, app *app.App) sdk.Context {
	t.Helper()
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(5 * time.Second))

	_, err := app.StakingKeeper.Keeper.EndBlocker(ctx)
	require.NoError(t, err)
	err = mint.BeginBlocker(ctx, app.MintKeeper, minttypes.DefaultInflationCalculationFn)
	require.NoError(t, err)

	err = distribution.BeginBlocker(ctx, app.DistrKeeper)
	require.NoError(t, err)

	return ctx
}
