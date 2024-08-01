package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/contract"
	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/staking/precompile"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

func TestStakingUndelegateABI(t *testing.T) {
	undelegateMethod := precompile.NewUndelegateMethod(nil)

	require.Equal(t, 2, len(undelegateMethod.Method.Inputs))
	require.Equal(t, 3, len(undelegateMethod.Method.Outputs))

	require.Equal(t, 5, len(undelegateMethod.Event.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestUndelegate() {
	undelegateMethod := precompile.NewUndelegateMethod(nil)
	undelegateV2Method := precompile.NewUndelegateV2Method(nil)
	testCases := []struct {
		name     string
		isV2     bool
		malleate func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				return types.UndelegateArgs{
					Validator: val.String(),
					Shares:    shares.TruncateInt().BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := val.String() + "1"
				return types.UndelegateArgs{
					Validator: newVal,
					Shares:    shares.TruncateInt().BigInt(),
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				return types.UndelegateArgs{
					Validator: newVal,
					Shares:    shares.TruncateInt().BigInt(),
				}, fmt.Errorf("validator not found: %s", newVal)
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				return types.UndelegateArgs{
					Validator: val.String(),
					Shares:    shares.TruncateInt().BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := val.String() + "1"
				return types.UndelegateArgs{
					Validator: newVal,
					Shares:    shares.TruncateInt().BigInt(),
				}, fmt.Errorf("undelegate failed: invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "contract - failed - validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				return types.UndelegateArgs{
					Validator: newVal,
					Shares:    shares.TruncateInt().BigInt(),
				}, fmt.Errorf("undelegate failed: validator not found: %s", newVal)
			},
			result: false,
		},

		{
			name: "ok v2",
			isV2: true,
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				return types.UndelegateV2Args{
					Validator: val.String(),
					Amount:    delAmt.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator address",
			isV2: true,
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := val.String() + "1"
				return types.UndelegateV2Args{
					Validator: newVal,
					Amount:    delAmt.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},

		{
			name: "contract - ok v2",
			isV2: true,
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				return types.UndelegateV2Args{
					Validator: val.String(),
					Amount:    delAmt.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - v2 invalid validator address",
			isV2: true,
			malleate: func(val sdk.ValAddress, shares sdk.Dec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := val.String() + "1"
				return types.UndelegateV2Args{
					Validator: newVal,
					Amount:    delAmt.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()
			delAmt := helpers.NewRandAmount()

			stakingContract := precompile.GetAddress()
			stakingABI := precompile.GetABI()
			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delAddr = suite.staking
			}

			pack, err := stakingABI.Pack(TestDelegateName, val.GetOperator().String())
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, delAmt.BigInt(), pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			delegation := suite.GetDelegation(delAddr.Bytes(), val.GetOperator())

			undelegations := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, delAddr.Bytes())
			suite.Require().Equal(0, len(undelegations))

			var packData []byte
			args, errResult := tc.malleate(val.GetOperator(), delegation.Shares, delAmt)
			if !tc.isV2 {
				packData, err = undelegateMethod.PackInput(args.(types.UndelegateArgs))
			} else {
				packData, err = undelegateV2Method.PackInput(args.(types.UndelegateV2Args))
			}
			suite.Require().NoError(err)

			if strings.HasPrefix(tc.name, "contract") {
				if !tc.isV2 {
					argsV1 := args.(types.UndelegateArgs)
					packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestUndelegateName, argsV1.Validator, argsV1.Shares)
				} else {
					argsV2 := args.(types.UndelegateV2Args)
					packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestUndelegateV2Name, argsV2.Validator, argsV2.Amount)
				}
				suite.Require().NoError(err)
			}

			res = suite.EthereumTx(suite.signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				if !tc.isV2 {
					unpack, err := stakingABI.Unpack(TestUndelegateName, res.Ret)
					suite.Require().NoError(err)
					// amount,reward,completionTime
					reward := unpack[1].(*big.Int)
					suite.Require().True(reward.Cmp(big.NewInt(0)) == 1, reward.String())
				}

				undelegations := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, delAddr.Bytes())
				suite.Require().Equal(1, len(undelegations))
				suite.Require().Equal(1, len(undelegations[0].Entries))
				suite.Require().Equal(sdk.AccAddress(delAddr.Bytes()).String(), undelegations[0].DelegatorAddress)
				suite.Require().Equal(val.GetOperator().String(), undelegations[0].ValidatorAddress)
				suite.Require().Equal(delAmt, undelegations[0].Entries[0].Balance)

				suite.CheckUndelegateLogs(res.Logs, delAddr, val.GetOperator().String(), delegation.Shares.TruncateInt().BigInt(),
					undelegations[0].Entries[0].Balance.BigInt(), undelegations[0].Entries[0].CompletionTime)

				suite.CheckUndeledateEvents(suite.Ctx, val.GetOperator().String(), undelegations[0].Entries[0].Balance.BigInt(),
					undelegations[0].Entries[0].CompletionTime)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
