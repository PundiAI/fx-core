package precompile_test

import (
	"math/big"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
	"github.com/functionx/fx-core/v8/x/staking/testutil"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type PrecompileTestSuite struct {
	helpers.BaseSuite
	testutil.StakingSuite
	signer          *helpers.Signer
	stakingTestAddr common.Address
	stakingAddr     common.Address

	allowanceSharesMethod   *precompile.AllowanceSharesMethod
	approveSharesMethod     *precompile.ApproveSharesMethod
	delegateV2Method        *precompile.DelegateV2Method
	delegationMethod        *precompile.DelegationMethod
	delegationRewardsMethod *precompile.DelegationRewardsMethod
	redelegateV2Method      *precompile.RedelegateV2Method
	// slashingInfoMethod       *precompile.SlashingInfoMethod
	transferSharesMethod     *precompile.TransferSharesMethod
	transferFromSharesMethod *precompile.TransferFromSharesMethod
	undelegateV2Method       *precompile.UndelegateV2Method
	validatorListMethod      *precompile.ValidatorListMethod
	withdrawMethod           *precompile.WithdrawMethod
}

func TestPrecompileTestSuite(t *testing.T) {
	fxtypes.SetConfig(true)
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *PrecompileTestSuite) SetupTest() {
	suite.MintValNumber = 2
	suite.BaseSuite.SetupTest()
	suite.Commit(10)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.signer = helpers.NewSigner(priv)
	suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18))))

	stakingContract, err := suite.App.EvmKeeper.DeployContract(suite.Ctx, suite.signer.Address(), contract.MustABIJson(testscontract.StakingTestMetaData.ABI), contract.MustDecodeHex(testscontract.StakingTestMetaData.Bin))
	suite.Require().NoError(err)
	suite.stakingTestAddr = stakingContract

	suite.stakingAddr = common.HexToAddress(contract.StakingAddress)

	suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18))))

	suite.StakingSuite.Init(suite.Require(), suite.Ctx, suite.App.StakingKeeper)

	suite.allowanceSharesMethod = precompile.NewAllowanceSharesMethod(nil)
	suite.approveSharesMethod = precompile.NewApproveSharesMethod(nil)
	suite.delegateV2Method = precompile.NewDelegateV2Method(nil)
	suite.delegationMethod = precompile.NewDelegationMethod(nil)
	suite.delegationRewardsMethod = precompile.NewDelegationRewardsMethod(nil)
	suite.redelegateV2Method = precompile.NewRedelegateV2Method(nil)
	// suite.slashingInfoMethod = precompile.NewSlashingInfoMethod(nil)
	suite.transferSharesMethod = precompile.NewTransferSharesMethod(nil)
	suite.transferFromSharesMethod = precompile.NewTransferFromSharesMethod(nil)
	suite.undelegateV2Method = precompile.NewUndelegateV2Method(nil)
	suite.validatorListMethod = precompile.NewValidatorListMethod(nil)
	suite.withdrawMethod = precompile.NewWithdrawMethod(nil)
}

func (suite *PrecompileTestSuite) DistributionQueryClient(ctx sdk.Context) distributiontypes.QueryClient {
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, suite.App.InterfaceRegistry())
	distributiontypes.RegisterQueryServer(queryHelper, distributionkeeper.NewQuerier(suite.App.DistrKeeper))
	return distributiontypes.NewQueryClient(queryHelper)
}

func (suite *PrecompileTestSuite) EthereumTx(signer *helpers.Signer, to common.Address, amount *big.Int, data []byte) *evmtypes.MsgEthereumTxResponse {
	ethTx := evmtypes.NewTx(
		fxtypes.EIP155ChainID(suite.Ctx.ChainID()),
		suite.App.EvmKeeper.GetNonce(suite.Ctx, signer.Address()),
		&to,
		amount,
		contract.DefaultGasCap,
		nil,
		nil,
		nil,
		data,
		nil,
	)
	ethTx.From = signer.Address().Bytes()
	err := ethTx.Sign(ethtypes.LatestSignerForChainID(fxtypes.EIP155ChainID(suite.Ctx.ChainID())), signer)
	suite.Require().NoError(err)

	res, err := suite.App.EvmKeeper.EthereumTx(suite.Ctx, ethTx)
	suite.Require().NoError(err)
	return res
}

func (suite *PrecompileTestSuite) RandSigner() *helpers.Signer {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	account := suite.App.AccountKeeper.NewAccountWithAddress(suite.Ctx, signer.AccAddress())
	suite.App.AccountKeeper.SetAccount(suite.Ctx, account)
	return signer
}

func (suite *PrecompileTestSuite) delegateFromFunc(val sdk.ValAddress, from, _ common.Address, delAmount sdkmath.Int) {
	suite.MintToken(from.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmount))
	_, err := stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) undelegateToFunc(val sdk.ValAddress, _, to common.Address, _ sdkmath.Int) {
	toDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) delegateFromToFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
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

