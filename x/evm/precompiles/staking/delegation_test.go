package staking_test

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

func (suite *PrecompileTestSuite) TestDelegation() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegationMethodName, val.String(), suite.precompileStaking)
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "ok - zero",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegationMethodName, val.String(), suite.signer.Address())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegationMethodName, newVal, suite.precompileStaking)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("delegation failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.AccAddress()).String()
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegationMethodName, newVal, suite.precompileStaking)
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

			delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
			pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
			suite.Require().NoError(err)
			delegateEthTx := suite.PackEthereumTx(suite.signer, suite.precompileStaking, delAmt.BigInt(), pack)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.precompileStaking.Bytes(), val0.GetOperator())
			suite.Require().True(found)

			suite.Commit()

			pack, errArgs := tc.malleate(val0.GetOperator(), delegation.Shares)
			res, err = suite.app.EvmKeeper.CallEVMWithoutGas(suite.ctx, suite.signer.Address(), &suite.precompileStaking, pack, false)
			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				var delRes struct {
					Value *big.Int
				}
				stakingABI := fxtypes.MustABIJson(StakingTestABI)
				err = stakingABI.UnpackIntoInterface(&delRes, staking.DelegationMethodName, res.Ret)
				suite.Require().NoError(err)
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
