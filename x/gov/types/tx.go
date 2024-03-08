package types

import v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

type MsgServerPro interface {
	MsgServer
	v1.MsgServer
}
