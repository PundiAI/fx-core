package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	fxstakingtypes "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func TestStakingApproveSharesABI(t *testing.T) {
	approveSharesMethod := precompile.NewApproveSharesMethod(nil)

	require.Len(t, approveSharesMethod.Method.Inputs, 3)
	require.Len(t, approveSharesMethod.Method.Outputs, 1)

	require.Len(t, approveSharesMethod.Event.Inputs, 4)
}

func (suite *PrecompileTestSuite) TestApproveShares() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error) {
				return fxstakingtypes.ApproveSharesArgs{
					Validator: val.String(),
					Spender:   spender.Address(),
					Shares:    allowance.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - approve zero",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error) {
				return fxstakingtypes.ApproveSharesArgs{
					Validator: val.String(),
					Spender:   spender.Address(),
					Shares:    big.NewInt(0),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error) {
				valStr := val.String() + "1"
				return fxstakingtypes.ApproveSharesArgs{
					Validator: valStr,
					Spender:   spender.Address(),
					Shares:    allowance.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error) {
				return fxstakingtypes.ApproveSharesArgs{
					Validator: val.String(),
					Spender:   spender.Address(),
					Shares:    allowance.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - approve zero",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error) {
				return fxstakingtypes.ApproveSharesArgs{
					Validator: val.String(),
					Spender:   spender.Address(),
					Shares:    big.NewInt(0),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, spender *helpers.Signer, allowance sdkmath.Int) (fxstakingtypes.ApproveSharesArgs, error) {
				valStr := val.String() + "1"
				return fxstakingtypes.ApproveSharesArgs{
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

			delAddr := suite.signer.Address()
			suite.WithContract(suite.stakingAddr)
			if strings.HasPrefix(tc.name, "contract") {
				suite.WithContract(suite.stakingTestAddr)
				delAddr = suite.stakingTestAddr
			}

			res := suite.WithError(expectErr).ApproveShares(suite.Ctx, args)
			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				allowance = suite.App.StakingKeeper.GetAllowance(suite.Ctx, valAddr, delAddr.Bytes(), spender.AccAddress())
				if allowance.Cmp(big.NewInt(0)) != 0 {
					suite.Require().Equal(0, allowance.Cmp(allowanceAmt.BigInt()))
				}

				existLog := false
				for _, log := range res.Logs {
					abi := precompile.NewApproveSharesABI()
					if log.Topics[0] == abi.Event.ID.String() {
						suite.Require().Equal(log.Address, suite.stakingAddr.String())
						event, err := abi.UnpackEvent(log.ToEthereum())
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
