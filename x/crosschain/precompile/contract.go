package precompile

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/contract"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

type Contract struct {
	methods   []contract.PrecompileMethod
	govKeeper GovKeeper
}

func NewPrecompiledContract(
	bankKeeper BankKeeper,
	erc20Keeper Erc20Keeper,
	ibcTransferKeeper IBCTransferKeeper,
	accountKeeper AccountKeeper,
	govKeeper GovKeeper,
	router *Router,
) *Contract {
	keeper := &Keeper{
		bankKeeper:        bankKeeper,
		erc20Keeper:       erc20Keeper,
		ibcTransferKeeper: ibcTransferKeeper,
		accountKeeper:     accountKeeper,
		router:            router,
	}
	return &Contract{
		govKeeper: govKeeper,
		methods: []contract.PrecompileMethod{
			NewBridgeCoinAmountMethod(keeper),

			NewCancelSendToExternalMethod(keeper),
			NewIncreaseBridgeFeeMethod(keeper),
			NewFIP20CrossChainMethod(keeper),
			NewCrossChainMethod(keeper),
			NewBridgeCallMethod(keeper),
			NewCancelPendingBridgeCallMethod(keeper),
			NewAddPendingPoolRewardsMethod(keeper),
			NewExecuteClaimMethod(keeper),
		},
	}
}

func (c *Contract) Address() common.Address {
	return crosschaintypes.GetAddress()
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

func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(contract.Input) <= 4 {
		return evmtypes.PackRetErrV2(errors.New("invalid input"))
	}

	for _, method := range c.methods {
		if bytes.Equal(method.GetMethodId(), contract.Input[:4]) {
			if readonly && !method.IsReadonly() {
				return evmtypes.PackRetErrV2(errors.New("write protection"))
			}

			stateDB := evm.StateDB.(evmtypes.ExtStateDB)
			if err = c.govKeeper.CheckDisabledPrecompiles(stateDB.CacheContext(), c.Address(), method.GetMethodId()); err != nil {
				return evmtypes.PackRetErrV2(err)
			}

			ret, err = method.Run(evm, contract)
			if err != nil {
				return evmtypes.PackRetErrV2(err)
			}
			return ret, nil
		}
	}
	return evmtypes.PackRetErrV2(errors.New("unknown method"))
}

func EmitEvent(evm *vm.EVM, data []byte, topics []common.Hash) {
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     crosschaintypes.GetAddress(),
		Topics:      topics,
		Data:        data,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
}
