package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	v3types "github.com/functionx/fx-core/v3/x/evm/legacy/v3/types"
)

// MigrateRejectUnprotectedTx used by ethermint version before v0.17.0
func MigrateRejectUnprotectedTx(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) error {
	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(evmtypes.ModuleName), '/'))

	bzR := paramsStore.Get(v3types.ParamStoreKeyRejectUnprotectedTx)
	var rejectUnprotectedTx bool
	if err := legacyAmino.UnmarshalJSON(bzR, &rejectUnprotectedTx); err != nil {
		return fmt.Errorf("legacy amino unmarshal %s: %s", err.Error(), v3types.ParamStoreKeyRejectUnprotectedTx)
	}

	allowUnprotectedTxs := !rejectUnprotectedTx
	bzA, err := legacyAmino.MarshalJSON(allowUnprotectedTxs)
	if err != nil {
		return fmt.Errorf("legacy amino marshal %s: %s", err.Error(), v3types.ParamStoreKeyRejectUnprotectedTx)
	}

	ctx.Logger().Info("migrate params", "module", evmtypes.ModuleName,
		"from", fmt.Sprintf("%s:%v", v3types.ParamStoreKeyRejectUnprotectedTx, rejectUnprotectedTx),
		"to", fmt.Sprintf("%s:%v", evmtypes.ParamStoreKeyAllowUnprotectedTxs, allowUnprotectedTxs))

	paramsStore.Delete(v3types.ParamStoreKeyRejectUnprotectedTx)
	paramsStore.Set(evmtypes.ParamStoreKeyAllowUnprotectedTxs, bzA)
	return nil
}
