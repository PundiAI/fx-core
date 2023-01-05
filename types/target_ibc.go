package types

import (
	"strings"
)

type TargetIBC struct {
	Prefix        string
	SourcePort    string
	SourceChannel string
}

func ParseTargetIBC(targetIbc string) (TargetIBC, bool) {
	// px/transfer/channel-0
	ibcData := strings.Split(targetIbc, "/")
	if len(ibcData) < 3 {
		return TargetIBC{}, false
	}
	return TargetIBC{
		Prefix:        ibcData[0],
		SourcePort:    ibcData[1],
		SourceChannel: ibcData[2],
	}, true
}
