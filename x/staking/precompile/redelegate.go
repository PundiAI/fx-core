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
)

type RedelegateV2Method struct {
	*Keeper
	RedelegateABI
}

func NewRedelegateV2Method(keeper *Keeper) *RedelegateV2Method {
	return &RedelegateV2Method{
		Keeper:        keeper,
		RedelegateABI: NewRedelegateV2ABI(),
	}
}

func (m *RedelegateV2Method) IsReadonly() bool {
	return false
}

func (m *RedelegateV2Method) GetMethodId() []byte {
	return m.Method.ID
}

func (m *RedelegateV2Method) RequiredGas() uint64 {
	return 60_000
}

func (m *RedelegateV2Method) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		resp, err := m.stakingMsgServer.BeginRedelegate(ctx, &stakingtypes.MsgBeginRedelegate{
			DelegatorAddress:    sdk.AccAddress(contract.Caller().Bytes()).String(),
			ValidatorSrcAddress: args.ValidatorSrc,
			ValidatorDstAddress: args.ValidatorDst,
			Amount:              m.NewStakingCoin(args.Amount),
		})
		if err != nil {
			return err
		}

		// add redelegate log
		data, topic, err := m.NewRedelegationEvent(contract.Caller(), args.ValidatorSrc, args.ValidatorDst, args.Amount, resp.CompletionTime.Unix())
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

type RedelegateABI struct {
	abi.Method
	abi.Event
}

func NewRedelegateV2ABI() RedelegateABI {
	return RedelegateABI{
		Method: stakingABI.Methods["redelegateV2"],
		Event:  stakingABI.Events["RedelegateV2"],
	}
}

func (m RedelegateABI) NewRedelegationEvent(sender common.Address, validatorSrc, validatorDst string, amount *big.Int, completionTime int64) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validatorSrc, validatorDst, amount, big.NewInt(completionTime))
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m RedelegateABI) PackInput(args fxcontract.RedelegateV2Args) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.ValidatorSrc, args.ValidatorDst, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m RedelegateABI) UnpackInput(data []byte) (*fxcontract.RedelegateV2Args, error) {
	args := new(fxcontract.RedelegateV2Args)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m RedelegateABI) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m RedelegateABI) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m RedelegateABI) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingRedelegateV2, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseRedelegateV2(*log)
}
