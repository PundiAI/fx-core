package precompile

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type UndelegateV2Method struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewUndelegateV2Method(keeper *Keeper) *UndelegateV2Method {
	return &UndelegateV2Method{
		Keeper: keeper,
		Method: stakingABI.Methods["undelegateV2"],
		Event:  stakingABI.Events["UndelegateV2"],
	}
}

func (m *UndelegateV2Method) IsReadonly() bool {
	return false
}

func (m *UndelegateV2Method) GetMethodId() []byte {
	return m.Method.ID
}

func (m *UndelegateV2Method) RequiredGas() uint64 {
	return 45_000
}

func (m *UndelegateV2Method) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		resp, err := m.stakingMsgServer.Undelegate(ctx, &stakingtypes.MsgUndelegate{
			DelegatorAddress: sdk.AccAddress(contract.Caller().Bytes()).String(),
			ValidatorAddress: args.Validator,
			Amount:           m.NewStakingCoin(args.Amount),
		})
		if err != nil {
			return err
		}

		// add undelegate log
		data, topic, err := m.NewUndelegateEvent(contract.Caller(), args.Validator, args.Amount, resp.CompletionTime.Unix())
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

func (m *UndelegateV2Method) NewUndelegateEvent(sender common.Address, validator string, amount *big.Int, completionTime int64) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, amount, big.NewInt(completionTime))
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *UndelegateV2Method) PackInput(args fxstakingtypes.UndelegateV2Args) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *UndelegateV2Method) UnpackInput(data []byte) (*fxstakingtypes.UndelegateV2Args, error) {
	args := new(fxstakingtypes.UndelegateV2Args)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *UndelegateV2Method) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m *UndelegateV2Method) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *UndelegateV2Method) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingUndelegateV2, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseUndelegateV2(*log)
}
