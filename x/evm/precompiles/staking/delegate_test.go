package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
)

func TestStakingDelegateABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.DelegateMethodName]
	require.Equal(t, method, staking.DelegateMethod)
	require.Equal(t, 1, len(staking.DelegateMethod.Inputs))
	require.Equal(t, 2, len(staking.DelegateMethod.Outputs))

	event := stakingABI.Events[staking.DelegateEventName]
	require.Equal(t, event, staking.DelegateEvent)
	require.Equal(t, 4, len(staking.DelegateEvent.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestDelegate() {
	testCases := []struct {
		name     string
		malleate func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := staking.GetABI().Pack(staking.DelegateMethodName, val.String())
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "ok - delegate - multiple",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount.Mul(sdk.NewInt(2)))))

				validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val)
				suite.Require().True(found)
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				pack, err := staking.GetABI().Pack(staking.DelegateMethodName, val.String())
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := staking.GetABI().Pack(staking.DelegateMethodName, val.String()+"1")
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), []string{val.String() + "1"}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - invalid value",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := staking.GetABI().Pack(staking.DelegateMethodName, val.String())
				suite.Require().NoError(err)
				return pack, big.NewInt(0), []string{big.NewInt(0).String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("invalid delegate amount: %s", args[0])
			},
			result: false,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := staking.GetABI().Pack(staking.DelegateMethodName, sdk.ValAddress(signer.AccAddress()).String())
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), []string{sdk.ValAddress(signer.AccAddress()).String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("validator not found: %s", args[0])
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegateName, val.String())
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "contract - ok - delegate - multiple",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount.Mul(sdk.NewInt(2)))))

				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegateName, val.String())
				suite.Require().NoError(err)

				tx, err := suite.PackEthereumTx(signer, suite.staking, delAmount.BigInt(), pack)
				suite.Require().NoError(err)

				res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				return pack, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegateName, val.String()+"1")
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), []string{val.String() + "1"}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: invalid validator address: %s", args[0])
			},
			result: false,
		},
		{
			name: "contract - failed - invalid value",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegateName, val.String())
				suite.Require().NoError(err)
				return pack, big.NewInt(0), []string{big.NewInt(0).String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: invalid delegate amount: %s", args[0])
			},
			result: false,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, delAmount sdkmath.Int) ([]byte, *big.Int, []string) {
				helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(StakingTestDelegateName, sdk.ValAddress(signer.AccAddress()).String())
				suite.Require().NoError(err)
				return pack, delAmount.BigInt(), []string{sdk.ValAddress(signer.AccAddress()).String()}
			},
			error: func(args []string) string {
				return fmt.Sprintf("execution reverted: delegate failed: validator not found: %s", args[0])
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val := vals[0]
			delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
			signer := suite.RandSigner()

			contract := staking.GetAddress()
			delAddr := signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				contract = suite.staking
				delAddr = suite.staking
			}

			chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())

			pack, value, errArgs := tc.malleate(signer, val.GetOperator(), delAmount)

			totalBefore, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)
			delBefore, delFound := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
			valBefore, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val.GetOperator())
			suite.Require().True(found)

			tx, err := suite.PackEthereumTx(signer, contract, value, pack)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				chainBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer.AccAddress())
				suite.Require().True(chainBalances.IsZero(), chainBalances.String())

				totalAfter, err := suite.app.BankKeeper.TotalSupply(suite.ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(totalAfter, totalBefore)

				delAfter, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
				suite.Require().True(found)

				vaAfter, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val.GetOperator())
				suite.Require().True(found)

				if !delFound {
					delBefore = stakingtypes.Delegation{Shares: sdk.ZeroDec()}
				}
				suite.Require().Equal(delAfter.GetShares().Sub(delBefore.GetShares()), vaAfter.GetDelegatorShares().Sub(valBefore.GetDelegatorShares()))
				suite.Require().Equal(delAmount, vaAfter.GetTokens().Sub(valBefore.GetTokens()))

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == staking.DelegateEvent.ID.String() {
						suite.Require().Equal(log.Address, staking.GetAddress().String())
						suite.Require().Equal(log.Topics[1], delAddr.Hash().String())
						unpack, err := staking.DelegateEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						unpackValidator := unpack[0].(string)
						suite.Require().Equal(unpackValidator, val.GetOperator().String())
						amount := unpack[1].(*big.Int)
						suite.Require().Equal(amount.String(), delAmount.BigInt().String())
						shares := unpack[2].(*big.Int)
						suite.Require().Equal(shares.String(), delAfter.GetShares().Sub(delBefore.GetShares()).TruncateInt().BigInt().String())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type == stakingtypes.TypeMsgDelegate {
						for _, attr := range event.Attributes {
							if string(attr.Key) == stakingtypes.AttributeKeyValidator {
								suite.Require().Equal(string(attr.Value), val.GetOperator().String())
								existEvent = true
							}
							if string(attr.Key) == sdk.AttributeKeyAmount {
								suite.Require().Equal(string(attr.Value), delAmount.String())
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

func (suite *PrecompileTestSuite) TestDelegateCompare() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val := vals[0]
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()

	helpers.AddTestAddr(suite.app, suite.ctx, signer1.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

	// signer1 chain delegate to val
	shares1, err := suite.app.StakingKeeper.Delegate(suite.ctx, signer1.AccAddress(), delAmount, stakingtypes.Unbonded, val, true)
	suite.Require().NoError(err)

	// signer2 evm delegate to val
	shares2 := suite.PrecompileStakingDelegate(signer2, val.GetOperator(), delAmount.BigInt())

	// shares1 should equal shares2
	suite.Require().EqualValues(shares1.TruncateInt().BigInt(), shares2)

	// generate block
	suite.Commit()

	// signer1 chain withdraw
	rewards1, err := suite.app.DistrKeeper.WithdrawDelegationRewards(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)

	// signer2 evm withdraw
	rewards2 := suite.PrecompileStakingWithdraw(signer2, val.GetOperator())

	// rewards1 should equal rewards2
	suite.Require().EqualValues(rewards1.AmountOf(fxtypes.DefaultDenom).BigInt(), rewards2)
}