func (suite *PrecompileTestSuite) delegateToFromFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
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

func (suite *PrecompileTestSuite) undelegateFromToFunc(val sdk.ValAddress, from, to common.Address, _ sdkmath.Int) {
	fromDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, from.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, from.Bytes(), val, fromDel.Shares)
	suite.Require().NoError(err)

	toDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) undelegateToFromFunc(val sdk.ValAddress, from, to common.Address, _ sdkmath.Int) {
	toDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, to.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)

	fromDel, err := suite.App.StakingKeeper.GetDelegation(suite.Ctx, from.Bytes(), val)
	suite.Require().NoError(err)
	_, _, err = suite.App.StakingKeeper.Undelegate(suite.Ctx, from.Bytes(), val, fromDel.Shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) approveFunc(val sdk.ValAddress, owner, spender common.Address, allowance *big.Int) {
	suite.App.StakingKeeper.SetAllowance(suite.Ctx, val, owner.Bytes(), spender.Bytes(), allowance)
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegation(val sdk.ValAddress, del common.Address) (*big.Int, *big.Int) {
	input, err := suite.delegationMethod.PackInput(fxstakingtypes.DelegationArgs{
		Validator: val.String(),
		Delegator: del,
	})
	suite.Require().NoError(err)
	res, err := suite.App.EvmKeeper.CallEVMWithoutGas(suite.Ctx, del, &suite.stakingAddr, nil, input, false)
	suite.Require().NoError(err)
	shares, amount, err := suite.delegationMethod.UnpackOutput(res.Ret)
	suite.Require().NoError(err)
	return shares, amount
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegateV2(signer *helpers.Signer, val sdk.ValAddress, amt *big.Int) *big.Int {
	suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(amt)))
	pack, err := suite.delegateV2Method.PackInput(fxstakingtypes.DelegateV2Args{
		Validator: val.String(),
		Amount:    amt,
	})
	suite.Require().NoError(err)

	_, amountBefore := suite.PrecompileStakingDelegation(val, signer.Address())

	res := suite.EthereumTx(signer, suite.stakingAddr, big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	shares, amount := suite.PrecompileStakingDelegation(val, signer.Address())
	suite.Require().Equal(amt.String(), big.NewInt(0).Sub(amount, amountBefore).String())
	return shares
}

func (suite *PrecompileTestSuite) PrecompileStakingWithdraw(signer *helpers.Signer, val sdk.ValAddress) *big.Int {
	balanceBefore := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	pack, err := suite.withdrawMethod.PackInput(fxstakingtypes.WithdrawArgs{Validator: val.String()})
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, suite.stakingAddr, big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	balanceAfter := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingTransferShares(signer *helpers.Signer, val sdk.ValAddress, receipt common.Address, shares *big.Int) (*big.Int, *big.Int) {
	balanceBefore := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	pack, err := suite.transferSharesMethod.PackInput(fxstakingtypes.TransferSharesArgs{
		Validator: val.String(),
		To:        receipt,
		Shares:    shares,
	})
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, suite.stakingAddr, big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	signerShares, _ := suite.PrecompileStakingDelegation(val, signer.Address())

	balanceAfter := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return signerShares, rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingUndelegateV2(signer *helpers.Signer, val sdk.ValAddress, shares *big.Int) *big.Int {
	balanceBefore := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	pack, err := suite.undelegateV2Method.PackInput(fxstakingtypes.UndelegateV2Args{
		Validator: val.String(),
		Amount:    shares,
	})
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, suite.stakingAddr, big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	balanceAfter := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingApproveShares(signer *helpers.Signer, val sdk.ValAddress, spender common.Address, shares *big.Int) {
	pack, err := suite.approveSharesMethod.PackInput(fxstakingtypes.ApproveSharesArgs{
		Validator: val.String(),
		Spender:   spender,
		Shares:    shares,
	})
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, suite.stakingAddr, big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *PrecompileTestSuite) PrecompileStakingTransferFromShares(signer *helpers.Signer, val sdk.ValAddress, from, receipt common.Address, shares *big.Int) {
	pack, err := suite.transferFromSharesMethod.PackInput(fxstakingtypes.TransferFromSharesArgs{
		Validator: val.String(),
		From:      from,
		To:        receipt,
		Shares:    shares,
	})
	suite.Require().NoError(err)
	res := suite.EthereumTx(signer, suite.stakingAddr, big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *PrecompileTestSuite) Delegate(val sdk.ValAddress, amount sdkmath.Int, dels ...sdk.AccAddress) {
	for _, del := range dels {
		suite.MintToken(del, sdk.NewCoin(fxtypes.DefaultDenom, amount))
		validator, err := suite.App.StakingKeeper.GetValidator(suite.Ctx, val)
		suite.Require().NoError(err)
		_, err = suite.App.StakingKeeper.Delegate(suite.Ctx, del, amount, stakingtypes.Unbonded, validator, true)
		suite.Require().NoError(err)
	}
}

func (suite *PrecompileTestSuite) Redelegate(valSrc, valDest sdk.ValAddress, del sdk.AccAddress, shares sdkmath.LegacyDec) {
	_, err := suite.App.StakingKeeper.BeginRedelegation(suite.Ctx, del, valSrc, valDest, shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) Error(res *evmtypes.MsgEthereumTxResponse, errResult error) {
	suite.Require().True(res.Failed())
	if res.VmError != vm.ErrExecutionReverted.Error() {
		suite.Require().Equal(errResult.Error(), res.VmError)
		return
	}

	if len(res.Ret) > 0 {
		reason, err := abi.UnpackRevert(res.Ret)
		suite.Require().NoError(err)
		suite.Require().Equal(errResult.Error(), reason)
		return
	}

	suite.Require().Equal(errResult.Error(), vm.ErrExecutionReverted.Error())
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

func (suite *PrecompileTestSuite) CheckRedelegateLogs(logs []*evmtypes.Log, delAddr common.Address, valSrc, valDst string, shares, amount *big.Int, completionTime int64) {
	redelegateV2Method := precompile.NewRedelegateV2Method(nil)
	existLog := false
	for _, log := range logs {
		if log.Topics[0] == redelegateV2Method.Event.ID.String() {
			suite.Require().Equal(log.Address, suite.stakingAddr.String())
			event, err := redelegateV2Method.UnpackEvent(log.ToEthereum())
			suite.Require().NoError(err)
			suite.Require().Equal(event.Sender, delAddr)
			suite.Require().Equal(event.ValSrc, valSrc)
			suite.Require().Equal(event.ValDst, valDst)
			suite.Require().Equal(event.Amount.String(), amount.String())
			suite.Require().Equal(event.CompletionTime.Int64(), completionTime)
			existLog = true
		}
	}
	suite.Require().True(existLog)
}

func (suite *PrecompileTestSuite) CheckRedelegateEvents(ctx sdk.Context, valSrc, valDst string, amount *big.Int, completionTime time.Time) {
	existEvent := false
	for _, event := range ctx.EventManager().Events() {
		if event.Type != stakingtypes.EventTypeRedelegate {
			continue
		}
		for _, attr := range event.Attributes {
			if attr.Key == stakingtypes.AttributeKeySrcValidator {
				suite.Require().Equal(attr.Value, valSrc)
			}
			if attr.Key == stakingtypes.AttributeKeyDstValidator {
				suite.Require().Equal(attr.Value, valDst)
			}
			if attr.Key == sdk.AttributeKeyAmount {
				suite.Require().Equal(strings.TrimSuffix(attr.Value, fxtypes.DefaultDenom), amount.String())
			}
			if attr.Key == stakingtypes.AttributeKeyCompletionTime {
				suite.Require().Equal(attr.Value, completionTime.Format(time.RFC3339))
			}
		}
		existEvent = true
	}
	suite.Require().True(existEvent)
}

func (suite *PrecompileTestSuite) CheckUndelegateLogs(logs []*evmtypes.Log, delAddr common.Address, valAddr string, shares, amount *big.Int, completionTime time.Time) {
	undelegateV2Method := precompile.NewUndelegateV2Method(nil)
	existLog := false
	for _, log := range logs {
		if log.Topics[0] == undelegateV2Method.Event.ID.String() {
			suite.Require().Equal(log.Address, suite.stakingAddr.String())
			event, err := undelegateV2Method.UnpackEvent(log.ToEthereum())
			suite.Require().NoError(err)
			suite.Require().Equal(event.Sender, delAddr)
			suite.Require().Equal(event.Validator, valAddr)
			suite.Require().Equal(event.Amount.String(), amount.String())
			suite.Require().Equal(event.CompletionTime.Int64(), completionTime.Unix())
			existLog = true
		}
	}
	suite.Require().True(existLog)
}

func (suite *PrecompileTestSuite) CheckUndeledateEvents(ctx sdk.Context, valAddr string, amount *big.Int, completionTime time.Time) {
	existEvent := false
	for _, event := range ctx.EventManager().Events() {
		if event.Type == stakingtypes.EventTypeUnbond {
			for _, attr := range event.Attributes {
				if attr.Key == stakingtypes.AttributeKeyValidator {
					suite.Require().Equal(attr.Value, valAddr)
					existEvent = true
				}
				if attr.Key == sdk.AttributeKeyAmount {
					suite.Require().Equal(strings.TrimSuffix(attr.Value, fxtypes.DefaultDenom), amount.String())
					existEvent = true
				}
				if attr.Key == stakingtypes.AttributeKeyCompletionTime {
					suite.Require().Equal(attr.Value, completionTime.Format(time.RFC3339))
					existEvent = true
				}
			}
		}
	}
	suite.Require().True(existEvent)
}
