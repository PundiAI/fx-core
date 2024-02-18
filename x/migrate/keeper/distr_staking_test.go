package keeper_test

import (
	"context"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v7/app"
	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	migratekeeper "github.com/functionx/fx-core/v7/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateStakingDelegate() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val1 := validators[0]

	// acc delegate
	_, err := suite.app.StakingKeeper.Delegate(suite.ctx, acc, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	// check acc delegate
	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().True(found)
	shares := delegation.Shares

	// check eth acc delegate
	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().False(found)

	// commit block
	suite.ctx = commitBlock(suite.T(), suite.ctx, suite.app)
	suite.ctx = commitBlock(suite.T(), suite.ctx, suite.app)

	// check reward
	rewards1, err := GetDelegateRewards(suite.ctx, suite.app, acc, val1.GetOperator())
	suite.Require().NoError(err)
	rewards2, err := GetDelegateRewards(suite.ctx, suite.app, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().Equal("delegation does not exist", err.Error())

	// migrate
	m := migratekeeper.NewDistrStakingMigrate(suite.app.GetKey(distritypes.StoreKey), suite.app.GetKey(stakingtypes.StoreKey), suite.app.StakingKeeper)
	err = m.Validate(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	// check eth acc delegate
	delegation, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(shares, delegation.Shares)

	// check acc delegate
	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().False(found)

	// check reward
	rewards3, err := GetDelegateRewards(suite.ctx, suite.app, acc, val1.GetOperator())
	suite.Require().Equal("delegation does not exist", err.Error())
	suite.Require().Equal(rewards2, rewards3)
	rewards4, err := GetDelegateRewards(suite.ctx, suite.app, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().NoError(err)
	suite.Require().Equal(rewards1, rewards4)
}

func (suite *KeeperTestSuite) TestMigrateStakingUnbonding() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val1 := validators[0]

	// delegate
	delegateAmount := sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000))
	_, err := suite.app.StakingKeeper.Delegate(suite.ctx, acc, delegateAmount, stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	del, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().True(found)

	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().False(found)

	// undelegate
	completionTime, err := suite.app.StakingKeeper.Undelegate(suite.ctx, acc, val1.GetOperator(), del.Shares.Quo(sdk.NewDec(10)))
	suite.Require().NoError(err)

	delegation2, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().True(found)

	unbondingDelegations := suite.app.StakingKeeper.GetAllUnbondingDelegations(suite.ctx, acc)
	suite.Require().Equal(1, len(unbondingDelegations))
	suite.Require().Equal(unbondingDelegations[0].Entries[0].CompletionTime, completionTime)
	suite.Require().Equal(unbondingDelegations[0].DelegatorAddress, acc.String())

	slice := suite.app.StakingKeeper.GetUBDQueueTimeSlice(suite.ctx, completionTime)
	suite.Require().Equal(1, len(slice))
	suite.Require().Equal(acc.String(), slice[0].DelegatorAddress)

	m := migratekeeper.NewDistrStakingMigrate(suite.app.GetKey(distritypes.StoreKey), suite.app.GetKey(stakingtypes.StoreKey), suite.app.StakingKeeper)
	err = m.Validate(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().False(found)

	delegation3, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(delegation2.Shares, delegation3.Shares)

	unbondingDelegations = suite.app.StakingKeeper.GetAllUnbondingDelegations(suite.ctx, acc)
	suite.Require().Equal(0, len(unbondingDelegations))

	unbondingDelegations = suite.app.StakingKeeper.GetAllUnbondingDelegations(suite.ctx, ethAcc.Bytes())
	suite.Require().Equal(1, len(unbondingDelegations))
	suite.Require().Equal(unbondingDelegations[0].Entries[0].CompletionTime, completionTime)
	suite.Require().Equal(unbondingDelegations[0].DelegatorAddress, sdk.AccAddress(ethAcc.Bytes()).String())

	slice = suite.app.StakingKeeper.GetUBDQueueTimeSlice(suite.ctx, completionTime)
	suite.Require().Equal(1, len(slice))
	suite.Require().Equal(sdk.AccAddress(ethAcc.Bytes()).String(), slice[0].DelegatorAddress)

	ethAccBalanceV1 := suite.app.BankKeeper.GetBalance(suite.ctx, ethAcc.Bytes(), fxtypes.DefaultDenom)
	suite.Require().True(ethAccBalanceV1.Amount.GT(sdkmath.NewInt(0)))

	suite.ctx = commitUnbonding(suite.T(), suite.ctx, suite.app)

	suite.ctx = commitBlock(suite.T(), suite.ctx, suite.app)

	ethAccBalanceV2 := suite.app.BankKeeper.GetBalance(suite.ctx, ethAcc.Bytes(), fxtypes.DefaultDenom)
	suite.Require().Equal(ethAccBalanceV2.Sub(ethAccBalanceV1).Amount, delegateAmount.Quo(sdkmath.NewInt(10)))
}

func (suite *KeeperTestSuite) TestMigrateStakingRedelegate() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val1, val2 := validators[0], validators[1]

	// delegate
	_, err := suite.app.StakingKeeper.Delegate(suite.ctx, acc, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
	suite.Require().NoError(err)

	_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().True(found)

	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().False(found)

	// redelegate
	var completionTime time.Time
	entries := suite.app.StakingKeeper.MaxEntries(suite.ctx)
	for i := 0; i < int(entries); i++ {
		completionTime, err = suite.app.StakingKeeper.BeginRedelegation(suite.ctx, acc, val1.GetOperator(), val2.GetOperator(), sdk.NewDec(1))
		suite.Require().NoError(err)
	}

	delegation2, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().True(found)

	delegation3, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val2.GetOperator())
	suite.Require().True(found)

	queue := suite.app.StakingKeeper.GetRedelegationQueueTimeSlice(suite.ctx, completionTime)
	suite.Require().Equal(int(entries), len(queue))
	suite.Require().Equal(queue[0].DelegatorAddress, acc.String())

	m := migratekeeper.NewDistrStakingMigrate(suite.app.GetKey(distritypes.StoreKey), suite.app.GetKey(stakingtypes.StoreKey), suite.app.StakingKeeper)
	err = m.Validate(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, acc, val1.GetOperator())
	suite.Require().False(found)

	delegation4, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val1.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(delegation2.Shares, delegation4.Shares)

	delegation5, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, ethAcc.Bytes(), val2.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(delegation3.Shares, delegation5.Shares)

	queue = suite.app.StakingKeeper.GetRedelegationQueueTimeSlice(suite.ctx, completionTime)
	suite.Require().Equal(int(entries), len(queue))
	suite.Require().Equal(queue[0].DelegatorAddress, sdk.AccAddress(ethAcc.Bytes()).String())
}

func GetDelegateRewards(ctx sdk.Context, app *app.App, delegate []byte, validator sdk.ValAddress) (sdk.DecCoins, error) {
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	distritypes.RegisterQueryServer(queryHelper, app.DistrKeeper)
	queryClient := distritypes.NewQueryClient(queryHelper)
	rewards, err := queryClient.DelegationRewards(context.Background(), &distritypes.QueryDelegationRewardsRequest{
		DelegatorAddress: sdk.AccAddress(delegate).String(),
		ValidatorAddress: validator.String(),
	})
	if err != nil {
		return nil, err
	}
	return rewards.Rewards, nil
}

func commitUnbonding(t *testing.T, ctx sdk.Context, app *app.App) sdk.Context {
	i := 0
	for i < 70 {
		ctx = commitBlock(t, ctx, app)
		i++
	}
	return ctx
}

func commitBlock(t *testing.T, ctx sdk.Context, app *app.App) sdk.Context {
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(5 * time.Second))

	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	mint.BeginBlocker(ctx, app.MintKeeper, minttypes.DefaultInflationCalculationFn)

	distribution.BeginBlocker(ctx, abcitypes.RequestBeginBlock{
		Hash:   nil,
		Header: tmproto.Header{},
		LastCommitInfo: abcitypes.LastCommitInfo{
			Round: 0,
			Votes: buildCommitVotes(t, ctx, app.StakingKeeper.Keeper, app.AppCodec()),
		},
		ByzantineValidators: nil,
	}, app.DistrKeeper)

	return ctx
}

func buildCommitVotes(t *testing.T, ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, codec codec.Codec) []abcitypes.VoteInfo {
	t.Helper()
	validators := stakingKeeper.GetAllValidators(ctx)

	result := make([]abcitypes.VoteInfo, 0, len(validators))
	for _, validator := range validators {
		if !validator.IsBonded() {
			continue
		}

		var pubkey cryptotypes.PubKey
		err := codec.UnpackAny(validator.ConsensusPubkey, &pubkey)
		require.NoError(t, err)
		result = append(result, abcitypes.VoteInfo{
			Validator: abcitypes.Validator{
				Address: pubkey.Address(),
				Power:   validator.GetConsensusPower(sdk.DefaultPowerReduction),
			},
			SignedLastBlock: true,
		})
	}
	return result
}
