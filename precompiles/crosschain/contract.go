package crosschain

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/types"
	evmtypes "github.com/pundiai/fx-core/v8/x/evm/types"
)

var (
	crosschainABI     = contract.MustABIJson(contract.ICrosschainMetaData.ABI)
	crosschainAddress = common.HexToAddress(contract.CrosschainAddress)
)

type Contract struct {
	methods   []contract.PrecompileMethod
	govKeeper types.GovKeeper
}

func NewPrecompiledContract(
	bankKeeper types.BankKeeper,
	govKeeper types.GovKeeper,
	erc20Keeper types.Erc20Keeper,
	router *Router,
) *Contract {
	keeper := &Keeper{
		bankKeeper:  bankKeeper,
		erc20Keeper: erc20Keeper,
		router:      router,
	}
	return &Contract{
		govKeeper: govKeeper,
		methods: []contract.PrecompileMethod{
			NewBridgeCoinAmountMethod(keeper),
			NewHasOracleMethod(keeper),
			NewIsOracleOnlineMethod(keeper),
			NewGetERC20TokenMethod(keeper),

			NewCrosschainMethod(keeper),
			NewBridgeCallMethod(keeper),
			NewExecuteClaimMethod(keeper),
		},
	}
}

func (c *Contract) Address() common.Address {
	return crosschainAddress
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
			if err = c.govKeeper.CheckDisabledPrecompiles(stateDB.Context(), c.Address(), method.GetMethodId()); err != nil {
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
