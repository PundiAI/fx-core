package precompile

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type TransferSharesMethod struct {
	*Keeper
	TransferSharesABI
}

func NewTransferSharesMethod(keeper *Keeper) *TransferSharesMethod {
	return &TransferSharesMethod{
		Keeper:            keeper,
		TransferSharesABI: NewTransferSharesABI(),
	}
}

func (m *TransferSharesMethod) IsReadonly() bool {
	return false
}

func (m *TransferSharesMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *TransferSharesMethod) RequiredGas() uint64 {
	return 50_000
}

func (m *TransferSharesMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	stateDB := evm.StateDB.(types.ExtStateDB)
	var result []byte
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		token, reward, err := m.handlerTransferShares(ctx, evm, valAddr, contract.Caller(), args.To, args.Shares)
		if err != nil {
			return err
		}

		data, topic, err := m.NewTransferShareEvent(contract.Caller(), args.To, valAddr.String(), args.Shares, token)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, stakingAddress, data, topic)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
				sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(contract.Caller().Bytes()).String()),
			),
		)
		result, err = m.PackOutput(token, reward)
		return err
	}); err != nil {
		return nil, err
	}
	return result, nil
}

type TransferSharesABI struct {
	transferShareABI
}

func NewTransferSharesABI() TransferSharesABI {
	return TransferSharesABI{
		transferShareABI: transferShareABI{
			Method: stakingABI.Methods["transferShares"],
			Event:  stakingABI.Events["TransferShares"],
		},
	}
}

func (m TransferSharesABI) PackInput(args fxstakingtypes.TransferSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.To, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m TransferSharesABI) UnpackInput(data []byte) (*fxstakingtypes.TransferSharesArgs, error) {
	args := new(fxstakingtypes.TransferSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

type TransferFromSharesMethod struct {
	*Keeper
	TransferFromSharesABI
}

func NewTransferFromSharesMethod(keeper *Keeper) *TransferFromSharesMethod {
	return &TransferFromSharesMethod{
		Keeper:                keeper,
		TransferFromSharesABI: NewTransferFromSharesABI(),
	}
}

func (m *TransferFromSharesMethod) IsReadonly() bool {
	return false
}

func (m *TransferFromSharesMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *TransferFromSharesMethod) RequiredGas() uint64 {
	return 60_000
}

func (m *TransferFromSharesMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		valAddr := args.GetValidator()
		spender := contract.Caller()
		if err := m.decrementAllowance(ctx, valAddr, args.From.Bytes(), spender.Bytes(), args.Shares); err != nil {
			return err
		}
		token, reward, err := m.handlerTransferShares(ctx, evm, valAddr, args.From, args.To, args.Shares)
		if err != nil {
			return err
		}

		data, topic, err := m.NewTransferShareEvent(args.From, args.To, valAddr.String(), args.Shares, token)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, stakingAddress, data, topic)

		result, err = m.PackOutput(token, reward)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
				sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(contract.Caller().Bytes()).String()),
			),
		)
		return nil
	})
	return result, err
}

type TransferFromSharesABI struct {
	transferShareABI
}

func NewTransferFromSharesABI() TransferFromSharesABI {
	return TransferFromSharesABI{
		transferShareABI: transferShareABI{
			Method: stakingABI.Methods["transferFromShares"],
			Event:  stakingABI.Events["TransferShares"],
		},
	}
}

func (m TransferFromSharesABI) PackInput(args fxstakingtypes.TransferFromSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.From, args.To, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m TransferFromSharesABI) UnpackInput(data []byte) (*fxstakingtypes.TransferFromSharesArgs, error) {
	args := new(fxstakingtypes.TransferFromSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

type transferShareABI struct {
	abi.Method
	abi.Event
}

func (m transferShareABI) PackOutput(amount, reward *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount, reward)
}

func (m transferShareABI) UnpackOutput(data []byte) (*big.Int, *big.Int, error) {
	unpacks, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, nil, err
	}
	return unpacks[0].(*big.Int), unpacks[1].(*big.Int), nil
}

func (m transferShareABI) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingTransferShares, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseTransferShares(*log)
}

func (m transferShareABI) NewTransferShareEvent(sender, to common.Address, validator string, shares, amount *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash(), to.Hash()}, validator, shares, amount)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}
