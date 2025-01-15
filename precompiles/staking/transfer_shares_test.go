package staking_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/staking"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func TestStakingTransferSharesABI(t *testing.T) {
	transferSharesABI := staking.NewTransferSharesABI()

	require.Len(t, transferSharesABI.Method.Inputs, 3)
	require.Len(t, transferSharesABI.Method.Outputs, 2)

	require.Len(t, transferSharesABI.Event.Inputs, 5)
}

func (suite *StakingPrecompileTestSuite) TestTransferShares() {
	testCases := []struct {
		name        string
		pretransfer func(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int)
		malleate    func(val sdk.ValAddress, to common.Address, shares *big.Int) (fxcontract.TransferSharesArgs, *big.Int, []string)
		suftransfer func(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int)
		error       func(errArgs []string) string
		result      bool
	}{
		{
			name:        "ok - from delegated",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			result:      true,
		},
		{
			name:        "ok - from delegated - undelegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - undelegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - delegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - undelegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - undelegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - delegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()
			delAmt := sdkmath.NewInt(int64(tmrand.Intn(10000) + 1000)).Mul(sdkmath.NewInt(1e18))
			suite.signer = suite.NewSigner()
			toSigner := suite.NewSigner()

			// contract := suite.stakingAddr
			delAddr := suite.GetDelAddr()
			// if strings.HasPrefix(tc.name, "contract") {
			// 	contract = suite.stakingTestAddr
			// 	delAddr = suite.stakingTestAddr
			// }

			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)

			tc.pretransfer(operator, delAddr, toSigner.Address(), delAmt)

			fromWithdrawAddr := helpers.GenHexAddress()
			err = suite.App.DistrKeeper.SetWithdrawAddr(suite.Ctx, delAddr.Bytes(), fromWithdrawAddr.Bytes())
			suite.Require().NoError(err)
			toWithdrawAddr := helpers.GenHexAddress()
			err = suite.App.DistrKeeper.SetWithdrawAddr(suite.Ctx, toSigner.AccAddress(), toWithdrawAddr.Bytes())
			suite.Require().NoError(err)

			suite.Commit()

			fromBalance := suite.App.BankKeeper.GetBalance(suite.Ctx, fromWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
			suite.Require().True(fromBalance.IsZero())
			toBalance := suite.App.BankKeeper.GetBalance(suite.Ctx, toWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
			suite.Require().True(toBalance.IsZero())

			fromDelBefore := suite.GetDelegation(delAddr.Bytes(), operator)
			found2 := true
			toDelBefore, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, toSigner.Address().Bytes(), operator)
			if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
				toDelBefore.Shares = sdkmath.LegacyZeroDec()
				found2 = false
			} else {
				suite.Require().NoError(err)
			}

			fromBeforeRewards, err := distributionkeeper.NewQuerier(suite.App.DistrKeeper).DelegationRewards(suite.Ctx, &distributiontypes.QueryDelegationRewardsRequest{
				DelegatorAddress: sdk.AccAddress(delAddr.Bytes()).String(),
				ValidatorAddress: sdk.ValAddress(operator).String(),
			})
			suite.Require().NoError(err)

			args, shares, _ := tc.malleate(operator, toSigner.Address(), fromDelBefore.GetShares().TruncateInt().BigInt())
			suite.WithSigner(suite.signer)
			res, _ := suite.TransferShares(suite.Ctx, args)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				fromDelAfter, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator)
				if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
					fromDelAfter.Shares = sdkmath.LegacyZeroDec()
				} else {
					suite.Require().NoError(err)
				}
				toDelAfter, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, toSigner.Address().Bytes(), operator)
				suite.Require().NoError(err)
				if !found2 {
					toDelBefore.Shares = sdkmath.LegacyZeroDec()
				}

				suite.Require().Equal(fromDelBefore.GetShares().TruncateInt().Sub(fromDelAfter.GetShares().TruncateInt()).BigInt(), shares)
				suite.Require().Equal(toDelAfter.GetShares().TruncateInt().Sub(toDelBefore.GetShares().TruncateInt()).BigInt(), shares)

				if tc.suftransfer != nil {
					tc.suftransfer(operator, delAddr, toSigner.Address(), delAmt)
				}

				fromBalance = suite.App.BankKeeper.GetBalance(suite.Ctx, fromWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
				suite.Require().Equal(fromBeforeRewards.Rewards.String(), sdk.NewDecCoinFromCoin(fromBalance).String())

				existLog := false
				for _, log := range res.Logs {
					transferSharesABI := staking.NewTransferSharesABI()
					if log.Topics[0] == transferSharesABI.Event.ID.String() {
						suite.Require().Len(log.Topics, 3)
						event, err := transferSharesABI.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.From, delAddr)
						suite.Require().Equal(event.To, toSigner.Address())
						suite.Require().Equal(event.Validator, sdk.ValAddress(operator).String())
						suite.Require().Equal(event.Shares.String(), shares.String())
						existLog = true
					}
				}
				suite.Require().True(existLog)
			} else {
				suite.Require().True(err != nil || res.Failed())
			}
		})
	}
}

