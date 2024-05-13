package staking_test

import (
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
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/precompiles/staking"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
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

//gocyclo:ignore
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

			fromWithdrawAddr := helpers.GenHexAddress()
			err := suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, delAddr.Bytes(), fromWithdrawAddr.Bytes())
			suite.Require().NoError(err)
			toWithdrawAddr := helpers.GenHexAddress()
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

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == staking.TransferSharesEvent.ID.String() {
						suite.Require().Equal(3, len(log.Topics))
						suite.Require().Equal(log.Topics[1], delAddr.Hash().String())
						suite.Require().Equal(log.Topics[2], toSigner.Address().Hash().String())
						unpack, err := staking.TransferSharesEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						unpackValStr := unpack[0].(string)
						unpackShares := unpack[1].(*big.Int)
						suite.Require().Equal(val.GetOperator().String(), unpackValStr)
						suite.Require().Equal(shares.String(), unpackShares.String())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type == fxstakingtypes.EventTypeTransferShares {
						for _, attr := range event.Attributes {
							if string(attr.Key) == fxstakingtypes.AttributeKeyFrom {
								suite.Require().Equal(string(attr.Value), sdk.AccAddress(delAddr.Bytes()).String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeyRecipient {
								suite.Require().Equal(string(attr.Value), toSigner.AccAddress().String())
							}
							if string(attr.Key) == stakingtypes.AttributeKeyValidator {
								suite.Require().Equal(string(attr.Value), val.GetOperator().String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeyShares {
								suite.Require().Equal(string(attr.Value), shares.String())
							}
						}
						existEvent = true
					}
				}
				suite.Require().True(existEvent)
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

//gocyclo:ignore
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

			fromWithdrawAddr := helpers.GenHexAddress()
			err := suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, delAddr.Bytes(), fromWithdrawAddr.Bytes())
			suite.Require().NoError(err)
			toWithdrawAddr := helpers.GenHexAddress()
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

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == staking.TransferSharesEvent.ID.String() {
						suite.Require().Equal(3, len(log.Topics))
						suite.Require().Equal(log.Topics[1], delAddr.Hash().String())
						suite.Require().Equal(log.Topics[2], toSigner.Address().Hash().String())
						unpack, err := staking.TransferSharesEvent.Inputs.NonIndexed().Unpack(log.Data)
						suite.Require().NoError(err)
						unpackValStr := unpack[0].(string)
						unpackShares := unpack[1].(*big.Int)
						suite.Require().Equal(val.GetOperator().String(), unpackValStr)
						suite.Require().Equal(shares.String(), unpackShares.String())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.ctx.EventManager().Events() {
					if event.Type == fxstakingtypes.EventTypeTransferShares {
						for _, attr := range event.Attributes {
							if string(attr.Key) == fxstakingtypes.AttributeKeyFrom {
								suite.Require().Equal(string(attr.Value), sdk.AccAddress(delAddr.Bytes()).String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeyRecipient {
								suite.Require().Equal(string(attr.Value), toSigner.AccAddress().String())
							}
							if string(attr.Key) == stakingtypes.AttributeKeyValidator {
								suite.Require().Equal(string(attr.Value), val.GetOperator().String())
							}
							if string(attr.Key) == fxstakingtypes.AttributeKeyShares {
								suite.Require().Equal(string(attr.Value), shares.String())
							}
						}
						existEvent = true
					}
				}
				suite.Require().True(existEvent)
			} else {
				suite.Require().True(err != nil || res.Failed())
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestTransferSharesCompare() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val := vals[0]
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()
	signer3 := suite.RandSigner()

	helpers.AddTestAddr(suite.app, suite.ctx, signer1.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))

	// starting info 1,2,3
	startingInfo := suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod)
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer2.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod)
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer3.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod)

	// signer1 chain delegate to val
	shares1, err := suite.app.StakingKeeper.Delegate(suite.ctx, signer1.AccAddress(), delAmount, stakingtypes.Unbonded, val, true)
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(2, startingInfo.PreviousPeriod)

	// signer2 evm delegate to val
	shares2 := suite.PrecompileStakingDelegate(signer2, val.GetOperator(), delAmount.BigInt())
	suite.Require().Equal(shares1.TruncateInt().BigInt(), shares2)

	// signer2 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer2.AccAddress())
	suite.Require().EqualValues(3, startingInfo.PreviousPeriod)

	// generate block
	suite.Commit()

	// signer1 withdraw
	rewards1, err := suite.app.DistrKeeper.WithdrawDelegationRewards(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(4, startingInfo.PreviousPeriod)

	// signer2 transfer shares to singer3
	halfShares := big.NewInt(0).Div(shares2, big.NewInt(2))
	surplusShares, rewards2 := suite.PrecompileStakingTransferShares(signer2, val.GetOperator(), signer3.Address(), halfShares)

	// starting info 2,3
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer2.AccAddress())
	suite.Require().EqualValues(5, startingInfo.PreviousPeriod)
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer3.AccAddress())
	suite.Require().EqualValues(6, startingInfo.PreviousPeriod)

	// rewards1 equal rewards2
	suite.Require().EqualValues(rewards1.AmountOf(fxtypes.DefaultDenom).BigInt(), rewards2)
	// surplus shares equal to half shares
	suite.Require().EqualValues(halfShares, surplusShares)

	// signer1 undelegate half shares
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, signer1.AccAddress(), val.GetOperator(), sdk.NewDecFromBigInt(halfShares))
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(7, startingInfo.PreviousPeriod)

	// signer1 shares equal to half shares
	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().True(found)
	suite.Require().EqualValues(halfShares, delegation.GetShares().TruncateInt().BigInt())

	// generate block
	suite.Commit()

	// transfer shares surplus
	surplusShares, rewards2 = suite.PrecompileStakingTransferShares(signer2, val.GetOperator(), signer3.Address(), surplusShares)

	// starting info 2,3
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer2.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transfer all shares, starting info removed
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer3.AccAddress())
	suite.Require().EqualValues(9, startingInfo.PreviousPeriod)

	// surplus shares equal to zero
	suite.Require().EqualValues(big.NewInt(0).String(), surplusShares.String())

	// get signer3 fx balance
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, signer3.AccAddress())
	suite.Require().EqualValues(rewards2, balances.AmountOf(fxtypes.DefaultDenom).BigInt())

	// signer1 withdraw
	rewards1, err = suite.app.DistrKeeper.WithdrawDelegationRewards(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(10, startingInfo.PreviousPeriod)

	// rewards1 equal rewards2
	suite.Require().EqualValues(rewards1.AmountOf(fxtypes.DefaultDenom).BigInt(), rewards2)

	// signer1 undelegate all shares
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, signer1.AccAddress(), val.GetOperator(), delegation.GetShares())
	suite.Require().NoError(err)

	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // undelegate all, starting info removed

	// signer1 shares equal to zero
	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().False(found)

	// signer3 delegation
	shares3, amount3 := suite.PrecompileStakingDelegation(val.GetOperator(), signer3.Address())
	suite.Require().EqualValues(big.NewInt(0).Mul(halfShares, big.NewInt(2)), shares3)
	suite.Require().EqualValues(delAmount.BigInt(), amount3)

	// generate block
	suite.Commit()

	// signer3 evm withdraw
	_ = suite.PrecompileStakingWithdraw(signer3, val.GetOperator())

	// signer3 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer3.AccAddress())
	suite.Require().EqualValues(12, startingInfo.PreviousPeriod)

	// signer3 transferFrom shares
	suite.PrecompileStakingApproveShares(signer3, val.GetOperator(), signer2.Address(), shares3)
	suite.PrecompileStakingTransferFromShares(signer2, val.GetOperator(), signer3.Address(), signer1.Address(), shares3)
	newShares1, _ := suite.PrecompileStakingDelegation(val.GetOperator(), signer1.Address())
	suite.Require().EqualValues(newShares1, shares3)

	// signer3 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer3.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transferFrom all shares, starting info removed
	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(14, startingInfo.PreviousPeriod)

	// generate block
	suite.Commit()

	// singer1 evm undelegate
	_ = suite.PrecompileStakingUndelegate(signer1, val.GetOperator(), halfShares)
	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(16, startingInfo.PreviousPeriod) // withdraw +1, undelegate +1

	_ = suite.PrecompileStakingUndelegate(signer1, val.GetOperator(), halfShares)
	// signer1 starting info
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // undelegate all, starting info removed

	// signer1 shares equal to zero
	_, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().False(found)

	// starting info 1,2,3
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer1.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // undelegate all, starting info removed
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer2.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transfer all shares, starting info removed
	startingInfo = suite.app.DistrKeeper.GetDelegatorStartingInfo(suite.ctx, val.GetOperator(), signer3.AccAddress())
	suite.Require().EqualValues(0, startingInfo.PreviousPeriod) // transfer all shares, starting info removed
}

