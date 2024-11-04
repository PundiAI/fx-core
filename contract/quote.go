package contract

import (
	"math/big"
	"time"
)

func (i IBridgeFeeQuoteQuoteInfo) IsTimeout(blockTime time.Time) bool {
	return new(big.Int).Sub(i.Expiry, big.NewInt(blockTime.Unix())).Sign() <= 0
}
