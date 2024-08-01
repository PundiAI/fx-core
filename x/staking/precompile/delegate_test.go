package precompile_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/contract"
	testscontract "github.com/functionx/fx-core/v7/tests/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/staking/precompile"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

func TestStakingDelegateABI(t *testing.T) {
	delegateMethod := precompile.NewDelegateMethod(nil)

	require.Equal(t, 1, len(delegateMethod.Method.Inputs))
	require.Equal(t, 2, len(delegateMethod.Method.Outputs))

	require.Equal(t, 4, len(delegateMethod.Event.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestDelegate() {
	delegateMethod := precompile.NewDelegateMethod(nil)
	delegateV2Method := precompile.NewDelegateV2Method(nil)
	testCases := []struct {
		name     string
		isV2     bool
		malleate func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: val.String()}, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "ok - delegate - multiple",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
				validator, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
				suite.Require().True(found)
				_, err := suite.App.StakingKeeper.Delegate(suite.Ctx, suite.signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				return types.DelegateArgs{Validator: val.String()}, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: val.String() + "1"},
					delAmount.BigInt(),
					fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},
		{
			name: "failed - invalid value",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: val.String()},
					big.NewInt(0),
					fmt.Errorf("invalid delegate amount: %s", big.NewInt(0).String())
			},
			result: false,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: sdk.ValAddress(suite.signer.AccAddress()).String()},
					delAmount.BigInt(),
					fmt.Errorf("validator not found: %s", sdk.ValAddress(suite.signer.AccAddress()).String())
			},
			error: func(args []string) string {
				return fmt.Sprintf("validator not found: %s", args[0])
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: val.String()}, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "contract - ok - delegate - multiple",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateName, val.String())
				suite.Require().NoError(err)

				res := suite.EthereumTx(suite.signer, suite.staking, delAmount.BigInt(), pack)
				suite.Require().False(res.Failed(), res.VmError)

				return types.DelegateArgs{Validator: val.String()}, delAmount.BigInt(), nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: val.String() + "1"}, delAmount.BigInt(),
					fmt.Errorf("delegate failed: invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},
		{
			name: "contract - failed - invalid value",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: val.String()}, big.NewInt(0),
					fmt.Errorf("delegate failed: invalid delegate amount: %s", big.NewInt(0).String())
			},
			result: false,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateArgs{Validator: sdk.ValAddress(suite.signer.AccAddress()).String()}, delAmount.BigInt(),
					fmt.Errorf("delegate failed: validator not found: %s", sdk.ValAddress(suite.signer.AccAddress()).String())
			},
			result: false,
		},

		{
			name: "ok - v2",
			isV2: true,
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "ok - v2 delegate - multiple",
			isV2: true,
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				validator, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
				suite.Require().True(found)
				_, err := suite.App.StakingKeeper.Delegate(suite.Ctx, suite.signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "failed - v2 invalid validator address",
			isV2: true,
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateV2Args{
					Validator: val.String() + "1",
					Amount:    delAmount.BigInt(),
				}, big.NewInt(0), fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},

		{
			name: "contract - ok v2",
			isV2: true,
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "contract - ok - v2 delegate - multiple",
			isV2: true,
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.staking.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateV2Name, val.String(), delAmount.BigInt())
				suite.Require().NoError(err)

				res := suite.EthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				suite.Require().False(res.Failed(), res.VmError)

				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, big.NewInt(0), nil
			},
			result: true,
		},
		{
			name: "contract - failed - v2 invalid validator address",
			isV2: true,
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				return types.DelegateV2Args{Validator: val.String() + "1", Amount: delAmount.BigInt()},
					big.NewInt(0),
					fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()

			delAmount := helpers.NewRandAmount()

			stakingContract := precompile.GetAddress()

			var packData []byte
			var err error

			args, value, errResult := tc.malleate(val.GetOperator(), delAmount)
			if !tc.isV2 {
				packData, err = delegateMethod.PackInput(args.(types.DelegateArgs))
			} else {
				packData, err = delegateV2Method.PackInput(args.(types.DelegateV2Args))
			}
			suite.Require().NoError(err)

			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				delAddr = suite.staking

				if !tc.isV2 {
					packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateName, args.(types.DelegateArgs).Validator)
				} else {
					v2Args := args.(types.DelegateV2Args)
					packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateV2Name, v2Args.Validator, v2Args.Amount)
					suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.staking.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
				}
				suite.Require().NoError(err)
			}

			delBefore, delFound := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val.GetOperator())
			valBefore, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val.GetOperator())
			suite.Require().True(found)

			res := suite.EthereumTx(suite.signer, stakingContract, value, packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				delAfter := suite.GetDelegation(delAddr.Bytes(), val.GetOperator())

				vaAfter, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val.GetOperator())
				suite.Require().True(found)

				if !delFound {
					delBefore = stakingtypes.Delegation{Shares: sdkmath.LegacyZeroDec()}
				}
				suite.Require().Equal(delAfter.GetShares().Sub(delBefore.GetShares()), vaAfter.GetDelegatorShares().Sub(valBefore.GetDelegatorShares()))
				suite.Require().Equal(delAmount, vaAfter.GetTokens().Sub(valBefore.GetTokens()))

				suite.CheckDelegateLogs(res.Logs, delAddr, val.GetOperator().String(),
					delAmount.BigInt(), delAfter.GetShares().Sub(delBefore.GetShares()).TruncateInt().BigInt())

				suite.CheckDelegateEvents(suite.Ctx, val.GetOperator(), delAmount)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestDelegateCompare() {
	val := suite.GetFirstValidator()
	delAmount := helpers.NewRandAmount()
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()

	helpers.AddTestAddr(suite.App, suite.Ctx, signer1.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

	// signer1 chain delegate to val
	shares1, err := suite.App.StakingKeeper.Delegate(suite.Ctx, signer1.AccAddress(), delAmount, stakingtypes.Unbonded, val, true)
	suite.Require().NoError(err)

	// signer2 evm delegate to val
	shares2 := suite.PrecompileStakingDelegate(signer2, val.GetOperator(), delAmount.BigInt())

	// shares1 should equal shares2
	suite.Require().EqualValues(shares1.TruncateInt().BigInt(), shares2)

	// generate block
	suite.Commit()

	// signer1 chain withdraw
	rewards1, err := suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)

	// signer2 evm withdraw
	rewards2 := suite.PrecompileStakingWithdraw(signer2, val.GetOperator())

	// rewards1 should equal rewards2
	suite.Require().EqualValues(rewards1.AmountOf(fxtypes.DefaultDenom).BigInt(), rewards2)
}

