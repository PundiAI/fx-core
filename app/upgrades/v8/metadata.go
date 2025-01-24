package v8

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func updateMetadata(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)

	removeMetadata := make([]string, 0, 2)
	for _, md := range mds {
		if md.Base == fxtypes.LegacyFXDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != pundixSymbol {
			continue
		}
		// remove alias
		md.DenomUnits[0].Aliases = []string{}

		newBase := strings.ToLower(md.Symbol)
		// update pundix/purse base denom
		if md.Base != newBase && !strings.Contains(md.Base, newBase) && !strings.HasPrefix(md.Display, ibctransfertypes.ModuleName+"/"+ibcchanneltypes.ChannelPrefix) {
			removeMetadata = append(removeMetadata, md.Base)

			md.Base = newBase
			md.Display = newBase
			md.DenomUnits[0].Denom = newBase
		}

		bankKeeper.SetDenomMetaData(ctx, md)
	}

	bk, ok := bankKeeper.(bankkeeper.BaseKeeper)
	if !ok {
		return errors.New("bank keeper not implement bank.BaseKeeper")
	}
	for _, base := range removeMetadata {
		if !bankKeeper.HasDenomMetaData(ctx, base) {
			continue
		}
		ctx.Logger().Info("remove metadata", "base", base, "module", "upgrade")
		if err := bk.BaseViewKeeper.DenomMetadata.Remove(ctx, base); err != nil {
			return err
		}
	}
	return nil
}

func migrateMetadataFXToPundiAI(ctx sdk.Context, keeper bankkeeper.Keeper) error {
	// add pundiai metadata
	metadata := fxtypes.NewDefaultMetadata()
	keeper.SetDenomMetaData(ctx, metadata)

	// remove FX metadata
	bk, ok := keeper.(bankkeeper.BaseKeeper)
	if !ok {
		return errors.New("bank keeper not implement bank.BaseKeeper")
	}
	return bk.BaseViewKeeper.DenomMetadata.Remove(ctx, fxtypes.LegacyFXDenom)
}

func migrateMetadataDisplay(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Display != md.Base || len(md.DenomUnits) <= 1 {
			continue
		}
		for _, dus := range md.DenomUnits {
			if dus.Denom != md.Base {
				md.Display = dus.Denom
				break
			}
		}
		if err := md.Validate(); err != nil {
			return err
		}
		bankKeeper.SetDenomMetaData(ctx, md)
	}
	return nil
}
