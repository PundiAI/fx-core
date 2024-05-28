package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

func PackBridgeCallback(sender common.Address, receiver common.Address, tokens []common.Address, amounts []*big.Int, data, memo []byte) ([]byte, error) {
	args, err := contract.GetBridgeCallBridgeCallback().Pack("bridgeCallback",
		sender,
		receiver,
		tokens,
		amounts,
		data,
		memo,
	)
	return args, err
}
