package types

import (
	"strings"
)

type FxTarget struct {
	isIBC         bool
	target        string
	Prefix        string
	SourcePort    string
	SourceChannel string
}

func ParseFxTarget(targetStr string) FxTarget {
	// px/transfer/channel-0
	ibcData := strings.Split(targetStr, "/")
	if len(ibcData) < 3 {
		return FxTarget{}
	}
	return FxTarget{
		isIBC:         true,
		target:        targetStr,
		Prefix:        ibcData[0],
		SourcePort:    ibcData[1],
		SourceChannel: ibcData[2],
	}
}

func (i FxTarget) GetTarget() string {
	if i.isIBC {
		return "ibc"
	}
	return i.target
}

func (i FxTarget) IsIBC() bool {
	return i.isIBC
}
