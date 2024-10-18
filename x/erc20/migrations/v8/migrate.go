package v8

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
)

func (m Migrator) MigrateToken(ctx sdk.Context) error {
	// add FX bridge token
	if err := m.addToken(ctx, fxtypes.DefaultDenom, ""); err != nil {
		return err
	}

	mds := m.bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0 {
			continue
		}
		baseDenom := strings.ToLower(md.Symbol)
		// add other bridge/ibc token
		for _, alias := range md.DenomUnits[0].Aliases {
			if err := m.addToken(ctx, baseDenom, alias); err != nil {
				return err
			}
		}
		// add pundix/purse token
		if md.Base == baseDenom {
			continue
		}
		if err := m.addToken(ctx, baseDenom, md.Base); err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) addToken(
	ctx sdk.Context,
	base, alias string,
) error {
	if strings.HasPrefix(alias, ibctransfertypes.DenomPrefix+"/") {
		return m.addIBCToken(ctx, base, alias)
	}
	return m.addBridgeToken(ctx, base, alias)
}

func (m Migrator) addIBCToken(ctx sdk.Context, base, alias string) error {
	channel, found := getIBCDenomTrace(ctx, alias)
	if !found {
		return sdkerrors.ErrInvalidCoins.Wrapf("ibc denom hash not found: %s %s", base, alias)
	}
	ctx.Logger().Info("add ibc token", "base-denom", base, "alias", alias, "channel", channel)
	return m.keeper.AddIBCToken(ctx, base, channel, alias)
}

func (m Migrator) addBridgeToken(
	ctx sdk.Context,
	base, alias string,
) error {
	if getExcludeBridgeToken(ctx, alias) {
		return nil
	}
	for _, ck := range m.crosschainKeepers {
		canAddFxBridgeToken := base == fxtypes.DefaultDenom && ck.ModuleName() == ethtypes.ModuleName
		canAddBridgeToken := strings.HasPrefix(alias, ck.ModuleName())
		excludeModule := ck.ModuleName() != arbitrumtypes.ModuleName && ck.ModuleName() != optimismtypes.ModuleName
		if ctx.ChainID() == fxtypes.MainnetChainId {
			canAddBridgeToken = canAddBridgeToken && excludeModule
		}
		if !canAddFxBridgeToken && !canAddBridgeToken {
			continue
		}

		if alias == "" { // FX token
			alias = base
		}
		legacyBridgeToken, found := ck.LegacyGetDenomBridgeToken(ctx, alias)
		if !found {
			return sdkerrors.ErrKeyNotFound.Wrapf("module %s bridge token: %s", ck.ModuleName(), alias)
		}
		ctx.Logger().Info("add bridge token", "base-denom", base, "alias", alias, "module", ck.ModuleName(), "contract", legacyBridgeToken.Token)
		isNativeErc20 := LegacyIsNativeERC20(ctx, m.storeKey, m.cdc, base)
		if err := m.keeper.AddBridgeToken(ctx, base, ck.ModuleName(), legacyBridgeToken.Token, isNativeErc20); err != nil {
			return err
		}
		break
	}
	return nil
}