func TestStakingDelegateV2ABI(t *testing.T) {
	delegateV2Method := precompile.NewDelegateV2Method(nil)

	require.Equal(t, 2, len(delegateV2Method.Method.Inputs))
	require.Equal(t, 1, len(delegateV2Method.Method.Outputs))

	require.Equal(t, 3, len(delegateV2Method.Event.Inputs))
}

func (suite *PrecompileTestSuite) TestDelegateV2() {
	delegateV2Method := precompile.NewDelegateV2Method(nil)
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error) {
				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - delegate - multiple",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				validator, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
				suite.Require().True(found)
				_, err := suite.App.StakingKeeper.Delegate(suite.Ctx, suite.signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error) {
				return types.DelegateV2Args{
					Validator: val.String() + "1",
					Amount:    delAmount.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},

		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error) {
				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - delegate - multiple",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.staking.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
				pack, err := contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateV2Name, val.String(), delAmount.BigInt())
				suite.Require().NoError(err)

				res := suite.EthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				suite.Require().False(res.Failed(), res.VmError)

				return types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (types.DelegateV2Args, error) {
				return types.DelegateV2Args{Validator: val.String() + "1", Amount: delAmount.BigInt()},
					fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()
			delAmount := helpers.NewRandAmount()

			args, errResult := tc.malleate(val.GetOperator(), delAmount)

			stakingContract := precompile.GetAddress()
			delAddr := suite.signer.Address()

			helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
			packData, err := delegateV2Method.PackInput(args)
			suite.Require().NoError(err)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				delAddr = suite.staking

				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateV2Name, args.Validator, args.Amount)
				suite.Require().NoError(err)

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.staking.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
			}

			delBefore, delFound := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), val.GetOperator())
			valBefore, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val.GetOperator())
			suite.Require().True(found)

			res := suite.EthereumTx(suite.signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				delAfter := suite.GetDelegation(delAddr.Bytes(), val.GetOperator())

				vaAfter, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, val.GetOperator())
				suite.Require().True(found)

				if !delFound {
					delBefore = stakingtypes.Delegation{Shares: sdkmath.LegacyZeroDec()}
				}
				suite.Require().Equal(delAfter.GetShares().Sub(delBefore.GetShares()), vaAfter.GetDelegatorShares().Sub(valBefore.GetDelegatorShares()))
				suite.Require().Equal(delAmount, vaAfter.GetTokens().Sub(valBefore.GetTokens()))

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == delegateV2Method.Event.ID.String() {
						suite.Require().Equal(log.Address, precompile.GetAddress().String())

						event, err := delegateV2Method.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.Delegator, delAddr)
						suite.Require().Equal(event.Validator, val.GetOperator().String())
						suite.Require().Equal(event.Amount.String(), delAmount.BigInt().String())
						existLog = true
					}
				}
				suite.Require().True(existLog)
			} else {
				suite.Error(res, errResult)
			}
		})
	}
}
