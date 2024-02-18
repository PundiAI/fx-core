package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestStakingDelegationABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.DelegationMethod.Name]
	require.Equal(t, method, staking.DelegationMethod)
	require.Equal(t, 2, len(staking.DelegationMethod.Inputs))
	require.Equal(t, 2, len(staking.DelegationMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestDelegation() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, del common.Address) ([]byte, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.DelegationMethodName, val.String(), del)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "ok - zero",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.DelegationMethodName, val.String(), del)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := staking.GetABI().Pack(staking.DelegationMethodName, newVal, del)
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
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := staking.GetABI().Pack(staking.DelegationMethodName, newVal, del)
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
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationName, val.String(), del)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - ok - zero",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationName, val.String(), del)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationName, newVal, del)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("delegation failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(val sdk.ValAddress, del common.Address) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegationName, newVal, del)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("delegation failed: validator not found: %s", errArgs[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]

			signer := suite.RandSigner()
			delAmount := sdk.NewInt(int64(tmrand.Int() + 100)).Mul(sdk.NewInt(1e18))
			helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

			stakingContract := staking.GetAddress()
			delAddr := signer.Address()
			stakingABI := staking.GetABI()
			delegateFunc := staking.DelegateMethodName
			delegationFunc := staking.DelegationMethodName
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				delAddr = suite.staking
				stakingABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delegateFunc = StakingTestDelegateName
				delegationFunc = StakingTestDelegationName
			}

			pack, err := stakingABI.Pack(delegateFunc, val0.GetOperator().String())
			suite.Require().NoError(err)
			delegateEthTx, err := suite.PackEthereumTx(signer, stakingContract, delAmount.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)
			unpack, err := stakingABI.Methods[delegateFunc].Outputs.Unpack(res.Ret)
			suite.Require().NoError(err)
			delShares := unpack[0].(*big.Int)

			suite.Commit()

			pack, errArgs := tc.malleate(val0.GetOperator(), delAddr)
			res, err = suite.app.EvmKeeper.CallEVMWithoutGas(suite.ctx, suite.signer.Address(), &stakingContract, nil, pack, false)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, sdk.AccAddress(delAddr.Bytes()), val0.GetOperator())
			suite.Require().True(found)

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				delValue, err := stakingABI.Methods[delegationFunc].Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal(delegation.GetShares().TruncateInt().String(), delValue[0].(*big.Int).String())
				suite.Require().Equal(delShares.String(), delValue[1].(*big.Int).String())
			} else {
				suite.Require().True(err != nil || res.Failed())
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