func (suite *StakingPrecompileTestSuite) packTransferRand(val sdk.ValAddress, to common.Address, shares *big.Int) (fxcontract.TransferSharesArgs, *big.Int, []string) {
	randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
	args := fxcontract.TransferSharesArgs{
		Validator: val.String(),
		To:        to,
		Shares:    randShares,
	}
	return args, randShares, nil
}

func (suite *StakingPrecompileTestSuite) packTransferAll(val sdk.ValAddress, to common.Address, shares *big.Int) (fxcontract.TransferSharesArgs, *big.Int, []string) {
	args := fxcontract.TransferSharesArgs{
		Validator: val.String(),
		To:        to,
		Shares:    shares,
	}
	return args, shares, nil
}

func TestStakingTransferFromSharesABI(t *testing.T) {
	transferFromSharesABI := staking.NewTransferFromSharesABI()

	require.Len(t, transferFromSharesABI.Method.Inputs, 4)
	require.Len(t, transferFromSharesABI.Method.Outputs, 2)
}

func (suite *StakingPrecompileTestSuite) TestTransferFromShares() {
	testCases := []struct {
		name        string
		pretransfer func(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int)
		malleate    func(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) (fxcontract.TransferFromSharesArgs, *big.Int, []string)
		suftransfer func(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int)
		error       func(errArgs []string) string
		result      bool
	}{
		{
			name:        "ok - from delegated",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			result:      true,
		},
		{
			name:        "ok - from delegated - undelegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - undelegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - delegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - undelegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - undelegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromAll,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - delegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			val := suite.GetFirstValidator()
			delAmt := sdkmath.NewInt(int64(tmrand.Intn(10000) + 1000)).Mul(sdkmath.NewInt(1e18))
			fromSigner := suite.NewSigner()
			toSigner := suite.NewSigner()
			suite.signer = suite.NewSigner()

			// from delegate, approve sender, sender send tx, transferFrom to toSigner
			// from delegate, approve contract, sender call contract, transferFrom to toSigner
			delAddr := fromSigner.Address()
			spender := suite.GetDelAddr()

			operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
			suite.Require().NoError(err)

			tc.pretransfer(operator, delAddr, toSigner.Address(), delAmt)

			fromWithdrawAddr := helpers.GenHexAddress()
			err = suite.App.DistrKeeper.SetWithdrawAddr(suite.Ctx, delAddr.Bytes(), fromWithdrawAddr.Bytes())
			suite.Require().NoError(err)
			toWithdrawAddr := helpers.GenHexAddress()
			err = suite.App.DistrKeeper.SetWithdrawAddr(suite.Ctx, toSigner.AccAddress(), toWithdrawAddr.Bytes())
			suite.Require().NoError(err)

			suite.Commit()

			fromBalance := suite.App.BankKeeper.GetBalance(suite.Ctx, fromWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
			suite.Require().True(fromBalance.IsZero())
			toBalance := suite.App.BankKeeper.GetBalance(suite.Ctx, toWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
			suite.Require().True(toBalance.IsZero())

			fromDelBefore := suite.GetDelegation(delAddr.Bytes(), operator)
			found2 := true
			toDelBefore, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, toSigner.Address().Bytes(), operator)
			if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
				toDelBefore.Shares = sdkmath.LegacyZeroDec()
				found2 = false
			} else {
				suite.Require().NoError(err)
			}

			fromBeforeRewards, err := distributionkeeper.NewQuerier(suite.App.DistrKeeper).DelegationRewards(suite.Ctx, &distributiontypes.QueryDelegationRewardsRequest{
				DelegatorAddress: sdk.AccAddress(delAddr.Bytes()).String(),
				ValidatorAddress: sdk.ValAddress(operator).String(),
			})
			suite.Require().NoError(err)

			toBeforeRewards, err := distributionkeeper.NewQuerier(suite.App.DistrKeeper).DelegationRewards(suite.Ctx, &distributiontypes.QueryDelegationRewardsRequest{
				DelegatorAddress: sdk.AccAddress(toSigner.Address().Bytes()).String(),
				ValidatorAddress: sdk.ValAddress(operator).String(),
			})
			if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
				toBeforeRewards = &distributiontypes.QueryDelegationRewardsResponse{Rewards: sdk.DecCoins{}}
			} else {
				suite.Require().NoError(err)
			}

			// NOTE: if contract test, spender is staking test contract
			args, shares, _ := tc.malleate(operator, spender, delAddr, toSigner.Address(), fromDelBefore.GetShares().TruncateInt().BigInt())

			suite.WithSigner(suite.signer)
			res, _ := suite.TransferFromShares(suite.Ctx, args)

			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				fromDelAfter, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, delAddr.Bytes(), operator)
				if err != nil && errors.Is(err, stakingtypes.ErrNoDelegation) {
					fromDelAfter.Shares = sdkmath.LegacyZeroDec()
				}
				toDelAfter, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, toSigner.Address().Bytes(), operator)
				suite.Require().NoError(err)

				suite.Require().Equal(fromDelBefore.GetShares().TruncateInt().Sub(fromDelAfter.GetShares().TruncateInt()).BigInt(), shares)
				suite.Require().Equal(toDelAfter.GetShares().TruncateInt().Sub(toDelBefore.GetShares().TruncateInt()).BigInt(), shares)

				allowance := suite.App.StakingKeeper.GetAllowance(suite.Ctx, operator, delAddr.Bytes(), spender.Bytes())
				suite.Require().EqualValues(big.NewInt(0), allowance)

				if tc.suftransfer != nil {
					tc.suftransfer(operator, delAddr, toSigner.Address(), delAmt)
				}

				fromBalance = suite.App.BankKeeper.GetBalance(suite.Ctx, fromWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
				suite.Require().Equal(fromBeforeRewards.Rewards.String(), sdk.NewDecCoinFromCoin(fromBalance).String())

				toBalance = suite.App.BankKeeper.GetBalance(suite.Ctx, toWithdrawAddr.Bytes(), fxtypes.DefaultDenom)
				if found2 {
					suite.Require().Equal(toBeforeRewards.Rewards.String(), sdk.NewDecCoinFromCoin(toBalance).String())
				} else {
					suite.Require().True(toBalance.IsZero())
				}

				existLog := false
				for _, log := range res.Logs {
					transferFromSharesABI := staking.NewTransferFromSharesABI()
					if log.Topics[0] == transferFromSharesABI.Event.ID.String() {
						suite.Require().Len(log.Topics, 3)
						event, err := transferFromSharesABI.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.From, delAddr)
						suite.Require().Equal(event.To, toSigner.Address())
						suite.Require().Equal(event.Validator, sdk.ValAddress(operator).String())
						suite.Require().Equal(event.Shares.String(), shares.String())
						existLog = true
					}
				}
				suite.Require().True(existLog)
			} else {
				suite.Require().True(err != nil || res.Failed())
			}
		})
	}
}

