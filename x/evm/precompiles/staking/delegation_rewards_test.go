package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
)

func TestStakingDelegationRewardsABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.DelegationRewardsMethod.Name]
	require.Equal(t, method, staking.DelegationRewardsMethod)
	require.Equal(t, 2, len(staking.DelegationRewardsMethod.Inputs))
	require.Equal(t, 1, len(staking.DelegationRewardsMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestDelegationRewards() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, del common.Address) ([]byte, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.DelegationRewardsMethodName, val.String(), del)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := staking.GetABI().Pack(staking.DelegationRewardsMethodName, newVal, del)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(_ sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := staking.GetABI().Pack(staking.DelegationRewardsMethodName, newVal, del)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("validator not found: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationRewardsName, val.String(), del)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationRewardsName, newVal, del)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("delegationRewards failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(_ sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationRewardsName, newVal, del)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("delegationRewards failed: validator not found: %s", errArgs[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]

			delAmt := sdkmath.NewInt(int64(tmrand.Intn(1000) + 100)).MulRaw(1e18)
			signer := suite.RandSigner()
			helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmt)))

			stakingContract := staking.GetAddress()
			stakingABI := staking.GetABI()
			delegateMethodName := staking.DelegateMethodName
			delegationRewardsMethodName := staking.DelegationRewardsMethodName
			delAddr := signer.Address()

			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delegateMethodName = StakingTestDelegateName
				delegationRewardsMethodName = StakingTestDelegationRewardsName
				delAddr = suite.staking
			}

			pack, err := stakingABI.Pack(delegateMethodName, val0.GetOperator().String())
			suite.Require().NoError(err)
			tx, err := suite.PackEthereumTx(signer, stakingContract, delAmt.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			resp, err := suite.app.DistrKeeper.DelegationRewards(suite.ctx, &distrtypes.QueryDelegationRewardsRequest{DelegatorAddress: sdk.AccAddress(delAddr.Bytes()).String(), ValidatorAddress: val0.GetOperator().String()})
			suite.Require().NoError(err)
			evmDenom := suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom

			pack, errArgs := tc.malleate(val0.GetOperator(), delAddr)
			res, err = suite.app.EvmKeeper.CallEVMWithoutGas(suite.ctx, suite.signer.Address(), &stakingContract, nil, pack, false)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				rewardsValue, err := stakingABI.Methods[delegationRewardsMethodName].Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				suite.Require().EqualValues(rewardsValue[0].(*big.Int), resp.Rewards.AmountOf(evmDenom).TruncateInt().BigInt())
			} else {
				if res.Failed() {
					if res.VmError != vm.ErrExecutionReverted.Error() {
						suite.Require().Equal(tc.error(errArgs), res.VmError)
					} else {
						if len(res.Ret) > 0 {
							reason, err := abi.UnpackRevert(common.CopyBytes(res.Ret))
							suite.Require().NoError(err)

							suite.Require().Equal(tc.error(errArgs), reason)
						} else {
							suite.Require().Equal(tc.error(errArgs), vm.ErrExecutionReverted.Error())
						}
					}
				} else {
					suite.Require().Equal(tc.error(errArgs), err)
				}
			}
		})
	}
}
