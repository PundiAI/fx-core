package contract

import (
	"time"
)

func (i IBridgeFeeQuoteQuoteInfo) IsTimeout(blockTime time.Time) bool {
	return i.Expiry-uint64(blockTime.Unix()) <= 0
}

func (i IBridgeFeeQuoteQuoteInfo) GetChainName() string {
	return Byte32ToString(i.ChainName)
}

func (i IBridgeFeeQuoteQuoteInfo) GetTokenName() string {
	return Byte32ToString(i.TokenName)
}
