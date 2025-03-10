package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func updateMetadataDesc(ctx sdk.Context, bankKeeper bankkeeper.Keeper) {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Description == "The cross chain token of the Function X" ||
			md.Description == "Cross chain token of Function X" ||
			md.Description == "The cross chain token of Function X" ||
			md.Description == "Cross chain token Function X" ||
			md.Description == "Function X coin token representation of 0x934B9f502dcED1eBf0594c7384Eb299bC3ca2bE6" {
			md.Description = "The crosschain token of the Pundi AIFX"
		}
		bankKeeper.SetDenomMetaData(ctx, md)
	}
}
