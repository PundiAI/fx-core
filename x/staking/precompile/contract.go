package precompile

import (
	"bytes"
	"errors"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/contract"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingkeeper "github.com/functionx/fx-core/v8/x/staking/keeper"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type Contract struct {
	methods   []contract.PrecompileMethod
	govKeeper GovKeeper
}

func NewPrecompiledContract(
	bankKeeper BankKeeper,
	stakingKeeper *fxstakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	stakingDenom string,
	govKeeper GovKeeper,
	slashingKeeper SlashingKeeper,
) *Contract {
	keeper := &Keeper{
		bankKeeper:       bankKeeper,
		distrKeeper:      distrKeeper,
		distrMsgServer:   distrkeeper.NewMsgServerImpl(distrKeeper),
		stakingKeeper:    stakingKeeper,
		stakingMsgServer: stakingkeeper.NewMsgServerImpl(stakingKeeper.Keeper),
		stakingDenom:     stakingDenom,
		slashingKeeper:   slashingKeeper,
	}

	delegateV2 := NewDelegateV2Method(keeper)
	redelegateV2 := NewRedelegateV2Method(keeper)
	undelegateV2 := NewUndelegateV2Method(keeper)
	slashingInfo := NewSlashingInfoMethod(keeper)
	validatorList := NewValidatorListMethod(keeper)
	return &Contract{
		methods: []contract.PrecompileMethod{
			NewAllowanceSharesMethod(keeper),
			NewDelegationMethod(keeper),
			NewDelegationRewardsMethod(keeper),

			NewApproveSharesMethod(keeper),
			NewTransferSharesMethod(keeper),
			NewTransferFromSharesMethod(keeper),
			NewWithdrawMethod(keeper),

			delegateV2,
			redelegateV2,
			undelegateV2,

			slashingInfo,
			validatorList,
		},
		govKeeper: govKeeper,
	}
}

func (c *Contract) Address() common.Address {
	return fxstakingtypes.GetAddress()
}

func (c *Contract) IsStateful() bool {
	return true
}

func (c *Contract) RequiredGas(input []byte) uint64 {
	if len(input) <= 4 {
		return 0
	}
	for _, method := range c.methods {
		if bytes.Equal(method.GetMethodId(), input[:4]) {
			return method.RequiredGas()
		}
	}
	return 0
}

func (c *Contract) Run(evm *vm.EVM, vmContract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(vmContract.Input) <= 4 {
		return contract.PackRetErrV2(errors.New("invalid input"))
	}

	for _, method := range c.methods {
		if bytes.Equal(method.GetMethodId(), vmContract.Input[:4]) {
			if readonly && !method.IsReadonly() {
				return contract.PackRetErrV2(errors.New("write protection"))
			}

			stateDB := evm.StateDB.(evmtypes.ExtStateDB)
			if err = c.govKeeper.CheckDisabledPrecompiles(stateDB.Context(), c.Address(), vmContract.Input[:4]); err != nil {
				return contract.PackRetError(err)
			}

			ret, err = method.Run(evm, vmContract)
			if err != nil {
				return contract.PackRetErrV2(err)
			}
			return ret, nil
		}
	}
	return contract.PackRetErrV2(errors.New("unknown method"))
}

func EmitEvent(evm *vm.EVM, data []byte, topics []common.Hash) {
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     fxstakingtypes.GetAddress(),
		Topics:      topics,
		Data:        data,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
}
