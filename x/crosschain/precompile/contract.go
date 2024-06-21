package precompile

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

type Contract struct {
	ctx               sdk.Context
	router            *Router
	bankKeeper        BankKeeper
	erc20Keeper       Erc20Keeper
	ibcTransferKeeper IBCTransferKeeper
	accountKeeper     AccountKeeper
}

func NewPrecompiledContract(
	ctx sdk.Context,
	bankKeeper BankKeeper,
	erc20Keeper Erc20Keeper,
	ibcTransferKeeper IBCTransferKeeper,
	accountKeeper AccountKeeper,
	router *Router,
) *Contract {
	return &Contract{
		ctx:               ctx,
		bankKeeper:        bankKeeper,
		erc20Keeper:       erc20Keeper,
		ibcTransferKeeper: ibcTransferKeeper,
		accountKeeper:     accountKeeper,
		router:            router,
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
	switch string(input[:4]) {
	case string(crosschaintypes.FIP20CrossChainMethod.ID):
		return crosschaintypes.FIP20CrossChainGas
	case string(crosschaintypes.CrossChainMethod.ID):
		return crosschaintypes.CrossChainGas
	case string(crosschaintypes.CancelSendToExternalMethod.ID):
		return crosschaintypes.CancelSendToExternalGas
	case string(crosschaintypes.IncreaseBridgeFeeMethod.ID):
		return crosschaintypes.IncreaseBridgeFeeGas
	case string(crosschaintypes.BridgeCoinAmountMethod.ID):
		return crosschaintypes.BridgeCoinAmountFeeGas
	case string(crosschaintypes.BridgeCallMethod.ID):
		return crosschaintypes.BridgeCallFeeGas
	default:
		return 0
	}
}

func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(contract.Input) <= 4 {
		return evmtypes.PackRetErrV2("invalid input")
	}

	cacheCtx, commit := c.ctx.CacheContext()
	snapshot := evm.StateDB.Snapshot()

	// parse input
	switch string(contract.Input[:4]) {
	case string(crosschaintypes.FIP20CrossChainMethod.ID):
		ret, err = c.FIP20CrossChain(cacheCtx, evm, contract, readonly)
	case string(crosschaintypes.CrossChainMethod.ID):
		ret, err = c.CrossChain(cacheCtx, evm, contract, readonly)
	case string(crosschaintypes.CancelSendToExternalMethod.ID):
		ret, err = c.CancelSendToExternal(cacheCtx, evm, contract, readonly)
	case string(crosschaintypes.IncreaseBridgeFeeMethod.ID):
		ret, err = c.IncreaseBridgeFee(cacheCtx, evm, contract, readonly)
	case string(crosschaintypes.BridgeCoinAmountMethod.ID):
		ret, err = c.BridgeCoinAmount(cacheCtx, evm, contract, readonly)
	case string(crosschaintypes.BridgeCallMethod.ID):
		ret, err = c.BridgeCall(cacheCtx, evm, contract, readonly)

	default:
		err = errors.New("unknown method")
	}

	if err != nil {
		// revert evm state
		evm.StateDB.RevertToSnapshot(snapshot)
		return evmtypes.PackRetErrV2(err.Error())
	}

	// commit and append events
	commit()

	return ret, nil
}

func (c *Contract) AddLog(evm *vm.EVM, event abi.Event, topics []common.Hash, args ...interface{}) error {
	data, newTopic, err := evmtypes.PackTopicData(event, topics, args...)
	if err != nil {
		return err
	}
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     c.Address(),
		Topics:      newTopic,
		Data:        data,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
	return nil
}

func (c *Contract) GetBlockGasLimit() uint64 {
	params := c.ctx.ConsensusParams()
	if params != nil && params.Block != nil && params.Block.MaxGas > 0 {
		return uint64(params.Block.MaxGas)
	}
	return 0
}
