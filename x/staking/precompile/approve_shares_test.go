package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingApproveSharesABI(t *testing.T) {
	approveSharesMethod := precompile.NewApproveSharesMethod(nil)

	require.Len(t, approveSharesMethod.Method.Inputs, 3)
	require.Len(t, approveSharesMethod.Method.Outputs, 1)

	require.Len(t, approveSharesMethod.Event.Inputs, 4)
}

func (suite *PrecompileTestSuite) TestApproveShares() {
	approveSharesMethod := precompile.NewApproveSharesMethod(nil)
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
			val := suite.GetFirstValidator()
			owner := suite.RandSigner()
			spender := suite.RandSigner()
			allowanceAmt := helpers.NewRandAmount()

			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)

			args, errResult := tc.malleate(operator, spender, allowanceAmt)

			packData, err := approveSharesMethod.PackInput(args)
			suite.Require().NoError(err)
			stakingContract := precompile.GetAddress()
			sender := owner.Address()

			if strings.HasPrefix(tc.name, "contract") {
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestApproveSharesName, args.Validator, args.Spender, args.Shares)
				suite.Require().NoError(err)

				stakingContract = suite.staking
				sender = suite.staking
			}

			allowance := suite.App.StakingKeeper.GetAllowance(suite.Ctx, operator, owner.AccAddress(), spender.AccAddress())
			suite.Require().Equal(0, allowance.Cmp(big.NewInt(0)))

			res := suite.EthereumTx(owner, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				allowance = suite.App.StakingKeeper.GetAllowance(suite.Ctx, operator, sender.Bytes(), spender.AccAddress())
				if allowance.Cmp(big.NewInt(0)) != 0 {
					suite.Require().Equal(0, allowance.Cmp(allowanceAmt.BigInt()))
				}

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == approveSharesMethod.Event.ID.String() {
						suite.Require().Equal(log.Address, precompile.GetAddress().String())
						event, err := approveSharesMethod.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.Owner, sender)
						suite.Require().Equal(event.Spender, spender.Address())
						suite.Require().Equal(event.Validator, val.GetOperator())
						if allowance.Cmp(big.NewInt(0)) != 0 {
							suite.Require().Equal(event.Shares.String(), allowanceAmt.BigInt().String())
						}
						existLog = true
					}
				}
				suite.Require().True(existLog)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
