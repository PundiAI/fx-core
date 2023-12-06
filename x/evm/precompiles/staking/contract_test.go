package staking_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/server/config"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v6/app"
	testscontract "github.com/functionx/fx-core/v6/tests/contract"
	"github.com/functionx/fx-core/v6/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v6/types"
	"github.com/functionx/fx-core/v6/x/evm/precompiles/staking"
)

const (
	StakingTestDelegateName           = "delegate"
	StakingTestUndelegateName         = "undelegate"
	StakingTestRedelegateName         = "redelegate"
	StakingTestWithdrawName           = "withdraw"
	StakingTestDelegationName         = "delegation"
	StakingTestDelegationRewardsName  = "delegationRewards"
	StakingTestAllowanceSharesName    = "allowanceShares"
	StakingTestApproveSharesName      = "approveShares"
	StakingTestTransferSharesName     = "transferShares"
	StakingTestTransferFromSharesName = "transferFromShares"
)

type PrecompileTestSuite struct {
	suite.Suite
	ctx     sdk.Context
	app     *app.App
	signer  *helpers.Signer
	staking common.Address
}

func TestPrecompileTestSuite(t *testing.T) {
	fxtypes.SetConfig(true)
	suite.Run(t, new(PrecompileTestSuite))
}

// Test helpers
func (suite *PrecompileTestSuite) SetupTest() {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(suite.T(), err)
	suite.signer = helpers.NewSigner(priv)

	set, accs, balances := helpers.GenerateGenesisValidator(tmrand.Intn(10)+3, nil)
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), set, accs, balances...)

	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight() + 1,
		ChainID:         fxtypes.ChainId(),
		ProposerAddress: set.Proposer.Address,
		Time:            time.Now().UTC(),
	})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	for _, validator := range set.Validators {
		signingInfo := slashingtypes.NewValidatorSigningInfo(
			validator.Address.Bytes(),
			suite.ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, validator.Address.Bytes(), signingInfo)
	}

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18)))))
	stakingContract, err := suite.app.EvmKeeper.DeployContract(suite.ctx, suite.signer.Address(), fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI), fxtypes.MustDecodeHex(testscontract.StakingTestMetaData.Bin))
	suite.Require().NoError(err)
	suite.staking = stakingContract

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18)))))
}

func (suite *PrecompileTestSuite) PackEthereumTx(signer *helpers.Signer, contract common.Address, amount *big.Int, data []byte) (*evmtypes.MsgEthereumTx, error) {
	fromAddr := signer.Address()
	value := hexutil.Big(*amount)
	args, err := json.Marshal(&evmtypes.TransactionArgs{To: &contract, From: &fromAddr, Data: (*hexutil.Bytes)(&data), Value: &value})
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	res, err := evmtypes.NewQueryClient(queryHelper).EstimateGas(sdk.WrapSDKContext(suite.ctx),
		&evmtypes.EthCallRequest{
			Args:    args,
			GasCap:  config.DefaultGasCap,
			ChainId: suite.app.EvmKeeper.ChainID().Int64(),
		},
	)
	if err != nil {
		return nil, err
	}

	ethTx := evmtypes.NewTx(
		fxtypes.EIP155ChainID(),
		suite.app.EvmKeeper.GetNonce(suite.ctx, signer.Address()),
		&contract,
		amount,
		res.Gas,
		nil,
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx),
		big.NewInt(1),
		data,
		nil,
	)
	ethTx.From = signer.Address().Hex()
	err = ethTx.Sign(ethtypes.LatestSignerForChainID(fxtypes.EIP155ChainID()), signer)
	return ethTx, err
}

func (suite *PrecompileTestSuite) Commit() {
	header := suite.ctx.BlockHeader()

	suite.app.EndBlock(abci.RequestEndBlock{Height: header.Height})
	suite.app.Commit()
	// begin block
	header.Time = time.Now().UTC()
	header.Height += 1

	vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	infos := make([]abci.VoteInfo, 0, len(vals))
	for _, val := range vals {
		addr, err := val.GetConsAddr()
		suite.Require().NoError(err)
		infos = append(infos, abci.VoteInfo{Validator: abci.Validator{Address: addr, Power: 100}})
	}

	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
		LastCommitInfo: abci.LastCommitInfo{
			Votes: infos,
		},
	})
	suite.ctx = suite.app.NewContext(false, header)
}

func (suite *PrecompileTestSuite) RandSigner() *helpers.Signer {
	privKey := helpers.NewEthPrivKey()
	// helpers.AddTestAddr(suite.app, suite.ctx, privKey.PubKey().Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))))
	signer := helpers.NewSigner(privKey)
	suite.app.AccountKeeper.SetAccount(suite.ctx, suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.AccAddress()))
	return signer
}

