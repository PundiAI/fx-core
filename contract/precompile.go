package contract

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type PrecompileMethod interface {
	GetMethodId() []byte
	RequiredGas() uint64
	IsReadonly() bool
	Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error)
}

func EmitEvent(evm *vm.EVM, address common.Address, data []byte, topics []common.Hash) {
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     address,
		Topics:      topics,
		Data:        data,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
}

type ERC20Call struct {
	abi      abi.ABI
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
		abi:      GetWFX().ABI,
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
	data, err := e.abi.Pack("burn", account, amount)
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
	data, err := e.abi.Pack("transferFrom", from, to, amount)
	if err != nil {
		return err
	}
	ret, err := e.call(data)
	if err != nil {
		return fmt.Errorf("call transferFrom: %s", err.Error())
	}
	var unpackedRet struct{ Value bool }
	if err = e.abi.UnpackIntoInterface(&unpackedRet, "transferFrom", ret); err != nil {
		return fmt.Errorf("unpack transferFrom: %s", err.Error())
	}
	if !unpackedRet.Value {
		return errors.New("transferFrom failed")
	}
	return nil
}

func (e *ERC20Call) TotalSupply() (*big.Int, error) {
	data, err := e.abi.Pack("totalSupply")
	if err != nil {
		return nil, err
	}
	ret, err := e.staticCall(data)
	if err != nil {
		return nil, fmt.Errorf("StaticCall totalSupply: %s", err.Error())
	}
	var unpackedRet struct{ Value *big.Int }
	if err = e.abi.UnpackIntoInterface(&unpackedRet, "totalSupply", ret); err != nil {
		return nil, fmt.Errorf("unpack totalSupply: %s", err.Error())
	}
	return unpackedRet.Value, nil
}