func (suite *StakingPrecompileTestSuite) packTransferFromRand(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) (fxcontract.TransferFromSharesArgs, *big.Int, []string) {
	randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
	suite.approveFunc(val, from, spender, randShares)
	return fxcontract.TransferFromSharesArgs{
		Validator: val.String(),
		From:      from,
		To:        to,
		Shares:    randShares,
	}, randShares, nil
}

func (suite *StakingPrecompileTestSuite) packTransferFromAll(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) (fxcontract.TransferFromSharesArgs, *big.Int, []string) {
	suite.approveFunc(val, from, spender, shares)
	return fxcontract.TransferFromSharesArgs{
		Validator: val.String(),
		From:      from,
		To:        to,
		Shares:    shares,
	}, shares, nil
}

func (suite *StakingPrecompileTestSuite) TestTransferSharesCompare() {
	if !suite.IsCallPrecompile() {
		suite.T().Skip()
	}
	val := suite.GetFirstValidator()
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).MulRaw(1e16)
	signer1 := suite.NewSigner()
	signer2 := suite.NewSigner()
	signer3 := suite.NewSigner()

	suite.MintToken(signer1.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))

	operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	suite.Require().NoError(err)

	// starting info 1,2,3
	startingInfo, err := suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod)
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer2.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod)
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer3.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod)

	// signer1 chain delegate to val
	shares1, err := suite.App.StakingKeeper.Delegate(suite.Ctx, signer1.AccAddress(), delAmount, stakingtypes.Unbonded, val, true)
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(2, startingInfo.PreviousPeriod)

	// signer2 evm delegate to val
	shares2 := suite.PrecompileStakingDelegateV2(signer2, operator, delAmount.BigInt())
	suite.Require().Equal(shares1.TruncateInt().BigInt(), shares2)

	// signer2 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer2.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(3, startingInfo.PreviousPeriod)

	// generate block
	suite.Commit()

	// signer1 withdraw
	rewards1, err := suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(4, startingInfo.PreviousPeriod)

	// signer2 transfer shares to singer3
	halfShares := big.NewInt(0).Div(shares2, big.NewInt(2))
	surplusShares, rewards2 := suite.PrecompileStakingTransferShares(signer2, operator, signer3.Address(), halfShares)

	// starting info 2,3
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer2.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(5, startingInfo.PreviousPeriod)
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer3.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(6, startingInfo.PreviousPeriod)

	// rewards1 equal rewards2
	suite.Require().EqualValues(rewards1.AmountOf(fxtypes.DefaultDenom).BigInt(), rewards2)
	// surplus shares equal to half shares
	suite.Require().EqualValues(halfShares, surplusShares)

	// signer1 undelegate half shares
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, signer1.AccAddress(), operator, sdkmath.LegacyNewDecFromBigInt(halfShares))
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(7, startingInfo.PreviousPeriod)

	// signer1 shares equal to half shares
	delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)
	suite.Require().EqualValues(halfShares, delegation.GetShares().TruncateInt().BigInt())

	// generate block
	suite.Commit()

	// transfer shares surplus
	surplusShares, rewards2 = suite.PrecompileStakingTransferShares(signer2, operator, signer3.Address(), surplusShares)

	// starting info 2,3
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer2.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transfer all shares, starting info removed
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer3.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(9, startingInfo.PreviousPeriod)

	// surplus shares equal to zero
	suite.Require().EqualValues(big.NewInt(0).String(), surplusShares.String())

	// get signer3 fx balance
	balances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, signer3.AccAddress())
	suite.Require().EqualValues(rewards2, balances.AmountOf(fxtypes.DefaultDenom).BigInt())

	// signer1 withdraw
	rewards1, err = suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(10, startingInfo.PreviousPeriod)

	// rewards1 equal rewards2
	suite.Require().EqualValues(rewards1.AmountOf(fxtypes.DefaultDenom).BigInt(), rewards2)

	// signer1 undelegate all shares
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, signer1.AccAddress(), operator, delegation.GetShares())
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // undelegate all, starting info removed

	// signer1 shares equal to zero
	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	// signer3 delegation
	shares3, amount3 := suite.PrecompileStakingDelegation(operator, signer3.Address())
	suite.Require().EqualValues(big.NewInt(0).Mul(halfShares, big.NewInt(2)), shares3)
	suite.Require().EqualValues(delAmount.BigInt(), amount3)

	// generate block
	suite.Commit()

	// signer3 evm withdraw
	_ = suite.PrecompileStakingWithdraw(signer3, operator)

	// signer3 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer3.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(12, startingInfo.PreviousPeriod)

	// signer3 transferFrom shares
	suite.PrecompileStakingApproveShares(signer3, operator, signer2.Address(), shares3)
	suite.PrecompileStakingTransferFromShares(signer2, operator, signer3.Address(), signer1.Address(), shares3)
	newShares1, _ := suite.PrecompileStakingDelegation(operator, signer1.Address())
	suite.Require().EqualValues(newShares1, shares3)

	// signer3 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer3.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transferFrom all shares, starting info removed
	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(14, startingInfo.PreviousPeriod)

	// generate block
	suite.Commit()

	// singer1 evm undelegate
	_ = suite.PrecompileStakingUndelegateV2(signer1, operator, halfShares)
	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(uint64(15), startingInfo.PreviousPeriod) // withdraw +1, undelegate +1

	_ = suite.PrecompileStakingUndelegateV2(signer1, operator, halfShares)
	// signer1 starting info
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // undelegate all, starting info removed

	// signer1 shares equal to zero
	_, err = suite.App.StakingKeeper.GetDelegation(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().ErrorIs(err, stakingtypes.ErrNoDelegation)

	// starting info 1,2,3
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer1.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // undelegate all, starting info removed
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer2.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transfer all shares, starting info removed
	startingInfo, err = suite.App.DistrKeeper.GetDelegatorStartingInfo(suite.Ctx, operator, signer3.AccAddress())
	suite.Require().NoError(err)
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transfer all shares, starting info removed
}