func (suite *PrecompileTestSuite) delegateFromFunc(val sdk.ValAddress, from, _ common.Address, delAmount sdkmath.Int) {
	helpers.AddTestAddr(suite.app, suite.ctx, from.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) undelegateToFunc(val sdk.ValAddress, _, to common.Address, _ sdkmath.Int) {
	toDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, to.Bytes(), val)
	suite.Require().True(found)
	_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) delegateFromToFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
	helpers.AddTestAddr(suite.app, suite.ctx, from.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)

	helpers.AddTestAddr(suite.app, suite.ctx, to.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) delegateToFromFunc(val sdk.ValAddress, from, to common.Address, delAmount sdkmath.Int) {
	helpers.AddTestAddr(suite.app, suite.ctx, to.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(to.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)

	helpers.AddTestAddr(suite.app, suite.ctx, from.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, delAmount)))
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Delegate(sdk.WrapSDKContext(suite.ctx), &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(from.Bytes()).String(),
		ValidatorAddress: val.String(),
		Amount:           sdk.NewCoin(fxtypes.DefaultDenom, delAmount),
	})
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) undelegateFromToFunc(val sdk.ValAddress, from, to common.Address, _ sdkmath.Int) {
	fromDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, from.Bytes(), val)
	suite.Require().True(found)
	_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, from.Bytes(), val, fromDel.Shares)
	suite.Require().NoError(err)

	toDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, to.Bytes(), val)
	suite.Require().True(found)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) undelegateToFromFunc(val sdk.ValAddress, from, to common.Address, _ sdkmath.Int) {
	toDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, to.Bytes(), val)
	suite.Require().True(found)
	_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, to.Bytes(), val, toDel.Shares)
	suite.Require().NoError(err)

	fromDel, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, from.Bytes(), val)
	suite.Require().True(found)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, from.Bytes(), val, fromDel.Shares)
	suite.Require().NoError(err)
}

