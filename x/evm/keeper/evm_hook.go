package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v4/x/evm/types"
)

// LogProcessEvmHook is an evm hook that convert specific contract logs into native module calls
type LogProcessEvmHook struct {
	handlers map[common.Hash]types.EvmLogHandler
}

// NewLogProcessEvmHook todo: remove unused code
func NewLogProcessEvmHook(handlers ...types.EvmLogHandler) LogProcessEvmHook {
	handlerMap := make(map[common.Hash]types.EvmLogHandler)
	for _, h := range handlers {
		handlerMap[h.EventID()] = h
	}
	return LogProcessEvmHook{handlerMap}
}

// PostTxProcessing delegate the call to underlying hooks
func (lh LogProcessEvmHook) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		if handler, ok := lh.handlers[log.Topics[0]]; ok {
			if err := handler.Handle(ctx, msg, log); err != nil {
				return err
			}
		}
	}
	return nil
}
