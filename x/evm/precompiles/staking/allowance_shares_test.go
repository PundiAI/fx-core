package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
)

func TestStakingAllowanceSharesABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.AllowanceSharesMethodName]
	require.Equal(t, method, staking.AllowanceSharesMethod)
	require.Equal(t, 3, len(staking.AllowanceSharesMethod.Inputs))
	require.Equal(t, 1, len(staking.AllowanceSharesMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestAllowanceShares() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.AllowanceSharesMethodName, val.String(), owner.Address(), spender.Address())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "ok - default allowance zero",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.AllowanceSharesMethodName, val.String(), suite.RandSigner().Address(), spender.Address())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string) {
				valStr := val.String() + "1"
				pack, err := staking.GetABI().Pack(staking.AllowanceSharesMethodName, valStr, suite.RandSigner().Address(), spender.Address())
				suite.Require().NoError(err)
				return pack, []string{valStr}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string) {
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestAllowanceSharesName, val.String(), owner.Address(), spender.Address())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - ok - default allowance zero",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string) {
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestAllowanceSharesName, val.String(), suite.RandSigner().Address(), spender.Address())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) ([]byte, []string) {
				valStr := val.String() + "1"
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestAllowanceSharesName, valStr, suite.RandSigner().Address(), spender.Address())
				suite.Require().NoError(err)
				return pack, []string{valStr}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: allowance shares failed: invalid validator address: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val := vals[0]
			owner := suite.RandSigner()
			spender := suite.RandSigner()
			allowanceAmt := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))

			// set allowance
			suite.app.StakingKeeper.SetAllowance(suite.ctx, val.GetOperator(), owner.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			stakingContract := staking.GetAddress()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
			}

			pack, errArgs := tc.malleate(val.GetOperator(), owner, spender)
			tx, err := suite.PackEthereumTx(owner, stakingContract, big.NewInt(0), pack)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				unpack, err := staking.AllowanceSharesMethod.Outputs.Unpack(res.Ret)
				suite.Require().NoError(err)
				shares := unpack[0].(*big.Int)
				if shares.Cmp(big.NewInt(0)) != 0 {
					suite.Require().Equal(allowanceAmt.BigInt(), shares)
				}
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
