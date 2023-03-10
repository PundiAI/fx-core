package staking_test

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

func (suite *PrecompileTestSuite) TestDelegate() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val0 := vals[0]

	_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.precompileStaking.Bytes(), val0.GetOperator())
	suite.False(found)

	delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
	pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
	suite.Require().NoError(err)

	ethTx := suite.PackEthereumTx(suite.signer, suite.precompileStaking, delAmt.BigInt(), pack)

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	del, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.precompileStaking.Bytes(), val0.GetOperator())
	suite.True(found)

	val0Del, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val0.GetOperator())
	suite.Require().True(found)

	suite.Require().Equal(del.Shares, val0Del.DelegatorShares.Sub(val0.DelegatorShares))
}

func (suite *PrecompileTestSuite) TestDelegates() {
	testCases := []struct {
		name     string
		malleate func() (*evmtypes.MsgEthereumTx, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func() (*evmtypes.MsgEthereumTx, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
				suite.Require().NoError(err)
				return suite.PackEthereumTx(suite.signer, suite.precompileStaking, delAmt.BigInt(), pack), nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func() (*evmtypes.MsgEthereumTx, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				val := val0.GetOperator().String() + "1"
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val, delAmt.BigInt())
				suite.Require().NoError(err)
				return suite.PackEthereumTx(suite.signer, suite.precompileStaking, delAmt.BigInt(), pack), []string{val}
			},
			error: func(args []string) string {
				return fmt.Sprintf("delegate failed: invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid amount",
			malleate: func() (*evmtypes.MsgEthereumTx, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(0)
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
				suite.Require().NoError(err)
				return suite.PackEthereumTx(suite.signer, suite.precompileStaking, delAmt.BigInt(), pack), []string{
					delAmt.String(),
				}
			},
			error: func(args []string) string {
				return fmt.Sprintf("delegate failed: invalid amount: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed invalid value",
			malleate: func() (*evmtypes.MsgEthereumTx, []string) {
				vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val0 := vals[0]

				delAmt := sdkmath.NewInt(1)
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
				suite.Require().NoError(err)
				return suite.PackEthereumTx(suite.signer, suite.precompileStaking, big.NewInt(0), pack), []string{}
			},
			error: func(args []string) string {
				return "execution reverted"
			},
			result: false,
		},
		{
			name: "failed invalid validator address",
			malleate: func() (*evmtypes.MsgEthereumTx, []string) {
				delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
				val := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val, delAmt.BigInt())
				suite.Require().NoError(err)
				return suite.PackEthereumTx(suite.signer, suite.precompileStaking, delAmt.BigInt(), pack), []string{val}
			},
			error: func(args []string) string {
				return fmt.Sprintf("delegate failed: validator not found: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ethTx, errArgs := tc.malleate()
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
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
				}
			}
		})
	}
}
