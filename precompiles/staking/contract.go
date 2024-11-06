package staking

import (
	"bytes"
	"errors"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/precompiles/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingkeeper "github.com/functionx/fx-core/v8/x/staking/keeper"
)

var (
	stakingAddress = common.HexToAddress(contract.StakingAddress)
	stakingABI     = contract.MustABIJson(contract.IStakingMetaData.ABI)
)

type Contract struct {
	methods   []contract.PrecompileMethod
	govKeeper types.GovKeeper
}

func NewPrecompiledContract(
	bankKeeper types.BankKeeper,
	stakingKeeper *fxstakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	stakingDenom string,
	govKeeper types.GovKeeper,
	slashingKeeper types.SlashingKeeper,
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

	return &Contract{
		methods: []contract.PrecompileMethod{
			NewAllowanceSharesMethod(keeper),
			NewDelegationMethod(keeper),
			NewDelegationRewardsMethod(keeper),

			NewApproveSharesMethod(keeper),
			NewTransferSharesMethod(keeper),
			NewTransferFromSharesMethod(keeper),
			NewWithdrawMethod(keeper),

			NewDelegateV2Method(keeper),
			NewRedelegateV2Method(keeper),
			NewUndelegateV2Method(keeper),

			NewSlashingInfoMethod(keeper),
			NewValidatorListMethod(keeper),
		},
		govKeeper: govKeeper,
	}
}

func (c *Contract) Address() common.Address {
	return stakingAddress
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
				return contract.PackRetErrV2(err)
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
