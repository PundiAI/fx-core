package keeper_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	fxevmtypes "github.com/pundiai/fx-core/v8/x/evm/types"
)

func (s *KeeperTestSuite) ethereumTx(signer *helpers.Signer, to *common.Address, data []byte, value *big.Int, gasLimit uint64) (*types.MsgEthereumTxResponse, error) {
	chanId := s.App.EvmKeeper.ChainID()
	s.Equal(fxtypes.EIP155ChainID(s.Ctx.ChainID()), chanId)
	if value == nil {
		value = big.NewInt(0)
	}

	nonce := s.App.EvmKeeper.GetNonce(s.Ctx, signer.Address())
	tx := types.NewTx(
		chanId,
		nonce,
		to,
		value,
		gasLimit,
		big.NewInt(fxtypes.DefaultGasPrice),
		nil,
		nil,
		data,
		nil,
	)
	tx.From = signer.Address().Bytes()
	s.NoError(tx.Sign(ethtypes.LatestSignerForChainID(chanId), signer))

	return s.App.EvmKeeper.EthereumTx(s.Ctx, tx)
}

func (s *KeeperTestSuite) TestKeeper_EthereumTx_Data() {
	signer := s.NewSigner()
	erc20Suite := s.NewERC20TokenSuite()
	erc20Suite.WithSigner(signer)

	contractAddr := erc20Suite.DeployERC20Token(s.Ctx, "TEST")
	erc20Suite.Mint(s.Ctx, erc20Suite.HexAddress(), big.NewInt(100))

	gasLimit := uint64(71000)

	totalSupplyBefore := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	gasPrice := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx).Int64()
	mintAmount := gasPrice * int64(gasLimit)
	s.MintFeeCollector(helpers.NewStakingCoins(mintAmount, 0))

	recipient := helpers.GenHexAddress()
	amount := big.NewInt(10)
	data, err := helpers.PackERC20Mint(recipient, amount)
	s.Require().NoError(err)

	res, err := s.ethereumTx(signer, &contractAddr, data, nil, gasLimit)
	s.Require().NoError(err)
	s.Require().False(res.Failed(), res)

	// s.Commit() call evm endBlock to setter the EVM Virtual balance
	s.Commit()

	refundAmount := gasPrice * int64(gasLimit-res.GasUsed)
	s.BurnEvmRefundFee(erc20Suite.AccAddress(), helpers.NewStakingCoins(refundAmount, 0))

	totalSupplyAfter := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	s.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	s.Require().Equal(amount, erc20Suite.BalanceOf(s.Ctx, recipient))
}

func (s *KeeperTestSuite) TestKeeper_EthereumTx_Value() {
	recipient := helpers.GenHexAddress()
	amount := big.NewInt(10)

	gasLimit := uint64(71000)

	signer := s.AddTestSigner()
	totalSupplyBefore := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	mintAmount := sdkmath.NewInt(s.App.FeeMarketKeeper.GetBaseFee(s.Ctx).Int64() * int64(gasLimit))
	s.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, mintAmount)))

	res, err := s.ethereumTx(signer, &recipient, nil, amount, gasLimit)
	s.Require().NoError(err)
	s.Require().False(res.Failed(), res)

	// s.Commit() call evm endBlock to setter the EVM Virtual balance
	s.Commit()

	refundAmount := sdkmath.NewInt(s.App.FeeMarketKeeper.GetBaseFee(s.Ctx).Int64() * int64(gasLimit-res.GasUsed))
	s.BurnEvmRefundFee(signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, refundAmount)))

	totalSupplyAfter := s.App.BankKeeper.GetSupply(s.Ctx, fxtypes.DefaultDenom)
	s.Require().Equal(totalSupplyBefore.String(), totalSupplyAfter.String())

	balance := s.App.EvmKeeper.GetEVMDenomBalance(s.Ctx, recipient)
	s.Equal(balance, amount)

	balance = s.App.BankKeeper.GetBalance(s.Ctx, recipient.Bytes(), fxtypes.DefaultDenom).Amount.BigInt()
	s.Equal(balance, amount)
}

func (s *KeeperTestSuite) TestKeeper_CallContract() {
	erc20Suite := s.NewERC20TokenSuite()
	contract := erc20Suite.DeployERC20Token(s.Ctx, "USD")

	amount := big.NewInt(100)
	recipient := erc20Suite.HexAddress()
	data, err := helpers.PackERC20Mint(recipient, amount)
	s.Require().NoError(err)

	// failed: not authorized
	_, err = s.App.EvmKeeper.CallContract(s.Ctx, &fxevmtypes.MsgCallContract{
		Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ContractAddress: contract.String(),
		Data:            common.Bytes2Hex(data),
	})
	s.Require().EqualError(err, "Ownable: caller is not the owner: evm transaction execution failed")

	// transfer erc20 token owner to evm module
	evmModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName))
	erc20Suite.TransferOwnership(s.Ctx, evmModuleAddr)

	// success
	_, err = s.App.EvmKeeper.CallContract(s.Ctx, &fxevmtypes.MsgCallContract{
		Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ContractAddress: contract.String(),
		Data:            common.Bytes2Hex(data),
	})
	s.Require().NoError(err)

	s.Equal(amount, erc20Suite.BalanceOf(s.Ctx, recipient))
}