func (suite *StakingPrecompileTestSuite) TestPrecompileStakingSteps() {
	if !suite.IsCallPrecompile() {
		suite.T().Skip()
	}
	val := suite.GetFirstValidator()
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).MulRaw(1e16)
	signer1 := suite.NewSigner()
	signer2 := suite.NewSigner()
	signer3 := suite.NewSigner()

	operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	suite.Require().NoError(err)

	// delegate 1 and 2
	suite.Delegate(operator, delAmount, signer1.AccAddress(), signer2.AccAddress())
	suite.Commit()

	// get 1 shares
	shares1, _ := suite.PrecompileStakingDelegation(operator, signer1.Address())

	// 1 transfer shares to 2
	suite.PrecompileStakingTransferShares(signer1, operator, signer2.Address(), shares1)
	suite.Commit()

	// 2 transfer shares to 3
	suite.PrecompileStakingTransferShares(signer2, operator, signer3.Address(), shares1)
	suite.Commit()

	// delegate 1,2,3
	suite.Delegate(operator, delAmount, signer1.AccAddress(), signer2.AccAddress(), signer3.AccAddress())
	suite.Commit()

	// precompile delegate 1,2,3
	suite.PrecompileStakingDelegateV2(signer1, operator, delAmount.BigInt())
	suite.PrecompileStakingDelegateV2(signer2, operator, delAmount.BigInt())
	suite.PrecompileStakingDelegateV2(signer3, operator, delAmount.BigInt())
	suite.Commit()

	// precompile withdraw 1,2,3
	suite.PrecompileStakingWithdraw(signer1, operator)
	suite.PrecompileStakingWithdraw(signer2, operator)
	suite.PrecompileStakingWithdraw(signer3, operator)
	suite.Commit()

	// withdraw 1,2,3
	_, err = suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)
	_, err = suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer2.AccAddress(), operator)
	suite.Require().NoError(err)
	_, err = suite.App.DistrKeeper.WithdrawDelegationRewards(suite.Ctx, signer3.AccAddress(), operator)
	suite.Require().NoError(err)
	suite.Commit()

	// precompile undelegate 1,2,3
	suite.PrecompileStakingUndelegateV2(signer1, operator, shares1)
	suite.PrecompileStakingUndelegateV2(signer2, operator, shares1)
	suite.PrecompileStakingUndelegateV2(signer3, operator, shares1)
	suite.Commit()

	// undelegate 1,2,3
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, signer1.AccAddress(), operator, sdkmath.LegacyNewDecFromBigInt(shares1))
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, signer2.AccAddress(), operator, sdkmath.LegacyNewDecFromBigInt(shares1))
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, signer3.AccAddress(), operator, sdkmath.LegacyNewDecFromBigInt(shares1))
	suite.Require().NoError(err)
	suite.Commit()

	// delegate 1,2,3
	suite.Delegate(operator, delAmount, signer1.AccAddress(), signer2.AccAddress(), signer3.AccAddress())
	suite.Commit()

	// transfer shares 1,2,3
	suite.PrecompileStakingTransferShares(signer1, operator, signer2.Address(), shares1)
	suite.PrecompileStakingTransferShares(signer2, operator, signer3.Address(), shares1)
	suite.PrecompileStakingTransferShares(signer3, operator, signer1.Address(), shares1)
	suite.Commit()

	// transferFrom shares 1,2,3
	suite.PrecompileStakingApproveShares(signer1, operator, signer2.Address(), shares1)
	suite.PrecompileStakingTransferFromShares(signer2, operator, signer1.Address(), signer2.Address(), shares1)
	suite.PrecompileStakingApproveShares(signer2, operator, signer3.Address(), shares1)
	suite.PrecompileStakingTransferFromShares(signer3, operator, signer2.Address(), signer3.Address(), shares1)
	suite.PrecompileStakingApproveShares(signer3, operator, signer1.Address(), shares1)
	suite.PrecompileStakingTransferFromShares(signer1, operator, signer3.Address(), signer1.Address(), shares1)
	suite.Commit()
}

