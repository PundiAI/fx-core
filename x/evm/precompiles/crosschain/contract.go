package crosschain

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/evm/types"
)

type Contract struct {
	ctx               sdk.Context
	evm               *vm.EVM
	router            *fxtypes.Router
	bankKeeper        BankKeeper
	evmKeeper         EvmKeeper
	erc20Keeper       Erc20Keeper
	ibcTransferKeeper IBCTransferKeeper
}

func NewPrecompiledContract(
	ctx sdk.Context,
	evm *vm.EVM,
	bankKeeper BankKeeper,
	evmKeeper EvmKeeper,
	erc20Keeper Erc20Keeper,
	ibcTransferKeeper IBCTransferKeeper,
	router *fxtypes.Router,
) *Contract {
	return &Contract{
		ctx:               ctx,
		evm:               evm,
		bankKeeper:        bankKeeper,
		evmKeeper:         evmKeeper,
		erc20Keeper:       erc20Keeper,
		ibcTransferKeeper: ibcTransferKeeper,
		router:            router,
	}
}

func (c *Contract) Address() common.Address {
	return precompileAddress
}

func (c *Contract) IsStateful() bool {
	return true
}

func (c *Contract) RequiredGas(input []byte) uint64 {
	if len(input) <= 4 {
		return 0
	}
	switch string(input[:4]) {
	case string(FIP20CrossChainMethod.ID):
		return FIP20CrossChainGas
	case string(CrossChainMethod.ID):
		return CrossChainGas
	case string(CancelSendToExternalMethod.ID):
		return CancelSendToExternalGas
	default:
		return 0
	}
}

func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(contract.Input) <= 4 {
		return types.PackRetError("invalid input")
	}

	cacheCtx, commit := c.ctx.CacheContext()
	snapshot := evm.StateDB.Snapshot()

	// parse input
	switch string(contract.Input[:4]) {
	case string(FIP20CrossChainMethod.ID):
		ret, err = c.FIP20CrossChain(cacheCtx, evm, contract, readonly)
	case string(CrossChainMethod.ID):
		ret, err = c.CrossChain(cacheCtx, evm, contract, readonly)
	case string(CancelSendToExternalMethod.ID):
		ret, err = c.CancelSendToExternal(cacheCtx, evm, contract, readonly)
	default:
		err = errors.New("unknown method")
	}

	if err != nil {
		// revert evm state
		evm.StateDB.RevertToSnapshot(snapshot)
		return types.PackRetError(err.Error())
	}

	// commit and append events
	commit()
	c.ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

	return ret, nil
}
