package precompile

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/contract"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

var _ contract.PrecompileMethod = (*ExecuteClaimMethod)(nil)

type ExecuteClaimMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewExecuteClaimMethod(keeper *Keeper) *ExecuteClaimMethod {
	return &ExecuteClaimMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["executeClaim"],
		Event:  crosschaintypes.GetABI().Events["ExecuteClaimEvent"],
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

	crosschainKeeper, has := m.router.GetRoute(args.DstChain)
	if !has {
		return nil, errors.New("chain not support")
	}
	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		if err = crosschainKeeper.ExecuteClaim(ctx, args.EventNonce.Uint64()); err != nil {
			return err
		}
		data, topic, err := m.NewExecuteClaimEvent(contract.Caller(), args.EventNonce.Uint64(), args.DstChain)
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return m.PackOutput(true)
}

func (m *ExecuteClaimMethod) NewExecuteClaimEvent(sender common.Address, eventNonce uint64, dstChain string) (data []byte, topic []common.Hash, err error) {
	data, topic, err = evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash()}, eventNonce, dstChain)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *ExecuteClaimMethod) PackInput(args crosschaintypes.ExecuteClaimArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.DstChain, args.EventNonce)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *ExecuteClaimMethod) UnpackInput(data []byte) (*crosschaintypes.ExecuteClaimArgs, error) {
	args := new(crosschaintypes.ExecuteClaimArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *ExecuteClaimMethod) PackOutput(success bool) ([]byte, error) {
	pack, err := m.Method.Outputs.Pack(success)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (m *ExecuteClaimMethod) UnpackOutput(data []byte) (bool, error) {
	success, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return success[0].(bool), nil
}
