package precompile

import (
	"errors"
	"fmt"
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

type IncreaseBridgeFeeMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewIncreaseBridgeFeeMethod(keeper *Keeper) *IncreaseBridgeFeeMethod {
	return &IncreaseBridgeFeeMethod{
		Keeper: keeper,
		Method: crosschaintypes.GetABI().Methods["increaseBridgeFee"],
		Event:  crosschaintypes.GetABI().Events["IncreaseBridgeFee"],
	}
}

func (m *IncreaseBridgeFeeMethod) IsReadonly() bool {
	return false
}

func (m *IncreaseBridgeFeeMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *IncreaseBridgeFeeMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *IncreaseBridgeFeeMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if m.router == nil {
		return nil, errors.New("cross chain router empty")
	}

	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	fxTarget := fxtypes.ParseFxTarget(args.Chain)
	route, has := m.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		value := contract.Value()
		sender := contract.Caller()
		totalCoin := sdk.Coin{}
		if value.Cmp(big.NewInt(0)) == 1 && fxcontract.IsZeroEthAddress(args.Token) {
			if args.Fee.Cmp(value) != 0 {
				return errors.New("add bridge fee not equal msg.value")
			}
			totalCoin, err = m.handlerOriginToken(ctx, evm, sender, args.Fee)
			if err != nil {
				return err
			}
		} else {
			totalCoin, err = m.handlerERC20Token(ctx, evm, sender, args.Token, args.Fee)
			if err != nil {
				return err
			}
		}

		// convert token to bridge fee token
		feeCoin := sdk.NewCoin(totalCoin.Denom, sdkmath.NewIntFromBigInt(args.Fee))
		addBridgeFee, err := m.erc20Keeper.ConvertDenomToTarget(ctx, sender.Bytes(), feeCoin, fxTarget)
		if err != nil {
			return err
		}

		if err = route.AddUnbatchedTxBridgeFee(ctx, args.TxID.Uint64(), sender.Bytes(), addBridgeFee); err != nil {
			return err
		}

		data, topic, err := m.NewIncreaseBridgeFeeEvent(args, sender)
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

func (m *IncreaseBridgeFeeMethod) NewIncreaseBridgeFeeEvent(args *crosschaintypes.IncreaseBridgeFeeArgs, sender common.Address) (data []byte, topic []common.Hash, err error) {
	return evmtypes.PackTopicData(m.Event, []common.Hash{sender.Hash(), args.Token.Hash()}, args.Chain, args.TxID, args.Fee)
}

func (m *IncreaseBridgeFeeMethod) PackInput(chainName string, txId *big.Int, token common.Address, fee *big.Int) ([]byte, error) {
	data, err := m.Method.Inputs.Pack(chainName, txId, token, fee)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), data...), nil
}

func (m *IncreaseBridgeFeeMethod) UnpackInput(data []byte) (*crosschaintypes.IncreaseBridgeFeeArgs, error) {
	args := new(crosschaintypes.IncreaseBridgeFeeArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *IncreaseBridgeFeeMethod) PackOutput(success bool) ([]byte, error) {
	return m.Method.Outputs.Pack(success)
}
