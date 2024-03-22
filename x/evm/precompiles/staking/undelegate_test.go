package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
)

func TestStakingUndelegateABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.UndelegateMethodName]
	require.Equal(t, method, staking.UndelegateMethod)
	require.Equal(t, 2, len(staking.UndelegateMethod.Inputs))
	require.Equal(t, 3, len(staking.UndelegateMethod.Outputs))

	event := stakingABI.Events[staking.UndelegateEventName]
	require.Equal(t, event, staking.UndelegateEvent)
	require.Equal(t, 5, len(staking.UndelegateEvent.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestUndelegate() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack, err := staking.GetABI().Pack(staking.UndelegateMethodName, val.String(), shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := staking.GetABI().Pack(staking.UndelegateMethodName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed - validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := staking.GetABI().Pack(staking.UndelegateMethodName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("validator not found: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestUndelegateName, val.String(), shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := val.String() + "1"
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestUndelegateName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: undelegate failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - failed - validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestUndelegateName, newVal, shares.TruncateInt().BigInt())
				suite.Require().NoError(err)
				return pack, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: undelegate failed: validator not found: %s", errArgs[0])
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
			undelegateMethodName := staking.UndelegateMethodName
			delAddr := signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
				delegateMethodName = StakingTestDelegateName
				undelegateMethodName = StakingTestUndelegateName
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
			undelegations := suite.app.StakingKeeper.GetAllUnbondingDelegations(suite.ctx, delAddr.Bytes())
			suite.Require().Equal(0, len(undelegations))

			pack, errArgs := tc.malleate(val.GetOperator(), delegation.Shares)
			tx, err = suite.PackEthereumTx(signer, stakingContract, big.NewInt(0), pack)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				unpack, err := stakingABI.Unpack(undelegateMethodName, res.Ret)
				suite.Require().NoError(err)
				// amount,reward,completionTime
				reward := unpack[1].(*big.Int)
				suite.Require().True(reward.Cmp(big.NewInt(0)) == 1, reward.String())

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, delAddr.Bytes())
				suite.Require().True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(totalAfter, totalBefore)

				undelegations := suite.app.StakingKeeper.GetAllUnbondingDelegations(suite.ctx, delAddr.Bytes())
				suite.Require().Equal(1, len(undelegations))
				suite.Require().Equal(1, len(undelegations[0].Entries))
				suite.Require().Equal(sdk.AccAddress(delAddr.Bytes()).String(), undelegations[0].DelegatorAddress)
				suite.Require().Equal(val.GetOperator().String(), undelegations[0].ValidatorAddress)
				suite.Require().Equal(delAmt, undelegations[0].Entries[0].Balance)

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == staking.UndelegateEvent.ID.String() {
						suite.Require().Equal(log.Address, staking.GetAddress().String())
						suite.Require().Equal(log.Topics[1], delAddr.Hash().String())
						unpack, err := staking.UndelegateEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						unpackValidator := unpack[0].(string)
						suite.Require().Equal(unpackValidator, val.GetOperator().String())
						shares := unpack[1].(*big.Int)
						suite.Require().Equal(shares.String(), delegation.Shares.TruncateInt().BigInt().String())
						amount := unpack[2].(*big.Int)
						suite.Require().Equal(amount.String(), undelegations[0].Entries[0].Balance.BigInt().String())
						completionTime := unpack[3].(*big.Int)
						suite.Require().Equal(completionTime.Int64(), undelegations[0].Entries[0].CompletionTime.Unix())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type == stakingtypes.EventTypeUnbond {
						for _, attr := range event.Attributes {
							if string(attr.Key) == stakingtypes.AttributeKeyValidator {
								suite.Require().Equal(string(attr.Value), val.GetOperator().String())
								existEvent = true
							}
							if string(attr.Key) == sdk.AttributeKeyAmount {
								suite.Require().Equal(string(attr.Value), undelegations[0].Entries[0].Balance.String())
								existEvent = true
							}
							if string(attr.Key) == stakingtypes.AttributeKeyCompletionTime {
								suite.Require().Equal(string(attr.Value), undelegations[0].Entries[0].CompletionTime.Format(time.RFC3339))
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
				} else {
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
				}
			}
		})
	}
}
