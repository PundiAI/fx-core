package staking_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

func TestStakingTransferSharesABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.TransferSharesMethodName]
	require.Equal(t, method, staking.TransferSharesMethod)
	require.Equal(t, 3, len(staking.TransferSharesMethod.Inputs))
	require.Equal(t, 2, len(staking.TransferSharesMethod.Outputs))

	event := stakingABI.Events[staking.TransferSharesEventName]
	require.Equal(t, event, staking.TransferSharesEvent)
	require.Equal(t, 5, len(staking.TransferSharesEvent.Inputs))
}

func (suite *PrecompileTestSuite) TestTransferShares() {
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
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val := vals[0]
			delAmt := sdkmath.NewInt(int64(tmrand.Intn(10000) + 1000)).Mul(sdkmath.NewInt(1e18))
			fromSigner := suite.RandSigner()
			toSigner := suite.RandSigner()

			contract := staking.GetAddress()
			delAddr := fromSigner.Address()
			if strings.HasPrefix(tc.name, "contract") {
				contract = suite.staking
				delAddr = suite.staking
			}

			tc.pretransfer(val.GetOperator(), delAddr, toSigner.Address(), delAmt)

			fromWithdrawAddr := helpers.GenerateAddress()
			err := suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, delAddr.Bytes(), fromWithdrawAddr.Bytes())
			suite.Require().NoError(err)
			toWithdrawAddr := helpers.GenerateAddress()
			err = suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, toSigner.AccAddress(), toWithdrawAddr.Bytes())
			suite.Require().NoError(err)

			suite.Commit()

			fromBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, fromWithdrawAddr.Bytes())
			suite.Require().True(fromBalances.Empty())
			toBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, toWithdrawAddr.Bytes())
			suite.Require().True(toBalances.Empty())

			fromDelBefore, found1 := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
			suite.Require().True(found1)
			toDelBefore, found2 := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.Address().Bytes(), val.GetOperator())

			pack, shares, _ := tc.malleate(val.GetOperator(), contract, toSigner.Address(), fromDelBefore.GetShares().TruncateInt().BigInt())
			tx, err := suite.PackEthereumTx(fromSigner, contract, big.NewInt(0), pack)
			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				fromDelAfter, found3 := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
				toDelAfter, found4 := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.Address().Bytes(), val.GetOperator())
				suite.Require().True(found4)
				if !found3 {
					fromDelAfter.Shares = sdk.ZeroDec()
				}
				if !found2 {
					toDelBefore.Shares = sdk.ZeroDec()
				}

				suite.Require().Equal(fromDelBefore.GetShares().TruncateInt().Sub(fromDelAfter.GetShares().TruncateInt()).BigInt(), shares)
				suite.Require().Equal(toDelAfter.GetShares().TruncateInt().Sub(toDelBefore.GetShares().TruncateInt()).BigInt(), shares)

				if tc.suftransfer != nil {
					tc.suftransfer(val.GetOperator(), delAddr, toSigner.Address(), delAmt)
				}

				fromBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, fromWithdrawAddr.Bytes())
				suite.Require().True(fromBalances.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.ZeroInt()))
				toBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, toWithdrawAddr.Bytes())
				if found2 {
					suite.Require().True(toBalances.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.ZeroInt()))
				} else {
					suite.Require().True(toBalances.Empty())
				}
			} else {
				suite.Require().True(err != nil || res.Failed())
			}
		})
	}
}

func TestStakingTransferFromSharesABI(t *testing.T) {
	stakingABI := staking.GetABI()

	method := stakingABI.Methods[staking.TransferFromSharesMethodName]
	require.Equal(t, method, staking.TransferFromSharesMethod)
	require.Equal(t, 4, len(staking.TransferFromSharesMethod.Inputs))
	require.Equal(t, 2, len(staking.TransferFromSharesMethod.Outputs))
}

func (suite *PrecompileTestSuite) TestTransferFromShares() {
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
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val := vals[0]
			delAmt := sdkmath.NewInt(int64(tmrand.Intn(10000) + 1000)).Mul(sdkmath.NewInt(1e18))
			fromSigner := suite.RandSigner()
			toSigner := suite.RandSigner()
			sender := suite.RandSigner()

			// from delegate, approve sender, sender send tx, transferFrom to toSigner
			// from delegate, approve contract, sender call contract, transferFrom to toSigner
			contract := staking.GetAddress()
			delAddr := fromSigner.Address()
			spender := sender.Address()
			if strings.HasPrefix(tc.name, "contract") {
				contract = suite.staking
				spender = suite.staking
			}

			tc.pretransfer(val.GetOperator(), delAddr, toSigner.Address(), delAmt)

			fromWithdrawAddr := helpers.GenerateAddress()
			err := suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, delAddr.Bytes(), fromWithdrawAddr.Bytes())
			suite.Require().NoError(err)
			toWithdrawAddr := helpers.GenerateAddress()
			err = suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, toSigner.AccAddress(), toWithdrawAddr.Bytes())
			suite.Require().NoError(err)

			suite.Commit()

			fromBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, fromWithdrawAddr.Bytes())
			suite.Require().True(fromBalances.Empty())
			toBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, toWithdrawAddr.Bytes())
			suite.Require().True(toBalances.Empty())

			fromDelBefore, found1 := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
			suite.Require().True(found1)
			toDelBefore, found2 := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.Address().Bytes(), val.GetOperator())

			// NOTE: if contract test, spender is staking test contract
			pack, shares, _ := tc.malleate(val.GetOperator(), spender, delAddr, toSigner.Address(), fromDelBefore.GetShares().TruncateInt().BigInt())
			tx, err := suite.PackEthereumTx(sender, contract, big.NewInt(0), pack)

			var res *evmtypes.MsgEthereumTxResponse
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				fromDelAfter, found3 := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr.Bytes(), val.GetOperator())
				toDelAfter, found4 := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.Address().Bytes(), val.GetOperator())
				suite.Require().True(found4)
				if !found3 {
					fromDelAfter.Shares = sdk.ZeroDec()
				}
				if !found2 {
					toDelBefore.Shares = sdk.ZeroDec()
				}

				suite.Require().Equal(fromDelBefore.GetShares().TruncateInt().Sub(fromDelAfter.GetShares().TruncateInt()).BigInt(), shares)
				suite.Require().Equal(toDelAfter.GetShares().TruncateInt().Sub(toDelBefore.GetShares().TruncateInt()).BigInt(), shares)

				allowance := suite.app.StakingKeeper.GetAllowance(suite.ctx, val.GetOperator(), delAddr.Bytes(), spender.Bytes())
				suite.Require().EqualValues(big.NewInt(0), allowance)

				if tc.suftransfer != nil {
					tc.suftransfer(val.GetOperator(), delAddr, toSigner.Address(), delAmt)
				}

				fromBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, fromWithdrawAddr.Bytes())
				suite.Require().True(fromBalances.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.ZeroInt()))
				toBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, toWithdrawAddr.Bytes())
				if found2 {
					suite.Require().True(toBalances.AmountOf(fxtypes.DefaultDenom).GT(sdkmath.ZeroInt()))
				} else {
					suite.Require().True(toBalances.Empty())
				}
			} else {
				suite.Require().True(err != nil || res.Failed())
			}
		})
	}
}
