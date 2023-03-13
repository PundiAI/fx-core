package crosschain

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

type CallerRef struct {
	addr common.Address
}

func (cr CallerRef) Address() common.Address {
	return cr.addr
}

type ContractCall struct {
	evm      *vm.EVM
	caller   CallerRef
	contract common.Address
}

func NewContractCall(evm *vm.EVM, caller, contract common.Address) *ContractCall {
	return &ContractCall{
		evm:      evm,
		caller:   CallerRef{addr: caller},
		contract: contract,
	}
}

func (cc *ContractCall) ERC20Burn(amount *big.Int) error {
	data, err := fxtypes.GetERC20().ABI.Pack("burn", cc.caller.Address(), amount)
	if err != nil {
		return fmt.Errorf("pack burn: %s", err.Error())
	}
	_, _, err = cc.evm.Call(cc.caller, cc.contract, data, math.MaxInt64, big.NewInt(0))
	if err != nil {
		return fmt.Errorf("call burn: %s", err.Error())
	}
	return nil
}

func (cc *ContractCall) ERC20TransferFrom(from, to common.Address, amount *big.Int) error {
	data, err := fxtypes.GetERC20().ABI.Pack("transferFrom", from, to, amount)
	if err != nil {
		return fmt.Errorf("pack transferFrom: %s", err.Error())
	}
	ret, _, err := cc.evm.Call(cc.caller, cc.contract, data, math.MaxInt64, big.NewInt(0))
	if err != nil {
		return fmt.Errorf("call transferFrom: %s", err.Error())
	}
	var unpackedRet struct{ Value bool }
	if err := fxtypes.GetERC20().ABI.UnpackIntoInterface(&unpackedRet, "transferFrom", ret); err != nil {
		return fmt.Errorf("unpack transferFrom: %s", err.Error())
	}
	if !unpackedRet.Value {
		return errors.New("transferFrom failed")
	}
	return nil
}
