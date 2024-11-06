package staking_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/precompiles/staking"
	"github.com/functionx/fx-core/v8/testutil/helpers"
)

func TestSlashingInfoABI(t *testing.T) {
	slashingInfoABI := staking.NewSlashingInfoABI()

	require.Len(t, slashingInfoABI.Method.Inputs, 1)
	require.Len(t, slashingInfoABI.Method.Outputs, 2)
}

func (suite *StakingPrecompileTestSuite) TestSlashingInfo() {
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
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator := suite.GetFirstValAddr()
			spender := suite.NewSigner()
			allowanceAmt := helpers.NewRandAmount()

			suite.SetAllowance(operator, suite.signer.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			args, expectErr := tc.malleate(operator)

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
