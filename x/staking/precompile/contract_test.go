package precompile_test

import (
	"bytes"
	"math/big"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
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
)

const (
	TestDelegateV2Name         = "delegateV2"
	TestUndelegateV2Name       = "undelegateV2"
	TestRedelegateV2Name       = "redelegateV2"
	TestWithdrawName           = "withdraw"
	TestDelegationName         = "delegation"
	TestDelegationRewardsName  = "delegationRewards"
	TestAllowanceSharesName    = "allowanceShares"
	TestApproveSharesName      = "approveShares"
	TestTransferSharesName     = "transferShares"
	TestTransferFromSharesName = "transferFromShares"
	TestSlashingInfoName       = "slashingInfo"
	TestValidatorListName      = "validatorList"
)

type PrecompileTestSuite struct {
	helpers.BaseSuite
	testutil.StakingSuite
	signer  *helpers.Signer
	staking common.Address
}

func TestPrecompileTestSuite(t *testing.T) {
	fxtypes.SetConfig(true)
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *PrecompileTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.Commit(100)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.signer = helpers.NewSigner(priv)
	helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18)))))

	stakingContract, err := suite.App.EvmKeeper.DeployContract(suite.Ctx, suite.signer.Address(), contract.MustABIJson(testscontract.StakingTestMetaData.ABI), contract.MustDecodeHex(testscontract.StakingTestMetaData.Bin))
	suite.Require().NoError(err)
	suite.staking = stakingContract

	helpers.AddTestAddr(suite.App, suite.Ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18)))))

	suite.StakingSuite.Init(suite.Require(), suite.Ctx, suite.App.StakingKeeper)
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
	helpers.AddTestAddr(suite.App, suite.Ctx, from.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
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
	helpers.AddTestAddr(suite.App, suite.Ctx, from.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err := stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)

	helpers.AddTestAddr(suite.App, suite.Ctx, to.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err = stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) delegateToFromFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
	helpers.AddTestAddr(suite.App, suite.Ctx, to.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err := stakingkeeper.NewMsgServerImpl(suite.App.StakingKeeper.Keeper).Delegate(suite.Ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)

	helpers.AddTestAddr(suite.App, suite.Ctx, from.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
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

func (suite *PrecompileTestSuite) packTransferRand(val sdk.ValAddress, contractAddr, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
	callABI := precompile.GetABI()
	if bytes.Equal(contractAddr.Bytes(), suite.staking.Bytes()) {
		callABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(TestTransferSharesName, val.String(), to, randShares)
	suite.Require().NoError(err)
	return pack, randShares, nil
}

func (suite *PrecompileTestSuite) packTransferAll(val sdk.ValAddress, contractAddr, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	callABI := precompile.GetABI()
	if bytes.Equal(contractAddr.Bytes(), suite.staking.Bytes()) {
		callABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(TestTransferSharesName, val.String(), to, shares)
	suite.Require().NoError(err)
	return pack, shares, nil
}

func (suite *PrecompileTestSuite) approveFunc(val sdk.ValAddress, owner, spender common.Address, allowance *big.Int) {
	suite.App.StakingKeeper.SetAllowance(suite.Ctx, val, owner.Bytes(), spender.Bytes(), allowance)
}

func (suite *PrecompileTestSuite) packTransferFromRand(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
	suite.approveFunc(val, from, spender, randShares)
	callABI := precompile.GetABI()
	if spender == suite.staking {
		callABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(TestTransferFromSharesName, val.String(), from, to, randShares)
	suite.Require().NoError(err)
	return pack, randShares, nil
}

func (suite *PrecompileTestSuite) packTransferFromAll(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	suite.approveFunc(val, from, spender, shares)
	callABI := precompile.GetABI()
	if spender == suite.staking {
		callABI = contract.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(TestTransferFromSharesName, val.String(), from, to, shares)
	suite.Require().NoError(err)
	return pack, shares, nil
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegation(val sdk.ValAddress, del common.Address) (*big.Int, *big.Int) {
	var res struct {
		Shares *big.Int `abi:"_shares"`
		Amount *big.Int `abi:"_delegateAmount"`
	}
	err := suite.App.EvmKeeper.QueryContract(suite.Ctx, del, precompile.GetAddress(), precompile.GetABI(),
		TestDelegationName, &res, val.String(), del)
	suite.Require().NoError(err)
	return res.Shares, res.Amount
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegateV2(signer *helpers.Signer, val sdk.ValAddress, amt *big.Int) *big.Int {
	helpers.AddTestAddr(suite.App, suite.Ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(amt))))
	pack, err := precompile.GetABI().Pack(TestDelegateV2Name, val.String(), amt)
	suite.Require().NoError(err)

	_, amountBefore := suite.PrecompileStakingDelegation(val, signer.Address())

	res := suite.EthereumTx(signer, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	shares, amount := suite.PrecompileStakingDelegation(val, signer.Address())
	suite.Require().Equal(amt.String(), big.NewInt(0).Sub(amount, amountBefore).String())
	return shares
}

func (suite *PrecompileTestSuite) PrecompileStakingWithdraw(signer *helpers.Signer, val sdk.ValAddress) *big.Int {
	balanceBefore := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	pack, err := precompile.GetABI().Pack(TestWithdrawName, val.String())
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	balanceAfter := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingTransferShares(signer *helpers.Signer, val sdk.ValAddress, receipt common.Address, shares *big.Int) (*big.Int, *big.Int) {
	balanceBefore := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	pack, err := precompile.GetABI().Pack(TestTransferSharesName, val.String(), receipt, shares)
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	signerShares, _ := suite.PrecompileStakingDelegation(val, signer.Address())

	balanceAfter := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return signerShares, rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingUndelegateV2(signer *helpers.Signer, val sdk.ValAddress, shares *big.Int) *big.Int {
	balanceBefore := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	pack, err := precompile.GetABI().Pack(TestUndelegateV2Name, val.String(), shares)
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)

	balanceAfter := suite.App.EvmKeeper.GetEVMDenomBalance(suite.Ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingApproveShares(signer *helpers.Signer, val sdk.ValAddress, spender common.Address, shares *big.Int) {
	pack, err := precompile.GetABI().Pack(TestApproveSharesName, val.String(), spender, shares)
	suite.Require().NoError(err)

	res := suite.EthereumTx(signer, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *PrecompileTestSuite) PrecompileStakingTransferFromShares(signer *helpers.Signer, val sdk.ValAddress, from, receipt common.Address, shares *big.Int) {
	pack, err := precompile.GetABI().Pack(TestTransferFromSharesName, val.String(), from, receipt, shares)
	suite.Require().NoError(err)
	res := suite.EthereumTx(signer, precompile.GetAddress(), big.NewInt(0), pack)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *PrecompileTestSuite) Delegate(val sdk.ValAddress, amount sdkmath.Int, dels ...sdk.AccAddress) {
	for _, del := range dels {
		helpers.AddTestAddr(suite.App, suite.Ctx, del, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)))
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
			suite.Require().Equal(log.Address, precompile.GetAddress().String())

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
			suite.Require().Equal(log.Address, precompile.GetAddress().String())
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
			suite.Require().Equal(log.Address, precompile.GetAddress().String())
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
