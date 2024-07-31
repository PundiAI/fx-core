package precompile

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

type CancelSendToExternalMethod struct {
	*Keeper
	Method abi.Method
	Event  abi.Event
}

func NewCancelSendToExternalMethod(keeper *Keeper) *CancelSendToExternalMethod {
	return &CancelSendToExternalMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["cancelSendToExternal"],
		Event:  crosschaintypes.GetABI().Events["CancelSendToExternal"],
	}
}

func (m *CancelSendToExternalMethod) IsReadonly() bool {
	return false
}

func (m *CancelSendToExternalMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *CancelSendToExternalMethod) RequiredGas() uint64 {
	return 30_000
}

func (m *CancelSendToExternalMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if m.router == nil {
		return nil, errors.New("cross chain router is empty")
	}

	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	sender := contract.Caller()
	route, has := m.router.GetRoute(args.Chain)
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	// NOTE: must be get relation before cancel, cancel will delete it if relation exist

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		if _, err = route.PrecompileCancelSendToExternal(ctx, args.TxID.Uint64(), sender.Bytes()); err != nil {
			return err
		}

		data, topic, err := m.NewCancelSendToExternalEvent(sender, args.Chain, args.TxID)
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

		return nil
	}); err != nil {
		return nil, err
	}

	return m.PackOutput(true)
}

func (m *CancelSendToExternalMethod) NewCancelSendToExternalEvent(sender common.Address, chainName string, txId *big.Int) (data []byte, topic []common.Hash, err error) {
	return evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash()}, chainName, txId)
}

func (m *CancelSendToExternalMethod) PackInput(chainName string, txId *big.Int) ([]byte, error) {
	data, err := m.Method.Inputs.Pack(chainName, txId)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), data...), nil
}

func (m *CancelSendToExternalMethod) UnpackInput(data []byte) (*crosschaintypes.CancelSendToExternalArgs, error) {
	args := new(crosschaintypes.CancelSendToExternalArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *CancelSendToExternalMethod) PackOutput(success bool) ([]byte, error) {
	return m.Method.Outputs.Pack(success)
}