func (suite *PrecompileTestSuite) TestPrecompileStakingSteps() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val := vals[0]
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()
	signer3 := suite.RandSigner()

	// delegate 1 and 2
	suite.Delegate(val.GetOperator(), delAmount, signer1.AccAddress(), signer2.AccAddress())
	suite.Commit()

	// get 1 shares
	shares1, _ := suite.PrecompileStakingDelegation(val.GetOperator(), signer1.Address())

	// 1 transfer shares to 2
	suite.PrecompileStakingTransferShares(signer1, val.GetOperator(), signer2.Address(), shares1)
	suite.Commit()

	// 2 transfer shares to 3
	suite.PrecompileStakingTransferShares(signer2, val.GetOperator(), signer3.Address(), shares1)
	suite.Commit()

	// delegate 1,2,3
	suite.Delegate(val.GetOperator(), delAmount, signer1.AccAddress(), signer2.AccAddress(), signer3.AccAddress())
	suite.Commit()

	// precompile delegate 1,2,3
	suite.PrecompileStakingDelegate(signer1, val.GetOperator(), delAmount.BigInt())
	suite.PrecompileStakingDelegate(signer2, val.GetOperator(), delAmount.BigInt())
	suite.PrecompileStakingDelegate(signer3, val.GetOperator(), delAmount.BigInt())
	suite.Commit()

	// precompile withdraw 1,2,3
	suite.PrecompileStakingWithdraw(signer1, val.GetOperator())
	suite.PrecompileStakingWithdraw(signer2, val.GetOperator())
	suite.PrecompileStakingWithdraw(signer3, val.GetOperator())
	suite.Commit()

	// withdraw 1,2,3
	_, err := suite.app.DistrKeeper.WithdrawDelegationRewards(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)
	_, err = suite.app.DistrKeeper.WithdrawDelegationRewards(suite.ctx, signer2.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)
	_, err = suite.app.DistrKeeper.WithdrawDelegationRewards(suite.ctx, signer3.AccAddress(), val.GetOperator())
	suite.Require().NoError(err)
	suite.Commit()

	// precompile undelegate 1,2,3
	suite.PrecompileStakingUndelegate(signer1, val.GetOperator(), shares1)
	suite.PrecompileStakingUndelegate(signer2, val.GetOperator(), shares1)
	suite.PrecompileStakingUndelegate(signer3, val.GetOperator(), shares1)
	suite.Commit()

	// undelegate 1,2,3
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, signer1.AccAddress(), val.GetOperator(), sdk.NewDecFromBigInt(shares1))
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, signer2.AccAddress(), val.GetOperator(), sdk.NewDecFromBigInt(shares1))
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, signer3.AccAddress(), val.GetOperator(), sdk.NewDecFromBigInt(shares1))
	suite.Require().NoError(err)
	suite.Commit()

	// delegate 1,2,3
	suite.Delegate(val.GetOperator(), delAmount, signer1.AccAddress(), signer2.AccAddress(), signer3.AccAddress())
	suite.Commit()

	// transfer shares 1,2,3
	suite.PrecompileStakingTransferShares(signer1, val.GetOperator(), signer2.Address(), shares1)
	suite.PrecompileStakingTransferShares(signer2, val.GetOperator(), signer3.Address(), shares1)
	suite.PrecompileStakingTransferShares(signer3, val.GetOperator(), signer1.Address(), shares1)
	suite.Commit()

	// transferFrom shares 1,2,3
	suite.PrecompileStakingApproveShares(signer1, val.GetOperator(), signer2.Address(), shares1)
	suite.PrecompileStakingTransferFromShares(signer2, val.GetOperator(), signer1.Address(), signer2.Address(), shares1)
	suite.PrecompileStakingApproveShares(signer2, val.GetOperator(), signer3.Address(), shares1)
	suite.PrecompileStakingTransferFromShares(signer3, val.GetOperator(), signer2.Address(), signer3.Address(), shares1)
	suite.PrecompileStakingApproveShares(signer3, val.GetOperator(), signer1.Address(), shares1)
	suite.PrecompileStakingTransferFromShares(signer1, val.GetOperator(), signer3.Address(), signer1.Address(), shares1)
	suite.Commit()
}

