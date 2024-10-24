package precompile

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

type BridgeCallMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewBridgeCallMethod(keeper *Keeper) *BridgeCallMethod {
	return &BridgeCallMethod{
		Keeper: keeper,
		Method: crosschainABI.Methods["bridgeCall"],
		Event:  crosschainABI.Events["BridgeCallEvent"],
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

	sender := contract.Caller()
	value := contract.Value()

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		baseCoins := make([]sdk.Coin, 0, len(args.Tokens)+1)
		originTokenAmount := sdkmath.ZeroInt()
		if value.Sign() > 0 {
			baseCoin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(value))
			if err = m.bankKeeper.SendCoins(ctx, crosschainAddress.Bytes(), sender.Bytes(), sdk.NewCoins(baseCoin)); err != nil {
				return err
			}
			baseCoins = append(baseCoins, baseCoin)
			originTokenAmount = baseCoin.Amount
		}
		fxTarget, err := crosschaintypes.ParseFxTarget(args.DstChain)
		if err != nil {
			return err
		}
		crosschainKeeper, ok := m.router.GetRoute(fxTarget.GetModuleName())
		if !ok {
			return errors.New("invalid router")
		}

		for i, token := range args.Tokens {
			baseCoin, err := m.EvmTokenToBaseCoin(ctx, evm, crosschainKeeper, sender, token, args.Amounts[i])
			if err != nil {
				return err
			}
			baseCoins = append(baseCoins, baseCoin)
		}
		nonce, err := crosschainKeeper.BridgeCallBaseCoin(ctx, sender, args.Refund, args.To,
			baseCoins, args.Data, args.Memo, fxTarget, originTokenAmount)
		if err != nil {
			return err
		}
		eventNonce := new(big.Int).SetUint64(nonce)
		data, topic, err := m.NewBridgeCallEvent(args, sender, evm.Origin, eventNonce)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, crosschainAddress, data, topic)

		result, err = m.PackOutput(eventNonce)
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
