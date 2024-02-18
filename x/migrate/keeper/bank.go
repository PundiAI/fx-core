package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	migratetypes "github.com/functionx/fx-core/v7/x/migrate/types"
)

type BankMigrate struct {
	bankKeeper migratetypes.BankKeeper
}

func NewBankMigrate(bk migratetypes.BankKeeper) MigrateI {
	return &BankMigrate{bankKeeper: bk}
}

func (m *BankMigrate) Validate(_ sdk.Context, _ codec.BinaryCodec, _ sdk.AccAddress, _ common.Address) error {
	return nil
}

func (m *BankMigrate) Execute(ctx sdk.Context, _ codec.BinaryCodec, from sdk.AccAddress, to common.Address) error {
	balances := m.bankKeeper.GetAllBalances(ctx, from)
	if balances.IsZero() {
		return nil
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		migratetypes.EventTypeMigrateBankSend,
		sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
	))
	return m.bankKeeper.SendCoins(ctx, from, to.Bytes(), balances)
}
