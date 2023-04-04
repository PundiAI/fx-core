package staking_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

func TestStakingTransferABI(t *testing.T) {
	stakingABI := fxtypes.MustABIJson(staking.JsonABI)

	method := stakingABI.Methods[staking.TransferMethodName]
	require.Equal(t, method, staking.TransferMethod)
	require.Equal(t, 3, len(staking.TransferMethod.Inputs))
	require.Equal(t, 2, len(staking.TransferMethod.Outputs))

	event := stakingABI.Events[staking.TransferEventName]
	require.Equal(t, event, staking.TransferEvent)
	require.Equal(t, 5, len(staking.TransferEvent.Inputs))
}

func (suite *PrecompileTestSuite) TestTransfer() {
	delegateFromFunc := func(val sdk.ValAddress, fromSigner, _ *helpers.Signer, delAmount sdkmath.Int) {
		helpers.AddTestAddr(suite.app, suite.ctx, fromSigner.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
		_, err := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: fromSigner.AccAddress().String(),
			ValidatorAddress: val.String(),
			Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
		})
		suite.Require().NoError(err)
	}
	undelegateToFunc := func(val sdk.ValAddress, _, toSigner *helpers.Signer, delAmount sdkmath.Int) {
		toDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.AccAddress(), val)
		suite.Require().True(found)
		_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, toSigner.AccAddress(), val, toDel.Shares)
		suite.Require().NoError(err)
	}
	delegateFromToFunc := func(val sdk.ValAddress, fromSigner, toSigner *helpers.Signer, delAmount sdkmath.Int) {
		helpers.AddTestAddr(suite.app, suite.ctx, fromSigner.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
		_, err := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: fromSigner.AccAddress().String(),
			ValidatorAddress: val.String(),
			Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
		})
		suite.Require().NoError(err)

		helpers.AddTestAddr(suite.app, suite.ctx, toSigner.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
		_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: toSigner.AccAddress().String(),
			ValidatorAddress: val.String(),
			Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
		})
		suite.Require().NoError(err)
	}
	delegateToFromFunc := func(val sdk.ValAddress, fromSigner, toSigner *helpers.Signer, delAmount sdkmath.Int) {
		helpers.AddTestAddr(suite.app, suite.ctx, toSigner.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
		_, err := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: toSigner.AccAddress().String(),
			ValidatorAddress: val.String(),
			Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
		})
		suite.Require().NoError(err)

		helpers.AddTestAddr(suite.app, suite.ctx, fromSigner.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
		_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: fromSigner.AccAddress().String(),
			ValidatorAddress: val.String(),
			Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
		})
		suite.Require().NoError(err)
	}
	undelegateFromToFunc := func(val sdk.ValAddress, fromSigner, toSigner *helpers.Signer, delAmount sdkmath.Int) {
		fromDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, fromSigner.AccAddress(), val)
		suite.Require().True(found)
		_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, fromSigner.AccAddress(), val, fromDel.Shares)
		suite.Require().NoError(err)

		toDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.AccAddress(), val)
		suite.Require().True(found)
		_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, toSigner.AccAddress(), val, toDel.Shares)
		suite.Require().NoError(err)
	}
	undelegateToFromFunc := func(val sdk.ValAddress, fromSigner, toSigner *helpers.Signer, delAmount sdkmath.Int) {
		toDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, toSigner.AccAddress(), val)
		suite.Require().True(found)
		_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, toSigner.AccAddress(), val, toDel.Shares)
		suite.Require().NoError(err)

		fromDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, fromSigner.AccAddress(), val)
		suite.Require().True(found)
		_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, fromSigner.AccAddress(), val, fromDel.Shares)
		suite.Require().NoError(err)
	}

	packTransferRand := func(val sdk.ValAddress, to *helpers.Signer, shares *big.Int) ([]byte, *big.Int, []string) {
		randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
		pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.TransferMethodName, val.String(), to.Address(), randShares)
		suite.Require().NoError(err)
		return pack, randShares, nil
	}
	packTransferAll := func(val sdk.ValAddress, to *helpers.Signer, shares *big.Int) ([]byte, *big.Int, []string) {
		pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.TransferMethodName, val.String(), to.Address(), shares)
		suite.Require().NoError(err)
		return pack, shares, nil
	}

	testCases := []struct {
		name        string
		pretransfer func(val sdk.ValAddress, fromSigner, toSigner *helpers.Signer, delAmount sdkmath.Int)
		malleate    func(val sdk.ValAddress, to *helpers.Signer, shares *big.Int) ([]byte, *big.Int, []string)
		suftransfer func(val sdk.ValAddress, fromSigner, toSigner *helpers.Signer, delAmount sdkmath.Int)
		error       func(errArgs []string) string
		result      bool
	}{
		{
			name:        "ok - from delegated",
			pretransfer: delegateFromFunc,
			malleate:    packTransferRand,
			result:      true,
		},
		{
			name:        "ok - from delegated - undelegate from and to",
			pretransfer: delegateFromFunc,
			malleate:    packTransferRand,
			suftransfer: undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - undelegate to and from",
			pretransfer: delegateFromFunc,
			malleate:    packTransferRand,
			suftransfer: undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - delegate from and to",
			pretransfer: delegateFromFunc,
			malleate:    packTransferRand,
			suftransfer: delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - delegate to and from",
			pretransfer: delegateFromFunc,
			malleate:    packTransferRand,
			suftransfer: delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - undelegate to",
			pretransfer: delegateFromFunc,
			malleate:    packTransferAll,
			suftransfer: undelegateToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - delegate from and to",
			pretransfer: delegateFromFunc,
			malleate:    packTransferAll,
			suftransfer: delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from delegated - transfer all - delegate to and from",
			pretransfer: delegateFromFunc,
			malleate:    packTransferAll,
			suftransfer: delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferRand,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - delegate from and to",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferRand,
			suftransfer: delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - delegate to and from",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferRand,
			suftransfer: delegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - undelegate from and to",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferRand,
			suftransfer: undelegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - undelegate to and from",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferRand,
			suftransfer: undelegateToFromFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferAll,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - undelegate to",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferAll,
			suftransfer: undelegateToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - delegate from and to",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferAll,
			suftransfer: delegateFromToFunc,
			result:      true,
		},
		{
			name:        "ok - from and to delegated - transfer all - delegate to and from",
			pretransfer: delegateFromToFunc,
			malleate:    packTransferAll,
			suftransfer: delegateToFromFunc,
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

			contract := staking.GetPrecompileAddress()
			delAddr := fromSigner.Address()

			tc.pretransfer(val.GetOperator(), fromSigner, toSigner, delAmt)

			fromWithdrawAddr := helpers.GenerateAddress()
			err := suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, fromSigner.AccAddress(), fromWithdrawAddr.Bytes())
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

			pack, shares, _ := tc.malleate(val.GetOperator(), toSigner, fromDelBefore.GetShares().TruncateInt().BigInt())
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
					tc.suftransfer(val.GetOperator(), fromSigner, toSigner, delAmt)
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

func TestStakingTransferFromABI(t *testing.T) {
	stakingABI := fxtypes.MustABIJson(staking.JsonABI)

	method := stakingABI.Methods[staking.TransferFromMethodName]
	require.Equal(t, method, staking.TransferFromMethod)
	require.Equal(t, 4, len(staking.TransferFromMethod.Inputs))
	require.Equal(t, 2, len(staking.TransferFromMethod.Outputs))
}
