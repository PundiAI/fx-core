package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetGravityID returns the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a oracle set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) GetGravityID(ctx sdk.Context) string {
	return k.GetParams(ctx).GravityId
}

func (k Keeper) GetOracleDelegateThreshold(ctx sdk.Context) sdk.Coin {
	return k.GetParams(ctx).DelegateThreshold
}

func (k Keeper) GetOracleDelegateMultiple(ctx sdk.Context) int64 {
	return k.GetParams(ctx).DelegateMultiple
}

func (k Keeper) GetSlashFraction(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).SlashFraction
}

func (k Keeper) GetSignedWindow(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).SignedWindow
}

func (k Keeper) GetIbcTransferTimeoutHeight(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).IbcTransferTimeoutHeight
}

func (k Keeper) GetOracleSetUpdatePowerChangePercent(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).OracleSetUpdatePowerChangePercent
}

func (k Keeper) GetBridgeCallTimeout(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).BridgeCallTimeout
}

func (k Keeper) GetBridgeCallMaxGasLimit(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).BridgeCallMaxGasLimit
}
