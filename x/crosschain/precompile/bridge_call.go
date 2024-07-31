package precompile

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

type BridgeCallMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewBridgeCallMethod(keeper *Keeper) *BridgeCallMethod {
	return &BridgeCallMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["bridgeCall"],
		Event:  crosschaintypes.GetABI().Events["BridgeCallEvent"],
	}
}

func (m *BridgeCallMethod) IsReadonly() bool {
	return false
}

func (m *BridgeCallMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *BridgeCallMethod) RequiredGas() uint64 {
	return 50_000
}

func (m *BridgeCallMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if m.router == nil {
		return nil, errors.New("bridge call router is empty")
	}

	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	route, has := m.router.GetRoute(args.DstChain)
	if !has {
		return nil, errors.New("invalid dstChain")
	}
	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		sender := contract.Caller()
		coins := make([]sdk.Coin, 0, len(args.Tokens)+1)
		value := contract.Value()
		if value.Cmp(big.NewInt(0)) == 1 {
			totalCoin, err := m.handlerOriginToken(ctx, evm, sender, value)
			if err != nil {
				return err
			}
			coins = append(coins, totalCoin)
		}
		for i, token := range args.Tokens {
			coin, err := m.handlerERC20Token(ctx, evm, sender, token, args.Amounts[i])
			if err != nil {
				return err
			}
			coins = append(coins, coin)
		}

		bridgeCallNonce, err := route.PrecompileBridgeCall(ctx, sender, args.Refund, coins, args.To, args.Data, args.Memo)
		if err != nil {
			return err
		}

		emitNonce := new(big.Int).SetUint64(bridgeCallNonce)
		data, topic, err := m.NewBridgeCallEvent(args, sender, evm.Origin, emitNonce)
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

		result, err = m.PackOutput(emitNonce)
		return err
	})
	return result, err
}

func (m *BridgeCallMethod) NewBridgeCallEvent(args *crosschaintypes.BridgeCallArgs, sender, origin common.Address, eventNonce *big.Int) (data []byte, topic []common.Hash, err error) {
	return evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash(), args.Refund.Hash(), args.To.Hash()}, origin, args.Value, eventNonce, args.DstChain, args.Tokens, args.Amounts, args.Data, args.Memo)
}

func (m *BridgeCallMethod) UnpackInput(data []byte) (*crosschaintypes.BridgeCallArgs, error) {
	args := new(crosschaintypes.BridgeCallArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *BridgeCallMethod) PackOutput(nonceNonce *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(nonceNonce)
}

func (m *BridgeCallMethod) PackInput(args crosschaintypes.BridgeCallArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.DstChain, args.Refund, args.Tokens, args.Amounts, args.To, args.Data, args.Value, args.Memo)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}
