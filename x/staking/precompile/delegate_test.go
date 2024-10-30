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
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func (suite *PrecompileTestSuite) TestDelegateCompare() {
	val := suite.GetFirstValidator()
	delAmount := helpers.NewRandAmount()
	signer1 := suite.NewSigner()
	signer2 := suite.NewSigner()

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
	suite.Require().EqualValues(rewards1.String(), rewards2.String())
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
		malleate func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error)
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error) {
				return contract.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "ok - delegate - multiple",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error) {
				suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

				validator, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
				suite.Require().NoError(err)
				_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, suite.signer.AccAddress(), delAmount, stakingtypes.Unbonded, validator, true)
				suite.Require().NoError(err)

				return contract.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error) {
				return contract.DelegateV2Args{
					Validator: val.String() + "1",
					Amount:    delAmount.BigInt(),
				}, fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error) {
				return contract.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - ok - delegate - multiple",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error) {
				suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

				suite.Require().NoError(suite.App.BankKeeper.SendCoins(suite.Ctx, suite.signer.AccAddress(), suite.stakingTestAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount))))

				suite.WithContract(suite.stakingTestAddr)
				res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				})
				suite.Require().False(res.Failed(), res.VmError)

				return contract.DelegateV2Args{
					Validator: val.String(),
					Amount:    delAmount.BigInt(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed - invalid validator address",
			malleate: func(val sdk.ValAddress, delAmount sdkmath.Int) (contract.DelegateV2Args, error) {
				return contract.DelegateV2Args{Validator: val.String() + "1", Amount: delAmount.BigInt()},
					fmt.Errorf("invalid validator address: %s", val.String()+"1")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator := suite.GetFirstValAddr()
			delAmount := helpers.NewRandAmount()

			args, expectErr := tc.malleate(operator, delAmount)

			suite.WithContract(suite.stakingAddr)
			delAddr := suite.signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				suite.WithContract(suite.stakingTestAddr)
				delAddr = suite.stakingTestAddr
				suite.MintToken(delAddr.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
			}

			delFound := true
			delBefore, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator)
			if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
				delFound = false
			} else {
				suite.Require().NoError(err)
			}

			valBefore := suite.GetValidator(operator)

			res := suite.WithError(expectErr).DelegateV2(suite.Ctx, args)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				delAfter := suite.GetDelegation(delAddr.Bytes(), operator)

				valAfter := suite.GetValidator(operator)

				if !delFound {
					delBefore = stakingtypes.Delegation{Shares: sdkmath.LegacyZeroDec()}
				}
				suite.Require().Equal(delAfter.GetShares().Sub(delBefore.GetShares()), valAfter.GetDelegatorShares().Sub(valBefore.GetDelegatorShares()))
				suite.Require().Equal(delAmount, valAfter.GetTokens().Sub(valBefore.GetTokens()))

				existLog := false
				for _, log := range res.Logs {
					abi := precompile.NewDelegateV2ABI()
					if log.Topics[0] == abi.Event.ID.String() {
						suite.Require().Equal(log.Address, suite.stakingAddr.String())

						event, err := abi.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.Delegator, delAddr)
						suite.Require().Equal(event.Validator, operator.String())
						suite.Require().Equal(event.Amount.String(), delAmount.BigInt().String())
						existLog = true
					}
				}
				suite.Require().True(existLog)
			}
		})
	}
}

func (suite *PrecompileTestSuite) CheckDelegateLogs(logs []*evmtypes.Log, delAddr common.Address, valAddr string, amount, shares *big.Int) {
	delegateV2Method := precompile.NewDelegateV2Method(nil)
	existLog := false
	for _, log := range logs {
		if log.Topics[0] == delegateV2Method.Event.ID.String() {
			suite.Require().Equal(log.Address, suite.stakingAddr.String())

			event, err := delegateV2Method.UnpackEvent(log.ToEthereum())
			suite.Require().NoError(err)
			suite.Require().Equal(event.Delegator, delAddr)
			suite.Require().Equal(event.Validator, valAddr)
			suite.Require().Equal(event.Amount.String(), amount.String())
			existLog = true
		}
	}
	suite.Require().True(existLog)
}

func (suite *PrecompileTestSuite) CheckDelegateEvents(ctx sdk.Context, valAddr sdk.ValAddress, delAmount sdkmath.Int) {
	existEvent := false
	for _, event := range ctx.EventManager().Events() {
		if event.Type == stakingtypes.EventTypeDelegate {
			for _, attr := range event.Attributes {
				if attr.Key == stakingtypes.AttributeKeyValidator {
					suite.Require().Equal(attr.Value, valAddr.String())
					existEvent = true
				}
				if attr.Key == sdk.AttributeKeyAmount {
					suite.Require().Equal(strings.TrimSuffix(attr.Value, fxtypes.DefaultDenom), delAmount.String())
					existEvent = true
				}
			}
		}
	}
	suite.Require().True(existEvent)
}
