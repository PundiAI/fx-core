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

	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
)

func TestStakingRedelegateABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.RedelegateMethodName]
	require.Equal(t, method, staking.RedelegateMethod)
	require.Equal(t, 3, len(staking.RedelegateMethod.Inputs))
	require.Equal(t, 3, len(staking.RedelegateMethod.Outputs))

	event := stakingABI.Events[staking.RedelegateEventName]
	require.Equal(t, event, staking.RedelegateEvent)
	require.Equal(t, 6, len(staking.RedelegateEvent.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestRedelegate() {
	mustPackRedelegateFunc := func(valSrc, valDst sdk.ValAddress, shares sdk.Dec) []byte {
		pack, err := staking.GetABI().Pack(staking.RedelegateMethodName, valSrc.String(), valDst.String(), shares.TruncateInt().BigInt())
		suite.Require().NoError(err)
		return pack
	}
	testABI := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
	mustPackTestABIFunc := func(valSrc, valDst sdk.ValAddress, shares sdk.Dec) []byte {
		pack, err := testABI.Pack(StakingTestRedelegateName, valSrc.String(), valDst.String(), shares.TruncateInt().BigInt())
		suite.Require().NoError(err)
		return pack
	}

	testCases := []struct {
		name     string
		malleate func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(_ *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack := mustPackRedelegateFunc(valSrc, valDst, shares)
				return pack, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator src",
			malleate: func(signer *helpers.Signer, _, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				pack := mustPackRedelegateFunc(valSrc, valDst, shares)
				return pack, []string{valSrc.String()}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("validator src not found: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed - invalid validator dst",
			malleate: func(signer *helpers.Signer, valSrc, _ sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				valDst := sdk.ValAddress(suite.signer.Address().Bytes())
				pack := mustPackRedelegateFunc(valSrc, valDst, shares)
				return pack, []string{valDst.String()}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("validator dst not found: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed - no delegation",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				// undelegate all before redelegate
				_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, signer.AccAddress(), valSrc, shares)
				suite.Require().NoError(err)

				pack := mustPackRedelegateFunc(valSrc, valDst, shares)
				return pack, []string{}
			},
			error: func(errArgs []string) string {
				return "delegation not found"
			},
			result: false,
		},
		{
			name: "failed - insufficient redelegate shares",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack := mustPackRedelegateFunc(valSrc, valDst, shares.Add(sdk.NewDec(1e18)))
				return pack, []string{}
			},
			error: func(errArgs []string) string {
				return "insufficient shares to redelegate"
			},
			result: false,
		},
		{
			name: "failed - redelegate limit",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				entries := suite.app.StakingKeeper.MaxEntries(suite.ctx)
				for i := uint32(0); i < entries; i++ {
					_, err := suite.app.StakingKeeper.BeginRedelegation(suite.ctx,
						signer.AccAddress(), valSrc, valDst, shares.QuoInt64(10))
					suite.Require().NoError(err)
				}

				pack := mustPackRedelegateFunc(valSrc, valDst, shares.QuoInt64(10))
				return pack, []string{}
			},
			error: func(errArgs []string) string {
				return "too many redelegation entries for (delegator, src-validator, dst-validator) tuple"
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack := mustPackTestABIFunc(valSrc, valDst, shares)
				return pack, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator src",
			malleate: func(signer *helpers.Signer, _, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				valSrc := sdk.ValAddress(suite.signer.Address().Bytes())
				pack := mustPackTestABIFunc(valSrc, valDst, shares)
				return pack, []string{valSrc.String()}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: redelegate failed: validator src not found: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - failed - invalid validator dst",
			malleate: func(signer *helpers.Signer, valSrc, _ sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				valDst := sdk.ValAddress(suite.signer.Address().Bytes())
				pack := mustPackTestABIFunc(valSrc, valDst, shares)
				return pack, []string{valDst.String()}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: redelegate failed: validator dst not found: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "contract - failed - no delegation",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				// undelegate all before redelegate
				_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, suite.staking.Bytes(), valSrc, shares)
				suite.Require().NoError(err)

				pack := mustPackTestABIFunc(valSrc, valDst, shares)
				return pack, []string{}
			},
			error: func(errArgs []string) string {
				return "execution reverted: redelegate failed: delegation not found"
			},
			result: false,
		},
		{
			name: "contract - failed - insufficient redelegate shares",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				pack := mustPackTestABIFunc(valSrc, valDst, shares.Add(sdk.NewDec(1e18)))
				return pack, []string{}
			},
			error: func(errArgs []string) string {
				return "execution reverted: redelegate failed: insufficient shares to redelegate"
			},
			result: false,
		},
		{
			name: "contract - failed - redelegate limit",
			malleate: func(signer *helpers.Signer, valSrc, valDst sdk.ValAddress, shares sdk.Dec) ([]byte, []string) {
				entries := suite.app.StakingKeeper.MaxEntries(suite.ctx)
				for i := uint32(0); i < entries; i++ {
					_, err := suite.app.StakingKeeper.BeginRedelegation(suite.ctx,
						suite.staking.Bytes(), valSrc, valDst, shares.QuoInt64(10))
					suite.Require().NoError(err)
				}

				pack := mustPackTestABIFunc(valSrc, valDst, shares.QuoInt64(10))
				return pack, []string{}
			},
			error: func(errArgs []string) string {
				return "execution reverted: redelegate failed: too many redelegation entries for (delegator, src-validator, dst-validator) tuple"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]
			val1 := vals[1]
			delAmt := sdkmath.NewInt(int64(tmrand.Intn(1000) + 100)).Mul(sdkmath.NewInt(1e18))
			signer := suite.RandSigner()
			helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(),
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmt)))

			stakingContract := staking.GetAddress()
			stakingABI := staking.GetABI()
			delegateMethodName := staking.DelegateMethodName
			redelegateMethodName := staking.RedelegateMethodName
			delAddr := signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				stakingABI = testABI
				delegateMethodName = StakingTestDelegateName
				redelegateMethodName = StakingTestRedelegateName
				delAddr = suite.staking
			}

			// delegate to val0
			pack, err := stakingABI.Pack(delegateMethodName, val0.GetOperator().String())
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

			delegation0, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val0.GetOperator())
			suite.Require().True(found)
			_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val1.GetOperator())
			suite.Require().False(found)

			val0, found = suite.app.StakingKeeper.GetValidator(suite.ctx, val0.GetOperator())
			suite.Require().True(found)

			pack, errArgs := tc.malleate(signer, val0.GetOperator(), val1.GetOperator(), delegation0.Shares)
			tx, err = suite.PackEthereumTx(signer, stakingContract, big.NewInt(0), pack)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.NoError(err)
				suite.False(res.Failed(), res.VmError)

				unpack, err := stakingABI.Unpack(redelegateMethodName, res.Ret)
				suite.NoError(err)
				// amount,reward,completionTime
				reward := unpack[1].(*big.Int)
				suite.True(reward.Cmp(big.NewInt(0)) == 1, reward.String())

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, delAddr.Bytes())
				suite.True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.NoError(err)
				suite.Equal(totalAfter, totalBefore)

				_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val0.GetOperator())
				suite.False(found)
				delegation1New, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val1.GetOperator())
				suite.True(found)
				suite.Equal(delegation0.Shares, delegation1New.Shares)

				redelegates := suite.app.StakingKeeper.GetAllRedelegations(suite.ctx, delAddr.Bytes(), val0.GetOperator(), val1.GetOperator())
				suite.Equal(1, len(redelegates))

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] != staking.RedelegateEvent.ID.String() {
						continue
					}

					suite.Equal(log.Address, staking.GetAddress().String())
					suite.Equal(log.Topics[1], delAddr.Hash().String())
					unpack, err := staking.RedelegateEvent.Inputs.NonIndexed().Unpack(log.Data)
					suite.NoError(err)
					unpackValidatorSrc := unpack[0].(string)
					suite.Equal(unpackValidatorSrc, val0.GetOperator().String())
					unpackValidatorDst := unpack[1].(string)
					suite.Equal(unpackValidatorDst, val1.GetOperator().String())
					shares := unpack[2].(*big.Int)
					suite.Require().Equal(shares.String(), delegation0.Shares.TruncateInt().BigInt().String())
					amount := unpack[3].(*big.Int)
					suite.Require().Equal(amount.String(), val0.TokensFromShares(delegation0.Shares).TruncateInt().String())
					completionTime := unpack[4].(*big.Int)
					suite.Require().Equal(completionTime.Int64(), redelegates[0].Entries[0].CompletionTime.Unix())

					existLog = true
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type != stakingtypes.EventTypeRedelegate {
						continue
					}
					for _, attr := range event.Attributes {
						if string(attr.Key) == stakingtypes.AttributeKeySrcValidator {
							suite.Equal(string(attr.Value), val0.GetOperator().String())
						}
						if string(attr.Key) == stakingtypes.AttributeKeyDstValidator {
							suite.Equal(string(attr.Value), val1.GetOperator().String())
						}
						if string(attr.Key) == sdk.AttributeKeyAmount {
							suite.Require().Equal(string(attr.Value), val0.TokensFromShares(delegation0.Shares).TruncateInt().String())
						}
						if string(attr.Key) == stakingtypes.AttributeKeyCompletionTime {
							suite.Require().Equal(string(attr.Value), redelegates[0].Entries[0].CompletionTime.Format(time.RFC3339))
						}
					}
					existEvent = true
				}
				suite.True(existEvent)
			} else {
				suite.True(err != nil || res.Failed())
				if err != nil {
					suite.Equal(tc.error(errArgs), err.Error())
					return
				}
				if res.VmError != vm.ErrExecutionReverted.Error() {
					suite.Equal(tc.error(errArgs), res.VmError)
					return
				}
				if len(res.Ret) > 0 {
					reason, err := abi.UnpackRevert(common.CopyBytes(res.Ret))
					suite.NoError(err)
					suite.Equal(tc.error(errArgs), reason)
					return
				}
				suite.Equal(tc.error(errArgs), vm.ErrExecutionReverted.Error())
			}
		})
	}
}
