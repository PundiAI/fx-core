package precompile_test

import (
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func (suite *PrecompileTestSuite) TestRedelegate() {
	testCases := []struct {
		name     string
		malleate func(valSrc, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (interface{}, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok v2",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator src",
			malleate: func(_, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (interface{}, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},

		{
			name: "contract - ok v2",
			malleate: func(valSrc, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (interface{}, error) {
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - v2 invalid validator src",
			malleate: func(_, valDst sdk.ValAddress, shares sdkmath.LegacyDec, delAmount sdkmath.Int) (interface{}, error) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				return types.RedelegateV2Args{
					ValidatorSrc: valSrc.String(),
					ValidatorDst: valDst.String(),
					Amount:       delAmount.BigInt(),
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			vals := suite.GetValidators()
			val0 := vals[0]
			val1 := vals[1]

			delAmt := helpers.NewRandAmount()

			stakingContract := suite.stakingAddr
			delAddr := suite.signer.Address()
			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
				delAddr = suite.stakingTestAddr
				value = delAmt.BigInt()
			}

			operator0, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val0.GetOperator())
			suite.Require().NoError(err)

			operator1, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val1.GetOperator())
			suite.Require().NoError(err)

			// delegate to val0
			pack, err := suite.delegateV2Method.PackInput(types.DelegateV2Args{
				Validator: val0.GetOperator(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			delegation0, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator0)
			suite.Require().NoError(err)
			_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator1)
			suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

			val0, err = suite.App.StakingKeeper.GetValidator(suite.Ctx, operator0)
			suite.Require().NoError(err)

			var packData []byte
			args, errResult := tc.malleate(operator0, operator1, delegation0.Shares, delAmt)
			packData, err = suite.redelegateV2Method.PackInput(args.(types.RedelegateV2Args))
			suite.Require().NoError(err)

			res = suite.EthereumTx(suite.signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator0)
				suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)
				delegation1New, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator1)
				suite.Require().NoError(err)
				suite.Require().Equal(delegation0.Shares, delegation1New.Shares)

				redelegates, err := suite.App.StakingKeeper.GetAllRedelegations(suite.Ctx, delAddr.Bytes(), operator0, operator1)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(redelegates))

				suite.CheckRedelegateLogs(res.Logs, delAddr, val0.GetOperator(), val1.GetOperator(),
					delegation0.Shares.TruncateInt().BigInt(), val0.TokensFromShares(delegation0.Shares).TruncateInt().BigInt(),
					redelegates[0].Entries[0].CompletionTime.Unix())

				suite.CheckRedelegateEvents(suite.Ctx, val0.GetOperator(), val1.GetOperator(),
					val0.TokensFromShares(delegation0.Shares).TruncateInt().BigInt(),
					redelegates[0].Entries[0].CompletionTime)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
