package testutil

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
)

type ERC20Suite struct {
	contract.ERC20ABI
	*EVMSuite
}

func NewERC20Suite(evmSuite *EVMSuite) *ERC20Suite {
	return &ERC20Suite{
		ERC20ABI: contract.NewERC20ABI(),
		EVMSuite: evmSuite,
	}
}

func (s *ERC20Suite) Call(method string, res interface{}, args ...interface{}) {
	s.EVMSuite.Call(s.ABI, method, res, args...)
}

func (s *ERC20Suite) Send(method string, args ...interface{}) *evmtypes.MsgEthereumTxResponse {
	return s.EVMSuite.Send(s.ABI, method, args...)
}

func (s *ERC20Suite) Deploy(symbol string) common.Address {
	data := contract.GetFIP20().Bin
	nonce := s.evmKeeper.GetNonce(s.ctx, s.signer.Address())
	msg := core.Message{
		To:                nil,
		From:              s.signer.Address(),
		Nonce:             nonce,
		Value:             big.NewInt(0),
		GasLimit:          1_700_000,
		GasPrice:          big.NewInt(500 * 1e9),
		GasFeeCap:         nil,
		GasTipCap:         nil,
		Data:              data,
		AccessList:        nil,
		SkipAccountChecks: false,
	}
	rsp, err := s.evmKeeper.ApplyMessage(s.ctx, msg, nil, true)
	s.Require().NoError(err)
	s.Require().False(rsp.Failed(), rsp.VmError)
	s.Require().Equal(uint64(1_600_010), rsp.GasUsed)
	addr := crypto.CreateAddress(s.signer.Address(), nonce)
	s.contractAddr = addr
	s.Initialize(symbol, true)
	return addr
}

func (s *ERC20Suite) Initialize(symbol string, result bool) {
	response := s.Send("initialize", symbol+" Token", symbol, uint8(18), helpers.GenHexAddress())
	s.Require().Equal(response.Failed(), !result)
}

func (s *ERC20Suite) Owner() common.Address {
	var ownerRes struct {
		Value common.Address
	}
	s.Call("owner", &ownerRes)
	return ownerRes.Value
}

func (s *ERC20Suite) Name() string {
	var nameRes struct {
		Value string
	}
	s.Call("name", &nameRes)
	return nameRes.Value
}

func (s *ERC20Suite) Symbol() string {
	var symbolRes struct {
		Value string
	}
	s.Call("symbol", &symbolRes)
	return symbolRes.Value
}

func (s *ERC20Suite) Decimals() uint8 {
	var decimalsRes struct {
		Value uint8
	}
	s.Call("decimals", &decimalsRes)
	return decimalsRes.Value
}

func (s *ERC20Suite) TotalSupply() *big.Int {
	var totalSupplyRes struct {
		Value *big.Int
	}
	s.Call("totalSupply", &totalSupplyRes)
	return totalSupplyRes.Value
}

func (s *ERC20Suite) BalanceOf(account common.Address) *big.Int {
	var balanceRes struct {
		Value *big.Int
	}
	s.Call("balanceOf", &balanceRes, account)
	return balanceRes.Value
}

func (s *ERC20Suite) Allowance(owner, spender common.Address) *big.Int {
	var allowanceRes struct {
		Value *big.Int
	}
	s.Call("allowance", &allowanceRes, owner, spender)
	return allowanceRes.Value
}

func (s *ERC20Suite) Approve(spender common.Address, amount *big.Int, result bool) {
	before := s.Allowance(s.signer.Address(), spender)
	response := s.Send("approve", spender, amount)
	after := s.Allowance(s.signer.Address(), spender)
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Add(before, amount))
	}
}

func (s *ERC20Suite) Transfer(recipient common.Address, amount *big.Int, result bool) {
	before := s.BalanceOf(s.signer.Address())
	response := s.Send("transfer", recipient, amount)
	after := s.BalanceOf(s.signer.Address())
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Sub(before, amount))
	}
}

func (s *ERC20Suite) TransferFrom(sender, recipient common.Address, amount *big.Int, result bool) {
	before := s.BalanceOf(recipient)
	response := s.Send("transferFrom", sender, recipient, amount)
	after := s.BalanceOf(recipient)
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Add(before, amount))
	}
}

func (s *ERC20Suite) Mint(to common.Address, amount *big.Int, result bool) {
	before := s.TotalSupply()
	response := s.Send("mint", to, amount)
	after := s.TotalSupply()
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Add(before, amount))
	}
}

func (s *ERC20Suite) Burn(from common.Address, amount *big.Int, result bool) {
	before := s.BalanceOf(from)
	response := s.Send("burn", from, amount)
	after := s.BalanceOf(from)
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Sub(before, amount))
	}
}

func (s *ERC20Suite) TransferOwnership(newOwner common.Address, result bool) {
	before := s.Owner()
	response := s.Send("transferOwnership", newOwner)
	after := s.Owner()
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().NotEqual(before, after)
	}
}

func (s *ERC20Suite) WithdrawSelf(amount *big.Int, result bool) {
	before := s.BalanceOf(s.signer.Address())
	response := s.Send("withdraw", amount)
	after := s.BalanceOf(s.signer.Address())
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Sub(before, amount))
	}
}

func (s *ERC20Suite) Withdraw(to common.Address, amount *big.Int, result bool) {
	before := s.BalanceOf(to)
	response := s.Send("withdraw0", to, amount)
	after := s.BalanceOf(to)
	s.Require().Equal(response.Failed(), !result)
	if result {
		s.Require().Equal(after, before)
	}
}

func (s *ERC20Suite) Deposit(value *big.Int, result bool) {
	data, err := s.ABI.Pack("deposit")
	s.Require().NoError(err)

	msg := core.Message{
		To:                &s.contractAddr,
		From:              s.signer.Address(),
		Nonce:             s.evmKeeper.GetNonce(s.ctx, s.signer.Address()),
		Value:             value,
		GasLimit:          80000,
		GasPrice:          big.NewInt(500 * 1e9),
		GasFeeCap:         nil,
		GasTipCap:         nil,
		Data:              data,
		AccessList:        nil,
		SkipAccountChecks: false,
	}

	before := s.BalanceOf(s.signer.Address())
	rsp, err := s.evmKeeper.ApplyMessage(s.ctx, msg, evmtypes.NewNoOpTracer(), true)
	after := s.BalanceOf(s.signer.Address())
	s.Require().NoError(err)
	s.Require().Equal(rsp.Failed(), !result)
	if result {
		s.Require().Equal(after, new(big.Int).Add(before, value))
	}
}

func (s *ERC20Suite) OnTest(name, symbol string, decimals uint8, totalSupply *big.Int, owner common.Address) {
	s.Require().Equal(name, s.Name())
	s.Require().Equal(symbol, s.Symbol())
	s.Require().Equal(decimals, s.Decimals())
	s.Require().Equal(totalSupply.String(), s.TotalSupply().String())
	s.Require().Equal(owner.String(), s.Owner().String())

	s.Require().Equal("0", s.Allowance(s.HexAddr(), s.HexAddr()).String())
	s.Approve(s.HexAddr(), big.NewInt(100), true)
	s.Require().Equal("100", s.Allowance(s.HexAddr(), s.HexAddr()).String())

	newSigner := helpers.NewSigner(helpers.NewEthPrivKey())

	s.Mint(s.signer.Address(), big.NewInt(200), true)
	s.Transfer(newSigner.Address(), big.NewInt(100), true)
	s.TransferFrom(s.signer.Address(), newSigner.Address(), big.NewInt(100), true)
	s.Burn(newSigner.Address(), big.NewInt(200), true)
}
