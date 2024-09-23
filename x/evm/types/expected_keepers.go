package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/x/evm/types"
)

// AccountKeeper defines the expected interface needed to retrieve account info.
type AccountKeeper interface {
	types.AccountKeeper

	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}
