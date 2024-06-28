package precompile

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/legacy"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

type CrossChainMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewCrossChainMethod(keeper *Keeper) *CrossChainMethod {
	return &CrossChainMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["crossChain"],
		Event:  crosschaintypes.GetABI().Events["CrossChain"],
	}
}

func (m *CrossChainMethod) IsReadonly() bool {
	return false
}

func (m *CrossChainMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *CrossChainMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *CrossChainMethod) Run(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	value := contract.Value()
	sender := contract.Caller()

	originToken := false
	totalCoin := sdk.Coin{}

	// cross-chain origin token
	if value.Cmp(big.NewInt(0)) == 1 && fxcontract.IsZeroEthAddress(args.Token) {
		totalAmount := big.NewInt(0).Add(args.Amount, args.Fee)
		if totalAmount.Cmp(value) != 0 {
			return nil, errors.New("amount + fee not equal msg.value")
		}

		totalCoin, err = m.handlerOriginToken(ctx, evm, sender, totalAmount)
		if err != nil {
			return nil, err
		}

		// origin token flag is true when cross chain evm denom
		originToken = true
	} else {
		totalCoin, err = m.handlerERC20Token(ctx, evm, sender, args.Token, big.NewInt(0).Add(args.Amount, args.Fee))
		if err != nil {
			return nil, err
		}
	}

	fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target))
	amountCoin := sdk.NewCoin(totalCoin.Denom, sdkmath.NewIntFromBigInt(args.Amount))
	feeCoin := sdk.NewCoin(totalCoin.Denom, sdkmath.NewIntFromBigInt(args.Fee))

	if err = m.handlerCrossChain(ctx, sender.Bytes(), args.Receipt, amountCoin, feeCoin, fxTarget, args.Memo, originToken); err != nil {
		return nil, err
	}

	data, topic, err := m.NewCrossChainEvent(sender, args.Token, amountCoin.Denom, args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
	if err != nil {
		return nil, err
	}
	EmitEvent(evm, data, topic)

	legacy.CrossChainEvents(ctx, sender, args.Token, args.Receipt, fxtypes.Byte32ToString(args.Target),
		amountCoin.Denom, args.Memo, args.Amount, args.Fee)

	return m.PackOutput(true)
}

func (m *CrossChainMethod) NewCrossChainEvent(sender common.Address, token common.Address, denom, receipt string, amount, fee *big.Int, target [32]byte, memo string) (data []byte, topic []common.Hash, err error) {
	data, topic, err = evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash(), token.Hash()}, denom, receipt, amount, fee, target, memo)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *CrossChainMethod) UnpackInput(data []byte) (*crosschaintypes.CrossChainArgs, error) {
	args := new(crosschaintypes.CrossChainArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *CrossChainMethod) PackInput(args crosschaintypes.CrossChainArgs) ([]byte, error) {
	data, err := m.Method.Inputs.Pack(args.Token, args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), data...), nil
}

func (m *CrossChainMethod) PackOutput(success bool) ([]byte, error) {
	pack, err := m.Method.Outputs.Pack(success)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
