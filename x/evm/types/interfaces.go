package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

// EvmLogHandler defines the interface for evm log handler
type EvmLogHandler interface {
	// EventID Return the id of the log signature it handles
	EventID() common.Hash
	// Handle Process the log
	Handle(ctx sdk.Context, msg core.Message, log *ethtypes.Log) error
}

type MethodArgs interface {
	Validate() error
}

// ExtStateDB defines extra methods of statedb to support stateful precompiled contracts
type ExtStateDB interface {
	vm.StateDB
	ExecuteNativeAction(contract common.Address, converter statedb.EventConverter, action func(ctx sdk.Context) error) error
	Context() sdk.Context
}
