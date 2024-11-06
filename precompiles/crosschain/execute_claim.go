package crosschain

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

var _ fxcontract.PrecompileMethod = (*ExecuteClaimMethod)(nil)

type ExecuteClaimMethod struct {
	*Keeper
	ExecuteClaimABI
}

func NewExecuteClaimMethod(keeper *Keeper) *ExecuteClaimMethod {
	return &ExecuteClaimMethod{
		Keeper:          keeper,
		ExecuteClaimABI: NewExecuteClaimABI(),
	}
}

func (m *ExecuteClaimMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *ExecuteClaimMethod) RequiredGas() uint64 {
	return 50_000
}

func (m *ExecuteClaimMethod) IsReadonly() bool {
	return false
}

func (m *ExecuteClaimMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if m.router == nil {
		return nil, errors.New("bridge call router is empty")
	}

	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	crosschainKeeper, has := m.router.GetRoute(args.Chain)
	if !has {
		return nil, errors.New("chain not support")
	}
	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		if err = crosschainKeeper.ExecuteClaim(ctx, args.EventNonce.Uint64()); err != nil {
			return err
		}
		data, topic, err := m.NewExecuteClaimEvent(contract.Caller(), args.EventNonce, args.Chain)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, crosschainAddress, data, topic)
		return nil
	}); err != nil {
		return nil, err
	}

	return m.PackOutput(true)
}

type ExecuteClaimABI struct {
	abi.Method
	abi.Event
}

func NewExecuteClaimABI() ExecuteClaimABI {
	return ExecuteClaimABI{
		Method: crosschainABI.Methods["executeClaim"],
		Event:  crosschainABI.Events["ExecuteClaimEvent"],
	}
}

func (m ExecuteClaimABI) NewExecuteClaimEvent(sender common.Address, eventNonce *big.Int, dstChain string) (data []byte, topic []common.Hash, err error) {
	return evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash()}, eventNonce, dstChain)
}

func (m ExecuteClaimABI) PackInput(args fxcontract.ExecuteClaimArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Chain, args.EventNonce)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m ExecuteClaimABI) UnpackInput(data []byte) (*fxcontract.ExecuteClaimArgs, error) {
	args := new(fxcontract.ExecuteClaimArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m ExecuteClaimABI) PackOutput(success bool) ([]byte, error) {
	return m.Method.Outputs.Pack(success)
}

func (m ExecuteClaimABI) UnpackOutput(data []byte) (bool, error) {
	success, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return success[0].(bool), nil
}
