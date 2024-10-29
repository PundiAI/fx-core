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

func TestSlashingInfoABI(t *testing.T) {
	slashingInfoMethod := precompile.NewSlashingInfoMethod(nil)

	require.Len(t, slashingInfoMethod.Method.Inputs, 1)
	require.Len(t, slashingInfoMethod.Method.Outputs, 2)
}

func (suite *PrecompileTestSuite) TestSlashingInfo() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress) (contract.SlashingInfoArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress) (contract.SlashingInfoArgs, error) {
				return contract.SlashingInfoArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress) (contract.SlashingInfoArgs, error) {
				valStr := val.String() + "1"
				return contract.SlashingInfoArgs{
					Validator: valStr,
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress) (contract.SlashingInfoArgs, error) {
				return contract.SlashingInfoArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress) (contract.SlashingInfoArgs, error) {
				valStr := val.String() + "1"
				return contract.SlashingInfoArgs{
					Validator: valStr,
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator := suite.GetFirstValAddr()
			spender := suite.NewSigner()
			allowanceAmt := helpers.NewRandAmount()

			suite.SetAllowance(operator, suite.signer.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			args, expectErr := tc.malleate(operator)

			suite.WithContract(suite.stakingAddr)
			if strings.HasPrefix(tc.name, "contract") {
				suite.WithContract(suite.stakingTestAddr)
			}

			jailed, missed := suite.WithError(expectErr).SlashingInfo(suite.Ctx, args)

			if tc.result {
				validator, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, operator)
				suite.Require().NoError(err)
				suite.Require().Equal(validator.Jailed, jailed)
				consAddr, err := validator.GetConsAddr()
				suite.Require().NoError(err)
				signingInfo, err := suite.App.SlashingKeeper.GetValidatorSigningInfo(suite.Ctx, consAddr)
				suite.Require().NoError(err)
				suite.Require().Equal(signingInfo.MissedBlocksCounter, missed.Int64())
			}
		})
	}
}
