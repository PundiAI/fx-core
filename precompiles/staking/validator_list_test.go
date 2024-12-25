package staking_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/staking"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

func TestValidatorListABI(t *testing.T) {
	validatorListABI := staking.NewValidatorListABI()

	require.Len(t, validatorListABI.Method.Inputs, 1)
	require.Len(t, validatorListABI.Method.Outputs, 1)
}

func (suite *StakingPrecompileTestSuite) TestValidatorList() {
	testCases := []struct {
		name     string
		malleate func() (contract.ValidatorListArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func() (contract.ValidatorListArgs, error) {
				return contract.ValidatorListArgs{
					SortBy: uint8(contract.ValidatorSortByPower),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - missed",
			malleate: func() (contract.ValidatorListArgs, error) {
				return contract.ValidatorListArgs{
					SortBy: uint8(contract.ValidatorSortByMissed),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid order value",
			malleate: func() (contract.ValidatorListArgs, error) {
				return contract.ValidatorListArgs{
					SortBy: 100,
				}, fmt.Errorf("over the sort by limit")
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

			args, expectErr := tc.malleate()

			valAddrs := suite.WithError(expectErr).ValidatorList(suite.Ctx, args)
			if tc.result {
				valsByPower, err := suite.App.StakingKeeper.GetBondedValidatorsByPower(suite.Ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(len(valAddrs), len(valsByPower))

				if args.GetSortBy() == contract.ValidatorSortByPower {
					for index, addr := range valAddrs {
						suite.Require().Equal(addr, valsByPower[index].OperatorAddress)
					}
				}
				if args.GetSortBy() == contract.ValidatorSortByMissed {
					valList := make([]staking.Validator, 0, len(valsByPower))
					for _, validator := range valsByPower {
						consAddr, err := validator.GetConsAddr()
						suite.Require().NoError(err)
						info, err := suite.App.SlashingKeeper.GetValidatorSigningInfo(suite.Ctx, consAddr)
						suite.Require().NoError(err)
						valList = append(valList, staking.Validator{
							ValAddr:      validator.OperatorAddress,
							MissedBlocks: info.MissedBlocksCounter,
						})
					}
					sort.Slice(valList, func(i, j int) bool {
						return valList[i].MissedBlocks > valList[j].MissedBlocks
					})
					for index, addr := range valAddrs {
						suite.Require().Equal(addr, valList[index].ValAddr)
					}
				}
			}
		})
	}
}
