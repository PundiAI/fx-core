package staking_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/precompiles/staking"
)

func TestStakingWithdrawABI(t *testing.T) {
	stakingABI := fxtypes.MustABIJson(staking.JsonABI)

	method := stakingABI.Methods[staking.WithdrawMethod.Name]
	require.Equal(t, method, staking.WithdrawMethod)
}

func (suite *PrecompileTestSuite) TestWithdraw() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestWithdrawName, val.String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestWithdrawName, newVal)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("execution reverted: withdraw failed: invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()
				pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestWithdrawName, newVal)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(suite.signer, suite.staking, big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return "execution reverted: withdraw failed: no validator distribution info"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]

			delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
			pack, err := fxtypes.MustABIJson(StakingTestABI).Pack(StakingTestDelegateName, val0.GetOperator().String(), delAmt.BigInt())
			suite.Require().NoError(err)
			delegateEthTx, err := suite.PackEthereumTx(suite.signer, suite.staking, delAmt.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.staking.Bytes(), val0.GetOperator())
			suite.Require().True(found)
			suite.Commit()

			bal1 := suite.app.BankKeeper.GetBalance(suite.ctx, suite.staking.Bytes(), fxtypes.DefaultDenom)

			ethTx, err, errArgs := tc.malleate(val0.GetOperator(), delegation.Shares)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)

				bal2 := suite.app.BankKeeper.GetBalance(suite.ctx, suite.staking.Bytes(), fxtypes.DefaultDenom)

				var reward struct {
					Value *big.Int
				}
				err = fxtypes.MustABIJson(staking.JsonABI).UnpackIntoInterface(&reward, "withdraw", res.Ret)
				suite.Require().NoError(err)
				suite.Require().True(reward.Value.Cmp(big.NewInt(0)) == 1)

				suite.Require().Equal(bal2.Sub(bal1).Amount.BigInt().String(), reward.Value.String())
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
				if res.Failed() {
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
				} else {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestAccountWithdraw() {
	testCases := []struct {
		name     string
		malleate func(signer *helpers.Signer, val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.WithdrawMethodName, val.String())
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(signer, staking.StakingAddress, big.NewInt(0), pack)
				return tx, err, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := val.String() + "1"
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.WithdrawMethodName, newVal)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(signer, staking.StakingAddress, big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return fmt.Sprintf("invalid validator address: %s", errArgs[0])
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(signer *helpers.Signer, val sdk.ValAddress, shares sdk.Dec) (*evmtypes.MsgEthereumTx, error, []string) {
				newVal := sdk.ValAddress(signer.Address().Bytes()).String()
				pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.WithdrawMethodName, newVal)
				suite.Require().NoError(err)
				tx, err := suite.PackEthereumTx(signer, staking.StakingAddress, big.NewInt(0), pack)
				return tx, err, []string{newVal}
			},
			error: func(errArgs []string) string {
				return "no validator distribution info"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
			val0 := vals[0]

			priv, err := ethsecp256k1.GenerateKey()
			require.NoError(suite.T(), err)
			newSigner := helpers.NewSigner(priv)

			helpers.AddTestAddr(suite.app, suite.ctx, newSigner.AccAddress(),
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18)))))

			delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
			pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
			suite.Require().NoError(err)
			delegateEthTx, err := suite.PackEthereumTx(newSigner, staking.StakingAddress, delAmt.BigInt(), pack)
			suite.Require().NoError(err)
			res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
			suite.Require().NoError(err)
			suite.Require().False(res.Failed(), res.VmError)

			delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, newSigner.AccAddress(), val0.GetOperator())
			suite.Require().True(found)
			suite.Commit()

			bal1 := suite.app.BankKeeper.GetBalance(suite.ctx, newSigner.AccAddress(), fxtypes.DefaultDenom)

			ethTx, err, errArgs := tc.malleate(newSigner, val0.GetOperator(), delegation.Shares)
			if err == nil {
				res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), ethTx)
			}

			if tc.result {
				suite.Require().NoError(err)
				suite.Require().False(res.Failed(), res.VmError)
				bal2 := suite.app.BankKeeper.GetBalance(suite.ctx, newSigner.AccAddress(), fxtypes.DefaultDenom)

				var reward struct {
					Value *big.Int
				}
				err = fxtypes.MustABIJson(staking.JsonABI).UnpackIntoInterface(&reward, "withdraw", res.Ret)
				suite.Require().NoError(err)
				suite.Require().True(reward.Value.Cmp(big.NewInt(0)) == 1)

				suite.Require().Equal(bal2.Sub(bal1).Amount.BigInt().String(), reward.Value.String())
			} else {
				suite.Require().True(err != nil || res.Failed())
				if err != nil {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
				if res.Failed() {
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
				} else {
					suite.Require().Equal(tc.error(errArgs), err.Error())
				}
			}
		})
	}
}

func (suite *PrecompileTestSuite) TestAccountWithdrawOtherAddress() {
	vals := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
	val0 := vals[0]

	delAmt := sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
	pack, err := fxtypes.MustABIJson(staking.JsonABI).Pack(staking.DelegateMethodName, val0.GetOperator().String(), delAmt.BigInt())
	suite.Require().NoError(err)
	delegateEthTx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, delAmt.BigInt(), pack)
	suite.Require().NoError(err)
	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), delegateEthTx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, suite.signer.AccAddress(), val0.GetOperator())
	suite.Require().True(found)

	rewardAddr := helpers.GenerateAddress()
	err = suite.app.DistrKeeper.SetWithdrawAddr(suite.ctx, suite.signer.AccAddress(), rewardAddr.Bytes())
	suite.Require().NoError(err)

	bal1 := suite.app.BankKeeper.GetBalance(suite.ctx, rewardAddr.Bytes(), fxtypes.DefaultDenom)

	suite.Commit()

	pack, err = fxtypes.MustABIJson(staking.JsonABI).Pack(staking.WithdrawMethodName, val0.GetOperator().String())
	suite.Require().NoError(err)
	withdrawTx, err := suite.PackEthereumTx(suite.signer, staking.StakingAddress, big.NewInt(0), pack)
	suite.Require().NoError(err)
	res, err = suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), withdrawTx)
	suite.Require().NoError(err)

	var reward struct {
		Value *big.Int
	}
	err = fxtypes.MustABIJson(staking.JsonABI).UnpackIntoInterface(&reward, "withdraw", res.Ret)
	suite.Require().NoError(err)
	suite.Require().True(reward.Value.Cmp(big.NewInt(0)) == 1)

	bal2 := suite.app.BankKeeper.GetBalance(suite.ctx, rewardAddr.Bytes(), fxtypes.DefaultDenom)
	suite.Require().Equal(bal2.Sub(bal1).Amount.BigInt().String(), reward.Value.String())
}
