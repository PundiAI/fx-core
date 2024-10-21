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

	require.Len(t, delegationMethod.Method.Inputs, 2)
	require.Len(t, delegationMethod.Method.Outputs, 2)
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

			stakingContract := precompile.GetAddress()
			delAddr := suite.signer.Address()
			stakingABI := precompile.GetABI()
			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				delAddr = suite.staking
				stakingABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
				value = delAmount.BigInt()
			}

			operator0, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val0.GetOperator())
			suite.Require().NoError(err)

			pack, err := stakingABI.Pack(TestDelegateV2Name, val0.GetOperator(), delAmount.BigInt())
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			args, errResult := tc.malleate(operator0, delAddr)
			packData, err := delegationMethod.PackInput(args)
			suite.Require().NoError(err)
			if strings.HasPrefix(tc.name, "contract") {
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegationName, args.Validator, args.Delegator)
				suite.Require().NoError(err)
			}
			delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator0)
			suite.Require().NoError(err)

			res, _ = suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, suite.signer.Address(), &stakingContract, nil, packData, false)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				delValue, err := stakingABI.Methods[TestDelegationName].Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal(delegation.GetShares().TruncateInt().String(), delValue[0].(*big.Int).String())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