func (suite *StakingPrecompileTestSuite) TestTransferSharesRedelegate() {
	vals := suite.GetValidators()
	val := vals[0]
	valTmp := vals[1]
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.NewSigner()
	signer2 := suite.NewSigner()

	operator, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	suite.Require().NoError(err)

	operatorTmp, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(valTmp.GetOperator())
	suite.Require().NoError(err)

	// delegate 1 and 2
	suite.Delegate(operatorTmp, delAmount, signer1.AccAddress())
	suite.Commit()

	delegationTmp, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, signer1.AccAddress(), operatorTmp)
	suite.Require().NoError(err)

	// redelegate
	_, err = suite.App.StakingKeeper.BeginRedelegation(suite.Ctx, signer1.AccAddress(), operatorTmp, operator, delegationTmp.Shares)
	suite.Require().NoError(err)
	suite.Commit()

	delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)
	suite.Require().Equal(delegationTmp.Shares, delegation.Shares)

	// transfer shares
	suite.WithSigner(signer1).WithContract(suite.stakingAddr)
	_, _ = suite.WithError(errors.New("from has receiving redelegation")).
		TransferShares(suite.Ctx, fxcontract.TransferSharesArgs{
			Validator: val.GetOperator(),
			To:        signer2.Address(),
			Shares:    delegation.Shares.TruncateInt().BigInt(),
		})
}

