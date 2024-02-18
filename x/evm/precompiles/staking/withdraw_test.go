package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
)

func TestStakingWithdrawABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.WithdrawMethodName]
	require.Equal(t, method, staking.WithdrawMethod)
	require.Equal(t, 1, len(staking.WithdrawMethod.Inputs))
	require.Equal(t, 1, len(staking.WithdrawMethod.Outputs))

	event := stakingABI.Events[staking.WithdrawEventName]
	require.Equal(t, event, staking.WithdrawEvent)
	require.Equal(t, 3, len(staking.WithdrawEvent.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestWithdraw() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.WithdrawMethodName, val.String())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := staking.GetABI().Pack(staking.WithdrawMethodName, newVal)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := staking.GetABI().Pack(staking.WithdrawMethodName, newVal)
				suite.Require().NoError(err)

				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return "no validator distribution info"
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestWithdrawName, val.String())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestWithdrawName, newVal)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: withdraw failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestWithdrawName, newVal)
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return "execution reverted: withdraw failed: no validator distribution info"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val := vals[0]

			delAmt := sdkmath.NewInt(int64(tmrand.Intn(1000) + 100)).Mul(sdkmath.NewInt(1e18))
			signer := suite.RandSigner()
			helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmt)))

			stakingContract := staking.GetAddress()
			stakingABI := staking.GetABI()
			delegateMethodName := staking.DelegateMethodName
			withdrawMethodName := staking.WithdrawMethodName
			delAddr := signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delegateMethodName = StakingTestDelegateName
				withdrawMethodName = StakingTestWithdrawName
				delAddr = suite.staking
			}

			pack, err := stakingABI.Pack(delegateMethodName, val.GetOperator().String())
			suite.Require().NoError(err)
			tx, err := suite.PackEthereumTx(signer, stakingContract, delAmt.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, delAddr.Bytes())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
			suite.Require().True(found)

			pack, errArgs := tc.malleate(val.GetOperator(), delegation.Shares)
			tx, err = suite.PackEthereumTx(signer, stakingContract, big.NewInt(0), pack)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(totalAfter, totalBefore)

				unpack, err := stakingABI.Unpack(withdrawMethodName, res.Ret)
				suite.Require().NoError(err)
				reward := unpack[0].(*big.Int)
				suite.Require().True(reward.Cmp(big.NewInt(0)) == 1, reward.String())
				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, delAddr.Bytes())
				suite.Require().True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == staking.WithdrawEvent.ID.String() {
						suite.Require().Equal(log.Address, staking.GetAddress().String())
						suite.Require().Equal(log.Topics[1], delAddr.Hash().String())
						unpack, err := staking.WithdrawEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						unpackValidator := unpack[0].(string)
						suite.Require().Equal(unpackValidator, val.GetOperator().String())
						reward := unpack[1].(*big.Int)
						suite.Require().Equal(reward.String(), chainBalances.AmountOf(fxtypes.DefaultDenom).BigInt().String())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type == distritypes.EventTypeWithdrawRewards {
						for _, attr := range event.Attributes {
							if string(attr.Key) == distritypes.AttributeKeyValidator {
								suite.Require().Equal(string(attr.Value), val.GetOperator().String())
								existEvent = true
							}
							if string(attr.Key) == sdk.AttributeKeyAmount {
								suite.Require().Equal(string(attr.Value), sdk.NewCoin(fxtypes.DefaultDenom, chainBalances.AmountOf(fxtypes.DefaultDenom)).String())
								existEvent = true
							}
							if string(attr.Key) == distritypes.AttributeKeyDelegator {
								suite.Require().Equal(string(attr.Value), sdk.AccAddress(delAddr.Bytes()).String())
								existEvent = true
							}
						}
					}
				}
				suite.Require().True(existEvent)
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
				if res.Failed() {
					if res.VmError != vm.ErrExecutionReverted.Error() {
						suite.Require().Equal(tc.error(errArgs), res.VmError)
					} else {
						if len(res.Ret) > 0 {
							reason, err := abi.UnpackRevert(common.CopyBytes(res.Ret))
							suite.Require().NoError(err)

							suite.Require().Equal(tc.error(errArgs), reason)
						} else {
							suite.Require().Equal(tc.error(errArgs), vm.ErrExecutionReverted.Error())
						}
					}
				} else {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
			}
		})
	}
}
