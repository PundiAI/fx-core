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

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/types"
)

func (suite *PrecompileTestSuite) TestDelegate() {
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
				suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

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
				suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.stakingTestAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
				pack, err := suite.delegateV2Method.PackInput(types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				})
				suite.Require().NoError(err)

				res := suite.EthereumTx(suite.signer, suite.stakingTestAddr, big.NewInt(0), pack)
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

			stakingContract := suite.stakingAddr

			var packData []byte
			var err error
			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)
			args, value, errResult := tc.malleate(operator, delAmount)
			packData, err = suite.delegateV2Method.PackInput(args.(types.DelegateV2Args))
			suite.Require().NoError(err)

			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
				delAddr = suite.stakingTestAddr

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.stakingTestAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
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

	suite.MintToken(signer1.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

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

	require.Len(t, delegateV2Method.Method.Inputs, 2)
	require.Len(t, delegateV2Method.Method.Outputs, 1)

	require.Len(t, delegateV2Method.Event.Inputs, 3)
}

func (suite *PrecompileTestSuite) TestDelegateV2() {
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
				suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

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
				suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.stakingTestAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
				pack, err := suite.delegateV2Method.PackInput(types.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				})
				suite.Require().NoError(err)

				res := suite.EthereumTx(suite.signer, suite.stakingTestAddr, big.NewInt(0), pack)
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

			stakingContract := suite.stakingAddr
			delAddr := suite.signer.Address()

			suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
			packData, err := suite.delegateV2Method.PackInput(args)
			suite.Require().NoError(err)
			if strings.HasPrefix(tc.name, "contract") {
				stakingContract = suite.stakingTestAddr
				delAddr = suite.stakingTestAddr

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.stakingTestAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))
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
					if log.Topics[0] == suite.delegateV2Method.Event.ID.String() {
						suite.Require().Equal(log.Address, suite.stakingAddr.String())

						event, err := suite.delegateV2Method.UnpackEvent(log.ToEthereum())
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