func (suite *StakingPrecompileTestSuite) delegateFromFunc(val sdk.ValAddress, from, _ common.Address, delAmount sdkmath.Int) {
	suite.MintToken(from.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
	_, err := stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *StakingPrecompileTestSuite) undelegateToFunc(val sdk.ValAddress, _, to common.Address, _ sdkmath.Int) {
	toDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)
}

func (suite *StakingPrecompileTestSuite) delegateFromToFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
	suite.MintToken(from.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
	_, err := stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)

	suite.MintToken(to.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
	_, err = stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *StakingPrecompileTestSuite) delegateToFromFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
	suite.MintToken(to.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
	_, err := stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)

	suite.MintToken(from.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
	_, err = stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *StakingPrecompileTestSuite) undelegateFromToFunc(val sdk.ValAddress, from, to common.Address, _ sdkmath.Int) {
	fromDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, from.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, from.Bytes(), val, fromDel.Shares)
	suite.Require().NoError(err)

	toDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)
}

func (suite *StakingPrecompileTestSuite) undelegateToFromFunc(val sdk.ValAddress, from, to common.Address, _ sdkmath.Int) {
	toDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)

	fromDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, from.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, from.Bytes(), val, fromDel.Shares)
	suite.Require().NoError(err)
}

func (suite *StakingPrecompileTestSuite) approveFunc(val sdk.ValAddress, owner, spender common.Address, allowance *big.Int) {
	suite.App.StakingKeeper.SetAllowance(suite.Ctx, val, owner.Bytes(), spender.Bytes(), allowance)
}

