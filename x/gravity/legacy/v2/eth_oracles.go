package v2

import (
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func ethInitOracles(chainId string) []string {
	if chainId == fxtypes.MainnetChainId {
		return []string{}
	} else if chainId == fxtypes.TestnetChainId {
		return []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
}
