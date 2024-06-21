package contract

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type ERC20Pack struct{}

func (e ERC20Pack) Name() (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("name")
	if err != nil {
		return nil, fmt.Errorf("pack name: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) UnpackName(ret []byte) (string, error) {
	var unpackedRet struct{ Value string }
	if err := GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "name", ret); err != nil {
		return "", fmt.Errorf("unpack name: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20Pack) Symbol() (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("symbol")
	if err != nil {
		return nil, fmt.Errorf("pack symbol: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) UnpackSymbol(ret []byte) (string, error) {
	var unpackedRet struct{ Value string }
	if err := GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "symbol", ret); err != nil {
		return "", fmt.Errorf("unpack symbol: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20Pack) Decimals() (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("decimals")
	if err != nil {
		return nil, fmt.Errorf("pack decimals: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) UnpackDecimals(ret []byte) (uint8, error) {
	var unpackedRet struct{ Value uint8 }
	if err := GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "decimals", ret); err != nil {
		return 0, fmt.Errorf("unpack decimals: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20Pack) BalanceOf(account common.Address) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("balanceOf", account)
	if err != nil {
		return nil, fmt.Errorf("pack balanceOf: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) UnpackBalanceOf(ret []byte) (*big.Int, error) {
	var unpackedRet struct{ Value *big.Int }
	if err := GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "balanceOf", ret); err != nil {
		return nil, fmt.Errorf("unpack balanceOf: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20Pack) TotalSupply() (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("pack totalSupply: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) UnpackTotalSupply(ret []byte) (*big.Int, error) {
	var unpackedRet struct{ Value *big.Int }
	if err := GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "totalSupply", ret); err != nil {
		return nil, fmt.Errorf("unpack totalSupply: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20Pack) Approve(spender common.Address, amount *big.Int) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("approve", spender, amount)
	if err != nil {
		return nil, fmt.Errorf("pack approve: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) Allowance(owner, spender common.Address) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("allowance", owner, spender)
	if err != nil {
		return nil, fmt.Errorf("pack allowance: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) TransferFrom(sender, to common.Address, amount *big.Int) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("transferFrom", sender, to, amount)
	if err != nil {
		return nil, fmt.Errorf("pack transferFrom: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) UnpackTransferFrom(ret []byte) (bool, error) {
	var unpackedRet struct{ Value bool }
	if err := GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "transferFrom", ret); err != nil {
		return false, fmt.Errorf("unpack transferFrom: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20Pack) Transfer(to common.Address, amount *big.Int) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("transfer", to, amount)
	if err != nil {
		return nil, fmt.Errorf("pack transfer: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) Burn(account common.Address, amount *big.Int) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("burn", account, amount)
	if err != nil {
		return nil, fmt.Errorf("pack burn: %s", err.Error())
	}
	return data, err
}

func (e ERC20Pack) Mint(account common.Address, amount *big.Int) (data []byte, err error) {
	data, err = GetFIP20().ABI.Pack("mint", account, amount)
	if err != nil {
		return nil, fmt.Errorf("pack mint: %s", err.Error())
	}
	return data, err
}

type ERC20Call struct {
	ERC20Pack
	evm      *vm.EVM
	caller   vm.AccountRef
	contract common.Address
	maxGas   uint64
}

func NewERC20Call(evm *vm.EVM, caller, contract common.Address, maxGas uint64) *ERC20Call {
	defMaxGas := DefaultGasCap
	if maxGas > 0 {
		defMaxGas = maxGas
	}
	return &ERC20Call{
		evm:      evm,
		caller:   vm.AccountRef(caller),
		contract: contract,
		maxGas:   defMaxGas,
	}
}

func (e *ERC20Call) call(data []byte) (ret []byte, err error) {
	ret, _, err = e.evm.Call(e.caller, e.contract, data, e.maxGas, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (e *ERC20Call) staticCall(data []byte) (ret []byte, err error) {
	ret, _, err = e.evm.StaticCall(e.caller, e.contract, data, e.maxGas)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (e *ERC20Call) Burn(account common.Address, amount *big.Int) error {
	data, err := e.ERC20Pack.Burn(account, amount)
	if err != nil {
		return err
	}
	_, err = e.call(data)
	if err != nil {
		return fmt.Errorf("call burn: %s", err.Error())
	}
	return nil
}

func (e *ERC20Call) TransferFrom(from, to common.Address, amount *big.Int) error {
	data, err := e.ERC20Pack.TransferFrom(from, to, amount)
	if err != nil {
		return err
	}
	ret, err := e.call(data)
	if err != nil {
		return fmt.Errorf("call transferFrom: %s", err.Error())
	}
	isSuccess, err := e.UnpackTransferFrom(ret)
	if err != nil {
		return err
	}
	if !isSuccess {
		return errors.New("transferFrom failed")
	}
	return nil
}

func (e *ERC20Call) TotalSupply() (*big.Int, error) {
	data, err := e.ERC20Pack.TotalSupply()
	if err != nil {
		return nil, err
	}
	ret, err := e.staticCall(data)
	if err != nil {
		return nil, fmt.Errorf("StaticCall totalSupply: %s", err.Error())
	}
	return e.ERC20Pack.UnpackTotalSupply(ret)
}
