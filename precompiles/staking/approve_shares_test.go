package staking_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/precompiles/staking"
	"github.com/functionx/fx-core/v8/testutil/helpers"
)

func TestStakingApproveSharesABI(t *testing.T) {
	approveSharesABI := staking.NewApproveSharesABI()

	require.Len(t, approveSharesABI.Method.Inputs, 3)
	require.Len(t, approveSharesABI.Method.Outputs, 1)

	require.Len(t, approveSharesABI.Event.Inputs, 4)
}

func (suite *StakingPrecompileTestSuite) TestApproveShares() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxcontract.ApproveSharesArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxcontract.ApproveSharesArgs, error) {
				return fxcontract.ApproveSharesArgs{
					Validator: val.String(),
					Spender:   spender.Address(),
					Shares:    allowance.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - approve zero",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxcontract.ApproveSharesArgs, error) {
				return fxcontract.ApproveSharesArgs{
					Validator: val.String(),
					Spender:   spender.Address(),
					Shares:    big.NewInt(0),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxcontract.ApproveSharesArgs, error) {
				valStr := val.String() + "1"
				return fxcontract.ApproveSharesArgs{
					Validator: valStr,
					Spender:   spender.Address(),
					Shares:    allowance.BigInt(),
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

			allowance := suite.App.StakingKeeper.GetAllowance(suite.Ctx, valAddr, suite.signer.AccAddress(), spender.AccAddress())
			suite.Require().Equal(0, allowance.Cmp(big.NewInt(0)))

			args, expectErr := tc.malleate(valAddr, spender, allowanceAmt)

			delAddr := suite.GetDelAddr()

			res := suite.WithError(expectErr).ApproveShares(suite.Ctx, args)
			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				allowance = suite.App.StakingKeeper.GetAllowance(suite.Ctx, valAddr, delAddr.Bytes(), spender.AccAddress())
				if allowance.Cmp(big.NewInt(0)) != 0 {
					suite.Require().Equal(0, allowance.Cmp(allowanceAmt.BigInt()))
				}

				existLog := false
				approveSharesABI := staking.NewApproveSharesABI()
				for _, log := range res.Logs {
					if log.Topics[0] == approveSharesABI.Event.ID.String() {
						suite.Require().Equal(fxcontract.StakingAddress, log.Address)
						event, err := approveSharesABI.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.Owner, delAddr)
						suite.Require().Equal(event.Spender, spender.Address())
						suite.Require().Equal(event.Validator, valAddr.String())
						if allowance.Cmp(big.NewInt(0)) != 0 {
							suite.Require().Equal(event.Shares.String(), allowanceAmt.BigInt().String())
						}
						existLog = true
					}
				}
				suite.Require().True(existLog)
			}
		})
	}
}
