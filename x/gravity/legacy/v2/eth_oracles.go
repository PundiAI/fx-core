package v2

import (
	"github.com/functionx/fx-core/v3/types"
)

func ethInitOracles(chainId string) []string {
	if chainId == types.MainnetChainId {
		return []string{}
	} else if chainId == types.TestnetChainId {
		return []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
}
