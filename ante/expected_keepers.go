package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Govkeeper interface {
	GetDisabledMsgs(ctx sdk.Context) []string
}
