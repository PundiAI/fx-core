package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestSlashingInfoABI(t *testing.T) {
	slashingInfoMethod := precompile.NewSlashingInfoMethod(nil)

	require.Equal(t, 1, len(slashingInfoMethod.Method.Inputs))
	require.Equal(t, 2, len(slashingInfoMethod.Method.Outputs))
}

func (suite *PrecompileTestSuite) TestSlashingInfo() {
	slashingInfoMethod := precompile.NewSlashingInfoMethod(nil)
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress) (types.SlashingInfoArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress) (types.SlashingInfoArgs, error) {
				return types.SlashingInfoArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress) (types.SlashingInfoArgs, error) {
				valStr := val.String() + "1"
				return types.SlashingInfoArgs{
					Validator: valStr,
				}, fmt.Errorf("invalid validator address: %s", valStr)
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress) (types.SlashingInfoArgs, error) {
				return types.SlashingInfoArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress) (types.SlashingInfoArgs, error) {
				valStr := val.String() + "1"
				return types.SlashingInfoArgs{
					Validator: valStr,
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
			suite.App.StakingKeeper.SetAllowance(suite.Ctx, val.GetOperator(), owner.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			args, errResult := tc.malleate(val.GetOperator())

			packData, err := slashingInfoMethod.PackInput(args)
			suite.Require().NoError(err)
			stakingContract := precompile.GetAddress()

			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestSlashingInfoName, args.Validator)
				suite.Require().NoError(err)
			}

			res := suite.EthereumTx(owner, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)
				jailed, missed, err := slashingInfoMethod.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				validator, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val.GetOperator())
				suite.True(found)
				suite.Equal(validator.Jailed, jailed)
				consAddr, err := validator.GetConsAddr()
				suite.NoError(err)
				signingInfo, found := suite.App.SlashingKeeper.GetValidatorSigningInfo(suite.Ctx, consAddr)
				suite.True(found)
				suite.Equal(signingInfo.MissedBlocksCounter, missed.Int64())
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
