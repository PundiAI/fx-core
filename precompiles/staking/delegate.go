package staking

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type DelegateV2Method struct {
	*Keeper
	DelegateV2ABI
}

func NewDelegateV2Method(keeper *Keeper) *DelegateV2Method {
	return &DelegateV2Method{
		Keeper:        keeper,
		DelegateV2ABI: NewDelegateV2ABI(),
	}
}

func (m *DelegateV2Method) IsReadonly() bool {
	return false
}

func (m *DelegateV2Method) GetMethodId() []byte {
	return m.Method.ID
}

func (m *DelegateV2Method) RequiredGas() uint64 {
	return 40_000
}

func (m *DelegateV2Method) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if contract.Value().Sign() != 0 {
		return nil, errors.New("msg.value must be zero")
	}

	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		if _, err = m.stakingMsgServer.Delegate(ctx, &stakingtypes.MsgDelegate{
			DelegatorAddress: sdk.AccAddress(contract.Caller().Bytes()).String(),
			ValidatorAddress: args.Validator,
			Amount:           m.NewStakingCoin(args.Amount),
		}); err != nil {
			return err
		}

		// add delegate log
		data, topic, err := m.NewDelegateEvent(contract.Caller(), args.Validator, args.Amount)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, stakingAddress, data, topic)
		return nil
	}); err != nil {
		return nil, err
	}
	return m.PackOutput(true)
}

type DelegateV2ABI struct {
	abi.Method
	abi.Event
}

func NewDelegateV2ABI() DelegateV2ABI {
	return DelegateV2ABI{
		Method: stakingABI.Methods["delegateV2"],
		Event:  stakingABI.Events["DelegateV2"],
	}
}

func (m DelegateV2ABI) NewDelegateEvent(sender common.Address, validator string, amount *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, amount)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m DelegateV2ABI) PackInput(args fxcontract.DelegateV2Args) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m DelegateV2ABI) UnpackInput(data []byte) (*fxcontract.DelegateV2Args, error) {
	args := new(fxcontract.DelegateV2Args)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m DelegateV2ABI) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m DelegateV2ABI) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m DelegateV2ABI) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingDelegateV2, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseDelegateV2(*log)
}
