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

func TestStakingDelegateABI(t *testing.T) {
	stakingABI := fxtypes.MustABIJson(staking.JsonABI)

	method := stakingABI.Methods[staking.DelegateMethod.Name]
	require.Equal(t, method, staking.DelegateMethod)
}

func (suite *PrecompileTestSuite) TestDelegate() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val0 := vals[0]

	_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.staking.Bytes(), val0.GetOperator())
	suite.False(found)

	delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
	pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val0.GetOperator().String())
	suite.Require().NoError(err)

	ethTx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
	suite.Require().NoError(err)

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	del, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.staking.Bytes(), val0.GetOperator())
	suite.True(found)

	val0Del, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val0.GetOperator())
	suite.Require().True(found)

	suite.Require().Equal(del.Shares, val0Del.DelegatorShares.Sub(val0.DelegatorShares))
}

func (suite *PrecompileTestSuite) TestAccountDelegate() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val0 := vals[0]

	_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.signer.AccAddress(), val0.GetOperator())
	suite.False(found)

	delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
	pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String())
	suite.Require().NoError(err)

	ethTx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
	suite.Require().NoError(err)

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	del, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.signer.AccAddress(), val0.GetOperator())
	suite.True(found)

	val0Del, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val0.GetOperator())
	suite.Require().True(found)

	suite.Require().Equal(del.Shares, val0Del.DelegatorShares.Sub(val0.DelegatorShares))
}

func (suite *PrecompileTestSuite) TestDelegates() {
	testCases := []struct {
		name     string
		malleate func() (*evmtypes.MsgEthereumTx, error, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				val := val0.GetOperator().String() + "1"
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
				return tx, err, []string{val}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid amount",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(0)
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
				return tx, err, []string{delAmt.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: invalid delegate amount: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid value",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := big.NewInt(0)
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt, pack)
				return tx, err, []string{delAmt.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: invalid delegate amount: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid validator address",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				val := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
				return tx, err, []string{val}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: validator not found: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			var res *evmtypes.MsgEthereumTxResponse
			ethTx, err, errArgs := tc.malleate()
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

func (suite *PrecompileTestSuite) TestAccountDelegates() {
	delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
	testCases := []struct {
		name     string
		malleate func() (*evmtypes.MsgEthereumTx, error, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "ok - delegate already",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
				suite.Require().NoError(err)
				res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				pack, err = fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err = suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				val := val0.GetOperator().String() + "1"
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
				return tx, err, []string{val}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid amount",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(0)
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
				return tx, err, []string{delAmt.String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid delegate amount: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid validator address",
			malleate: func() (*evmtypes.MsgEthereumTx, error, []string) {
				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				val := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
				return tx, err, []string{val}
			},
			error: func(args []string) string {
				return fmt.Sprintf("validator not found: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			var res *evmtypes.MsgEthereumTxResponse
			ethTx, err, errArgs := tc.malleate()
			if err == nil {
				// generate rewards
				suite.Commit()
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
