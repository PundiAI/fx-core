package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingDelegationABI(t *testing.T) {
	delegationMethod := precompile.NewDelegationMethod(nil)

	require.Len(t, delegationMethod.Method.Inputs, 2)
	require.Len(t, delegationMethod.Method.Outputs, 2)
}

func (suite *PrecompileTestSuite) TestDelegation() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				return types.DelegationArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "ok - zero",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				return types.DelegationArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				newVal := val.String() + "1"
				return types.DelegationArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()

				return types.DelegationArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				return types.DelegationArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - zero",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				return types.DelegationArgs{
					Validator: val.String(),
					Delegator: del,
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				newVal := val.String() + "1"
				return types.DelegationArgs{
					Validator: newVal,
					Delegator: del,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(val sdk.ValAddress, del common.Address) (types.DelegationArgs, error) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()

				return types.DelegationArgs{
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

			delAmount := helpers.NewRandAmount()

			stakingContract := suite.stakingAddr
			delAddr := suite.signer.Address()
			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
				delAddr = suite.stakingTestAddr
				value = delAmount.BigInt()
			}

			operator0, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val0.GetOperator())
			suite.Require().NoError(err)

			pack, err := suite.delegateV2Method.PackInput(types.DelegateV2Args{
				Validator: val0.GetOperator(),
				Amount:    delAmount.BigInt(),
			})
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			args, errResult := tc.malleate(operator0, delAddr)
			packData, err := suite.delegationMethod.PackInput(args)
			suite.Require().NoError(err)
			delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator0)
			suite.Require().NoError(err)

			res, _ = suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, suite.signer.Address(), &stakingContract, nil, packData, false)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				delValue, _, err := suite.delegationMethod.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal(delegation.GetShares().TruncateInt().String(), delValue.String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
