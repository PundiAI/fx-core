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

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
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
	delegationRewardMethod := precompile.NewDelegationRewardsMethod(nil)
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

			stakingContract := precompile.GetAddress()
			stakingABI := precompile.GetABI()
			delAddr := signer.Address()

			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delAddr = suite.staking
				value = delAmt.BigInt()
			}

			operator0, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val0.GetOperator())
			suite.Require().NoError(err)
			pack, err := stakingABI.Pack(TestDelegateV2Name, val0.GetOperator(), delAmt.BigInt())
			suite.Require().NoError(err)

			res := suite.EthereumTx(signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			resp, err := suite.DistributionQueryClient(suite.Ctx).DelegationRewards(suite.Ctx, &distrtypes.QueryDelegationRewardsRequest{DelegatorAddress: sdk.AccAddress(delAddr.Bytes()).String(), ValidatorAddress: val0.GetOperator()})
			suite.Require().NoError(err)
			evmDenom := suite.App.EvmKeeper.GetParams(suite.Ctx).EvmDenom

			args, errResult := tc.malleate(operator0, delAddr)
			packData, err := delegationRewardMethod.PackInput(args)
			suite.Require().NoError(err)
			if strings.HasPrefix(tc.name, "contract") {
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegationRewardsName, args.Validator, args.Delegator)
				suite.Require().NoError(err)
			}

			res, err = suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, suite.signer.Address(), &stakingContract, nil, packData, false)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				rewardsValue, err := stakingABI.Methods[TestDelegationRewardsName].Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				suite.Require().EqualValues(rewardsValue[0].(*big.Int).String(), resp.Rewards.AmountOf(evmDenom).TruncateInt().BigInt().String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