func (suite *PrecompileTestSuite) TestTransferSharesRedelegate() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val := vals[0]
	valTmp := vals[1]
	delAmount := sdkmath.NewInt(int64(tmrand.Int() + 100)).Mul(sdkmath.NewInt(1e18))
	signer1 := suite.RandSigner()
	signer2 := suite.RandSigner()

	// delegate 1 and 2
	suite.Delegate(valTmp.GetOperator(), delAmount, signer1.AccAddress())
	suite.Commit()

	delegationTmp, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, signer1.AccAddress(), valTmp.GetOperator())
	suite.Require().True(found)

	// redelegate
	suite.Redelegate(valTmp.GetOperator(), val.GetOperator(), signer1.AccAddress(), delegationTmp.Shares)
	suite.Commit()

	delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, signer1.AccAddress(), val.GetOperator())
	suite.Require().True(found)
	suite.Require().Equal(delegationTmp.Shares, delegation.Shares)

	// transfer shares
	pack, err := staking.GetABI().Pack(staking.TransferSharesMethodName, val.GetOperator().String(), signer2.Address(), delegation.Shares.TruncateInt().BigInt())
	suite.Require().NoError(err)
	_, err = suite.PackEthereumTx(signer1, staking.GetAddress(), big.NewInt(0), pack)
	suite.Require().EqualError(err, "from has receiving redelegation")
}
