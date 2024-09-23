package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func (s *KeeperTestSuite) TestKeeper_EthereumTx_Data() {
	erc20Suite := s.NewERC20Suite()
	contractAddr := erc20Suite.Deploy("TEST")

	erc20Suite.Mint(erc20Suite.HexAddr(), big.NewInt(100), true)

	gasLimit := uint64(71000)

	totalSupplyBefore := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	gasPrice := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx).Int64()
	mintAmount := gasPrice * int64(gasLimit)
	s.MintFeeCollector(helpers.NewStakingCoins(mintAmount, 0))

	recipient := helpers.GenHexAddress()
	amount := big.NewInt(10)
	data, err := erc20Suite.PackTransfer(recipient, amount)
	s.Require().NoError(err)

	res, err := s.EVMSuite.EthereumTx(&contractAddr, data, nil, gasLimit)
	s.Require().NoError(err)
	s.Require().False(res.Failed(), res)

	// s.Commit() call evm endBlock to setter the EVM Virtual balance
	s.Commit()

	refundAmount := gasPrice * int64(gasLimit-res.GasUsed)
	s.BurnEvmRefundFee(erc20Suite.AccAddr(), helpers.NewStakingCoins(refundAmount, 0))

	totalSupplyAfter := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	s.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	s.Require().Equal(amount, erc20Suite.BalanceOf(recipient))
}

func (s *KeeperTestSuite) TestKeeper_EthereumTx_Value() {
	recipient := helpers.GenHexAddress()
	amount := big.NewInt(10)

	gasLimit := uint64(71000)

	totalSupplyBefore := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	mintAmount := sdkmath.NewInt(s.App.FeeMarketKeeper.GetBaseFee(s.Ctx).Int64() * int64(gasLimit))
	s.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, mintAmount)))

	res, err := s.EVMSuite.EthereumTx(&recipient, nil, amount, gasLimit)
	s.Require().NoError(err)
	s.Require().False(res.Failed(), res)

	// s.Commit() call evm endBlock to setter the EVM Virtual balance
	s.Commit()

	refundAmount := sdkmath.NewInt(s.App.FeeMarketKeeper.GetBaseFee(s.Ctx).Int64() * int64(gasLimit-res.GasUsed))
	s.BurnEvmRefundFee(s.EVMSuite.AccAddr(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, refundAmount)))

	totalSupplyAfter := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	s.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	balance := s.App.EvmKeeper.GetEVMDenomBalance(s.Ctx, recipient)
	s.Equal(balance, amount)

	balance = s.App.BankKeeper.GetBalance(s.Ctx, recipient.Bytes(), fxtypes.DefaultDenom).Amount.BigInt()
	s.Equal(balance, amount)
}

func (s *KeeperTestSuite) TestKeeper_CallContract() {
	s.EVMSuite.DeployUpgradableERC20Logic("USD")

	erc20Suite := s.NewERC20Suite()

	amount := big.NewInt(100)
	recipient := erc20Suite.HexAddr()
	data, err := erc20Suite.PackMint(recipient, amount)
	s.Require().NoError(err)

	account := s.App.AccountKeeper.GetAccount(s.Ctx, authtypes.NewModuleAddress(types.ModuleName))
	s.Require().NotNil(account)

	// failed: not authorized
	err = s.EVMSuite.CallContract(data)
	s.Require().EqualError(err, "Ownable: caller is not the owner: evm transaction execution failed")

	// transfer erc20 token owner to evm module
	evmModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName))
	erc20Suite.TransferOwnership(evmModuleAddr, true)

	// success
	err = s.EVMSuite.CallContract(data)
	s.Require().NoError(err)

	s.Equal(amount, erc20Suite.BalanceOf(recipient))
}
