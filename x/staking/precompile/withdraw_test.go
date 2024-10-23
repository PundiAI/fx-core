package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingWithdrawABI(t *testing.T) {
	withdrawMethod := precompile.NewWithdrawMethod(nil)

	require.Len(t, withdrawMethod.Method.Inputs, 1)
	require.Len(t, withdrawMethod.Method.Outputs, 1)

	require.Len(t, withdrawMethod.Event.Inputs, 3)
}

func (suite *PrecompileTestSuite) TestWithdraw() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error) {
				return types.WithdrawArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error) {
				newVal := val.String() + "1"
				return types.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()

				return types.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error) {
				return types.WithdrawArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error) {
				newVal := val.String() + "1"
				return types.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (types.WithdrawArgs, error) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()

				return types.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()

			delAmt := helpers.NewRandAmount()
			signer := suite.RandSigner()
			suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmt))

			stakingContract := suite.stakingAddr
			delAddr := signer.Address()
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

			res := suite.EthereumTx(signer, stakingContract, value, pack)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, delAddr.Bytes())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			totalBefore, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			delegation := suite.GetDelegation(delAddr.Bytes(), operator)

			args, errResult := tc.malleate(operator, delegation.Shares)
			packData, err := suite.withdrawMethod.PackInput(args)
			suite.Require().NoError(err)
			res = suite.EthereumTx(signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				totalAfter, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(totalAfter, totalBefore)

				reward, err := suite.withdrawMethod.UnpackOutput(res.Ret)
				suite.Require().NoError(err)
				chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, delAddr.Bytes())
				suite.Require().True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == suite.withdrawMethod.Event.ID.String() {
						suite.Require().Equal(log.Address, suite.stakingAddr.String())

						event, err := suite.withdrawMethod.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.Sender, delAddr)
						suite.Require().Equal(event.Validator, val.GetOperator())
						suite.Require().Equal(event.Reward.String(), chainBalances.AmountOf(fxtypes.DefaultDenom).BigInt().String())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.Ctx.EventManager().Events() {
					if event.Type == distritypes.EventTypeWithdrawRewards {
						for _, attr := range event.Attributes {
							if attr.Key == distritypes.AttributeKeyValidator {
								suite.Require().Equal(attr.Value, val.GetOperator())
								existEvent = true
							}
							if attr.Key == sdk.AttributeKeyAmount {
								suite.Require().Equal(attr.Value, sdk.NewCoin(fxtypes.DefaultDenom, chainBalances.AmountOf(fxtypes.DefaultDenom)).String())
								existEvent = true
							}
							if attr.Key == distritypes.AttributeKeyDelegator {
								suite.Require().Equal(attr.Value, sdk.AccAddress(delAddr.Bytes()).String())
								existEvent = true
							}
						}
					}
				}
				suite.Require().True(existEvent)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
