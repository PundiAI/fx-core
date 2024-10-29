package precompile_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func TestStakingDelegationRewardsABI(t *testing.T) {
	delegationRewardMethod := precompile.NewDelegationRewardsMethod(nil)

	require.Len(t, delegationRewardMethod.Method.Inputs, 2)
	require.Len(t, delegationRewardMethod.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestDelegationRewards() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error) {
				return contract.DelegationRewardsArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error) {
				newVal := val.String() + "1"
				return contract.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(_ sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()

				return contract.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error) {
				return contract.DelegationRewardsArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error) {
				newVal := val.String() + "1"
				return contract.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(_ sdk.ValAddress, del common.Address) (contract.DelegationRewardsArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				return contract.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator0 := suite.GetFirstValAddr()
			delAmt := helpers.NewRandAmount()

			suite.WithContract(suite.stakingAddr)
			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				suite.WithContract(suite.stakingTestAddr)
				delAddr = suite.stakingTestAddr
				suite.MintToken(delAddr.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmt))
			}

			res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
				Validator: operator0.String(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			resp, err := suite.DistributionQueryClient(suite.Ctx).DelegationRewards(suite.Ctx,
				&distrtypes.QueryDelegationRewardsRequest{
					DelegatorAddress: sdk.AccAddress(delAddr.Bytes()).String(),
					ValidatorAddress: operator0.String(),
				})
			suite.Require().NoError(err)

			args, expectErr := tc.malleate(operator0, delAddr)

			rewards := suite.WithError(expectErr).DelegationRewards(suite.Ctx, args)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().EqualValues(resp.Rewards[0].Amount.TruncateInt().String(), rewards.String())
			}
		})
	}
}
