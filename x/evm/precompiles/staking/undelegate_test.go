package staking_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

func TestStakingUndelegateABI(t *testing.T) {
	stakingABI := fxtypes.MustABIJson(staking.JsonABI)

	method := stakingABI.Methods[staking.UndelegateMethod.Name]
	require.Equal(t, method, staking.UndelegateMethod)
	require.Equal(t, 2, len(staking.UndelegateMethod.Inputs))
	require.Equal(t, 3, len(staking.UndelegateMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestUndelegates() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestUndelegateName, val.String(), shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestUndelegateName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: undelegate failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestUndelegateName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: undelegate failed: validator not found: %s", errArgs[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]

			delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
			pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val0.GetOperator().String())
			suite.Require().NoError(err)
			delegateEthTx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.staking.Bytes(), val0.GetOperator())
			suite.Require().True(found)

			ethTx, err, errArgs := tc.malleate(val0.GetOperator(), delegation.Shares)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				} else {
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
				}
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestAccountUndelegates() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.UndelegateMethodName, val.String(), shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.GetPrecompileAddress(), big.NewInt(0), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.UndelegateMethodName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.GetPrecompileAddress(), big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.UndelegateMethodName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.GetPrecompileAddress(), big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("validator not found: %s", errArgs[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]

			delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
			pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String())
			suite.Require().NoError(err)
			delegateEthTx, err := suite.PackEthereumTx(suite.signer, staking.GetPrecompileAddress(), delAmt.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.signer.AccAddress(), val0.GetOperator())
			suite.Require().True(found)

			ethTx, err, errArgs := tc.malleate(val0.GetOperator(), delegation.Shares)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				} else {
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
				}
			}
		})
	}
}
