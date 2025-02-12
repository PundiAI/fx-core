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
	require  *require.Assertions
	contract common.Address
	err      error

	ERC20TokenKeeper contract.ERC20TokenKeeper
	evmKeeper        *fxevmkeeper.Keeper
}

func NewERC20Suite(require *require.Assertions, evmKeeper *fxevmkeeper.Keeper) ERC20TokenSuite {
	return ERC20TokenSuite{
		require: require,

		ERC20TokenKeeper: contract.NewERC20TokenKeeper(evmKeeper),
		evmKeeper:        evmKeeper,
	}
}

func (s ERC20TokenSuite) WithContract(addr common.Address) ERC20TokenSuite {
	suite := s
	suite.contract = addr
	return suite
}

func (s ERC20TokenSuite) WithError(err error) ERC20TokenSuite {
	suite := s
	suite.err = err
	return suite
}

func (s ERC20TokenSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s ERC20TokenSuite) DeployERC20Token(ctx sdk.Context, from common.Address, symbol string) common.Address {
	erc20Contract := contract.GetERC20()
	erc20ModuleAddress := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
	initializeArgs := []interface{}{symbol + " Token", symbol, uint8(18), erc20ModuleAddress}
	newContractAddr, err := s.evmKeeper.DeployUpgradableContract(ctx, from, erc20Contract.Address, nil, &erc20Contract.ABI, initializeArgs...)
	s.require.NoError(err)
	return newContractAddr
}

func (s ERC20TokenSuite) Owner(ctx context.Context) common.Address {
	owner, err := s.ERC20TokenKeeper.Owner(ctx, s.contract)
	s.requireError(err)
	return owner
}

func (s ERC20TokenSuite) Name(ctx context.Context) string {
	name, err := s.ERC20TokenKeeper.Name(ctx, s.contract)
	s.requireError(err)
	return name
}

func (s ERC20TokenSuite) Symbol(ctx context.Context) string {
	symbol, err := s.ERC20TokenKeeper.Symbol(ctx, s.contract)
	s.requireError(err)
	return symbol
}

func (s ERC20TokenSuite) Decimals(ctx context.Context) uint8 {
	decimals, err := s.ERC20TokenKeeper.Decimals(ctx, s.contract)
	s.requireError(err)
	return decimals
}

func (s ERC20TokenSuite) TotalSupply(ctx context.Context) *big.Int {
	totalSupply, err := s.ERC20TokenKeeper.TotalSupply(ctx, s.contract)
	s.requireError(err)
	return totalSupply
}

func (s ERC20TokenSuite) BalanceOf(ctx context.Context, address common.Address) *big.Int {
	balance, err := s.ERC20TokenKeeper.BalanceOf(ctx, s.contract, address)
	s.requireError(err)
	return balance
}

func (s ERC20TokenSuite) Allowance(ctx context.Context, owner, spender common.Address) *big.Int {
	allowance, err := s.ERC20TokenKeeper.Allowance(ctx, s.contract, owner, spender)
	s.requireError(err)
	return allowance
}

func (s ERC20TokenSuite) Approve(ctx context.Context, from, spender common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Approve(ctx, s.contract, from, spender, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) Transfer(ctx context.Context, from, recipient common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Transfer(ctx, s.contract, from, recipient, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) TransferFrom(ctx context.Context, from, sender, recipient common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.TransferFrom(ctx, s.contract, from, sender, recipient, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) Mint(ctx context.Context, from, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Mint(ctx, s.contract, from, to, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) MintFromERC20Module(ctx context.Context, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	minter := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
	res, err := s.ERC20TokenKeeper.Mint(ctx, s.contract, minter, to, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) Burn(ctx context.Context, from common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Burn(ctx, s.contract, from, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) TransferOwnership(ctx context.Context, from, newOwner common.Address) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.TransferOwnership(ctx, s.contract, from, newOwner)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) WithdrawSelf(ctx context.Context, from common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Withdraw(ctx, s.contract, from, from, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) Withdraw(ctx context.Context, from, to common.Address, amount *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Withdraw(ctx, s.contract, from, to, amount)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) Deposit(ctx context.Context, from common.Address, value *big.Int) *evmtypes.MsgEthereumTxResponse {
	res, err := s.ERC20TokenKeeper.Deposit(ctx, s.contract, from, value)
	s.requireError(err)
	return res
}

func (s ERC20TokenSuite) OnTest(ctx context.Context, from common.Address, name, symbol string, decimals uint8, totalSupply *big.Int, owner common.Address) {
	s.require.Equal(name, s.Name(ctx))
	s.require.Equal(symbol, s.Symbol(ctx))
	s.require.Equal(decimals, s.Decimals(ctx))
	s.require.Equal(totalSupply.String(), s.TotalSupply(ctx).String())
	s.require.Equal(owner.String(), s.Owner(ctx).String())

	s.require.Equal("0", s.Allowance(ctx, from, from).String())
	s.Approve(ctx, from, from, big.NewInt(100))
	s.require.Equal("100", s.Allowance(ctx, from, from).String())
	s.Mint(ctx, from, from, big.NewInt(200))
	s.TransferFrom(ctx, from, from, from, big.NewInt(100))
	s.Burn(ctx, from, big.NewInt(200))
}
