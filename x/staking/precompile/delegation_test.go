package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingDelegationABI(t *testing.T) {
	delegationMethod := precompile.NewDelegationMethod(nil)

	require.Equal(t, 2, len(delegationMethod.Method.Inputs))
	require.Equal(t, 2, len(delegationMethod.Method.Outputs))
}

func (suite *PrecompileTestSuite) TestDelegation() {
	delegationMethod := precompile.NewDelegationMethod(nil)
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
				}, fmt.Errorf("validator not found: %s", newVal)
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
				}, fmt.Errorf("delegation failed: invalid validator address: %s", newVal)
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
				}, fmt.Errorf("delegation failed: validator not found: %s", newVal)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val0 := suite.GetFirstValidator()

			delAmount := helpers.NewRandAmount()

			stakingContract := precompile.GetAddress()
			delAddr := suite.signer.Address()
			stakingABI := precompile.GetABI()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				delAddr = suite.staking
				stakingABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
			}

			pack, err := stakingABI.Pack(TestDelegateName, val0.GetOperator().String())
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, delAmount.BigInt(), pack)
			suite.Require().False(res.Failed(), res.VmError)

			unpack, err := stakingABI.Methods[TestDelegateName].Outputs.Unpack(res.Ret)
			suite.Require().NoError(err)
			delShares := unpack[0].(*big.Int)

			suite.Commit()

			args, errResult := tc.malleate(val0.GetOperator(), delAddr)
			packData, err := delegationMethod.PackInput(args)
			suite.Require().NoError(err)
			if strings.HasPrefix(tc.name, "contract") {
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegationName, args.Validator, args.Delegator)
				suite.Require().NoError(err)
			}
			res, err = suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, suite.signer.Address(), &stakingContract, nil, packData, false)

			delegation, found := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val0.GetOperator())
			suite.Require().True(found)

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				delValue, err := stakingABI.Methods[TestDelegationName].Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal(delegation.GetShares().TruncateInt().String(), delValue[0].(*big.Int).String())
				suite.Require().Equal(delShares.String(), delValue[1].(*big.Int).String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
