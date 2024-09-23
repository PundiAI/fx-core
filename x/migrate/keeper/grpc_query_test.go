package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"github.com/functionx/fx-core/v8/x/migrate/types"
)

func (suite *KeeperTestSuite) TestMigrateRecord() {
	var (
		req    *types.QueryMigrateRecordRequest
		expRes *types.QueryMigrateRecordResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"fail - no address",
			func() {
				req = &types.QueryMigrateRecordRequest{}
				expRes = &types.QueryMigrateRecordResponse{}
			},
			false,
		},
		{
			"success - address not migrate",
			func() {
				key := secp256k1.GenPrivKey()
				req = &types.QueryMigrateRecordRequest{
					Address: sdk.AccAddress(key.PubKey().Address()).String(),
				}
				expRes = &types.QueryMigrateRecordResponse{
					Found:         false,
					MigrateRecord: types.MigrateRecord{},
				}
			},
			true,
		},
		{
			"success - address from migrate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				suite.App.MigrateKeeper.SetMigrateRecord(suite.Ctx, from, to)

				req = &types.QueryMigrateRecordRequest{
					Address: from.String(),
				}
				expRes = &types.QueryMigrateRecordResponse{
					Found: true,
					MigrateRecord: types.MigrateRecord{
						From:   from.String(),
						To:     to.String(),
						Height: suite.Ctx.BlockHeight(),
					},
				}
			},
			true,
		},
		{
			"success - address to migrate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				suite.App.MigrateKeeper.SetMigrateRecord(suite.Ctx, from, to)

				req = &types.QueryMigrateRecordRequest{
					Address: to.String(),
				}
				expRes = &types.QueryMigrateRecordResponse{
					Found: true,
					MigrateRecord: types.MigrateRecord{
						From:   from.String(),
						To:     to.String(),
						Height: suite.Ctx.BlockHeight(),
					},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			tc.malleate()

			res, err := suite.queryClient.MigrateRecord(suite.Ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.Found, res.Found)
				suite.Require().Equal(expRes.MigrateRecord, res.MigrateRecord)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMigrateCheckAccount() {
	var req *types.QueryMigrateCheckAccountRequest
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"fail - no address",
			func() {
				req = &types.QueryMigrateCheckAccountRequest{}
			},
			false,
		},
		{
			"fail - no from address",
			func() {
				toKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())
				req = &types.QueryMigrateCheckAccountRequest{
					To: to.String(),
				}
			},
			false,
		},
		{
			"fail - no to address",
			func() {
				fromKey := secp256k1.GenPrivKey()
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
				}
			},
			false,
		},
		{
			"success - can migrate",
			func() {
				fromKey := secp256k1.GenPrivKey()
				toKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)

				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())
				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"failed - has migrated",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				suite.App.MigrateKeeper.SetMigrateRecord(suite.Ctx, from, to)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has delegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
				suite.Require().NoError(err)
				val1 := validators[0]
				// delegate
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, from, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"fail - to has delegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
				suite.Require().NoError(err)
				val1 := validators[0]
				// delegate
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, to.Bytes(), sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has undelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
				suite.Require().NoError(err)
				val1 := validators[0]
				// delegate
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, from, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				valAddress, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val1.GetOperator())
				suite.Require().NoError(err)
				_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, from, valAddress, sdkmath.LegacyNewDec(1))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"fail - to has undelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
				suite.Require().NoError(err)
				val1 := validators[0]
				// delegate
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, to.Bytes(), sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				valAddress, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val1.GetOperator())
				suite.Require().NoError(err)
				delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), valAddress)
				suite.Require().NoError(err)

				_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), valAddress, delegation.Shares)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has redelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
				suite.Require().NoError(err)
				val1, val2 := validators[0], validators[1]
				// delegate
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, from, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				valAddress1, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val1.GetOperator())
				suite.Require().NoError(err)
				valAddress2, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val2.GetOperator())
				suite.Require().NoError(err)
				_, err = suite.App.StakingKeeper.BeginRedelegation(suite.Ctx, from, valAddress1, valAddress2, sdkmath.LegacyNewDec(1))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"fail - to has redelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators, err := suite.App.StakingKeeper.GetValidators(suite.Ctx, 10)
				suite.Require().NoError(err)
				val1, val2 := validators[0], validators[1]
				// delegate
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, to.Bytes(), sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				valAddress, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val1.GetOperator())
				suite.Require().NoError(err)
				delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), valAddress)
				suite.Require().NoError(err)

				valAddress1, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val1.GetOperator())
				suite.Require().NoError(err)
				valAddress2, err := suite.App.StakingKeeper.Keeper.ValidatorAddressCodec().StringToBytes(val2.GetOperator())
				suite.Require().NoError(err)
				_, err = suite.App.StakingKeeper.BeginRedelegation(suite.Ctx, to.Bytes(), valAddress1, valAddress2, delegation.Shares)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			tc.malleate()

			_, err := suite.queryClient.MigrateCheckAccount(suite.Ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err, err)
			}
		})
	}
}
