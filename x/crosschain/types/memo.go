package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

const MemoSendCallTo = "0000000000000000000000000000000000000000000000000000000000010000"

func IsMemoSendCallTo(memo []byte) bool {
	return len(memo) == 32 && bytes.Equal(memo, common.Hex2Bytes(MemoSendCallTo))
}
