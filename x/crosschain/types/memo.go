package types

import (
	"github.com/ethereum/go-ethereum/common"
)

var MemoSendCallTo = common.HexToHash("0000000000000000000000000000000000000000000000000000000000010000")

func IsMemoSendCallTo(memo []byte) bool {
	return len(memo) == 32 && common.BytesToHash(memo) == MemoSendCallTo
}
