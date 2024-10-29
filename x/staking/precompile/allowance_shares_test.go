package precompile_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func TestStakingAllowanceSharesABI(t *testing.T) {
	allowanceSharesMethod := precompile.NewAllowanceSharesMethod(nil)

	require.Len(t, allowanceSharesMethod.Method.Inputs, 3)
	require.Len(t, allowanceSharesMethod.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestAllowanceShares() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error) {
				return contract.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     owner.Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - default allowance zero",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error) {
				return contract.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     suite.NewSigner().Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error) {
				valStr := val.String() + "1"

				return contract.AllowanceSharesArgs{
					Validator: valStr,
					Owner:     suite.NewSigner().Address(),
					Spender:   spender.Address(),
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error) {
				return contract.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     owner.Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - default allowance zero",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error) {
				return contract.AllowanceSharesArgs{
					Validator: val.String(),
					Owner:     suite.NewSigner().Address(),
					Spender:   spender.Address(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, owner, spender *helpers.Signer) (contract.AllowanceSharesArgs, error) {
				valStr := val.String() + "1"

				return contract.AllowanceSharesArgs{
					Validator: valStr,
					Owner:     suite.NewSigner().Address(),
					Spender:   spender.Address(),
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			valAddr := suite.GetFirstValAddr()
			spender := suite.NewSigner()
			allowanceAmt := helpers.NewRandAmount()

			suite.SetAllowance(valAddr, suite.signer.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			args, expectErr := tc.malleate(valAddr, suite.signer, spender)

			suite.WithContract(suite.stakingAddr)
			if strings.HasPrefix(tc.name, "contract") {
				suite.WithContract(suite.stakingTestAddr)
			}
			shares := suite.WithError(expectErr).AllowanceShares(suite.Ctx, args)
			if tc.result && shares.Sign() > 0 {
				suite.Require().Equal(allowanceAmt.BigInt(), shares)
			}
		})
	}
}
