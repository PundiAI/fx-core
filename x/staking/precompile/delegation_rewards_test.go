package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingDelegationRewardsABI(t *testing.T) {
	delegationRewardMethod := precompile.NewDelegationRewardsMethod(nil)

	require.Len(t, delegationRewardMethod.Method.Inputs, 2)
	require.Len(t, delegationRewardMethod.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestDelegationRewards() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error) {
				return types.DelegationRewardsArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error) {
				newVal := val.String() + "1"
				return types.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(_ sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()

				return types.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error) {
				return types.DelegationRewardsArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error) {
				newVal := val.String() + "1"
				return types.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(_ sdk.ValAddress, del common.Address) (types.DelegationRewardsArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				return types.DelegationRewardsArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val0 := suite.GetFirstValidator()

			delAmt := sdkmath.NewInt(int64(tmrand.Intn(1000) + 100)).MulRaw(1e18)
			signer := suite.RandSigner()
			suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmt))

			stakingContract := suite.stakingAddr
			delAddr := signer.Address()

			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
				delAddr = suite.stakingTestAddr
				value = delAmt.BigInt()
			}

			operator0, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val0.GetOperator())
			suite.Require().NoError(err)
			pack, err := suite.delegateV2Method.PackInput(types.DelegateV2Args{
				Validator: val0.GetOperator(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().NoError(err)

			res := suite.EthereumTx(signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			resp, err := suite.DistributionQueryClient(suite.Ctx).DelegationRewards(suite.Ctx, &distrtypes.QueryDelegationRewardsRequest{DelegatorAddress: sdk.AccAddress(delAddr.Bytes()).String(), ValidatorAddress: val0.GetOperator()})
			suite.Require().NoError(err)
			evmDenom := suite.App.EvmKeeper.GetParams(suite.Ctx).EvmDenom

			args, errResult := tc.malleate(operator0, delAddr)
			packData, err := suite.delegationRewardsMethod.PackInput(args)
			suite.Require().NoError(err)

			res, err = suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, suite.signer.Address(), &stakingContract, nil, packData, false)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				rewardsValue, err := suite.delegationRewardsMethod.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				suite.Require().EqualValues(rewardsValue.String(), resp.Rewards.AmountOf(evmDenom).TruncateInt().BigInt().String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
