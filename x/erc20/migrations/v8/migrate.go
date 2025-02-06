package v8

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	arbitrumtypes "github.com/pundiai/fx-core/v8/x/arbitrum/types"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	optimismtypes "github.com/pundiai/fx-core/v8/x/optimism/types"
)

func (m Migrator) MigrateToken(ctx sdk.Context) error {
	// add FX bridge token
	if err := m.addToken(ctx, fxtypes.LegacyFXDenom, ""); err != nil {
		return err
	}

	mds := m.bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		// exclude FX and alias empty, except PUNDIX
		if md.Base == fxtypes.LegacyFXDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != "PUNDIX" {
			continue
		}

		newBaseDenom := md.Base
		if !strings.Contains(md.Base, strings.ToLower(md.Symbol)) {
			newBaseDenom = strings.ToLower(md.Symbol)
		}
		// add other bridge/ibc token
		for _, alias := range md.DenomUnits[0].Aliases {
			if err := m.addToken(ctx, newBaseDenom, alias); err != nil {
				return err
			}
		}
		// only add pundix/purse token
		if md.Base == newBaseDenom || strings.Contains(md.Base, newBaseDenom) {
			continue
		}
		if err := m.addToken(ctx, newBaseDenom, md.Base); err != nil {
			return err
		}

		// add purse bsc module bridge token
		if strings.HasPrefix(md.Base, ibctransfertypes.DenomPrefix+"/") {
			if err := m.addBscBridgePurse(ctx, newBaseDenom, md.Base); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m Migrator) addToken(ctx sdk.Context, base, alias string) error {
	if strings.HasPrefix(alias, ibctransfertypes.DenomPrefix+"/") {
		return m.addIBCToken(ctx, base, alias)
	}
	return m.addBridgeToken(ctx, base, alias)
}

func (m Migrator) addIBCToken(ctx sdk.Context, base, alias string) error {
	channel, found := GetIBCDenomTrace(ctx, alias)
	if !found {
		return sdkerrors.ErrInvalidCoins.Wrapf("ibc denom hash not found: %s %s", base, alias)
	}
	ctx.Logger().Info("add ibc token", "base-denom", base, "alias", alias, "channel", channel)
	return m.keeper.AddIBCToken(ctx, channel, base, alias)
}

func (m Migrator) addBridgeToken(ctx sdk.Context, base, alias string) error {
	if getExcludeBridgeToken(ctx, alias) {
		return nil
	}
	for _, ck := range m.crosschainKeepers {
		canAddFxBridgeToken := base == fxtypes.LegacyFXDenom && ck.ModuleName() == ethtypes.ModuleName

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
		erc20Token, err := m.keeper.GetERC20Token(ctx, base)
		if err != nil {
			return err
		}
		ctx.Logger().Info("add bridge token", "base-denom", base, "alias", alias, "module", ck.ModuleName(), "contract", legacyBridgeToken.Token)
		isNative := erc20Token.IsNativeERC20()

		// NOTE: purse bridge token can not mint/burn
		if base == "purse" {
			isNative = true
		}
		if err = m.keeper.AddBridgeToken(ctx, base, ck.ModuleName(), legacyBridgeToken.Token, isNative); err != nil {
			return err
		}
		break
	}
	return nil
}

func (m Migrator) addBscBridgePurse(ctx sdk.Context, newBaseDenom, base string) error {
	for _, ck := range m.crosschainKeepers {
		if ck.ModuleName() != bsctypes.ModuleName {
			continue
		}
		legacyBridgeToken, found := ck.LegacyGetDenomBridgeToken(ctx, base)
		if !found {
			return sdkerrors.ErrKeyNotFound.Wrapf("module %s bridge token: %s", ck.ModuleName(), base)
		}
		bridgeDenom := crosschaintypes.NewBridgeDenom(ck.ModuleName(), legacyBridgeToken.Token)
		ctx.Logger().Info("add bridge token", "base-denom", newBaseDenom, "alias", bridgeDenom, "module", ck.ModuleName(), "contract", legacyBridgeToken.Token)
		// NOTE: purse bridge token can not mint/burn
		if err := m.keeper.AddBridgeToken(ctx, newBaseDenom, ck.ModuleName(), legacyBridgeToken.Token, true); err != nil {
			return err
		}
	}
	return nil
}