func (suite *PrecompileTestSuite) packTransferRand(val sdk.ValAddress, contract, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
	callFunc := staking.TransferSharesMethodName
	callABI := staking.GetABI()
	if bytes.Equal(contract.Bytes(), suite.staking.Bytes()) {
		callFunc = StakingTestTransferSharesName
		callABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(callFunc, val.String(), to, randShares)
	suite.Require().NoError(err)
	return pack, randShares, nil
}

func (suite *PrecompileTestSuite) packTransferAll(val sdk.ValAddress, contract, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	callFunc := staking.TransferSharesMethodName
	callABI := staking.GetABI()
	if bytes.Equal(contract.Bytes(), suite.staking.Bytes()) {
		callFunc = StakingTestTransferSharesName
		callABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(callFunc, val.String(), to, shares)
	suite.Require().NoError(err)
	return pack, shares, nil
}

func (suite *PrecompileTestSuite) approveFunc(val sdk.ValAddress, owner, spender common.Address, allowance *big.Int) {
	suite.app.StakingKeeper.SetAllowance(suite.ctx, val, owner.Bytes(), spender.Bytes(), allowance)
}

func (suite *PrecompileTestSuite) packTransferFromRand(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	randShares := big.NewInt(0).Sub(shares, big.NewInt(0).Mul(big.NewInt(tmrand.Int63n(900)+100), big.NewInt(1e18)))
	suite.approveFunc(val, from, spender, randShares)
	callFunc := staking.TransferFromSharesMethodName
	callABI := staking.GetABI()
	if spender == suite.staking {
		callFunc = StakingTestTransferFromSharesName
		callABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(callFunc, val.String(), from, to, randShares)
	suite.Require().NoError(err)
	return pack, randShares, nil
}

func (suite *PrecompileTestSuite) packTransferFromAll(val sdk.ValAddress, spender, from, to common.Address, shares *big.Int) ([]byte, *big.Int, []string) {
	suite.approveFunc(val, from, spender, shares)
	callFunc := staking.TransferFromSharesMethodName
	callABI := staking.GetABI()
	if spender == suite.staking {
		callFunc = StakingTestTransferFromSharesName
		callABI = fxtypes.MustABIJson(testscontract.StakingTestMetaData.ABI)
	}
	pack, err := callABI.Pack(callFunc, val.String(), from, to, shares)
	suite.Require().NoError(err)
	return pack, shares, nil
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegation(val sdk.ValAddress, del common.Address) (*big.Int, *big.Int) {
	var res struct {
		Shares *big.Int `abi:"_shares"`
		Amount *big.Int `abi:"_delegateAmount"`
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, del, staking.GetAddress(), staking.GetABI(),
		staking.DelegationMethodName, &res, val.String(), del)
	suite.Require().NoError(err)
	return res.Shares, res.Amount
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegate(signer *helpers.Signer, val sdk.ValAddress, amt *big.Int) *big.Int {
	helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(amt))))
	pack, err := staking.GetABI().Pack(staking.DelegateMethodName, val.String())
	suite.Require().NoError(err)
	tx, err := suite.PackEthereumTx(signer, staking.GetAddress(), amt, pack)
	suite.Require().NoError(err)

	_, amountBefore := suite.PrecompileStakingDelegation(val, signer.Address())

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	shares, amount := suite.PrecompileStakingDelegation(val, signer.Address())
	suite.Require().Equal(amt.String(), big.NewInt(0).Sub(amount, amountBefore).String())
	return shares
}

func (suite *PrecompileTestSuite) PrecompileStakingWithdraw(signer *helpers.Signer, val sdk.ValAddress) *big.Int {
	balanceBefore := suite.app.EvmKeeper.GetBalance(suite.ctx, signer.Address())
	pack, err := staking.GetABI().Pack(staking.WithdrawMethodName, val.String())
	suite.Require().NoError(err)
	tx, err := suite.PackEthereumTx(signer, staking.GetAddress(), big.NewInt(0), pack)
	suite.Require().NoError(err)
	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
	balanceAfter := suite.app.EvmKeeper.GetBalance(suite.ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingTransferShares(signer *helpers.Signer, val sdk.ValAddress, receipt common.Address, shares *big.Int) (*big.Int, *big.Int) {
	balanceBefore := suite.app.EvmKeeper.GetBalance(suite.ctx, signer.Address())
	pack, err := staking.GetABI().Pack(staking.TransferSharesMethodName, val.String(), receipt, shares)
	suite.Require().NoError(err)
	tx, err := suite.PackEthereumTx(signer, staking.GetAddress(), big.NewInt(0), pack)
	suite.Require().NoError(err)
	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	signerShares, _ := suite.PrecompileStakingDelegation(val, signer.Address())

	balanceAfter := suite.app.EvmKeeper.GetBalance(suite.ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return signerShares, rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingUndelegate(signer *helpers.Signer, val sdk.ValAddress, shares *big.Int) *big.Int {
	balanceBefore := suite.app.EvmKeeper.GetBalance(suite.ctx, signer.Address())
	pack, err := staking.GetABI().Pack(staking.UndelegateMethodName, val.String(), shares)
	suite.Require().NoError(err)
	tx, err := suite.PackEthereumTx(signer, staking.GetAddress(), big.NewInt(0), pack)
	suite.Require().NoError(err)
	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
	balanceAfter := suite.app.EvmKeeper.GetBalance(suite.ctx, signer.Address())
	rewards := big.NewInt(0).Sub(balanceAfter, balanceBefore)
	return rewards
}

func (suite *PrecompileTestSuite) PrecompileStakingApproveShares(signer *helpers.Signer, val sdk.ValAddress, spender common.Address, shares *big.Int) {
	pack, err := staking.GetABI().Pack(staking.ApproveSharesMethodName, val.String(), spender, shares)
	suite.Require().NoError(err)
	tx, err := suite.PackEthereumTx(signer, staking.GetAddress(), big.NewInt(0), pack)
	suite.Require().NoError(err)
	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *PrecompileTestSuite) PrecompileStakingTransferFromShares(signer *helpers.Signer, val sdk.ValAddress, from, receipt common.Address, shares *big.Int) {
	pack, err := staking.GetABI().Pack(staking.TransferFromSharesMethodName, val.String(), from, receipt, shares)
	suite.Require().NoError(err)
	tx, err := suite.PackEthereumTx(signer, staking.GetAddress(), big.NewInt(0), pack)
	suite.Require().NoError(err)
	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *PrecompileTestSuite) Delegate(val sdk.ValAddress, amount sdkmath.Int, dels ...sdk.AccAddress) {
	for _, del := range dels {
		helpers.AddTestAddr(suite.app, suite.ctx, del, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)))
		validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val)
		suite.Require().True(found)
		_, err := suite.app.StakingKeeper.Delegate(suite.ctx, del, amount, stakingtypes.Unbonded, validator, true)
		suite.Require().NoError(err)
	}
}

func (suite *PrecompileTestSuite) Redelegate(valSrc, valDest sdk.ValAddress, del sdk.AccAddress, shares sdk.Dec) {
	_, err := suite.app.StakingKeeper.BeginRedelegation(suite.ctx, del, valSrc, valDest, shares)
	suite.Require().NoError(err)
}
