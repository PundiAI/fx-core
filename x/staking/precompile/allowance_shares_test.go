package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingAllowanceSharesABI(t *testing.T) {
	allowanceSharesMethod := precompile.NewAllowanceSharesMethod(nil)

	require.Len(t, allowanceSharesMethod.Method.Inputs, 3)
	require.Len(t, allowanceSharesMethod.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestAllowanceShares() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error) {
				return types.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     owner.Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - default allowance zero",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error) {
				return types.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     suite.RandSigner().Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error) {
				valStr := val.String() + "1"

				return types.AllowanceSharesArgs{
					Validator: valStr,
					Owner:     suite.RandSigner().Address(),
					Spender:   spender.Address(),
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error) {
				return types.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     owner.Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - default allowance zero",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error) {
				return types.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     suite.RandSigner().Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (types.AllowanceSharesArgs, error) {
				valStr := val.String() + "1"

				return types.AllowanceSharesArgs{
					Validator: valStr,
					Owner:     suite.RandSigner().Address(),
					Spender:   spender.Address(),
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()
			owner := suite.RandSigner()
			spender := suite.RandSigner()
			allowanceAmt := helpers.NewRandAmount()

			// set allowance
			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)
			suite.App.StakingKeeper.SetAllowance(suite.Ctx, operator, owner.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			args, errResult := tc.malleate(operator, owner, spender)

			packData, err := suite.allowanceSharesMethod.PackInput(args)
			suite.Require().NoError(err)
			stakingContract := suite.stakingAddr

			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
			}

			res := suite.EthereumTx(owner, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)
				shares, err := suite.allowanceSharesMethod.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				if shares.Cmp(big.NewInt(0)) != 0 {
					suite.Require().Equal(allowanceAmt.BigInt(), shares)
				}
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
