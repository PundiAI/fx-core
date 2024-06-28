package precompile

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

type CancelPendingBridgeCallMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewCancelPendingBridgeCallMethod(keeper *Keeper) *CancelPendingBridgeCallMethod {
	return &CancelPendingBridgeCallMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["cancelPendingBridgeCall"],
		Event:  crosschaintypes.GetABI().Events["CancelPendingBridgeCallEvent"],
	}
}

func (m *CancelPendingBridgeCallMethod) IsReadonly() bool {
	return false
}

func (m *CancelPendingBridgeCallMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *CancelPendingBridgeCallMethod) RequiredGas() uint64 {
	return 30_000
}

func (m *CancelPendingBridgeCallMethod) NewCancelPendingBridgeCallEvent(sender common.Address, chainName string, txId *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash()}, chainName, txId)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *CancelPendingBridgeCallMethod) Run(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
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

	conis, err := route.PrecompileCancelPendingBridgeCall(ctx, args.TxID.Uint64(), sender.Bytes())
	if err != nil {
		return nil, err
	}

	originAmount := conis.AmountOf(fxtypes.DefaultDenom)
	if originAmount.IsPositive() {
		// add refund + reward to sender in evm state db, because bank keeper add refund to sender
		evm.StateDB.AddBalance(sender, originAmount.BigInt())
	}

	data, topic, err := m.NewCancelPendingBridgeCallEvent(sender, args.Chain, args.TxID)
	if err != nil {
		return nil, err
	}
	EmitEvent(evm, data, topic)

	return m.PackOutput(true)
}

func (m *CancelPendingBridgeCallMethod) UnpackInput(data []byte) (*crosschaintypes.CancelPendingBridgeCallArgs, error) {
	args := new(crosschaintypes.CancelPendingBridgeCallArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *CancelPendingBridgeCallMethod) PackInput(args crosschaintypes.CancelPendingBridgeCallArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Chain, args.TxID)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *CancelPendingBridgeCallMethod) PackOutput(success bool) ([]byte, error) {
	pack, err := m.Method.Outputs.Pack(success)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
