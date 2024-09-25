package precompile_test

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func (suite *PrecompileTestSuite) TestDelegate() {
	delegateV2Method := precompile.NewDelegateV2Method(nil)
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - v2",
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
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (interface{}, *big.Int, error) {
				helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

				validator, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
				suite.Require().NoError(err)
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, suite.signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
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
			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)
			args, value, errResult := tc.malleate(operator, delAmount)
			packData, err = delegateV2Method.PackInput(args.(types.DelegateV2Args))
			suite.Require().NoError(err)

			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.staking
				delAddr = suite.staking

				v2Args := args.(types.DelegateV2Args)
				packData, err = contract.MustABIJson(testscontract.StakingTestMetaData.ABI).Pack(TestDelegateV2Name, v2Args.Validator, v2Args.Amount)
				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.staking.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
				suite.Require().NoError(err)
			}

			delFound := true
			delBefore, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator)
			if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
				delFound = false
			} else {
				suite.Require().NoError(err)
			}

			valBefore, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, operator)
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, value, packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				delAfter := suite.GetDelegation(delAddr.Bytes(), operator)

				vaAfter, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, operator)
				suite.Require().NoError(err)

				if !delFound {
					delBefore = stakingtypes.Delegation{Shares: sdkmath.LegacyZeroDec()}
				}
				suite.Require().Equal(delAfter.GetShares().Sub(delBefore.GetShares()), vaAfter.GetDelegatorShares().Sub(valBefore.GetDelegatorShares()))
				suite.Require().Equal(delAmount, vaAfter.GetTokens().Sub(valBefore.GetTokens()))

				suite.CheckDelegateLogs(res.Logs, delAddr, val.GetOperator(),
					delAmount.BigInt(), delAfter.GetShares().Sub(delBefore.GetShares()).TruncateInt().BigInt())

				suite.CheckDelegateEvents(suite.Ctx, operator, delAmount)
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

	operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	suite.Require().NoError(err)
	// signer2 evm delegate to val
	shares2 := suite.PrecompileStakingDelegateV2(signer2, operator, delAmount.BigInt())

	// shares1 should equal shares2
	suite.Require().EqualValues(shares1.TruncateInt().BigInt(), shares2)

	// generate block
	suite.Commit()

	// signer1 chain withdraw
	rewards1, err := suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)

	// signer2 evm withdraw
	rewards2 := suite.PrecompileStakingWithdraw(signer2, operator)

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

				validator, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
				suite.Require().NoError(err)
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, suite.signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
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

			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)

			args, errResult := tc.malleate(operator, delAmount)

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

			delFound := true
			delBefore, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator)
			if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
				delFound = false
			} else {
				suite.Require().NoError(err)
			}
			valBefore, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, operator)
			suite.Require().NoError(err)

			res := suite.EthereumTx(suite.signer, stakingContract, big.NewInt(0), packData)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				delAfter := suite.GetDelegation(delAddr.Bytes(), operator)

				vaAfter, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, operator)
				suite.Require().NoError(err)

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
						suite.Require().Equal(event.Validator, val.GetOperator())
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
