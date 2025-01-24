package helpers

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

type ERC20TokenSuite struct {
	*ContractBaseSuite
	err              error
	ERC20TokenKeeper contract.ERC20TokenKeeper
	evmKeeper        *fxevmkeeper.Keeper
}

func NewERC20Suite(require *require.Assertions, signer *Signer, evmKeeper *fxevmkeeper.Keeper) ERC20TokenSuite {
	return ERC20TokenSuite{
		ContractBaseSuite: NewContractBaseSuite(require, signer),
		ERC20TokenKeeper:  contract.NewERC20TokenKeeper(evmKeeper),
		evmKeeper:         evmKeeper,
	}
}

func (s ERC20TokenSuite) WithError(err error) ERC20TokenSuite {
	newErc20Suite := s
	newErc20Suite.err = err
	return newErc20Suite
}

func (s ERC20TokenSuite) Error(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s ERC20TokenSuite) DeployERC20Token(ctx sdk.Context, symbol string) common.Address {
	erc20Contract := contract.GetERC20()
	erc20ModuleAddress := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
	initializeArgs := []interface{}{symbol + " Token", symbol, uint8(18), erc20ModuleAddress}
	newContractAddr, err := s.evmKeeper.DeployUpgradableContract(ctx,
		s.signer.Address(), erc20Contract.Address, nil, &erc20Contract.ABI, initializeArgs...)
	s.require.NoError(err)
	s.WithContract(newContractAddr)
	return newContractAddr
}

func (s ERC20TokenSuite) Owner(ctx context.Context) common.Address {
	owner, err := s.ERC20TokenKeeper.Owner(ctx, s.contract)
	s.Error(err)
	return owner
}

func (s ERC20TokenSuite) Name(ctx context.Context) string {
	name, err := s.ERC20TokenKeeper.Name(ctx, s.contract)
	s.Error(err)
	return name
}

func (s ERC20TokenSuite) Symbol(ctx context.Context) string {
	symbol, err := s.ERC20TokenKeeper.Symbol(ctx, s.contract)
	s.Error(err)
	return symbol
}

func (s ERC20TokenSuite) Decimals(ctx context.Context) uint8 {
	decimals, err := s.ERC20TokenKeeper.Decimals(ctx, s.contract)
	s.Error(err)
	return decimals
}

func (s ERC20TokenSuite) TotalSupply(ctx context.Context) *big.Int {
	totalSupply, err := s.ERC20TokenKeeper.TotalSupply(ctx, s.contract)
	s.Error(err)
	return totalSupply
}

func (s ERC20TokenSuite) BalanceOf(ctx context.Context, address common.Address) *big.Int {
	balance, err := s.ERC20TokenKeeper.BalanceOf(ctx, s.contract, address)
	s.Error(err)
	return balance
}

func (s ERC20TokenSuite) Allowance(ctx context.Context, owner, spender common.Address) *big.Int {
	allowance, err := s.ERC20TokenKeeper.Allowance(ctx, s.contract, owner, spender)
	s.Error(err)
	return allowance
}

func (s ERC20TokenSuite) Approve(ctx context.Context, spender common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Approve(ctx, s.contract, s.HexAddress(), spender, amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) Transfer(ctx context.Context, recipient common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Transfer(ctx, s.contract, s.HexAddress(), recipient, amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) TransferFrom(ctx context.Context, sender, recipient common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.TransferFrom(ctx, s.contract, s.HexAddress(), sender, recipient, amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) Mint(ctx context.Context, from, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Mint(ctx, s.contract, from, to, amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) Burn(ctx context.Context, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Burn(ctx, s.contract, s.HexAddress(), amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) TransferOwnership(ctx context.Context, newOwner common.Address) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.TransferOwnership(ctx, s.contract, s.HexAddress(), newOwner)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) WithdrawSelf(ctx context.Context, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Withdraw(ctx, s.contract, s.HexAddress(), s.HexAddress(), amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) Withdraw(ctx context.Context, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Withdraw(ctx, s.contract, s.HexAddress(), to, amount)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) Deposit(ctx context.Context, value *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Deposit(ctx, s.contract, s.HexAddress(), value)
	s.Error(err)
	return res
}

func (s ERC20TokenSuite) OnTest(ctx context.Context, name, symbol string, decimals uint8, totalSupply *big.Int, owner common.Address) {
	s.require.Equal(name, s.Name(ctx))
	s.require.Equal(symbol, s.Symbol(ctx))
	s.require.Equal(decimals, s.Decimals(ctx))
	s.require.Equal(totalSupply.String(), s.TotalSupply(ctx).String())
	s.require.Equal(owner.String(), s.Owner(ctx).String())

	s.require.Equal("0", s.Allowance(ctx, s.HexAddress(), s.HexAddress()).String())
	s.Approve(ctx, s.HexAddress(), big.NewInt(100))
	s.require.Equal("100", s.Allowance(ctx, s.HexAddress(), s.HexAddress()).String())
	s.Mint(ctx, s.signer.Address(), s.signer.Address(), big.NewInt(200))
	s.TransferFrom(ctx, s.HexAddress(), s.HexAddress(), big.NewInt(100))
	s.Burn(ctx, big.NewInt(200))
}
