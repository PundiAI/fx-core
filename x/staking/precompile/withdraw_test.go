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

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingWithdrawABI(t *testing.T) {
	withdrawMethod := precompile.NewWithdrawMethod(nil)

	require.Equal(t, 1, len(withdrawMethod.Method.Inputs))
	require.Equal(t, 1, len(withdrawMethod.Method.Outputs))

	require.Equal(t, 3, len(withdrawMethod.Event.Inputs))
}

func (suite *PrecompileTestSuite) TestWithdraw() {
	withdrawMethod := precompile.NewWithdrawMethod(nil)
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
			helpers.AddTestAddr(suite.App, suite.Ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmt)))

			stakingContract := precompile.GetAddress()
			stakingABI := precompile.GetABI()
			delAddr := signer.Address()
			value := big.NewInt(0)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delAddr = suite.staking
				value = delAmt.BigInt()
			}

			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)

			pack, err := stakingABI.Pack(TestDelegateV2Name, val.GetOperator(), delAmt.BigInt())
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
			packData, err := withdrawMethod.PackInput(args)
			suite.Require().NoError(err)
			if strings.HasPrefix(tc.name, "contract") {
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestWithdrawName, args.Validator)
				suite.Require().NoError(err)
			}
			res = suite.EthereumTx(signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				totalAfter, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(totalAfter, totalBefore)

				unpack, err := stakingABI.Unpack(TestWithdrawName, res.Ret)
				suite.Require().NoError(err)
				reward := unpack[0].(*big.Int)
				chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, delAddr.Bytes())
				suite.Require().True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == withdrawMethod.Event.ID.String() {
						suite.Require().Equal(log.Address, precompile.GetAddress().String())

						event, err := withdrawMethod.UnpackEvent(log.ToEthereum())
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
