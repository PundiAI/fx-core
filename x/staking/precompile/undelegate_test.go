package precompile_test

import (
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func (suite *PrecompileTestSuite) TestUndelegate() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (interface{}, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok v2",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (interface{}, error) {
				return types.UndelegateV2Args{
					Validator: val.String(),
					Amount:    delAmt.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := val.String() + "1"
				return types.UndelegateV2Args{
					Validator: newVal,
					Amount:    delAmt.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},

		{
			name: "contract - ok v2",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (interface{}, error) {
				return types.UndelegateV2Args{
					Validator: val.String(),
					Amount:    delAmt.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - v2 invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec, delAmt sdkmath.Int) (interface{}, error) {
				newVal := val.String() + "1"
				return types.UndelegateV2Args{
					Validator: newVal,
					Amount:    delAmt.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()
			delAmt := helpers.NewRandAmount()

			stakingContract := suite.stakingAddr
			delAddr := suite.signer.Address()
			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
				delAddr = suite.stakingTestAddr
				value = delAmt.BigInt()
			}

			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)

			pack, err := suite.delegateV2Method.PackInput(types.DelegateV2Args{
				Validator: val.GetOperator(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			delegation := suite.GetDelegation(delAddr.Bytes(), operator)

			undelegations, err := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, delAddr.Bytes())
			suite.Require().NoError(err)
			suite.Require().Equal(0, len(undelegations))

			var packData []byte
			args, errResult := tc.malleate(operator, delegation.Shares, delAmt)
			packData, err = suite.undelegateV2Method.PackInput(args.(types.UndelegateV2Args))
			suite.Require().NoError(err)

			res = suite.EthereumTx(suite.signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				undelegations, err := suite.App.StakingKeeper.GetAllUnbondingDelegations(suite.Ctx, delAddr.Bytes())
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(undelegations))
				suite.Require().Equal(1, len(undelegations[0].Entries))
				suite.Require().Equal(sdk.AccAddress(delAddr.Bytes()).String(), undelegations[0].DelegatorAddress)
				suite.Require().Equal(val.GetOperator(), undelegations[0].ValidatorAddress)
				suite.Require().Equal(delAmt, undelegations[0].Entries[0].Balance)

				suite.CheckUndelegateLogs(res.Logs, delAddr, val.GetOperator(), delegation.Shares.TruncateInt().BigInt(),
					undelegations[0].Entries[0].Balance.BigInt(), undelegations[0].Entries[0].CompletionTime)

				suite.CheckUndeledateEvents(suite.Ctx, val.GetOperator(), undelegations[0].Entries[0].Balance.BigInt(),
					undelegations[0].Entries[0].CompletionTime)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
