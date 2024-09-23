package precompile_test

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

func TestStakingTransferSharesABI(t *testing.T) {
	transferSharesMethod := precompile.NewTransferSharesMethod(nil)

	require.Equal(t, 3, len(transferSharesMethod.Method.Inputs))
	require.Equal(t, 2, len(transferSharesMethod.Method.Outputs))

	require.Equal(t, 5, len(transferSharesMethod.Event.Inputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestTransferShares() {
	transferSharesMethod := precompile.NewTransferSharesMethod(nil)
	testCases := []struct {
		name        string
		pretransfer func(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int)
		malleate    func(val sdk.ValAddress, contract, to common.Address, shares *big.Int) ([]byte, *big.Int, []string)
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

		{
			name:        "contract - ok - from delegated",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - undelegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - undelegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - transfer all - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - delegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - undelegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - undelegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all - delegate to and from",
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
			fromSigner := suite.RandSigner()
			toSigner := suite.RandSigner()

			contract := precompile.GetAddress()
			delAddr := fromSigner.Address()
			if strings.HasPrefix(tc.name, "contract") {
				contract = suite.staking
				delAddr = suite.staking
			}

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

			fromBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, fromWithdrawAddr.Bytes())
			suite.Require().True(fromBalances.Empty())
			toBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, toWithdrawAddr.Bytes())
			suite.Require().True(toBalances.Empty())

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
			suite.NoError(err)

			pack, shares, _ := tc.malleate(operator, contract, toSigner.Address(), fromDelBefore.GetShares().TruncateInt().BigInt())
			res := suite.EthereumTx(fromSigner, contract, big.NewInt(0), pack)

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

				fromBalances = suite.App.BankKeeper.GetAllBalances(suite.Ctx, fromWithdrawAddr.Bytes())
				suite.Equal(fromBeforeRewards.Rewards.String(), fromBalances.String())

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == transferSharesMethod.Event.ID.String() {
						suite.Require().Equal(3, len(log.Topics))
						event, err := transferSharesMethod.UnpackEvent(log.ToEthereum())
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

func TestStakingTransferFromSharesABI(t *testing.T) {
	transferFromSharesMethod := precompile.NewTransferFromSharesMethod(nil)

	require.Equal(t, 4, len(transferFromSharesMethod.Method.Inputs))
	require.Equal(t, 2, len(transferFromSharesMethod.Method.Outputs))
}

//gocyclo:ignore
func (suite *PrecompileTestSuite) TestTransferFromShares() {
	transferFromSharesMethod := precompile.NewTransferFromSharesMethod(nil)
	testCases := []struct {
		name        string
		pretransfer func(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int)
		malleate    func(val sdk.ValAddress, spedner, from, to common.Address, shares *big.Int) ([]byte, *big.Int, []string)
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

		{
			name:        "contract - ok - from delegated",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - undelegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - undelegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - transfer all - delegate from and to",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from delegated - transfer all - delegate to and from",
			pretransfer: suite.delegateFromFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - delegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - delegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.delegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - undelegate from and to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - undelegate to and from",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromRand,
			suftransfer: suite.undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromAll,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all - undelegate to",
			pretransfer: suite.delegateFromToFunc,
			malleate:    suite.packTransferFromAll,
			suftransfer: suite.undelegateToFunc,
			result:      true,
		},
		{
			name:        "contract - ok - from and to delegated - transfer all - delegate from and to",
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
			fromSigner := suite.RandSigner()
			toSigner := suite.RandSigner()
			sender := suite.RandSigner()

			// from delegate, approve sender, sender send tx, transferFrom to toSigner
			// from delegate, approve contract, sender call contract, transferFrom to toSigner
			contract := precompile.GetAddress()
			delAddr := fromSigner.Address()
			spender := sender.Address()
			if strings.HasPrefix(tc.name, "contract") {
				contract = suite.staking
				spender = suite.staking
			}

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

			fromBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, fromWithdrawAddr.Bytes())
			suite.Require().True(fromBalances.Empty())
			toBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, toWithdrawAddr.Bytes())
			suite.Require().True(toBalances.Empty())

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
			suite.NoError(err)

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
			pack, shares, _ := tc.malleate(operator, spender, delAddr, toSigner.Address(), fromDelBefore.GetShares().TruncateInt().BigInt())

			res := suite.EthereumTx(sender, contract, big.NewInt(0), pack)

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

				fromBalances = suite.App.BankKeeper.GetAllBalances(suite.Ctx, fromWithdrawAddr.Bytes())
				suite.Require().Equal(fromBeforeRewards.Rewards.String(), fromBalances.String())

				toBalances = suite.App.BankKeeper.GetAllBalances(suite.Ctx, toWithdrawAddr.Bytes())
				if found2 {
					suite.Require().Equal(toBeforeRewards.Rewards.String(), toBalances.String())
				} else {
					suite.Require().True(toBalances.Empty())
				}

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == transferFromSharesMethod.Event.ID.String() {
						suite.Require().Equal(3, len(log.Topics))
						event, err := transferFromSharesMethod.UnpackEvent(log.ToEthereum())
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

func (suite *PrecompileTestSuite) TestTransferSharesCompare() {
	val := suite.GetFirstValidator()
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()
	signer3 := suite.RandSigner()

	helpers.AddTestAddr(suite.App, suite.Ctx, signer1.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

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

func (suite *PrecompileTestSuite) TestPrecompileStakingSteps() {
	val := suite.GetFirstValidator()
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()
	signer3 := suite.RandSigner()

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

func (suite *PrecompileTestSuite) TestTransferSharesRedelegate() {
	vals := suite.GetValidators()
	val := vals[0]
	valTmp := vals[1]
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()

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
	suite.Redelegate(operatorTmp, operator, signer1.AccAddress(), delegationTmp.Shares)
	suite.Commit()

	delegation, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, signer1.AccAddress(), operator)
	suite.Require().NoError(err)
	suite.Require().Equal(delegationTmp.Shares, delegation.Shares)

	// transfer shares
	transferSharesMethod := precompile.NewTransferSharesMethod(nil)
	pack, err := transferSharesMethod.PackInput(fxstakingtypes.TransferSharesArgs{
		Validator: val.GetOperator(),
		To:        signer2.Address(),
		Shares:    delegation.Shares.TruncateInt().BigInt(),
	})
	suite.Require().NoError(err)
	res := suite.EthereumTx(signer1, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Error(res, errors.New("from has receiving redelegation"))
}