func (suite *StakingPrecompileTestSuite) PrecompileStakingDelegation(val sdk.ValAddress, del common.Address) (*big.Int, *big.Int) {
	return suite.Delegation(suite.Ctx, fxcontract.DelegationArgs{
		Validator: val.String(),
		Delegator: del,
	})
}

func (suite *StakingPrecompileTestSuite) PrecompileStakingTransferShares(signer *helpers.Signer, val sdk.ValAddress, receipt common.Address, shares *big.Int) (*big.Int, *big.Int) {
	balanceBefore := suite.GetStakingBalance(signer.AccAddress())
	suite.WithSigner(signer)
	res, _ := suite.TransferShares(suite.Ctx, fxcontract.TransferSharesArgs{
		Validator: val.String(),
		To:        receipt,
		Shares:    shares,
	})
	suite.Require().False(res.Failed(), res.VmError)

	signerShares, _ := suite.PrecompileStakingDelegation(val, signer.Address())

	balanceAfter := suite.GetStakingBalance(signer.AccAddress())
	rewards := balanceAfter.Sub(balanceBefore)
	return signerShares, rewards.BigInt()
}

func (suite *StakingPrecompileTestSuite) PrecompileStakingUndelegateV2(signer *helpers.Signer, val sdk.ValAddress, shares *big.Int) *big.Int {
	balanceBefore := suite.GetStakingBalance(signer.AccAddress())
	suite.WithSigner(signer)
	res := suite.UndelegateV2(suite.Ctx, fxcontract.UndelegateV2Args{
		Validator: val.String(),
		Amount:    shares,
	})
	suite.Require().False(res.Failed(), res.VmError)

	balanceAfter := suite.GetStakingBalance(signer.AccAddress())
	rewards := balanceAfter.Sub(balanceBefore)
	return rewards.BigInt()
}

func (suite *StakingPrecompileTestSuite) PrecompileStakingApproveShares(signer *helpers.Signer, val sdk.ValAddress, spender common.Address, shares *big.Int) {
	suite.WithSigner(signer)
	suite.ApproveShares(suite.Ctx, fxcontract.ApproveSharesArgs{
		Validator: val.String(),
		Spender:   spender,
		Shares:    shares,
	})
}

func (suite *StakingPrecompileTestSuite) PrecompileStakingTransferFromShares(signer *helpers.Signer, val sdk.ValAddress, from, receipt common.Address, shares *big.Int) {
	suite.WithSigner(signer)
	suite.TransferFromShares(suite.Ctx, fxcontract.TransferFromSharesArgs{
		Validator: val.String(),
		From:      from,
		To:        receipt,
		Shares:    shares,
	})
}

func (suite *StakingPrecompileTestSuite) Delegate(val sdk.ValAddress, amount sdkmath.Int, dels ...sdk.AccAddress) {
	for _, del := range dels {
		suite.MintToken(del, sdk.NewCoin(fxtypes.DefaultDenom, amount))
		validator, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
		suite.Require().NoError(err)
		_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, del, amount, stakingtypes.Unbonded, validator, true)
		suite.Require().NoError(err)
	}
}
