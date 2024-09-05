package v4

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/functionx/fx-core/v8/x/gov/types"
)

func MigrateFXParams(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, params govv1.Params) error {
	store := ctx.KVStore(storeKey)
	fxParamsStore := prefix.NewStore(store, types.FxBaseParamsKeyPrefix)
	iter := fxParamsStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var fxParams types.Params
		err := cdc.Unmarshal(iter.Value(), &fxParams)
		if err != nil {
			return err
		}

		fxParams.MinInitialDepositRatio = params.MinInitialDepositRatio
		fxParams.BurnVoteQuorum = params.BurnVoteQuorum
		fxParams.BurnProposalDepositPrevote = params.BurnProposalDepositPrevote
		fxParams.BurnVoteVeto = params.BurnVoteVeto
		paramsBz, err := cdc.Marshal(&fxParams)
		if err != nil {
			return err
		}
		fxParamsStore.Set(iter.Key(), paramsBz)
	}
	return nil
}
