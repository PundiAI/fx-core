package precompile_test

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestValidatorListABI(t *testing.T) {
	validatorListMethod := precompile.NewValidatorListMethod(nil)

	require.Len(t, validatorListMethod.Method.Inputs, 1)
	require.Len(t, validatorListMethod.Method.Outputs, 1)
}

func (suite *PrecompileTestSuite) TestValidatorList() {
	validatroListMethod := precompile.NewValidatorListMethod(nil)
	testCases := []struct {
		name     string
		malleate func() (types.ValidatorListArgs, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func() (types.ValidatorListArgs, error) {
				return types.ValidatorListArgs{
					SortBy: uint8(types.ValidatorSortByPower),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - missed",
			malleate: func() (types.ValidatorListArgs, error) {
				return types.ValidatorListArgs{
					SortBy: uint8(types.ValidatorSortByMissed),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid order value",
			malleate: func() (types.ValidatorListArgs, error) {
				return types.ValidatorListArgs{
					SortBy: 100,
				}, fmt.Errorf("over the sort by limit")
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func() (types.ValidatorListArgs, error) {
				return types.ValidatorListArgs{
					SortBy: uint8(types.ValidatorSortByPower),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - missed",
			malleate: func() (types.ValidatorListArgs, error) {
				return types.ValidatorListArgs{
					SortBy: uint8(types.ValidatorSortByMissed),
				}, nil
			},
			result: true,
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

			// set allowance
			suite.App.StakingKeeper.SetAllowance(suite.Ctx, operator, owner.AccAddress(), spender.AccAddress(), allowanceAmt.BigInt())

			args, errResult := tc.malleate()

			packData, err := validatroListMethod.PackInput(args)
			suite.Require().NoError(err)
			stakingContract := precompile.GetAddress()

			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestValidatorListName, args.SortBy)
				suite.Require().NoError(err)
			}

			res := suite.EthereumTx(owner, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)
				valAddrs, err := validatroListMethod.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				valsByPower, err := suite.App.StakingKeeper.GetBondedValidatorsByPower(suite.Ctx)
				suite.Require().NoError(err)
				suite.Equal(len(valAddrs), len(valsByPower))

				if args.GetSortBy() == types.ValidatorSortByPower {
					for index, addr := range valAddrs {
						suite.Equal(addr, valsByPower[index].OperatorAddress)
					}
				}
				if args.GetSortBy() == types.ValidatorSortByMissed {
					valList := make([]precompile.ValidatorList, 0, len(valsByPower))
					for _, validator := range valsByPower {
						consAddr, err := validator.GetConsAddr()
						suite.Require().NoError(err)
						info, err := suite.App.SlashingKeeper.GetValidatorSigningInfo(suite.Ctx, consAddr)
						suite.Require().NoError(err)
						valList = append(valList, precompile.ValidatorList{
							ValAddr:      validator.OperatorAddress,
							MissedBlocks: info.MissedBlocksCounter,
						})
					}
					sort.Slice(valList, func(i, j int) bool {
						return valList[i].MissedBlocks > valList[j].MissedBlocks
					})
					for index, addr := range valAddrs {
						suite.Equal(addr, valList[index].ValAddr)
					}
				}
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
