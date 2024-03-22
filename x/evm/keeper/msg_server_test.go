package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	fxevmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

func (suite *KeeperTestSuite) TestKeeper_EthereumTx() {
	contractAddr := suite.DeployERC20Contract()

	recipient := helpers.GenerateAddress()
	amount := big.NewInt(10)
	data, err := contract.GetFIP20().ABI.Pack("transfer", recipient, amount)
	suite.Require().NoError(err)

	chanId := suite.app.EvmKeeper.ChainID()
	suite.Equal(fxtypes.EIP155ChainID(), chanId)

	gasLimit := uint64(71000)

	totalSupplyBefore := suite.app.BankKeeper.GetSupply(suite.ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	mintAmount := sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64() * int64(gasLimit))
	suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, mintAmount)))

	tx := types.NewTx(
		chanId,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		&contractAddr,
		big.NewInt(0),
		gasLimit,
		big.NewInt(500*1e9),
		nil,
		nil,
		data,
		nil,
	)
	tx.From = suite.signer.Address().String()
	suite.NoError(tx.Sign(ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID()), suite.signer))

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res)

	refundAmount := sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64() * int64(gasLimit-res.GasUsed))
	suite.BurnEvmRefundFee(suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, refundAmount)))

	totalSupplyAfter := suite.app.BankKeeper.GetSupply(suite.ctx, fxtypes.DefaultDenom)
	suite.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	suite.Equal(amount, suite.BalanceOf(contractAddr, recipient))
}

func (suite *KeeperTestSuite) TestKeeper_EthereumTx2() {
	recipient := helpers.GenerateAddress()
	amount := big.NewInt(10)

	chanId := suite.app.EvmKeeper.ChainID()
	suite.Equal(fxtypes.EIP155ChainID(), chanId)

	gasLimit := uint64(71000)

	totalSupplyBefore := suite.app.BankKeeper.GetSupply(suite.ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	mintAmount := sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64() * int64(gasLimit))
	suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, mintAmount)))

	tx := types.NewTx(
		chanId,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		&recipient,
		amount,
		gasLimit,
		big.NewInt(500*1e9),
		nil,
		nil,
		nil,
		nil,
	)
	tx.From = suite.signer.Address().String()
	suite.NoError(tx.Sign(ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID()), suite.signer))

	res, err := suite.app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(suite.ctx), tx)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res)

	refundAmount := sdkmath.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64() * int64(gasLimit-res.GasUsed))
	suite.BurnEvmRefundFee(suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, refundAmount)))

	totalSupplyAfter := suite.app.BankKeeper.GetSupply(suite.ctx, fxtypes.DefaultDenom)
	suite.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	balance := suite.app.EvmKeeper.GetBalance(suite.ctx, recipient)
	suite.Equal(balance, amount)

	balance = suite.app.BankKeeper.GetBalance(suite.ctx, recipient.Bytes(), fxtypes.DefaultDenom).Amount.BigInt()
	suite.Equal(balance, amount)
}

func (suite *KeeperTestSuite) TestKeeper_CallContract() {
	erc20 := contract.GetFIP20()
	initializeArgs := []interface{}{"FunctionX USD", "fxUSD", uint8(18), suite.app.Erc20Keeper.ModuleAddress()}

	// deploy contract
	contract, err := suite.app.EvmKeeper.DeployUpgradableContract(suite.ctx, suite.signer.Address(), erc20.Address, nil, &erc20.ABI, initializeArgs...)
	suite.NoError(err)
	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, suite.signer.Address().Bytes())
	suite.NoError(err)
	contractAddr := crypto.CreateAddress(suite.signer.Address(), nonce-1)
	suite.Equal(contractAddr, contract)
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
	args, err := erc20.ABI.Pack("mint", suite.signer.Address(), amount)
	suite.Require().NoError(err)

	failMsg := &fxevmtypes.MsgCallContract{
		Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ContractAddress: contractAddr.String(),
		Data:            common.Bytes2Hex(args),
	}
	_, err = suite.app.EvmKeeper.CallContract(sdk.WrapSDKContext(suite.ctx), failMsg)
	suite.Require().Error(err)
	// transferOwnership
	_, err = suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contract, nil, erc20.ABI, "transferOwnership", common.BytesToAddress(suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)))
	suite.Require().NoError(err)
	// CallContract
	msg := &fxevmtypes.MsgCallContract{
		Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ContractAddress: contractAddr.String(),
		Data:            common.Bytes2Hex(args),
	}
	_, err = suite.app.EvmKeeper.CallContract(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Equal(amount, suite.BalanceOf(contract, suite.signer.Address()))
}
