package precompile

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type ApproveSharesMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewApproveSharesMethod(keeper *Keeper) *ApproveSharesMethod {
	return &ApproveSharesMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["approveShares"],
		Event:  fxstakingtypes.GetABI().Events["ApproveShares"],
	}
}

func (m *ApproveSharesMethod) IsReadonly() bool {
	return false
}

func (m *ApproveSharesMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *ApproveSharesMethod) RequiredGas() uint64 {
	return 10_000
}

func (m *ApproveSharesMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		owner := contract.Caller()
		m.stakingKeeper.SetAllowance(ctx, args.GetValidator(), owner.Bytes(), args.Spender.Bytes(), args.Shares)

		ApproveSharesEmitEvents(ctx, args.GetValidator(), owner.Bytes(), args.Spender.Bytes(), sdkmath.NewIntFromBigInt(args.Shares))

		data, topic, err := m.NewApproveSharesEvent(owner, args.Spender, args.Validator, args.Shares)
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

func (m *ApproveSharesMethod) NewApproveSharesEvent(owner, spender common.Address, validator string, shares *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{owner.Hash(), spender.Hash()}, validator, shares)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *ApproveSharesMethod) PackInput(args fxstakingtypes.ApproveSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Spender, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *ApproveSharesMethod) UnpackInput(data []byte) (*fxstakingtypes.ApproveSharesArgs, error) {
	args := new(fxstakingtypes.ApproveSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *ApproveSharesMethod) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m *ApproveSharesMethod) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *ApproveSharesMethod) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingApproveShares, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseApproveShares(*log)
}

func ApproveSharesEmitEvents(ctx sdk.Context, validator sdk.ValAddress, owner, spender sdk.AccAddress, shares sdkmath.Int) {
	if shares.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "approve_shares")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(shares.Int64()),
				[]metrics.Label{telemetry.NewLabel("validator", validator.String())},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			fxstakingtypes.EventTypeApproveShares,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeySpender, spender.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyShares, shares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
}
