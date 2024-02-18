package crosschain

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxserverconfig "github.com/functionx/fx-core/v7/server/config"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

type CallerRef struct {
	addr common.Address
}

func (cr CallerRef) Address() common.Address {
	return cr.addr
}

type ContractCall struct {
	ctx      sdk.Context
	evm      *vm.EVM
	caller   CallerRef
	contract common.Address
	maxGas   uint64
}

func NewContractCall(ctx sdk.Context, evm *vm.EVM, caller, contract common.Address) *ContractCall {
	cc := &ContractCall{
		ctx:      ctx,
		evm:      evm,
		caller:   CallerRef{addr: caller},
		contract: contract,
		maxGas:   fxserverconfig.DefaultGasCap,
	}
	params := ctx.ConsensusParams()
	if params != nil && params.Block != nil && params.Block.MaxGas > 0 {
		cc.maxGas = uint64(params.Block.MaxGas)
	}
	return cc
}

func (cc *ContractCall) ERC20Burn(amount *big.Int) error {
	data, err := fxtypes.GetFIP20().ABI.Pack("burn", cc.caller.Address(), amount)
	if err != nil {
		return fmt.Errorf("pack burn: %s", err.Error())
	}
	_, _, err = cc.evm.Call(cc.caller, cc.contract, data, cc.maxGas, big.NewInt(0))
	if err != nil {
		return fmt.Errorf("call burn: %s", err.Error())
	}
	return nil
}

func (cc *ContractCall) ERC20TransferFrom(from, to common.Address, amount *big.Int) error {
	data, err := fxtypes.GetFIP20().ABI.Pack("transferFrom", from, to, amount)
	if err != nil {
		return fmt.Errorf("pack transferFrom: %s", err.Error())
	}
	ret, _, err := cc.evm.Call(cc.caller, cc.contract, data, cc.maxGas, big.NewInt(0))
	if err != nil {
		return fmt.Errorf("call transferFrom: %s", err.Error())
	}
	var unpackedRet struct{ Value bool }
	if err := fxtypes.GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "transferFrom", ret); err != nil {
		return fmt.Errorf("unpack transferFrom: %s", err.Error())
	}
	if !unpackedRet.Value {
		return errors.New("transferFrom failed")
	}
	return nil
}

func (cc *ContractCall) ERC20TotalSupply() (*big.Int, error) {
	data, err := fxtypes.GetFIP20().ABI.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("pack totalSupply: %s", err.Error())
	}
	ret, _, err := cc.evm.StaticCall(cc.caller, cc.contract, data, cc.maxGas)
	if err != nil {
		return nil, fmt.Errorf("StaticCall totalSupply: %s", err.Error())
	}
	var unpackedRet struct{ Value *big.Int }
	if err := fxtypes.GetFIP20().ABI.UnpackIntoInterface(&unpackedRet, "totalSupply", ret); err != nil {
		return nil, fmt.Errorf("unpack totalSupply: %s", err.Error())
	}
	return unpackedRet.Value, nil
}
